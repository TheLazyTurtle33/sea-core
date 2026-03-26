package context

import (
	"github.com/TheLazyTurtle33/sea-core/bot/internal/auth"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
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
		logger.Warn("Broadcaster could not be retrieved")
	}
	c.Bot = datatypes.NewUser(c.Auth.GetBotOauthToken(), c.Auth.GetClientId())
	if c.Bot == nil {
		logger.Warn("Bot could not be retrieved")
	}
	return c
}

func (c *Context) GetBroadcaster() *datatypes.User {
	if c.Broadcaster == nil {
		c.Broadcaster = datatypes.NewUser(c.Auth.GetUserOauthToken(), c.Auth.GetClientId())
		if c.Broadcaster == nil {
			logger.Warn("Broadcaster could not be retrieved")
		}
	}

	return c.Broadcaster
}

func (c *Context) GetBot() *datatypes.User {
	if c.Bot == nil {
		c.Bot = datatypes.NewUser(c.Auth.GetUserOauthToken(), c.Auth.GetClientId())
		if c.Bot == nil {
			logger.Warn("Broadcaster could not be retrieved")
		}
	}
	return c.Bot
}

func (c *Context) GetAuth() *auth.Auth {
	return c.Auth
}
