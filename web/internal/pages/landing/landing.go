package landing

import (
	"fmt"
	"net/http"
)

func Page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello from the web server! :3")
}
