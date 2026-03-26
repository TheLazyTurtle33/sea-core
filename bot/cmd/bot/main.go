package main

import (
	"github.com/TheLazyTurtle33/sea-core/bot/internal/api"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/eventsub"
)

func main() {
	logger.Init()
	logger.Log("bot starting")
	eventsub.SubscribeAll()
	context.Get().Auth.StartRefreshTokensWorker()
	api.Start()
}
