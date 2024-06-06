package app

import "github.com/charmbracelet/bubbles/key"

type createMigrationKeymap struct {
	Save key.Binding
	Exit key.Binding
}

var defaultCreateMigrationKeymap = createMigrationKeymap{
	Save: key.NewBinding(
		key.WithKeys("ctrl+s", "enter"),
		key.WithHelp("ctrl+s/enter", "save"),
	),
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "exit"),
	),
}

func (m createMigrationKeymap) ShortHelp() []key.Binding {
	return []key.Binding{m.Exit, m.Save}
}

func (m createMigrationKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}
