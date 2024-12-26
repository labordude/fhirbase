package cmd

import (
	// "context"
	"fmt"
	"os"

	// "os/signal"
	// "syscall"

	app "github.com/labordude/fhirbase/app"
	"github.com/labordude/fhirbase/shared"
	"github.com/spf13/cobra"
)

var AvailableSchemas = []string{
	"1.0.2", "1.1.0", "1.4.0",
	"1.6.0", "1.8.0", "3.0.1",
	"3.2.0", "3.3.0", "4.0.0",
}
var Version = "0.1."
var BuildDate = "2024-12-23"

const logo = ` 
(        )  (    (                   (
)\ )  ( /(  )\ ) )\ )   (     (      )\ )
(()/(  )\())(()/((()/( ( )\    )\    (()/( (
/(_))((_)\  /(_))/(_)))((_)((((_)(   /(_)))\
(_))_| _((_)(_)) (_)) ((_)_  )\ _ )\ (_)) ((_)
| |_  | || ||_ _|| _ \ | _ ) (_)_\(_)/ __|| __|
| __| | __ | | | |   / | _ \  / _ \  \__ \| _|
|_|   |_||_||___||_|_\ |___/ /_/ \_\ |___/|___|`

var data = struct {
	Logo      string
	Version   string
	BuildDate string
}{
	Logo:      logo,
	Version:   Version,
	BuildDate: BuildDate,
}

var welcomeMessage = fmt.Sprintf(`
%s
Version: %s
Build Date: %s
`, data.Logo, data.Version, data.BuildDate)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fhirbase",
	Short: "Welcome to fhirbase",

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		err := initializeConfig(cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		app.Start(shared.RootView)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Create a cancellable context
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel() // Ensure resources are released on program exit

	// // Add signal handling to cancel context on interrupt
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-signalChan
	// 	cancel()
	// }()

	// // Pass context to rootCmd
	// rootCmd.SetContext(ctx)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {

	// viper.Debug()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringP("fhir", "f", "4.0.0", "FHIR version to use")
	rootCmd.PersistentFlags().StringP("db", "d", "", "Database name to use")
	rootCmd.PersistentFlags().Uint64P("port", "p", 5432, "Port to listen on")
	rootCmd.PersistentFlags().StringP("host", "", "localhost", "Host to listen on")
	rootCmd.PersistentFlags().StringP("username", "U", "postgres", "Username to use")
	rootCmd.PersistentFlags().StringP("password", "W", "", "Password to use")
	rootCmd.PersistentFlags().StringP("sslmode", "s", "disable", "SSL mode to use")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
