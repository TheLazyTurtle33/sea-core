package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/file"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
)

type Auth struct {
	UserOauthToken        string `json:"user_oauth_token"`
	BotOauthToken         string `json:"bot_oauth_token"`
	UserOauthRefreshToken string `json:"user_oauth_refresh_token"`
	BotOauthRefreshToken  string `json:"bot_oauth_refresh_token"`
	ClientId              string `json:"client_id"`
	ClientSecret          string `json:"client_secret"`
	AppAccessToken        string `json:"app_access_token"`
	ExpectingToken        string `json:"expecting_token"`
}

func New() *Auth {
	a := &Auth{}
	secret := file.New("/app/data/secret.json").Read()
	if secret == nil {
		logger.Error(nil, "secret.json is required. look at secret.json.example for an example")
	}
	err := json.Unmarshal(secret, a)
	if err != nil {
		logger.Error(err, "failed to unmarshal secret.json")
	}
	return a
}

func (a *Auth) GetUserOauthToken() string {
	return a.UserOauthToken
}

func (a *Auth) SetOauthToken(code string) {
	if a.ExpectingToken == "" {
		logger.Warn("not expecting token, ignoring")
		return
	}
	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token", map[string][]string{
		"client_id":     {a.ClientId},
		"client_secret": {a.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectUrl},
	})
	if err != nil {
		logger.Error(err, "failed to get oauth token")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, "failed to read response body")
		return
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Error(err, "failed to unmarshal response body")
		return
	}
	token := data["access_token"].(string)
	refreshToken := data["refresh_token"].(string)
	switch a.ExpectingToken {
	case "user":
		a.UserOauthToken = token
		a.UserOauthRefreshToken = refreshToken
		a.ExpectingToken = ""
		a.Save()
	case "bot":
		a.BotOauthToken = token
		a.BotOauthRefreshToken = refreshToken
		a.ExpectingToken = ""
		a.Save()
	default:
		logger.Warn("not expecting token, ignoring")
	}

}

func (a *Auth) StartRefreshTokensWorker() {
	go func() {
		for {
			logger.Log("refreshing tokens")
			a.RefreshToken("user")
			a.RefreshToken("bot")
			a.RefreshAppAccessToken()
			time.Sleep(1 * time.Hour)
		}
	}()
}

func (a *Auth) RefreshToken(tokenType string) {
	var refreshToken string
	switch tokenType {
	case "user":
		refreshToken = a.UserOauthRefreshToken
	case "bot":
		refreshToken = a.BotOauthRefreshToken
	default:
		logger.Warn("invalid token type")
		return
	}
	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token", map[string][]string{
		"client_id":     {a.ClientId},
		"client_secret": {a.ClientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	})
	if err != nil {
		logger.Error(err, "failed to refresh token")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, "failed to read response body")
		return
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Error(err, "failed to unmarshal response body")
		return
	}
	token := data["access_token"].(string)
	refreshToken = data["refresh_token"].(string)
	switch tokenType {
	case "user":
		a.UserOauthToken = token
		a.UserOauthRefreshToken = refreshToken
	case "bot":
		a.BotOauthToken = token
		a.BotOauthRefreshToken = refreshToken
	}
	a.Save()
}
func (a *Auth) RefreshAppAccessToken() {
	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token", map[string][]string{
		"client_id":     {a.ClientId},
		"client_secret": {a.ClientSecret},
		"grant_type":    {"client_credentials"},
	})
	if err != nil {
		logger.Error(err, "failed to refresh app access token")
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, "failed to read response body")
		return
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Error(err, "failed to unmarshal response body")
		return
	}
	a.AppAccessToken = data["access_token"].(string)
}

func (a *Auth) GetAppAccessToken() string {
	if a.AppAccessToken == "" {
		a.RefreshAppAccessToken()
	}
	return a.AppAccessToken
}

func (a *Auth) GetBotOauthToken() string {
	return a.BotOauthToken
}

func (a *Auth) GetClientId() string {
	return a.ClientId
}

func (a *Auth) GetClientSecret() string {
	return a.ClientSecret
}

func (a *Auth) StartExpectingToken(tokenType string) {
	a.ExpectingToken = tokenType
}

func (a *Auth) Save() {
	file.New("/app/data/secret.json").Save(a.Export())
}

func (a *Auth) Export() []byte {
	out, err := json.Marshal(a)
	if err != nil {
		logger.Error(err, "failed to marshal auth")
		return nil
	}
	return out
}

const Scope = `
bits:read+
channel:manage:broadcast+
channel:manage:ads+
channel:manage:clips+
channel:read:charity+
channel:edit:commercial+
channel:manage:guest_star+
channel:read:hype_train+
channel:manage:moderators+
channel:manage:polls+
channel:manage:predictions+
channel:manage:raids+
channel:manage:redemptions+
channel:manage:schedule+
channel:read:subscriptions+
channel:manage:vips+
channel:moderate+
channel:bot+
clips:edit+
editor:manage:clips+
moderation:read+
moderator:manage:announcements+
moderator:manage:automod+
moderator:manage:automod_settings+
moderator:read:blocked_terms+
moderator:manage:banned_users+
moderator:manage:blocked_terms+
moderator:read:chat_messages+
moderator:manage:chat_settings+
moderator:manage:shoutouts+
moderator:manage:shield_mode+
moderator:manage:suspicious_users+
moderator:manage:unban_requests+
moderator:manage:warnings+
moderator:read:followers+
user:edit+
user:read:chat+
user:read:email+
user:read:emotes+
user:read:follows+
user:read:subscriptions+
user:manage:whispers+
user:write:chat+
user:bot+
chat:edit+
chat:read+
whispers:read
`
const redirectUrl = "https://lazyturtle33.live/auth/callback"

func (a *Auth) CreateUserOauthUrl() string {
	a.ExpectingToken = "user"
	scope := strings.Replace(Scope, "\n", "", -1)
	return fmt.Sprintf(
		"https://id.twitch.tv/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		a.ClientId, redirectUrl, scope)
}

func (a *Auth) CreateBotOauthUrl() string {
	a.ExpectingToken = "bot"
	scope := strings.Replace(Scope, "\n", "", -1)
	return fmt.Sprintf(
		"https://id.twitch.tv/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		a.ClientId, redirectUrl, scope)
}
