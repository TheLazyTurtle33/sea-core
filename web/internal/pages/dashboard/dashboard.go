package dashboard

import (
	"fmt"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "no dashboard yet x.x")
}
