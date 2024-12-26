package app

import (
	"fmt"

	help "github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/labordude/fhirbase/app/router"
	"github.com/labordude/fhirbase/internal/logger"
	"github.com/labordude/fhirbase/pkg/components"
	shared "github.com/labordude/fhirbase/shared"
)

type MainModel struct {
	help     help.Model
	selector tea.Model
	err      error
}

func NewMainModel() tea.Model {

	selectorComponent := components.NewSelectorModel(components.NewSelectorModelOpts{
		Filter: components.FilterOpts{
			Placeholder: "Search...",
		},
		Options: []components.SelectorOption{
			{Label: "Home", Value: shared.RootView},
			{Label: "Check Configuration", Value: shared.ConfigView},
			{Label: "Initialize Database", Value: shared.InitDbView},
			{Label: "Load Database", Value: shared.LoadDbView},
			{Label: "Bulk Download", Value: shared.BulkGetView},
			{Label: "Transform Resources", Value: shared.TransformView},
			{Label: "Start Web Server", Value: shared.WebServerView},
			{Label: "Update software", Value: shared.UpdateView},
			{Label: "Quit", Value: "quit"},
		},
	},
	)
	h := help.New()
	return MainModel{
		help:     h,
		selector: selectorComponent,
	}
}

func (m MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.selector.Init())
	return tea.Batch(cmds...)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.DebugMsg(m, msg)
	if m.err != nil {
		return m, tea.Quit
	}

	var cmds = []tea.Cmd{}

	/* Handle possible commands by selector */
	var cmd tea.Cmd
	m.selector, cmd = m.selector.Update(msg)
	cmds = append(cmds, cmd)

	/* All other events */
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
	case components.SelectMsg:
		if msg.Option.Value == "quit" {
			return m, tea.Quit
		}
		return m, router.ChangeView(msg.Option.Value)
	case tea.KeyMsg:
		switch msg.String() {
		case shared.FhirbaseOptions.Keys.Help:
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	base := appTitle
	base += m.selector.View()
	base += fmt.Sprintf("\n\n%s", m.help.View(newKeys()))
	return base
}
