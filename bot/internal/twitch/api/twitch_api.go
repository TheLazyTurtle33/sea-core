package api

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
)

const baseURL = "https://api.twitch.tv/helix"

type TwitchAPI struct {
	as string
}

func As(as string) *TwitchAPI {
	return &TwitchAPI{
		as: as,
	}
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

func (t *TwitchAPI) request(endpoint, method, bodyType string, data io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(method, baseURL+endpoint, data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Client-ID", context.Get().Auth.GetClientId())
	if t.as == "bot" {
		req.Header.Set("Authorization", "Bearer "+context.Get().Auth.GetBotOauthToken())
	} else {
		req.Header.Set("Authorization", "Bearer "+context.Get().Auth.GetUserOauthToken())
	}
	if bodyType != "" {
		req.Header.Set("Content-Type", bodyType)
	}
	return http.DefaultClient.Do(req)
}
