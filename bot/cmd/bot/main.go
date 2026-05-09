package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/queue"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/api"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/eventsub"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/redeems"
	"github.com/TheLazyTurtle33/sea-core/shared/cleanup"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
)

func main() {
	logger.Init()
	logger.Log("bot starting")
	trapSignals()

	context.Get().Auth.StartRefreshTokensWorker()

	// eventsub.UnsubAll()
	eventsub.SubscribeAll()
	queue.StartUp()

	redeems.RegisterRedeems()

	api.Start()

}

func trapSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	cleanup.RegisterCleaner(&cleaner{})
	go func() {
		sig := <-sigs
		_ = sig
		cleanup.Clean()
		os.Exit(0)
	}()
}

type cleaner struct {
	cleanup.Cleaner
}

func (c *cleaner) Clean() {
	logger.Log("bot stopping")
}
