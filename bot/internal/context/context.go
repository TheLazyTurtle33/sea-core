package context

import "github.com/TheLazyTurtle33/sea-core/bot/internal/auth"

type Context struct {
	Auth *auth.Auth
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
	return c
}
