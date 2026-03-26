package main

import (
	"log"
	"net/http"

	"github.com/TheLazyTurtle33/sea-core/web/internal/auth"
	"github.com/TheLazyTurtle33/sea-core/web/internal/webhook"

	"github.com/TheLazyTurtle33/sea-core/web/internal/pages/commands"
	"github.com/TheLazyTurtle33/sea-core/web/internal/pages/dashboard"
	"github.com/TheLazyTurtle33/sea-core/web/internal/pages/landing"
)

func main() {

	// back end
	http.HandleFunc("/auth/callback", auth.CallbackHandler)
	http.HandleFunc("/auth/user", auth.UserHandler)
	http.HandleFunc("/auth/bot", auth.BotHandler)
	http.HandleFunc("/eventsub", webhook.Handler)

	// front end
	http.HandleFunc("/", landing.Page)
	http.HandleFunc("/commands", commands.Page)
	http.HandleFunc("/dashboard", dashboard.Page)

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
