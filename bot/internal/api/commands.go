package api

import (
	"net/http"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/command"
)

func HandleCommands(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.Write(command.GetCommandsJson())
		return
	}

	// TODO: implement POST to add commands from the web interface
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	// TODO: implement PATCH to edit commands from the web interface
	if r.Method == http.MethodPatch {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)

}
