package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	// "github.com/TheLazyTurtle33/sea-core/bot/internal/cleanup"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/file"
)

var logFile *file.File

const logFilePath = "/app/data/logs/logs-%v.log"

var StanderOutLogger *slog.Logger
var FileLogger *slog.Logger

// var WebLogger *slog.Logger

var Loggers = []*slog.Logger{}

func Init() {

	logFile = file.New(fmt.Sprintf(logFilePath, time.Now().Format("2006-01-02")))
	StanderOutLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	Loggers = append(Loggers, StanderOutLogger)
	writer, err := logFile.GetFileWiter()
	if err != nil {
		log.Fatal(err)
		return
	}
	FileLogger = slog.New(slog.NewTextHandler(writer, nil))
	Loggers = append(Loggers, FileLogger)
	// WebLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// cleanup.RegisterCleaner(&cleaner{})
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

// type cleaner struct {
// 	cleanup.Cleaner
// }

// func (c *cleaner) clean() {

// }
