package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/TheLazyTurtle33/sea-core/shared/logger"
	"github.com/TheLazyTurtle33/sea-core/web/internal/botclient"
)

type EventNotification struct {
	Subscription struct {
		Type string `json:"type"`
	} `json:"subscription"`
	Event json.RawMessage `json:"event"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	if !verifySignature(r.Header, body) {
		http.Error(w, "invalid signature", http.StatusForbidden)
		return
	}

	messageType := r.Header.Get("Twitch-Eventsub-Message-Type")
	switch messageType {
	case "webhook_callback_verification":
		handleChallenge(w, body)
	case "notification":
		handleNotification(w, body)
	case "revocation":
		logger.Log("subscription revoked:", r.Header.Get("Twitch-Eventsub-Subscription-Type"))
		w.WriteHeader(http.StatusNoContent)
	}
}

func verifySignature(headers http.Header, body []byte) bool {
	msgId := headers.Get("Twitch-Eventsub-Message-Id")
	timestamp := headers.Get("Twitch-Eventsub-Message-Timestamp")
	signature := headers.Get("Twitch-Eventsub-Message-Signature")

	auth, err := botclient.GetAuth()
	if err != nil {
		logger.Error("webhook: faid to get auth", err)
		return false
	}
	mac := hmac.New(sha256.New, []byte(auth["client_secret"]))
	mac.Write([]byte(msgId + timestamp + string(body)))
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

func handleChallenge(w http.ResponseWriter, body []byte) {
	var payload struct {
		Challenge string `json:"challenge"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "bad challenge", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload.Challenge))
}

func handleNotification(w http.ResponseWriter, body []byte) {
	err := botclient.SendNotification(body)
	if err != nil {
		logger.Error("Webhook: faild to send nota", err)
		http.Error(w, "failed to send notification", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
