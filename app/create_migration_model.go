package app

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"path"
	"time"
)

type createMigrationModel struct {
	name   textinput.Model
	status string

	// prevM for returning to state
	prevM migrateModel

	help   help.Model
	keymap createMigrationKeymap
}

func newCreateMigrationModel(mm migrateModel) createMigrationModel {
	cm := createMigrationModel{
		prevM:  mm,
		keymap: defaultCreateMigrationKeymap,
		help:   help.New(),
	}

	cm.name = textinput.New()
	cm.name.Prompt = "> "
	cm.name.Placeholder = "create_table"
	cm.name.Focus()
	return cm
}

func (cmm createMigrationModel) Init() tea.Cmd {
	return textinput.Blink
}

func (cmm createMigrationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, cmm.keymap.Save):
			if cmm.status == "" {
				// input stage
				err := cmm.create()
				if err != nil {
					cmm.status = fmt.Sprintf("error: %s", err)
				} else {
					cmm.status = "created!"
				}
				cmm.name.Blur()
				return cmm, nil
			} else {
				// status stage
				return cmm.prevM, nil
			}
		case key.Matches(msg, cmm.keymap.Exit):
			return cmm.prevM, nil
		}
	}

	var cmd tea.Cmd
	cmm.name, cmd = cmm.name.Update(msg)
	return cmm, cmd
}

func (cmm createMigrationModel) View() string {
	view := lipgloss.JoinVertical(
		0,
		"Name of new migration:\n",
		cmm.name.View(),
	)
	if cmm.status != "" {
		view = lipgloss.JoinVertical(
			0,
			fmt.Sprintf("[status]: %s", cmm.status),
		)
	}

	view += fmt.Sprintf("\n\n%s", cmm.help.View(cmm.keymap))

	return view
}

func (cmm createMigrationModel) create() error {
	name := fmt.Sprintf(
		"%s_%s",
		cmm.formatTime(time.Now().UTC()),
		cmm.name.Value(),
	)
	if err := createFile(path.Join(cmm.prevM.sourceFile, fmt.Sprintf("%s.up.sql", name))); err != nil {
		return err
	}
	if err := createFile(path.Join(cmm.prevM.sourceFile, fmt.Sprintf("%s.down.sql", name))); err != nil {
		return err
	}
	return nil
}

func (cmm createMigrationModel) formatTime(t time.Time) string {
	return t.Format("200601021504")
}

func createFile(p string) error {
	_, err := os.Create(p)
	return err
}
