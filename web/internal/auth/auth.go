package auth

import (
	"net/http"

	"github.com/TheLazyTurtle33/sea-core/shared/logger"
	"github.com/TheLazyTurtle33/sea-core/web/internal/botclient"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("code")
	botclient.SendToken(token)
	http.Redirect(w, r, "/", http.StatusFound)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	url, err := botclient.CreateOauthUrl("user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("auth: faid to get user url", err)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func BotHandler(w http.ResponseWriter, r *http.Request) {
	url, err := botclient.CreateOauthUrl("bot")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error("auth: faid to get bot url", err)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
