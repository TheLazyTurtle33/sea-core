package actions

import (
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
)

type CreateDiscordInvite struct {
	action.Action
}

const InviteLink = "https://discord.gg/GJpybRAEq5"

func (a *CreateDiscordInvite) Run(passThough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThough}
	switch v[0].(type) {
	case datatypes.ChatMessageData:
		data := v[0].(datatypes.ChatMessageData)
		return command(flags, data)
	default:
		return schedualedMessage(flags)
	}
}

func command(flags action.Flags, chat datatypes.ChatMessageData) action.Flags {
	message := "YESSS join the server >:3 (" + InviteLink + ")"

	if chat.SourceBroadcasterUserID != "" {
		if chat.SourceBroadcasterUserID != context.Get().GetBroadcaster().Id {
			message = "Join my Discord as well while your at it ;3 (" + InviteLink + ")"
		}
	}
	flags.AddActions = action.Flag{
		Active: true,
		Actions: []action.Action{
			&ReplyToMessage{Message: message},
		},
		ActionData: []any{chat},
	}
	return flags
}

func schedualedMessage(flags action.Flags) action.Flags {
	logger.Debug("scheduald Discord mesage", "last chatter id", context.Get().LastChat.ChatterUserID, "bot id", context.Get().GetBot().Id)
	if context.Get().LastChat.ChatterUserID != context.Get().GetBot().Id {
		flags.AddActions = action.Flag{
			Active: true,
			Actions: []action.Action{
				&SendMessage{Message: "Come check out the discord if ya wanna hang after stream ^w^ (" + InviteLink + ")"},
			},
			ActionData: []any{nil},
		}
	} else {
		logger.Log("Last Message was from bot, sleeping")
		flags.AddActions = action.Flag{
			Active: true,
			Actions: []action.Action{
				&Delay{Duration: 60 * 5},
			},
			ActionData: []any{nil},
		}
	}
	return flags
}

func (a *CreateDiscordInvite) OnAdd(v ...any) action.Flags {
	flags := action.Flags{}
	return flags
}
