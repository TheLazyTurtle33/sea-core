package commands

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TheLazyTurtle33/sea-core/web/internal/botclient"
)

func Page(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "still under costructoin! sowy >.<")
	commands, err := botclient.GetCommandsJson()
	if err != nil {
		log.Printf("Error getting Json: %e", err)
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, string(commands))
}
