package app

import "github.com/charmbracelet/bubbles/key"

type migrateKeymap struct {
	Exit   key.Binding
	Up     key.Binding
	Down   key.Binding
	Goto   key.Binding
	Config key.Binding
	New    key.Binding
}

var defaultMigrateKeymap = migrateKeymap{
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "exit"),
	),
	Up: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "down"),
	),
	Goto: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "goto"),
	),
	Config: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "config"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
}

func (m migrateKeymap) ShortHelp() []key.Binding {
	return []key.Binding{m.Exit, m.Up, m.Down, m.Goto, m.Config, m.New}
}

func (m migrateKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}
