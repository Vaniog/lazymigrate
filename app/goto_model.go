package app

import (
	"github.com/Vaniog/lazymigrate/app/config"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type gotoModel struct {
	migration  textinput.Model
	migrations []migrationDto

	help   help.Model
	keymap gotoKeymap

	prevM migrateModel
}

type migrationDto struct {
	version uint
	name    string
}

func newGotoModel(prevM migrateModel) gotoModel {
	gm := gotoModel{
		migration: textinput.New(),
		help:      help.New(),
		keymap:    defaultGotoKeymap,
		prevM:     prevM,
	}

	gm.migration.Focus()
	gm.migration.Placeholder = ""
	gm.migration.Prompt = "> "

	var suggestions []string
	gm.migrations = availableMigrations()
	for _, m := range gm.migrations {
		suggestions = append(suggestions, m.name, strconv.Itoa(int(m.version)))
	}

	gm.migration.SetSuggestions(suggestions)
	gm.migration.ShowSuggestions = true
	return gm
}

func (gm gotoModel) Init() tea.Cmd {
	return nil
}

func (gm gotoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, gm.keymap.Exit):
			return gm.prevM, nil
		case key.Matches(msg, gm.keymap.Apply):
			if v, err := strconv.ParseInt(gm.migration.Value(), 10, 64); err == nil {
				return gm.prevM, func() tea.Msg {
					return gotoMigrationMsg{uint(v)}
				}
			}

			for _, m := range gm.migrations {
				if m.name == gm.migration.Value() {
					return gm.prevM, func() tea.Msg {
						return gotoMigrationMsg{m.version}
					}
				}
			}

			return gm.prevM, nil
		}
	}

	newMigration, cmd := gm.migration.Update(msg)
	gm.migration = newMigration
	return gm, cmd
}

func (gm gotoModel) View() string {
	return lipgloss.JoinVertical(
		0,
		"Goto (migration name or version)\n",
		gm.migration.View(),
		"\n",
		gm.help.View(gm.keymap),
	)
}

func availableMigrations() []migrationDto {
	var ms []migrationDto
	_ = filepath.WalkDir(config.Source, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".up.sql") {
			parts := strings.Split(d.Name(), "_")
			if len(parts) > 1 {
				versionStr := parts[0]
				version, err := strconv.ParseUint(versionStr, 10, 64)
				if err != nil {
					return nil
				}

				name := strings.Split(d.Name()[len(versionStr)+1:len(d.Name())], ".")[0]
				ms = append(ms, migrationDto{
					version: uint(version),
					name:    name,
				})
			}
		}

		return nil
	})
	return ms
}
