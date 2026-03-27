package command

import (
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/queue"
)

type Command struct {
	Name        string          `json:"name"`
	Triggers    []string        `json:"triggers"`
	Description string          `json:"description"`
	Usage       string          `json:"usage"`
	Actions     []action.Action `json:"-"`
	QueueName   string          `json:"queue"`
	Blocking    bool            `json:"blocking"`
	Active      bool            `json:"is_active"`
}

func (c *Command) AddActions(data any) {
	if !c.Active {
		return
	}
	q := queue.GetQueue(c.QueueName)
	for _, a := range c.Actions {
		q.AddActions(a, data)
	}

	if c.Blocking {
		q.Lock()
	}

	q.Start()
}
