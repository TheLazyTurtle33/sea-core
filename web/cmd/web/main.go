package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TheLazyTurtle33/sea-core/web/internal/auth"
)

func main() {

	http.HandleFunc("/auth/callback", auth.CallbackHandler)
	http.HandleFunc("/auth/user", auth.UserHandler)
	http.HandleFunc("/auth/bot", auth.BotHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello from the web server! :3")

	})

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
