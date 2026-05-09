package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/eventsub"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
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
	context.Get().Auth.SetOauthToken(body.Token)
	w.WriteHeader(http.StatusOK)
}

func HandleOauthUrl(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	url := ""
	switch body.Type {
	case "user":
		context.Get().Auth.StartExpectingToken("user")
		url = context.Get().Auth.CreateUserOauthUrl()
	case "bot":
		context.Get().Auth.StartExpectingToken("bot")
		url = context.Get().Auth.CreateBotOauthUrl()
	default:
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": url})

}

func HandleGetAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(context.Get().Auth.Export())

}

func HandleNotification(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	var notification eventsub.EventNotification
	if err := json.Unmarshal(body, &notification); err != nil {
		logger.Error("failed to unmarshal body", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	logger.Debug("received eventsub notification", "body", string(body))
	eventsub.HandleNotification(&notification)

	w.WriteHeader(http.StatusOK)
}

func Start() {
	http.HandleFunc("/internal/token", HandleToken)
	http.HandleFunc("/internal/oauth-url", HandleOauthUrl)
	http.HandleFunc("/internal/get-auth", HandleGetAuth)
	http.HandleFunc("/internal/notification", HandleNotification)
	http.HandleFunc("/internal/commands", HandleCommands)
	logger.Log("bot internal API listening on :9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		logger.Error("failed to start internal API", err)
	}
}
