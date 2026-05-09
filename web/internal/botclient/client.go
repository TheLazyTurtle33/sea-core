package botclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TheLazyTurtle33/sea-core/shared/logger"
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

func CreateOauthUrl(tokenType string) (string, error) {
	body, _ := json.Marshal(map[string]string{"type": tokenType})
	resp, err := http.Post(fmt.Sprintf("%s/internal/oauth-url", botURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to reach bot: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bot returned status: %d", resp.StatusCode)
	}
	var bodyReturn struct {
		Url string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&bodyReturn); err != nil {
		return "", fmt.Errorf("failed to decode bot response: %w", err)
	}
	return bodyReturn.Url, nil
}

func GetAuth() (map[string]string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/internal/get-auth", botURL))
	if err != nil {
		return nil, fmt.Errorf("failed to reach bot: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bot returned status: %d", resp.StatusCode)
	}
	var auth map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return nil, fmt.Errorf("failed to decode bot response: %w", err)
	}
	return auth, nil
}

func SendNotification(body []byte) error {
	resp, err := http.Post(fmt.Sprintf("%s/internal/notification", botURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to reach bot: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bot returned status: %d", resp.StatusCode)
	}
	return nil
}

func GetCommandsJson() ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("%s/internal/commands", botURL))
	if err != nil {
		return nil, fmt.Errorf("failed to reach bot: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bot returned status: %d", resp.StatusCode)
	}
	var data json.RawMessage
	logger.Log("client: ContentLength", "ContentLength", resp.ContentLength)
	// var json []byte = make([]byte, resp.ContentLength)
	// _, err = resp.Body.Read(json)
	// return json, nil
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decomp json: %w", err)
	}

	return []byte(data), nil

}
