package eventsub

import (
	"encoding/json"
	"log"
	"slices"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
)

var eventsToSubscribeTo []string = []string{
	"channel.chat.message",
}

type SubEvents struct {
	Data []struct {
		ID        string `json:"id"`
		Status    string `json:"status"`
		Type      string `json:"type"`
		Version   string `json:"version"`
		Cost      int    `json:"cost"`
		Condition struct {
			BroadcasterUserID string `json:"broadcaster_user_id"`
		} `json:"condition"`
		CreatedAt string `json:"created_at"`
		Transport struct {
			Method   string `json:"method"`
			Callback string `json:"callback"`
		} `json:"transport"`
	} `json:"data"`
	Total        int `json:"total"`
	TotalCost    int `json:"total_cost"`
	MaxTotalCost int `json:"max_total_cost"`
	Pagination   struct {
	} `json:"pagination"`
}

func SubscribeAll() {
	body, err := twitchapi.As("app").Get("/eventsub/subscriptions")
	if err != nil {
		log.Println(err)
		return
	}
	var subEvents SubEvents
	if err := json.Unmarshal(body, &subEvents); err != nil {
		log.Println(err)
		return
	}

	for _, sub := range subEvents.Data {
		if !slices.Contains(eventsToSubscribeTo, sub.Type) {
			switch sub.Type {
			case "channel.chat.message":
				subChat()
			}
		}
	}
}

func subChat() {
	if context.Get().Broadcaster == nil {
		log.Println("failed to subscribe to chat messages: no broadcaster")
		return
	}
	if context.Get().Bot == nil {
		log.Println("failed to subscribe to chat messages: no bot")
		return
	}
	err := Subscribe("channel.chat.message", "1", Condition{
		"broadcaster_user_id": context.Get().Broadcaster.Id,
		"user_id":             context.Get().Bot.Id,
	})
	if err != nil {
		log.Printf("failed to subscribe to chat messages: %s", err)
	}
}
