package app

import "github.com/charmbracelet/bubbles/key"

type configKeymap struct {
	Apply key.Binding
	Save  key.Binding
	Next  key.Binding
	Exit  key.Binding
}

var defaultConfigKeymap = configKeymap{
	Apply: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "apply"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save to .lazymigrate"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab", "up", "down"),
		key.WithHelp("tab/↑/↓", "edit next"),
	),
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "exit"),
	),
}

func (m configKeymap) ShortHelp() []key.Binding {
	return []key.Binding{m.Exit, m.Next, m.Apply, m.Save}
}

func (m configKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}
