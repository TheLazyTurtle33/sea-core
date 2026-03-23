package auth

import (
	"fmt"
	"log"
	"strings"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/file"
)

type Auth struct {
	userOauthToken string
	botOauthToken  string
	clientId       string
	clientSecret   string
	expectingToken string
}

func New() *Auth {
	a := &Auth{}
	secret := file.New("/app/data/secret.json").ReadAsJson()
	if secret == nil {
		log.Fatal("secret.json is required. look at secret.json.example for an example")
	}
	a.botOauthToken = secret["bot_oauth_token"].(string)
	a.userOauthToken = secret["user_oauth_token"].(string)

	a.clientId = secret["client_id"].(string)
	if a.clientId == "" {
		log.Fatal("client_id is required")
	}

	a.clientSecret = secret["client_secret"].(string)
	if a.clientSecret == "" {
		log.Fatal("client_secret is required")
	}
	return a
}

func (a *Auth) GetUserOauthToken() string {
	return a.userOauthToken
}

func (a *Auth) SetOauthToken(token string) {
	switch a.expectingToken {
	case "user":
		a.userOauthToken = token
		a.expectingToken = ""
		a.Save()
	case "bot":
		a.botOauthToken = token
		a.expectingToken = ""
		a.Save()
	default:
		log.Println("not expecting token, ignoring")
	}

}

func (a *Auth) GetBotOauthToken() string {
	return a.botOauthToken
}

func (a *Auth) GetClientId() string {
	return a.clientId
}

func (a *Auth) GetClientSecret() string {
	return a.clientSecret
}

func (a *Auth) ExpectingToken(tokenType string) {
	a.expectingToken = tokenType
}

func (a *Auth) Save() {
	file.New("/app/data/secret.json").Save([]byte(fmt.Sprintf(`{
	"bot_oauth_token": "%s",
	"user_oauth_token": "%s",
	"client_id": "%s",
	"client_secret": "%s"
}`, a.botOauthToken, a.userOauthToken, a.clientId, a.clientSecret)))
}

const userScope = `
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
user:edit+
user:read:chat+
user:read:email+
user:read:emotes+
user:read:follows+
user:read:subscriptions+
user:manage:whispers+
user:write:chat+
chat:edit+
chat:read+
whispers:read
`
const botScop = "user:bot+user:write:chat+user:manage:whispers+chat:edit+chat:read+whispers:read"
const redirectUrl = "https://lazyturtle33.live/auth/callback"

func (a *Auth) CreateUserOauthUrl() string {
	a.expectingToken = "user"
	scope := strings.Replace(userScope, "\n", "", -1)
	return fmt.Sprintf(
		"https://id.twitch.tv/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		a.clientId, redirectUrl, scope)
}

func (a *Auth) CreateBotOauthUrl() string {
	a.expectingToken = "bot"
	return fmt.Sprintf(
		"https://id.twitch.tv/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		a.clientId, redirectUrl, botScop)
}
