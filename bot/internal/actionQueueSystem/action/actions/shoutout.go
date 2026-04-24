package actions

import (
	"fmt"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
)

type ShoutOut struct {
	action.Action
}

func (a ShoutOut) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}

	if len(v) == 0 {
		flags.Error = fmt.Errorf("Expectided data")
		return flags
	}

	chat, ok := v[0].(datatypes.ChatMessageData)
	if !ok {
		flags.Error = fmt.Errorf("Expectided Chat data type")
		return flags
	}

	if len(chat.Message.Fragments) < 2 || chat.Message.Fragments[1].Mention.UserName == "" {
		flags.Error = fmt.Errorf("Expectided two fragments")
		message := "to use !so please @ person to shout out :3"
		flags.AddActions = action.Flag{
			Active: true,
			Actions: []action.Action{
				&ReplyToMessage{Message: message},
			},
			ActionData: [][]any{{chat}},
		}
		return flags
	}

	user := chat.Message.Fragments[1].Mention.UserName
	flags.AddActions = action.Flag{
		Active: true,
		Actions: []action.Action{
			&SendAnnouncement{
				Message: fmt.Sprintf("HAY EVERYONE! go checkout @%s channel :3", user),
			},
			&SendShoutOut{},
		},
	}
	flags.PassThrough = chat.Message.Fragments[1].Mention.UserID
	return flags
}

func (a ShoutOut) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}

type SendShoutOut struct {
	action.Action
	ShoutOutID string
}

func (a SendShoutOut) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}

	if a.ShoutOutID == "" {
		id, ok := passThrough.(string)
		if !ok {
			flags.Error = fmt.Errorf("faild to get id string from passthrough")
			return flags
		}
		a.ShoutOutID = id
	}

	_, err := twitchapi.AsBot().Post(
		fmt.Sprintf(
			"/chat/shoutouts?from_broadcaster_id=%s&to_broadcaster_id=%s&moderator_id=%s",
			context.Get().GetBroadcaster().Id,
			a.ShoutOutID,
			context.Get().GetBot().Id,
		),
		"",
	)

	if err != nil {
		flags.Error = err
	}

	return flags
}
func (a SendShoutOut) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}
