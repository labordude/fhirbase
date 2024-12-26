package app

import (
	"github.com/charmbracelet/bubbles/key"

	"github.com/labordude/fhirbase/shared"
)

type keyMap struct {
	Quit   key.Binding
	Up     key.Binding
	Down   key.Binding
	Back   key.Binding
	Toggle key.Binding
	Select key.Binding
	Help   key.Binding
}

func newKeys() keyMap {
	return keyMap{
		Quit: key.NewBinding(
			key.WithKeys(shared.FhirbaseOptions.Keys.Quit),
			key.WithHelp(shared.FhirbaseOptions.Keys.Quit, "quit"),
		),
		Up: key.NewBinding(
			key.WithKeys(shared.FhirbaseOptions.Keys.Up),
			key.WithHelp(shared.FhirbaseOptions.Keys.Up, "up"),
		),
		Down: key.NewBinding(
			key.WithKeys(shared.FhirbaseOptions.Keys.Down),
			key.WithHelp(shared.FhirbaseOptions.Keys.Down, "down"),
		),
		Back: key.NewBinding(
			key.WithKeys(shared.FhirbaseOptions.Keys.Back),
			key.WithHelp(shared.FhirbaseOptions.Keys.Back, "back"),
		),
		Toggle: key.NewBinding(
			key.WithKeys(shared.FhirbaseOptions.Keys.Toggle),
			key.WithHelp(shared.FhirbaseOptions.Keys.Toggle, "toggle"),
		),
		Select: key.NewBinding(
			key.WithKeys(shared.FhirbaseOptions.Keys.Select),
			key.WithHelp(shared.FhirbaseOptions.Keys.Select, "select"),
		),
		Help: key.NewBinding(
			key.WithKeys(shared.FhirbaseOptions.Keys.Help),
			key.WithHelp(shared.FhirbaseOptions.Keys.Help, "help"),
		),
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit,
			k.Up,
			k.Down,
			k.Back,
			k.Toggle,
			k.Select,
			k.Help},
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Quit,
		k.Help,
	}
}
