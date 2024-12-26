package app

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/labordude/fhirbase/app/router"

	shared "github.com/labordude/fhirbase/shared"
	"github.com/spf13/viper"
)

var appTitle = "Fhirbase\n\n"

/* Initializes the Main model and starts the TUI application */
func Start(view string) {
	if viper.GetBool("debug.messages") {
		f, err := tea.LogToFile(shared.FhirbaseOptions.Debug.FilePath, "debug")
		if err != nil {
			fmt.Printf("Error setting up logging: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	r := router.NewRouterModel(router.NewRouterModelOpts{
		Quit: shared.FhirbaseOptions.Keys.Quit,
		Views: map[string]tea.Model{
			shared.RootView: NewMainModel(),
			// shared.LoadDbView:    NewMainModel(),
			// shared.BulkGetView:   NewMainModel(),
			// shared.TransformView: NewMainModel(),
			// shared.WebServerView: NewMainModel(),
			// shared.UpdateView:    NewMainModel(),
			// shared.ConfigView:    NewMainModel(),
			// shared.InitDbView:    NewMainModel(),
		},
		View: view,
	})

	p := tea.NewProgram(r)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting BubbleTea: %v\n", err)
		os.Exit(1)
	}
}
