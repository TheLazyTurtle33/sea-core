package eventsub

import (
	"encoding/json"
	"slices"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
)

var eventsToSubscribeTo []string = []string{
	"channel.chat.message",
	"channel.channel_points_custom_reward_redemption.add",
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
		logger.Error(err, "failed to get subscriptions")
		return
	}
	var subEvents SubEvents
	if err := json.Unmarshal(body, &subEvents); err != nil {
		logger.Error(err, "failed to unmarshal subscriptions")
		return
	}

	for _, sub := range subEvents.Data {
		if slices.Contains(eventsToSubscribeTo, sub.Type) {
			eventsToSubscribeTo = slices.Delete(eventsToSubscribeTo, slices.Index(eventsToSubscribeTo, sub.Type), slices.Index(eventsToSubscribeTo, sub.Type)+1)
		}
	}
	for _, event := range eventsToSubscribeTo {
		switch event {
		case "channel.chat.message":
			subChat()
		case "channel.channel_points_custom_reward_redemption.add":
			subRedemptions()
		default:
			logger.Warn("unknown subscription type", "type", event)
		}
	}
}

func subChat() {
	if context.Get().GetBroadcaster() == nil {
		logger.Warn("failed to subscribe to chat messages: no broadcaster")
		return
	}
	if context.Get().GetBot() == nil {
		logger.Warn("failed to subscribe to chat messages: no bot")
		return
	}
	err := Subscribe("channel.chat.message", "1", Condition{
		"broadcaster_user_id": context.Get().GetBroadcaster().Id,
		"user_id":             context.Get().GetBot().Id,
	})
	if err != nil {
		logger.Error(err, "failed to subscribe to chat messages")
	}
}

func subRedemptions() {
	if context.Get().GetBroadcaster() == nil {
		logger.Warn("failed to subscribe to redemptions: no broadcaster")
		return
	}
	err := Subscribe("channel.channel_points_custom_reward_redemption.add", "1", Condition{
		"broadcaster_user_id": context.Get().GetBroadcaster().Id,
	})
	if err != nil {
		logger.Error(err, "failed to subscribe to redemptions")
	}
}
