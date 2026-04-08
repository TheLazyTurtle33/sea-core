package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/queue"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/api"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/cleanup"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/eventsub"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/redeems"
)

func main() {
	logger.Init()
	logger.Log("bot starting")
	trapSignals()

	context.Get().Auth.StartRefreshTokensWorker()

	eventsub.SubscribeAll()
	queue.StartUp()

	redeems.RegisterRedeems()

	// if err := obs.Init(); err != nil {
	// 	logger.Error(err, "failed to connect to OBS")
	// }

	// logger.Log("connected to OBS")
	// client, err := obs.Get()
	// if err != nil {
	// 	logger.Error(err, "failed to get OBS client")
	// }
	// result, err := client.Send("GetVersion", nil)
	// if err != nil {
	// 	logger.Error(err, "failed to get OBS version")
	// }
	// logger.Log("OBS version", result["obsWebSocketVersion"])
	// result, err = client.Send("SetCurrentProgramScene", map[string]any{
	// 	"sceneName": "BRB",
	// })
	// if err != nil {
	// 	logger.Error(err, "failed to set scene")
	// }
	// logger.Log("Set scene", result)
	api.Start()
}

func trapSignals() {
	// Create a channel to receive OS signals.
	sigs := make(chan os.Signal, 1)

	// Register the channel to receive SIGINT and SIGTERM signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cleanup.RegisterCleaner(&cleaner{})

	// Start a goroutine to wait for the signal.
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
