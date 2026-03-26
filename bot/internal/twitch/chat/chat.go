package chat

import (
	"slices"
	"strings"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/command"
)

func HandleMessage(data datatypes.ChatMessageData) {
	logger.Log("chat message received", "message", data.Message.Text, "user", data.ChatterUserName)
	context.Get().LastChat = &data
	parseForHello(data)
	parseForCommand(data)
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
	words := strings.Split(strings.ToLower(data.Message.Text), " ")
	var helloFound, botMentioned, mean bool

	for _, word := range words {
		if slices.Contains(hellos, word) {
			helloFound = true
		}
		for _, botName := range botNames {
			if strings.Contains(word, botName.name) {
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
		_, err := api.As("bot").SendReply(message, data.MessageID)
		if err != nil {
			logger.Error(err, "failed to send reply")
			return
		}
	}

}

func parseForCommand(data datatypes.ChatMessageData) {
	command.RegisterCommands()
	if com := command.GetCommand(data.Message.Fragments[0].Text); com != nil {
		logger.Log("found command", "command", com.Name)
		com.AddActions(data)
	}
}
