package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/twitch-bot-callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello from the web server callback! :3")
		log.Println(strings.Split(r.RequestURI, "access_token=")[1])
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello from the web server! :3")

	})

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
