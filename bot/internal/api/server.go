package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/store"
)

// example endpoint web can call to send the oauth token to the bot
func HandleToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// log.Printf("received token: %s", body.Token)
	store.Get().Set("user_oauth_token", body.Token)
	w.WriteHeader(http.StatusOK)
}

func Start() {
	http.HandleFunc("/internal/token", HandleToken)
	log.Println("bot internal API listening on :9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}
