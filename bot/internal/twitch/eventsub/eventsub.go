package eventsub

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/chat"
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
	log.Printf("subscribing to %s", eventType)
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
	log.Println(string(body2))
	return nil
}

func HandleNotification(notification *EventNotification) {
	switch notification.Subscription.Type {
	case "channel.chat.message":
		var data chat.EventData
		if err := json.Unmarshal(notification.Event, &data); err != nil {
			log.Println(err)
			return
		}
		chat.HandleMessage(data)
	default:
		log.Printf("unknown event type: %s", notification.Subscription.Type)
	}
}
