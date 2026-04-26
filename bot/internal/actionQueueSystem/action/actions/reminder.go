package actions

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
)

const path = "/app/shared/reminders.json"

type AddReminder struct {
	action.Action
}

func (a AddReminder) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	if len(v) == 0 {
		flags.Error = fmt.Errorf("reminider: expectid one var of type chat")
		return flags
	}
	chat, ok := v[0].(datatypes.ChatMessageData)
	if !ok {
		flags.Error = fmt.Errorf("reminider: expectid chat message")
		return flags
	}

	words := strings.Split(chat.Message.Text, " ")
	if len(words) <= 1 {
		flags.AddActions = action.Flag{
			Active:     true,
			Actions:    []action.Action{ReplyToMessage{Message: "use `!remind Message` to add a reminder"}},
			ActionData: [][]any{{chat}},
		}
		return flags
	}
	words = words[1:]
	var reminder strings.Builder
	for _, word := range words {
		reminder.WriteString(word + " ")
	}

	var reminders map[string][]string

	text, err := os.ReadFile(path)
	if err != nil {
		flags.Error = fmt.Errorf("reminder: faild to read file; %w", err)
		return flags
	}

	err = json.Unmarshal(text, &reminders)
	if err != nil {
		flags.Error = fmt.Errorf("reminder: faild to turn content to json; %w", err)
		return flags
	}

	reminders["reminders"] = append(reminders["reminders"], strings.TrimSpace(reminder.String()))

	js, err := json.Marshal(reminders)
	if err != nil {
		flags.Error = fmt.Errorf("reminder: faild to make json; %w", err)
		return flags
	}

	err = os.WriteFile(path, []byte(js), 0644)
	if err != nil {
		flags.Error = fmt.Errorf("reminder: faild to write file; %w", err)
		return flags
	}

	return flags
}

func (a AddReminder) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}
