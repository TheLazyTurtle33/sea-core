package actions

import (
	"fmt"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
)

type CreateDiscordInvite struct {
	action.Action
}

const InviteLink = "https://discord.gg/GJpybRAEq5"

func (a *CreateDiscordInvite) Run(v ...any) action.Flags {
	flags := action.Flags{}
	v = v[1:] // remove pervious action outup data from v
	if len(v) == 0 {
		flags.Error = fmt.Errorf("no data provided")
		return flags
	}
	switch v[0].(type) {
	case datatypes.ChatMessageData:
		data := v[0].(datatypes.ChatMessageData)
		if data.SourceBroadcasterUserID != "" {
			if data.SourceBroadcasterUserID != context.Get().GetBroadcaster().Id {
				flags.DataForNextAction = "Join my Discord as well while your at it ;3 (" + InviteLink + ")"
				return flags
			}
		}
		flags.DataForNextAction = "YESSS join the server >:3 (" + InviteLink + ")"
		return flags
	case nil:
		if context.Get().LastChat.ChatterUserID != context.Get().GetBot().Id {
			flags.DataForNextAction = "Come check out the discord if ya wanna hang after stream ^w^ (" + InviteLink + ")"
		} else {
			logger.Log("Last Message was from bot, sleeping")
			flags.Pause = action.Flag{Active: true, IntData: 60 * 5}
		}
		return flags
	default:
		flags.Error = fmt.Errorf("expected ChatMessageData or nil, got %T", v[0])
		flags.Skip = action.Flag{Active: true, IntData: 1}
		return flags
	}

}

func (a *CreateDiscordInvite) OnAdd(v ...any) action.Flags {
	flags := action.Flags{}
	return flags
}
