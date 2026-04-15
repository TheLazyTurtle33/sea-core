package command

import (
	"slices"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/queue"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
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
	AlowedUsers []string        `json:"alowed_users"`
}

func (c *Command) AddActions(data any) {
	if !c.Active {
		return
	}
	if c.QueueName == "" {
		c.QueueName = "default"
	}
	q := queue.GetQueue(c.QueueName)
	if len(c.AlowedUsers) == 0 || slices.Contains(c.AlowedUsers, "everyone") {
		for _, a := range c.Actions {
			q.AddActions(a, data)
		}
	} else {
		chatData, ok := data.(datatypes.ChatMessageData)
		if !ok {
			logger.Error("if command has allowed user feild, you must provide chat data", nil)
			return
		}
		badges := []string{}
		for _, badge := range chatData.Badges {
			badges = append(badges, badge.SetID)
		}

		if isAlowedUser(c.AlowedUsers, badges) {
			for _, a := range c.Actions {
				q.AddActions(a, data)
			}
		}

	}

	if c.Blocking {
		q.Lock()
	}

	q.Start()
}

func isAlowedUser(alowedUsers, badges []string) bool {
	if len(alowedUsers) == 0 || len(badges) == 0 {
		return false
	}
	for _, allowedUser := range alowedUsers {
		if slices.Contains(badges, allowedUser) {
			return true
		}
	}
	return slices.Contains(badges, "broadcaster")
}
