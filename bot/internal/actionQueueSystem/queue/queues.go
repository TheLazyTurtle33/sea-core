package queue

import (
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action/actions"
)

var DefaultQueue = Queue{
	name:        "default",
	locked:      false,
	repeating:   false,
	persistent:  true,
	repeatDelay: 0,
}

var RedeemsQueue = Queue{
	name:       "redeems",
	persistent: true,
}

var DiscordQueue = Queue{
	name:        "Discord",
	locked:      true,
	repeating:   true,
	persistent:  true,
	repeatDelay: 30 * time.Minute,
	actions:     []action.Action{&actions.CreateDiscordInvite{}, &actions.SendMessage{}},
}

var Queues = []*Queue{
	&DefaultQueue,
	&DiscordQueue,
	&RedeemsQueue,
}

func GetQueue(name string) *Queue {
	for _, q := range Queues {
		if q.name == name {
			return q
		}
	}
	return nil
}

func StartUp() {
	DefaultQueue.Start()
	DiscordQueue.puased += 30 * time.Minute
	DiscordQueue.Start()
}
