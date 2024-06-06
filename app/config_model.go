package app

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type configModel struct {
	source   textinput.Model
	database textinput.Model

	// which field is edit now
	edit string

	help   help.Model
	keymap configKeymap

	// to return
	prevM migrateModel
}

func newConfigModel(mm migrateModel) configModel {
	cm := configModel{
		source:   textinput.New(),
		database: textinput.New(),
		help:     help.New(),
		keymap:   defaultConfigKeymap,
	}

	cm.source.Prompt = "  Source       "
	cm.source.Placeholder = "path"
	cm.source.SetValue(mm.sourceFile)
	cm.database.Prompt = "  Database     "
	cm.database.Placeholder = "driver://..."
	cm.database.SetValue(mm.databaseUrl)
	cm.database.Focus()

	cm.edit = "database"

	cm.prevM = mm

	return cm
}

func (cm configModel) Init() tea.Cmd {
	return textinput.Blink
}

func (cm configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, cm.keymap.Next):
			if cm.edit == "source" {
				cm.edit = "database"
				cm.source.Blur()
				cm.database.Focus()
			} else {
				cm.edit = "source"
				cm.database.Blur()
				cm.source.Focus()
			}
			return cm, nil
		case key.Matches(msg, cm.keymap.Apply):
			mm := newMigrateModel(cm.source.Value(), cm.database.Value())
			return mm, mm.Init()
		case key.Matches(msg, cm.keymap.Save):
			cm.save()
			mm := newMigrateModel(cm.source.Value(), cm.database.Value())
			return mm, mm.Init()
		case key.Matches(msg, cm.keymap.Exit):
			return cm.prevM, nil
		}

	case tea.WindowSizeMsg:
		cm.help.Width = msg.Width
	}

	if cm.edit == "source" {
		var cmd tea.Cmd
		cm.source, cmd = cm.source.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		var cmd tea.Cmd
		cm.database, cmd = cm.database.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cm, tea.Batch(cmds...)
}

func (cm configModel) save() {
	_ = os.Remove(".lazymigrate")
	f, err := os.Create(".lazymigrate")
	if err != nil {
		return
	}
	cfg := fmt.Sprintf(
		"LAZYMIGRATE_URL=%s\nLAZYMIGRATE_SOURCE=%s\n",
		cm.database.Value(),
		cm.source.Value(),
	)
	_, _ = f.Write([]byte(cfg))
}

func (cm configModel) View() string {

	return fmt.Sprintf(lipgloss.JoinVertical(
		0,
		">> Edit Config",
		cm.database.View(),
		cm.source.View(),
		"",
		cm.help.View(cm.keymap),
	))
}
