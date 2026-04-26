package eventsub

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/chat"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/redeems"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
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

	_, err = twitchapi.AsApp().Post("/eventsub/subscriptions", string(body))
	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", eventType, err)
	}
	logger.Log("subscribed to event", "event", eventType)
	return nil
}

func HandleNotification(notification *EventNotification) {
	logToTempLogger(fmt.Sprintf("new %s event:\n %s\n\n", notification.Subscription.Type, notification.Event))
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
	case "stream.online":
		context.Get().IsLive = true
	case "stream.offline":
		context.Get().IsLive = false
	case "channel.raid":
		var data raidData
		if err := json.Unmarshal(notification.Event, &data); err != nil {
			logger.Error("failed to unmarshal event", err)
			return
		}
		HandleRaid(data)
	default:
		logger.Warn("unknown event type", "event", notification.Subscription.Type)
	}
}

func logToTempLogger(msg string) {
	file, err := os.OpenFile("/app/data/logs/log-temp.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Error("eventsub: faild to make temp log file", err)
	}
	defer file.Close()
	file.Write([]byte(msg))
}
