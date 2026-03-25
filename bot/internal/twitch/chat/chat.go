package chat

import (
	"log"
	"slices"
	"strings"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
)

type EventData struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	ChatterUserID        string `json:"chatter_user_id"`
	ChatterUserLogin     string `json:"chatter_user_login"`
	ChatterUserName      string `json:"chatter_user_name"`
	MessageID            string `json:"message_id"`
	Message              struct {
		Text      string `json:"text"`
		Fragments []struct {
			Type      string      `json:"type"`
			Text      string      `json:"text"`
			Cheermote interface{} `json:"cheermote"`
			Emote     interface{} `json:"emote"`
			Mention   interface{} `json:"mention"`
		} `json:"fragments"`
	} `json:"message"`
	Color  string `json:"color"`
	Badges []struct {
		SetID string `json:"set_id"`
		ID    string `json:"id"`
		Info  string `json:"info"`
	} `json:"badges"`
	MessageType                 string      `json:"message_type"`
	Cheer                       interface{} `json:"cheer"`
	Reply                       interface{} `json:"reply"`
	ChannelPointsCustomRewardID interface{} `json:"channel_points_custom_reward_id"`
	SourceBroadcasterUserID     string      `json:"source_broadcaster_user_id"`
	SourceBroadcasterUserLogin  string      `json:"source_broadcaster_user_login"`
	SourceBroadcasterUserName   string      `json:"source_broadcaster_user_name"`
	SourceMessageID             string      `json:"source_message_id"`
	SourceBadges                []struct {
		SetID string `json:"set_id"`
		ID    string `json:"id"`
		Info  string `json:"info"`
	} `json:"source_badges"`
	IsSourceOnly bool `json:"is_source_only"`
}

func HandleMessage(data EventData) {
	log.Println(data.Message.Text)
	parseForHello(data)
}

var hellos = []string{
	"hello",
	"hi",
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
	{name: context.Get().Bot.Login, mean: false},
	{name: "@" + context.Get().Bot.Login, mean: false},
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

func parseForHello(data EventData) {
	if data.ChatterUserLogin == context.Get().Bot.Login {
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
		log.Println("sending reply")
		_, err := api.As("bot").SendReply(message, data.MessageID)
		if err != nil {
			log.Println(err)
			return
		}
	}

}
