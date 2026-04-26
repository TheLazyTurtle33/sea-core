package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
)

const baseURL = "https://api.twitch.tv/helix"

type TwitchAPI struct {
	auth string
	id   string
}

func AsBot() *TwitchAPI {
	return &TwitchAPI{
		auth: context.Get().Auth.GetBotOauthToken(),
		id:   context.Get().GetBot().Id,
	}
}

func AsUser() *TwitchAPI {
	return &TwitchAPI{
		auth: context.Get().Auth.GetUserOauthToken(),
		id:   context.Get().GetBroadcaster().Id,
	}
}

func AsApp() *TwitchAPI {
	return &TwitchAPI{
		auth: context.Get().Auth.GetAppAccessToken(),
		id:   "",
	}
}

func (t *TwitchAPI) Get(endpoint string) ([]byte, error) {
	resp, err := t.request(endpoint, "GET", "", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", string(body))
	}

	return body, nil
}

func (t *TwitchAPI) Post(endpoint, json string) ([]byte, error) {
	resp, err := t.request(endpoint, "POST", "application/json", strings.NewReader(json))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", string(body))
	}

	return body, nil
}

func (t *TwitchAPI) Patch(endpoint, json string) ([]byte, error) {
	resp, err := t.request(endpoint, "PATCH", "application/json", strings.NewReader(json))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", string(body))
	}

	return body, nil
}

func (t *TwitchAPI) Delete(endpoint string) ([]byte, error) {
	resp, err := t.request(endpoint, "DELETE", "", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", string(body))
	}

	return body, nil
}

func (t *TwitchAPI) SendMessage(message string) ([]byte, error) {
	return t.Post("/chat/messages", fmt.Sprintf(`{"broadcaster_id": "%s", "sender_id": "%s", "message": "%s"}`, context.Get().GetBroadcaster().Id, t.id, message))
}

func (t *TwitchAPI) SendReply(message, parentMessageID string) ([]byte, error) {
	return t.Post("/chat/messages", fmt.Sprintf(`{"broadcaster_id": "%s", "sender_id": "%s", "message": "%s", "reply_parent_message_id": "%s"}`, context.Get().GetBroadcaster().Id, t.id, message, parentMessageID))
}

func (t *TwitchAPI) SendAnnouncement(message string, color string) ([]byte, error) {
	return t.Post(fmt.Sprintf("/chat/announcements?broadcaster_id=%s&moderator_id=%s", context.Get().GetBroadcaster().Id, t.id), fmt.Sprintf(`{"message": "%s", "color": "%s"}`, message, color))
}

func (t *TwitchAPI) request(endpoint, method, bodyType string, data io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(method, baseURL+endpoint, data)
	if err != nil {
		logger.Error("failed to create request", err)
		return nil, err
	}
	req.Header.Set("Client-ID", context.Get().Auth.GetClientId())

	req.Header.Set("Authorization", "Bearer "+t.auth)

	if bodyType != "" {
		req.Header.Set("Content-Type", bodyType)
	}
	return http.DefaultClient.Do(req)
}
