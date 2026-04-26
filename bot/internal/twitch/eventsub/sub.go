package eventsub

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
)

var eventsToSubscribeTo = []string{
	"channel.chat.message",
	"channel.channel_points_custom_reward_redemption.add",
	"stream.online",
	"stream.offline",
	"channel.raid",
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
	body, err := twitchapi.AsApp().Get("/eventsub/subscriptions")
	if err != nil {
		logger.Error("sub: failed to get subscriptions", err)
		return
	}
	var subEvents SubEvents
	if err := json.Unmarshal(body, &subEvents); err != nil {
		logger.Error("sub: failed to unmarshal subscriptions", err)
		return
	}

	for _, sub := range subEvents.Data {
		if slices.Contains(eventsToSubscribeTo, sub.Type) {
			if sub.Status == "enabled" {
				eventsToSubscribeTo = slices.Delete(eventsToSubscribeTo, slices.Index(eventsToSubscribeTo, sub.Type), slices.Index(eventsToSubscribeTo, sub.Type)+1)
			}
		}
	}

	for _, event := range eventsToSubscribeTo {
		switch event {
		case "channel.chat.message":
			subChat()
		case "channel.channel_points_custom_reward_redemption.add":
			subRedemptions()
		case "stream.online":
			if err := Subscribe("stream.online", "1", Condition{
				"broadcaster_user_id": context.Get().GetBroadcaster().Id,
			}); err != nil {
				logger.Error("sub: failed to subscribe to stream.online", err)
			}
		case "stream.offline":
			if err := Subscribe("stream.offline", "1", Condition{
				"broadcaster_user_id": context.Get().GetBroadcaster().Id,
			}); err != nil {
				logger.Error("sub: failed to subscribe to stream.offline", err)
			}
		case "channel.raid":
			if err := Subscribe("channel.raid", "1", Condition{
				// "to_broadcaster_user_id":   context.Get().GetBroadcaster().Id,
				"from_broadcaster_user_id": context.Get().GetBroadcaster().Id,
			}); err != nil {
				logger.Error("sub: failed to subscribe to channle.raid", err)
			}
		default:
			logger.Warn("sub: unknown subscription type", "type", event)
		}
	}

}

func UnsubAll() {
	body, err := twitchapi.AsApp().Get("/eventsub/subscriptions")
	if err != nil {
		logger.Error("sub: failed to get subscriptions", err)
		return
	}
	logger.Log("sub: list of subs", "body", body)
	var subEvents SubEvents
	if err := json.Unmarshal(body, &subEvents); err != nil {
		logger.Error("sub: failed to unmarshal subscriptions", err)
		return
	}
	for _, sub := range subEvents.Data {
		if slices.Contains(eventsToSubscribeTo, sub.Type) {
			if sub.Status == "enabled" {
				_, err := twitchapi.AsApp().Delete(fmt.Sprintf("/eventsub/subscriptions?id=%s", sub.ID))
				if err != nil {
					logger.Error("sub: Faild to delete old event", err, "eventID", sub.ID)
				}
			}
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
		logger.Error("failed to subscribe to chat messages", err)
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
		logger.Error("failed to subscribe to redemptions", err)
	}
}
