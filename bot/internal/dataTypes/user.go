package datatypes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
)

type User struct {
	Id              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`

	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
	CreatedAt       string `json:"created_at"`
}

func NewUser(auth, clientId string) *User {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		logger.Error("failed to create request", err)
		return nil
	}
	req.Header.Set("Client-ID", clientId)

	req.Header.Set("Authorization", "Bearer "+auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("failed to do request", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("failed to read response body", err)
		return nil
	}
	var data struct {
		Data []User `json:"data"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Error("failed to unmarshal response body", err)
		return nil
	}
	if len(data.Data) == 0 {
		logger.Warn("no user found")
		return nil
	}
	u := data.Data[0]
	return &u
}
