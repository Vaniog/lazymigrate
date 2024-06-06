package app

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

type logModel struct {
	lineLogger   *lineLogger
	logsLifetime time.Duration
}

type logModelTickMessage struct{}

var defaultLogModel = logModel{
	lineLogger:   newLineLogger(5),
	logsLifetime: 2000 * time.Millisecond,
}

func (lm logModel) active() bool {
	logs := lm.lineLogger.logs
	if len(logs) == 0 {
		return false
	}
	return logs[len(logs)-1].time.Add(lm.logsLifetime).After(time.Now())
}

func (lm logModel) Init() tea.Cmd {
	return nil
}

func (lm logModel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return lm, nil
}

func (lm logModel) View() string {
	view := ""
	lm.lineLogger.logsLock.Lock()
	defer lm.lineLogger.logsLock.Unlock()

	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	for _, log := range lm.lineLogger.logs {
		timeStr := timeStyle.Render(log.time.Format("15:04:05"))
		view += fmt.Sprintf(" > %s %s\n", timeStr, log.line)
	}
	return view
}
