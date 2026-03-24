package context

import (
	"log"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/auth"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
)

type Context struct {
	Auth        *auth.Auth
	Broadcaster *datatypes.User
	Bot         *datatypes.User
}

var instance *Context

func Get() *Context {
	if instance == nil {
		instance = new()
	}
	return instance
}

func new() *Context {
	c := &Context{}
	c.Auth = auth.New()
	c.Broadcaster = datatypes.NewUser(c.Auth.GetUserOauthToken(), c.Auth.GetClientId())
	if c.Broadcaster == nil {
		log.Print("failed to get broadcaster")
	}
	c.Bot = datatypes.NewUser(c.Auth.GetBotOauthToken(), c.Auth.GetClientId())
	if c.Bot == nil {
		log.Print("failed to get bot")
	}
	return c
}
