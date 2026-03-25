package main

import (
	"github.com/TheLazyTurtle33/sea-core/bot/internal/api"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/eventsub"
)

func main() {
	eventsub.SubscribeAll()
	context.Get().Auth.RefreshToken("bot")
	context.Get().Auth.RefreshToken("user")
	api.Start()
}
