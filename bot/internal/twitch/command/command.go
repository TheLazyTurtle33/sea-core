package command

import (
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/queue"
)

type Command struct {
	Name        string
	Triggers    []string
	Description string
	Usage       string
	Actions     []action.Action
	QueueName   string
	blocking    bool
}

func (c *Command) AddActions(data any) {
	q := queue.GetQueue(c.QueueName)
	for _, a := range c.Actions {
		q.AddActions(a, data)
	}

	if c.blocking {
		q.Lock()
	}

	q.Start()
}
