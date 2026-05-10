package actions

import (
	"fmt"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
)

var CreateLurkText = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		if len(actionData) == 0 || actionData[0].Type != action.ChatMessageType {
			flags.Error = fmt.Errorf("CreateLurkText: expected chat message")
			flags.Skip.Active = true
			return flags
		}

		chat, ok := actionData[0].Data.(datatypes.ChatMessageData)
		if !ok {
			flags.Error = fmt.Errorf("CreateLurkText: expected chat message, got %T", actionData[0].Data)
			flags.Skip.Active = true
			return flags
		}

		if chat.SourceBroadcasterUserID != "" && chat.SourceBroadcasterUserID != context.Get().GetBroadcaster().Id {
			flags.Skip.Active = true
			return flags
		}

		flags.PassThrough = action.ActionData{Type: action.StringType, Data: fmt.Sprintf("O: thx for the lurk %s you da best <3 ^w^", chat.ChatterUserName)}
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "CreateLurkText",
		Description: "Generate a lurk thank-you message.",
		RunData: []action.ActionData{
			{
				Type:        action.ChatMessageType,
				Description: "The chat message that triggered the lurk.",
				Required:    true,
			},
		},
	},
}

func init() {
	action.ActionMap[CreateLurkText.MetaData.Name] = CreateLurkText
}
