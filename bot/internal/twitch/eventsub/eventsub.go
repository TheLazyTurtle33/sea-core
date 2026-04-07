package eventsub

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/chat"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/redeems"
)

const webhookCallbackUrl = "https://lazyturtle33.live/eventsub"

type EventNotification struct {
	Subscription struct {
		ID        string `json:"id"`
		Status    string `json:"status"`
		Type      string `json:"type"`
		Version   string `json:"version"`
		Condition struct {
			BroadcasterUserID string `json:"broadcaster_user_id"`
			UserID            string `json:"user_id"`
		} `json:"condition"`
		Transport struct {
			Method   string `json:"method"`
			Callback string `json:"callback"`
		} `json:"transport"`
		CreatedAt time.Time `json:"created_at"`
		Cost      int       `json:"cost"`
	} `json:"subscription"`
	Event json.RawMessage `json:"event"`
}

type Condition map[string]string

type subscription struct {
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	Condition Condition `json:"condition"`
	Transport transport `json:"transport"`
}

type transport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

func Subscribe(eventType, version string, condition Condition) error {
	logger.Log("subscribing to %s", "event", eventType)
	sub := subscription{
		Type:      eventType,
		Version:   version,
		Condition: condition,
		Transport: transport{
			Method:   "webhook",
			Callback: webhookCallbackUrl,
			Secret:   context.Get().Auth.GetClientSecret(),
		},
	}
	body, err := json.Marshal(sub)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription: %w", err)
	}
	var body2 []byte
	body2, err = twitchapi.As("app").Post("/eventsub/subscriptions", string(body))
	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", eventType, err)
	}
	logger.Log("subscribed to %s", "event", eventType, "body", string(body2))
	return nil
}

func HandleNotification(notification *EventNotification) {
	switch notification.Subscription.Type {
	case "channel.chat.message":
		var data datatypes.ChatMessageData
		if err := json.Unmarshal(notification.Event, &data); err != nil {
			logger.Error("failed to unmarshal event", err)
			return
		}
		chat.HandleMessage(data)
	case "channel.channel_points_custom_reward_redemption.add":
		var data datatypes.RedemptionData
		if err := json.Unmarshal(notification.Event, &data); err != nil {
			logger.Error("failed to unmarshal event", err)
			return
		}
		redeems.HandleRedemption(data)
	default:
		logger.Warn("unknown event type", "event", notification.Subscription.Type)
	}
}
