package context

import (
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/auth"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
)

type TTSContext struct {
	IsActive  bool
	Delay     time.Duration
	Count     int
	YappyChat bool
}

const TTSDelayDefualt = 5 * time.Minute

type Context struct {
	Auth        *auth.Auth
	broadcaster *datatypes.User
	bot         *datatypes.User
	LastChat    *datatypes.ChatMessageData
	IsLive      bool
	CanShock    bool
	TTSContext  TTSContext
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
	c.LastChat = &datatypes.ChatMessageData{}
	c.broadcaster = datatypes.NewUser(c.Auth.GetUserOauthToken(), c.Auth.GetClientId())
	if c.broadcaster == nil {
		logger.Warn("Broadcaster could not be retrieved")
	}
	c.bot = datatypes.NewUser(c.Auth.GetBotOauthToken(), c.Auth.GetClientId())
	if c.bot == nil {
		logger.Warn("Bot could not be retrieved")
	}
	c.TTSContext = TTSContext{
		IsActive: true,
		Delay:    TTSDelayDefualt,
	}
	return c
}

func (c *Context) GetBroadcaster() *datatypes.User {
	if c.broadcaster == nil {
		c.broadcaster = datatypes.NewUser(c.Auth.GetUserOauthToken(), c.Auth.GetClientId())
		if c.broadcaster == nil {
			logger.Warn("Broadcaster could not be retrieved")
		}
	}

	return c.broadcaster
}

func (c *Context) GetBot() *datatypes.User {
	if c.bot == nil {
		c.bot = datatypes.NewUser(c.Auth.GetUserOauthToken(), c.Auth.GetClientId())
		if c.bot == nil {
			logger.Warn("Broadcaster could not be retrieved")
		}
	}
	return c.bot
}

func (c *Context) GetAuth() *auth.Auth {
	return c.Auth
}
