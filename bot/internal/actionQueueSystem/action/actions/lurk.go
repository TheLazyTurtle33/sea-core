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

func (a *CreateLurkText) Run(passThough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThough}
	if len(v) == 0 {
		flags.Error = fmt.Errorf("no data provided")
		return flags
	}
	data, ok := v[0].(datatypes.ChatMessageData)
	if !ok {
		flags.Error = fmt.Errorf("expected ChatMessageData or nil, got %T", v[0])
		flags.Skip = action.Flag{Active: true, IntData: 1}
		return flags
	}
	if data.SourceBroadcasterUserID != "" && data.SourceBroadcasterUserID != context.Get().GetBroadcaster().Id {
		flags.Skip = action.Flag{Active: true, IntData: 1}
		return flags
	}
	flags.PassThrough = "O: thx for the lurk " + data.ChatterUserName + " you da best <3 ^w^"
	return flags

}

func (a *CreateLurkText) OnAdd(v ...any) action.Flags {
	flags := action.Flags{}
	return flags
}
