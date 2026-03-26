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
	culling:     true,
}

var RepetingQueueExample = Queue{
	name:        "repeting",
	locked:      false,
	repeating:   true,
	persistent:  false,
	repeatDelay: 10 * time.Second,
	actions:     []action.Action{&actions.ExampleActoin{}, &actions.ExampleActoin{}},
	culling:     false,
}

var Queues = []*Queue{&DefaultQueue, &RepetingQueueExample}

func GetQueue(name string) *Queue {
	for _, q := range Queues {
		if q.name == name {
			return q
		}
	}
	return nil
}
