package chat

import (
	"os/exec"
	"slices"
	"strings"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/tts"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/command"
)

func HandleMessage(data datatypes.ChatMessageData) {
	logger.Log("chat message received", "message", data.Message.Text, "user", data.ChatterUserName)
	context.Get().LastChat = &data
	parseForHello(data)
	parseForCommand(data)

	if context.Get().YappyChat {
		tts.MakeTTS(data.Message.Text, data.ChatterUserLogin)

		cmd := exec.Command("scp", "-i", "/root/.ssh/tts_key", "-o", "StrictHostKeyChecking=no", "/app/data/tts/tts.wav", "turt@192.168.0.182:/home/turt/StreamShit/audio/tts/tts.wav")
		if err := cmd.Run(); err != nil {
			logger.Error("erros sending file", err)
		}
	}

}

var hellos = []string{
	"hello",
	"hi",
	"hai",
	"hey",
	"sup",
	"yo",
	"howdy",
	"morning",
	"afternoon",
	"evening",
	"night",
}

type botName struct {
	name string
	mean bool
}

var botNames = []botName{
	{name: "lazybot33", mean: false},
	{name: "@lazybot33", mean: false},
	{name: "thelazybot33", mean: false},
	{name: "@thelazybot33", mean: false},
	{name: "bots", mean: false},
	{name: "bot", mean: false},
	{name: "@bots", mean: false},
	{name: "@bot", mean: true},
	{name: "lazy", mean: false},
	{name: "@lazy", mean: true},
	{name: "lazybot", mean: false},
	{name: "@lazybot", mean: false},
	{name: "clanker", mean: true},
	{name: "@clanker", mean: true},
}

func parseForHello(data datatypes.ChatMessageData) {
	if data.ChatterUserLogin == context.Get().GetBot().Login {
		return
	}
	var helloFound, botMentioned, mean bool

	for _, word := range data.Message.Fragments {
		if slices.Contains(hellos, word.Text) {
			helloFound = true
		}
		for _, botName := range botNames {
			if word.Text == botName.name {
				botMentioned = true
				mean = botName.mean
				break
			}
		}
	}

	var message = "Hai " + data.ChatterUserName + " ^w^"

	if data.ChatterUserLogin == "riceball_129" {
		message = "RICEBALL! hi :3"
	}

	if mean {
		message = ":c"
	}

	if helloFound && botMentioned {
		_, err := api.AsBot().SendReply(message, data.MessageID)
		if err != nil {
			logger.Error("failed to send reply", err)
			return
		}
	}

}

func parseForCommand(data datatypes.ChatMessageData) {
	command.RegisterCommands()
	parts := strings.Split(data.Message.Text, " ")
	if com := command.GetCommand(strings.TrimSpace(strings.ToLower(parts[0]))); com != nil {
		logger.Log("found command", "command", com.Name)
		com.AddActions(data)
	}
}
