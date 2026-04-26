package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/TheLazyTurtle33/sea-core/shared/cleanup"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
	"github.com/TheLazyTurtle33/sea-core/web/internal/auth"
	"github.com/TheLazyTurtle33/sea-core/web/internal/data"
	"github.com/TheLazyTurtle33/sea-core/web/internal/webhook"

	"github.com/TheLazyTurtle33/sea-core/web/internal/pages/commands"
	"github.com/TheLazyTurtle33/sea-core/web/internal/pages/dashboard"
	"github.com/TheLazyTurtle33/sea-core/web/internal/pages/landing"
)

func main() {
	logger.Init()
	logger.Log("web starting")
	trapSignals()

	// back end
	http.HandleFunc("/auth/callback", auth.CallbackHandler)
	http.HandleFunc("/auth/user", auth.UserHandler)
	http.HandleFunc("/auth/bot", auth.BotHandler)
	http.HandleFunc("/eventsub", webhook.Handler)
	http.HandleFunc("/reminders", data.HandleReminder)

	// front end
	http.HandleFunc("/", landing.Page)
	http.HandleFunc("/commands", commands.Page)
	http.HandleFunc("/dashboard", dashboard.Page)

	logger.Log("main: starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

func trapSignals() {
	// Create a channel to receive OS signals.
	sigs := make(chan os.Signal, 1)

	// Register the channel to receive SIGINT and SIGTERM signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cleanup.RegisterCleaner(&cleaner{})

	// Start a goroutine to wait for the signal.
	go func() {
		sig := <-sigs
		_ = sig
		cleanup.Clean()
		os.Exit(0)
	}()
}

type cleaner struct {
	cleanup.Cleaner
}

func (c *cleaner) Clean() {
	logger.Log("web stopping")
}
