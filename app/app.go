package app

import (
	"github.com/Vaniog/lazymigrate/app/config"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	tea *tea.Program
}

func NewApp() *App {
	model := newMigrateModel(config.Source, config.URL)
	return &App{
		tea.NewProgram(model),
	}
}

func (a *App) Run() error {
	_, err := a.tea.Run()
	return err
}
