package log

import (
	"fmt"
	"time"

	"github.com/DilemaFixer/gog/internal/api"
)

type ChanLogger struct {
	clog chan string
}

func NewChanLogger() api.Logger {
	clog := make(chan string, 50)
	go log(clog)
	return &ChanLogger{clog: clog}
}

func (lc *ChanLogger) Log(level api.LogLevel, format string, args ...interface{}) {
	prefix := ""
	switch level {
	case api.LogLevelDebug:
		prefix = "[DEBUG]"
	case api.LogLevelInfo:
		prefix = "[INFO ]"
	case api.LogLevelWarn:
		prefix = "[WARN ]"
	default:
		prefix = "[LOG  ]"
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	entry := fmt.Sprintf("%s %s %s", prefix, timestamp, msg)

	select {
	case lc.clog <- entry:
	default:

	}
}

func log(clog chan string) {
	for msg := range clog {
		fmt.Println(msg)
	}
}
