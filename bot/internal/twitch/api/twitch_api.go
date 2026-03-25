package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
)

const baseURL = "https://api.twitch.tv/helix"

type TwitchAPI struct {
	auth string
	id   string
}

func As(as string) *TwitchAPI {
	switch as {
	case "bot":
		return &TwitchAPI{
			auth: context.Get().Auth.GetBotOauthToken(),
			id:   context.Get().Bot.Id,
		}
	case "user":
		return &TwitchAPI{
			auth: context.Get().Auth.GetUserOauthToken(),
			id:   context.Get().Broadcaster.Id,
		}
	case "app":
		return &TwitchAPI{
			auth: context.Get().Auth.GetAppAccessToken(),
			id:   "",
		}
	}
	return nil
}

func (t *TwitchAPI) Get(endpoint string) ([]byte, error) {
	resp, err := t.request(endpoint, "GET", "", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (t *TwitchAPI) Post(endpoint, json string) ([]byte, error) {
	resp, err := t.request(endpoint, "POST", "application/json", strings.NewReader(json))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (t *TwitchAPI) SendMessage(message string) ([]byte, error) {
	return t.Post("/chat/messages", fmt.Sprintf(`{"broadcaster_id": "%s", "sender_id": "%s", "message": "%s"}`, context.Get().Broadcaster.Id, t.id, message))
}

func (t *TwitchAPI) SendReply(message, parentMessageID string) ([]byte, error) {
	return t.Post("/chat/messages", fmt.Sprintf(`{"broadcaster_id": "%s", "sender_id": "%s", "message": "%s", "reply_parent_message_id": "%s"}`, context.Get().Broadcaster.Id, t.id, message, parentMessageID))
}

func (t *TwitchAPI) SendAnnouncement(message string, color string) ([]byte, error) {
	return t.Post(fmt.Sprintf("/chat/announcements?broadcaster_id=%s,moderator_id=%s", context.Get().Broadcaster.Id, t.id), fmt.Sprintf(`{"message": "%s", "color": "%s"}`, message, color))
}
func (t *TwitchAPI) request(endpoint, method, bodyType string, data io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(method, baseURL+endpoint, data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Client-ID", context.Get().Auth.GetClientId())

	req.Header.Set("Authorization", "Bearer "+t.auth)

	if bodyType != "" {
		req.Header.Set("Content-Type", bodyType)
	}
	return http.DefaultClient.Do(req)
}
