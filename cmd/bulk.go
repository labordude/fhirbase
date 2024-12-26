/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"math"

	"strconv"

	"net/http"
	urlPkg "net/url"
	"os"
	"path"
	"regexp"

	"io"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var numdl int = 5
var acceptHeader string = "application/fhir+json"

// bulkCmd represents the bulk command
var bulkCmd = &cobra.Command{
	Use:     "bulkget",
	Short:   "Downloads FHIR data from Bulk Data API endpoint and saves NDJSON files on local filesystem into specified directory",
	Example: "fhirbase bulkget [--numdl=10] http://some-fhir-server.com/fhir/Patient/$everything /output/dir/",
	Long: `
Downloads FHIR data from Bulk Data API endpoint and saves results into
specific directory on a local filesystem.

NDJSON files generated by remote server will be downloaded in
parallel and you can specify number of threads with "--numdl" flag.

To mitigate differences between Bulk Data API implementations, there
is an "--accept-header" option which sets the value for "Accept"
header. Most likely you won't need to set it, but if Bulk Data server
rejects queries because of "Accept" header value, consider explicitly
set it to something it expects.
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Not enough arguments")
			return
		}
		ctx := cmd.Context()
		err := BulkGetCommand(ctx, args)
		if err != nil {
			fmt.Printf("Error downloading files: %v", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(bulkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	bulkCmd.PersistentFlags().IntVarP(&numdl, "numdl", "n", 5, "Number of parallel downloads")
	bulkCmd.PersistentFlags().StringVar(&acceptHeader, "accept-header", "a", "application/fhir+json")
	viper.BindPFlag("numdl", bulkCmd.PersistentFlags().Lookup("numdl"))
	viper.BindPFlag("accept-header", bulkCmd.PersistentFlags().Lookup("accept-header"))
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bulkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var matchNonDigits, _ = regexp.Compile("[^\\d]")

func getBulkDataFiles(pingURL string, client *http.Client) ([]string, error) {
	maxRetries := 10
	baseDelay := 10 * time.Second
	var fileURLs []string
	fmt.Println("Waiting for Bulk Data API server to prepare files...")

	for i := 0; i < maxRetries; i++ {

		req, err := http.NewRequest("GET", pingURL, nil)
		if err != nil {
			return nil, fmt.Errorf("error while creating HTTP request: %w", err)
		}

		resp, err := client.Do(req)

		if err != nil {
			return nil, fmt.Errorf("error while performing polling request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 202 {
			fmt.Println("still waiting...")
			// read Retry-After header or use an exponential backoff
			retryAfter := resp.Header.Get("Retry-After")

			if retryAfter != "" {
				delay, err := strconv.Atoi(retryAfter)
				fmt.Printf("Retry in %d seconds...\r", delay)
				if err != nil {
					time.Sleep(time.Duration(delay) * time.Second)
					// show countdown
					continue
				}
			}

			// exponential backoff
			delay := time.Duration(math.Pow(2, float64(i))) * baseDelay
			time.Sleep(delay)
			continue
		} else if resp.StatusCode == 200 {
			//parse the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("error reading response body: %w", err)
			}
			fileURLs, err = parseFileURLs(body)
			if err != nil {
				return nil, fmt.Errorf("error parsing response body: %w", err)
			}
			break

		} else {
			respBody, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("expected 202 response, got %d; response body is: %s", resp.StatusCode, respBody)
		}
	}

	if fileURLs == nil {
		return nil, fmt.Errorf("Bulk Data API server did not return any files")
	}

	return fileURLs, nil
}

func parseFileURLs(body []byte) ([]string, error) {
	// fmt.Println("Bulk Data API server is ready", body)
	iter := jsoniter.ConfigDefault.BorrowIterator(body)
	defer jsoniter.ConfigDefault.ReturnIterator(iter)

	obj := iter.Read()

	if obj == nil {
		return nil, fmt.Errorf("cannot parse JSON from Bulk Data API server")
	}

	fileURLs := make([]string, 0)
	objMap, ok := obj.(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("expecting JSON object at the top level")
	}

	output := objMap["output"]

	if output == nil {
		return nil, fmt.Errorf("expecting to have 'output' attribute")
	}

	outputArr, ok := output.([]interface{})

	if !ok {
		return nil, fmt.Errorf("'output' attribute is not an JSON Array")
	}

	for _, v := range outputArr {
		item, ok := v.(map[string]interface{})

		if !ok {
			return nil, fmt.Errorf("got non-object in 'output' array")
		}

		url := item["url"]

		if url == nil {
			return nil, fmt.Errorf("cannot get 'url' attribute in item of 'output' array")
		}

		urlString, ok := url.(string)

		if !ok {
			return nil, fmt.Errorf("'url' attribute is not a string")
		}

		fileURLs = append(fileURLs, urlString)
	}

	return fileURLs, nil
}

func stripURL(url string, length int) string {
	if len(url) < length {
		return strings.Repeat(" ", length-len(url)) + url
	}

	return "..." + url[len(url)-length-3:]
}

func ensureDirectoryExists(dir string) error {
	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}
	return nil
}

func startDlWorker(jobs chan string, results chan interface{}, targetDir string, wg *sync.WaitGroup) {
	defer wg.Done() // Signal when the worker is done

	client := &http.Client{}

	for url := range jobs {
		parsedURL, err := urlPkg.Parse(url)
		fileName := path.Base(parsedURL.EscapedPath())

		if err != nil {
			results <- fmt.Errorf("cannot parse URL: %v", err)
			continue
		}

		// Ensure the file has the .ndjson extension
		if !strings.HasSuffix(fileName, ".ndjson") {
			fileName += ".ndjson"
		}

		// Create a file in the target directory
		targetPath := path.Join(targetDir, fileName)
		targetFile, err := os.Create(targetPath)

		if err != nil {
			results <- fmt.Errorf("cannot create file: %v", err)
			continue
		}

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Add("Accept-Encoding", "gzip")
		resp, err := client.Do(req)

		if err != nil {
			results <- fmt.Errorf("cannot download %s: %v", url, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			results <- fmt.Errorf("got non-200 response while downloading %s", url)
			continue
		}

		// Write response body to the file
		_, err = io.Copy(targetFile, resp.Body)
		if err != nil {
			results <- fmt.Errorf("error while downloading %s: %v", targetFile.Name(), err)
			continue
		}

		results <- targetFile
	}

}

func downloadAllFiles(fileURLs []string, numWorkers uint, targetDir string) ([]*os.File, error) {
	jobs := make(chan string, len(fileURLs))
	results := make(chan interface{}, len(fileURLs))
	files := make([]*os.File, 0)

	var wg sync.WaitGroup

	// Start workers
	for i := uint(0); i < numWorkers; i++ {
		wg.Add(1)
		go startDlWorker(jobs, results, targetDir, &wg)
	}

	// Send jobs to the channel
	for _, url := range fileURLs {
		jobs <- url
	}
	close(jobs) // Close jobs channel to signal workers no more jobs

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results) // Close results channel once all workers are done
	}()

	// Collect results
	for res := range results {
		switch r := res.(type) {
		case error:
			fmt.Printf("Got an error while downloading file: %s\n", r.Error())
		case *os.File:
			files = append(files, r)
		default:
			fmt.Printf("Got result of unknown type: %v\n", r)
		}
	}

	fmt.Printf("Finished downloading, got %d files\n", len(files))
	return files, nil
}

func getBulkData(url string, numWorkers uint, acceptHdr string, targetDir string) ([]*os.File, error) {
	client := &http.Client{}

	if strings.Contains(url, "$export-poll-status") {
		fileURLs, err := getBulkDataFiles(url, client)

		if err != nil {
			return nil, fmt.Errorf("error while getting files from Bulk Data API server: %v", err)
		}

		return downloadAllFiles(fileURLs, numWorkers, targetDir)
	}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("error while creating request to Bulk Data API server: %v", err)
	}

	// add headers for async response
	req.Header.Add("Prefer", "respond-async")
	req.Header.Add("Accept", acceptHdr)
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error while pinging Bulk Data API server: %v", err)
	}

	defer resp.Body.Close()

	// check if we got 20x response
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("expected 20x response, got %d; response body is: %s", resp.StatusCode, respBody)
	}

	pingURL := resp.Header.Get("Content-Location")

	if len(pingURL) == 0 {
		return nil, fmt.Errorf("No Content-Location header was returned by Bulk Data API server")
	}

	fileURLs, err := getBulkDataFiles(pingURL, client)

	if err != nil {
		return nil, fmt.Errorf("error while getting files from Bulk Data API server: %v", err)
	}

	return downloadAllFiles(fileURLs, numWorkers, targetDir)
}

// BulkGetCommand loads data from Bulk Data Endpoint and saves it to local filesystem
func BulkGetCommand(ctx context.Context, args []string) error {

	numWorkers := uint(viper.GetInt("numdl"))
	acceptHdr := viper.GetString("accept-header")
	bulkURL := args[0]
	destPath := args[1]

	// Ensure the destination directory exists before downloading files
	err := ensureDirectoryExists(destPath)

	if err != nil {
		return err
	}

	fileHndlrs, err := getBulkData(bulkURL, numWorkers, acceptHdr, destPath)

	if err != nil {
		return err
	}

	for _, f := range fileHndlrs {
		fn := f.Name()
		fbn := path.Base(fn)

		err := os.Rename(f.Name(), path.Join(destPath, fbn))

		if err != nil {
			fmt.Printf("Error moving %s to %s: %v", f.Name(), path.Join(destPath, fbn), err)
		}
	}

	return nil

}
