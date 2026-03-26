package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/cleanup"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/file"
)

const logDir = "/app/data/logs"
const logFilePath = "/app/data/logs/logs-%v.log"

var StanderOutLogger *slog.Logger
var FileLogger *slog.Logger

// var WebLogger *slog.Logger

var Loggers = []*slog.Logger{}

var logFileHandle *os.File

func Init() {
	// Ensure the logs directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}

	StanderOutLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	Loggers = append(Loggers, StanderOutLogger)

	path := fmt.Sprintf(logFilePath, time.Now().Format("2006-01-02"))
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	logFileHandle = f

	FileLogger = slog.New(slog.NewTextHandler(f, nil))
	Loggers = append(Loggers, FileLogger)
	// WebLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cleanup.RegisterCleaner(&cleaner{})
}

func Log(msg string, args ...any) {
	for _, l := range Loggers {
		l.Info(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	for _, l := range Loggers {
		l.Warn(msg, args...)
	}
}

func Error(err error, args ...any) {
	for _, l := range Loggers {
		l.Error(err.Error(), args...)
	}
}

func Debug(msg string, args ...any) {
	if file.Exists("/app/data/debug") {
		for _, l := range Loggers {
			l.Debug(msg, args...)
		}
	}
}

type cleaner struct {
	cleanup.Cleaner
}

func (c *cleaner) Clean() {
	if logFileHandle != nil {
		logFileHandle.Close()
	}
}
