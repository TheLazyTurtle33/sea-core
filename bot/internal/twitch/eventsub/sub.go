package eventsub

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
)

var eventsToSubscribeTo = []string{
	"channel.chat.message",
	"channel.channel_points_custom_reward_redemption.add",
	"stream.online",
	"stream.offline",
	"channel.raid",
}

var eventsToUnSubscribeFrom = []string{}

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
		logger.Error("failed to get subscriptions", err)
		return
	}
	logger.Log("sub events", "body", string(body))
	var subEvents SubEvents
	if err := json.Unmarshal(body, &subEvents); err != nil {
		logger.Error("failed to unmarshal subscriptions", err)
		return
	}

	for _, sub := range subEvents.Data {
		if slices.Contains(eventsToSubscribeTo, sub.Type) {
			if sub.Status == "enabled" {
				eventsToSubscribeTo = slices.Delete(eventsToSubscribeTo, slices.Index(eventsToSubscribeTo, sub.Type), slices.Index(eventsToSubscribeTo, sub.Type)+1)
			} else {
				logger.Log("removing non enabled event", "event", sub.Type, "status", sub.Status, "eventID", sub.ID)
				eventsToUnSubscribeFrom = append(eventsToUnSubscribeFrom, sub.ID)
			}
		}
	}

	for _, eventID := range eventsToUnSubscribeFrom {
		body, err := twitchapi.AsUser().Delete(fmt.Sprintf("/eventsub/subscriptions?id=%s", eventID))
		if err != nil {
			logger.Error("Faild to delete old event", err, "eventID", eventID)
		}
		logger.Log("event sub body", "body", string(body))
	}

	for _, event := range eventsToSubscribeTo {
		switch event {
		case "channel.chat.message":
			subChat()
		case "channel.channel_points_custom_reward_redemption.add":
			subRedemptions()
		case "strem.online":
			if err := Subscribe("strem.online", "1", Condition{
				"broadcaster_user_id": context.Get().GetBroadcaster().Id,
			}); err != nil {
				logger.Error("failed to subscribe to chat messages", err)
			}
		case "strem.offline":
			if err := Subscribe("strem.offline", "1", Condition{
				"broadcaster_user_id": context.Get().GetBroadcaster().Id,
			}); err != nil {
				logger.Error("failed to subscribe to chat messages", err)
			}
		case "channel.raid":
			if err := Subscribe("channel.raid", "1", Condition{
				"broadcaster_user_id": context.Get().GetBroadcaster().Id,
			}); err != nil {
				logger.Error("failed to subscribe to chat messages", err)
			}
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
