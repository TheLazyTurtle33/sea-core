package auth

import (
	"net/http"

	"github.com/TheLazyTurtle33/sea-core/web/internal/botclient"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("code")
	botclient.SendToken(token)
	http.Redirect(w, r, "/", http.StatusFound)
}
