package actions

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
)

const reminderFilePath = "/app/shared/reminders.json"

var AddReminder = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		if len(actionData) == 0 || actionData[0].Type != action.ChatMessageType {
			flags.Error = fmt.Errorf("AddReminder: expected chat message")
			return flags
		}

		chat, ok := actionData[0].Data.(datatypes.ChatMessageData)
		if !ok {
			flags.Error = fmt.Errorf("AddReminder: expected chat message, got %T", actionData[0].Data)
			return flags
		}

		words := strings.Split(chat.Message.Text, " ")
		if len(words) <= 1 {
			flags.AddActions = action.Flag{
				Active: true,
				Actions: []action.Action{
					ReplyToMessage.Make([]action.ActionData{
						{Type: action.ChatMessageType, Data: chat},
						{Type: action.StringType, Data: "use `!remind Message` to add a reminder"},
					}),
				},
			}
			return flags
		}

		reminderText := strings.TrimSpace(strings.Join(words[1:], " "))
		reminders := map[string][]string{}

		text, err := os.ReadFile(reminderFilePath)
		if err != nil {
			flags.Error = fmt.Errorf("reminder: failed to read file: %w", err)
			return flags
		}
		if len(text) > 0 {
			if err := json.Unmarshal(text, &reminders); err != nil {
				flags.Error = fmt.Errorf("reminder: failed to unmarshal file: %w", err)
				return flags
			}
		}

		reminders["reminders"] = append(reminders["reminders"], reminderText)
		js, err := json.Marshal(reminders)
		if err != nil {
			flags.Error = fmt.Errorf("reminder: failed to marshal reminders: %w", err)
			return flags
		}

		if err := os.WriteFile(reminderFilePath, js, 0644); err != nil {
			flags.Error = fmt.Errorf("reminder: failed to write file: %w", err)
			return flags
		}
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "AddReminder",
		Description: "Store a reminder from chat input.",
		RunData: []action.ActionData{
			{
				Type:        action.ChatMessageType,
				Description: "Reminder chat command.",
				Required:    true,
			},
		},
	},
}

func init() {
	action.ActionMap[AddReminder.MetaData.Name] = AddReminder
}
