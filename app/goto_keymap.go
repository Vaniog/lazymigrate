package app

import "github.com/charmbracelet/bubbles/key"

type gotoKeymap struct {
	Apply key.Binding
	Exit  key.Binding
}

var defaultGotoKeymap = gotoKeymap{
	Apply: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "apply"),
	),
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "exit"),
	),
}

func (m gotoKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.Exit,
		key.NewBinding(key.WithKeys("up", "down"), key.WithHelp("↑/↓", "scroll")),
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "complete")),
		m.Apply,
	}
}

func (m gotoKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}
