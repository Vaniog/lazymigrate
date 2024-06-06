package app

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type log struct {
	time time.Time
	line string
}

type lineLogger struct {
	logs     []log
	logsLock sync.Mutex
	maxSize  int
}

func newLineLogger(maxSize int) *lineLogger {
	return &lineLogger{
		logs:     nil,
		logsLock: sync.Mutex{},
		maxSize:  maxSize,
	}
}

func (ll *lineLogger) Printf(format string, v ...interface{}) {
	line := fmt.Sprintf(format, v...)

	line = strings.Split(line, "\n")[0]
	ll.logsLock.Lock()
	defer ll.logsLock.Unlock()
	ll.logs = append(ll.logs, log{time.Now(), line})
	if len(ll.logs) > ll.maxSize {
		ll.logs = ll.logs[1:]
	}
}

// Verbose is for implementing migrate.Logger
func (_ *lineLogger) Verbose() bool {
	return false
}
