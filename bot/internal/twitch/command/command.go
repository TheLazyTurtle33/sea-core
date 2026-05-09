package command

import (
	"fmt"
	"slices"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/queue"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
)

type Command struct {
	Name        string          `json:"name"`
	Triggers    []string        `json:"triggers"`
	Description string          `json:"description"`
	Usage       string          `json:"usage"`
	Actions     []action.Action `json:"actions"`
	QueueName   string          `json:"queue"`
	Blocking    bool            `json:"blocking"`
	Active      bool            `json:"is_active"`
	AlowedUsers []string        `json:"alowed_users"`
}

// func MakeCommand(jsonData []byte) Command {
// 	var cmd Command
// 	json.Unmarshal(jsonData, &cmd)
// 	for _, actionName := range cmd.ActionList {
// 		action, err := action.GetActionByName(actionName)
// 		if err != nil {
// 			logger.Error("failed to get action by name", err)
// 			continue
// 		}
// 		cmd.Actions = append(cmd.Actions, action)
// 	}
// 	return cmd
// }

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
			c.parseData(&a, data)
			q.AddAction(a)
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
				c.parseData(&a, data)
				q.AddAction(a)
			}
		}

	}

	if c.Blocking {
		q.Lock()
	}

	q.Start()
}

func (c *Command) parseData(a *action.Action, ChatMessageData any) {
	if a.ActionData != nil {
		for i := range a.ActionData {
			switch a.ActionData[i].Type {
			case action.ChatMessageType:
				if a.ActionData[i].Data == nil {
					a.ActionData[i].Data = ChatMessageData
				}
			case action.ActionsType:
				acts, ok := a.ActionData[i].Data.([]action.Action)
				if !ok {
					logger.Error("expected []action.Action for ActionsType ActionData, got", fmt.Errorf("expected []action.Action for ActionsType ActionData, got %T", a.ActionData[i].Data))
					continue
				}
				for j := range acts {
					c.parseData(&acts[j], ChatMessageData)
				}
				a.ActionData[i].Data = acts
			}
		}
	}
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
