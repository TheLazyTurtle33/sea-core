package actions

import (
	"fmt"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
)

type CreateLurkText struct {
	action.Action
}

func (a *CreateLurkText) Run(v ...any) action.Flags {
	flags := action.Flags{}
	v = v[1:] // remove pervious action outup data from v
	if len(v) == 0 {
		flags.Error = fmt.Errorf("no data provided")
		return flags
	}
	switch v[0].(type) {
	case datatypes.ChatMessageData:
		data := v[0].(datatypes.ChatMessageData)
		if data.SourceBroadcasterUserID != "" && data.SourceBroadcasterUserID != context.Get().GetBroadcaster().Id {
			flags.Skip = action.Flag{Active: true, IntData: 1}
			return flags
		}
		flags.DataForNextAction = "O: thx for the lurk " + data.ChatterUserName + " you da best <3 ^w^"
		return flags
	default:
		flags.Error = fmt.Errorf("expected ChatMessageData or nil, got %T", v[0])
		flags.Skip = action.Flag{Active: true, IntData: 1}
		return flags
	}

}

func (a *CreateLurkText) OnAdd(v ...any) action.Flags {
	flags := action.Flags{}
	return flags
}
