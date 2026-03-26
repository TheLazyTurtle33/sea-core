package commands

import (
	"fmt"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "still under costructoin! sowy >.<")
}
