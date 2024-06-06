package app

import (
	"fmt"
	_ "github.com/Vaniog/lazymigrate/app/build"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
)

type migrateModel struct {
	sourceFile  string
	databaseUrl string

	migrate *migrate.Migrate
	err     error

	help   help.Model
	keymap migrateKeymap

	logModel logModel
}

type setMigrateMsg struct {
	migrate *migrate.Migrate
}
type setMigrateErrMsg struct {
	err error
}
type gotoMigrationMsg struct {
	version uint
}

func newMigrateModel(sourceUrl, databaseUrl string) migrateModel {
	mm := migrateModel{
		sourceFile:  sourceUrl,
		databaseUrl: databaseUrl,
		logModel:    defaultLogModel,
		keymap:      defaultMigrateKeymap,
	}

	mm.help = help.New()
	return mm
}

func (mm migrateModel) Init() tea.Cmd {
	migrateCmd := func() tea.Msg {
		m, err := migrate.New(
			fmt.Sprintf("file://%s", mm.sourceFile),
			mm.databaseUrl,
		)
		if err != nil {
			return setMigrateErrMsg{err}
		}

		m.Log = mm.logModel.lineLogger
		return setMigrateMsg{m}
	}

	return tea.Batch(migrateCmd, mm.logModel.Init())
}

func (mm migrateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, mm.keymap.Exit):
			return mm, tea.Quit
		case key.Matches(msg, mm.keymap.Config):
			cm := newConfigModel(mm)
			return cm, cm.Init()
		case key.Matches(msg, mm.keymap.Up):
			if mm.migrate == nil {
				mm.logModel.lineLogger.Printf("no connection")
				return mm, nil
			}
			err := mm.migrate.Up()
			if err != nil {
				mm.logModel.lineLogger.Printf("up error: %s", err)
			} else {
				mm.logModel.lineLogger.Printf("up succeed")
			}
			return mm, nil
		case key.Matches(msg, mm.keymap.Down):
			if mm.migrate == nil {
				mm.logModel.lineLogger.Printf("no connection")
				return mm, nil
			}
			err := mm.migrate.Down()
			if err != nil {
				mm.logModel.lineLogger.Printf("down error: %s", err)
			} else {
				mm.logModel.lineLogger.Printf("down succeed")
			}
			return mm, nil
		case key.Matches(msg, mm.keymap.Goto):
			if mm.migrate == nil {
				mm.logModel.lineLogger.Printf("no connection")
				return mm, nil
			}
			gm := newGotoModel(mm)
			return gm, gm.Init()
		case key.Matches(msg, mm.keymap.New):
			cmm := newCreateMigrationModel(mm)
			return cmm, cmm.Init()
		}
	case setMigrateMsg:
		mm.migrate = msg.migrate
		return mm, nil
	case setMigrateErrMsg:
		mm.err = msg.err
		return mm, nil
	case logModelTickMessage:
		_, cmd := mm.logModel.Update(msg)
		return mm, cmd
	case gotoMigrationMsg:
		err := mm.migrate.Migrate(msg.version)
		if err != nil {
			mm.logModel.lineLogger.Printf("goto error: %s", err)
		} else {
			mm.logModel.lineLogger.Printf("goto succeed")
		}
	}

	return mm, nil
}

func (mm migrateModel) View() string {
	configTable := table.New()
	configTable.Row("Database", mm.databaseUrl)
	configTable.Row("Source", mm.sourceFile)
	configTable.Border(lipgloss.NormalBorder())

	statusTable := table.New()

	if mm.err != nil {
		statusTable.Row("Status", "Error")
		statusTable.Row("Error", mm.err.Error())
	} else if mm.migrate != nil {
		version, dirty, _ := mm.migrate.Version()
		statusTable.Row("Status", "Connected")
		statusTable.Row("Version", fmt.Sprintf("%d", version))
		statusTable.Row("Dirty", fmt.Sprintf("%t", dirty))
	} else {
		statusTable.Row("Status", "Connecting ...")
	}

	re := lipgloss.NewRenderer(os.Stdout)
	styleTable := func(t *table.Table) {
		leftStyle := re.NewStyle().Padding(0, 1).Width(12)
		rightStyle := re.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("7"))
		t.Border(lipgloss.NormalBorder())
		t.BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238")))
		t.StyleFunc(func(_, col int) lipgloss.Style {
			if col == 0 {
				return leftStyle
			}
			return rightStyle
		})
	}

	styleTable(configTable)
	styleTable(statusTable)

	logsStyle := lipgloss.NewStyle().Height(6)

	view := lipgloss.JoinVertical(
		0,
		configTable.Render(),
		statusTable.Render(),
		logsStyle.Render(mm.logModel.View()),
		mm.help.View(mm.keymap),
	)
	return view
}
