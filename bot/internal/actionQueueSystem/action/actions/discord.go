package actions

import (
	"fmt"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
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
	data, ok := v[0].(datatypes.ChatMessageData)
	if !ok {
		flags.Error = fmt.Errorf("expected ChatMessageData, got %T", v[2])
		return flags
	}
	if data.SourceBroadcasterUserID != "" {
		if data.SourceBroadcasterUserID != context.Get().Broadcaster.Id {
			flags.DataForNextAction = "Join my Discord as well while your at it ;3 (" + InviteLink + ")"
			return flags
		}
	}
	flags.DataForNextAction = "YESSS join the server >:3 (" + InviteLink + ")"
	return flags
}

func (a *CreateDiscordInvite) OnAdd(v ...any) action.Flags {
	flags := action.Flags{}
	return flags
}
