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

func (a CreateDiscordInvite) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
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
		ActionData: [][]any{{chat}},
	}
	return flags
}

func schedualedMessage(flags action.Flags) action.Flags {
	logger.Debug("scheduald Discord mesage", "last chatter id", context.Get().LastChat.ChatterUserID, "bot id", context.Get().GetBot().Id)
	if context.Get().LastChat.ChatterUserID != context.Get().GetBot().Id {
		flags.PassThrough = "Come check out the discord if ya wanna hang after stream ^w^ (" + InviteLink + ")"
	} else {
		logger.Log("Last Message was from bot, sleeping")
		flags.AddActions = action.Flag{
			Active: true,
			Actions: []action.Action{
				&Delay{Duration: 60 * 5},
			},
		}
		flags.Skip.Active = true
	}
	return flags
}

func (a CreateDiscordInvite) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}
