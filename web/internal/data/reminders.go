package data

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const path = "/app/shared/reminders.json"

func HandleReminder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := os.ReadFile(path)
		if err != nil {
			http.Error(w, fmt.Sprintf("faild to read file; %e", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	if r.Method == "DELETE" {
		text := `{"reminders":[]}`
		err := os.WriteFile(path, []byte(text), 0644)
		if err != nil {
			http.Error(w, fmt.Sprintf("faild to write file; %e", err), http.StatusInternalServerError)
			log.Printf("error writing file %e", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
		return
	}

	http.Error(w, "have to use GET or DELETE", http.StatusBadRequest)

}
