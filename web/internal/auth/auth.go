package auth

import (
	"log"
	"net/http"

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
		log.Println(err)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func BotHandler(w http.ResponseWriter, r *http.Request) {
	url, err := botclient.CreateOauthUrl("bot")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
