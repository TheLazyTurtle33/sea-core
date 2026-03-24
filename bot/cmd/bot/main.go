package main

import (
	"github.com/TheLazyTurtle33/sea-core/bot/internal/api"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/eventsub"
)

func main() {
	eventsub.SubscribeAll()
	api.Start()
}
