package botclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const botURL = "http://bot:9090"

func SendToken(token string) error {
	body, _ := json.Marshal(map[string]string{"token": token})
	resp, err := http.Post(fmt.Sprintf("%s/internal/token", botURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to reach bot: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bot returned status: %d", resp.StatusCode)
	}
	return nil
}
