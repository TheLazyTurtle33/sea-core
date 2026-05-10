package actions

import (
	"fmt"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
)

var ShoutOut = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		if len(actionData) == 0 || actionData[0].Type != action.ChatMessageType {
			flags.Error = fmt.Errorf("ShoutOut: expected chat message")
			return flags
		}

		chat, ok := actionData[0].Data.(datatypes.ChatMessageData)
		if !ok {
			flags.Error = fmt.Errorf("ShoutOut: expected chat message, got %T", actionData[0].Data)
			return flags
		}

		if len(chat.Message.Fragments) < 2 || chat.Message.Fragments[1].Mention.UserName == "" {
			flags.AddActions = action.Flag{
				Active: true,
				Actions: []action.Action{
					ReplyToMessage.Make([]action.ActionData{
						{Type: action.ChatMessageType, Data: chat},
						{Type: action.StringType, Data: "to use !so please @ person to shout out :3"},
					}),
				},
			}
			return flags
		}

		user := chat.Message.Fragments[1].Mention.UserName
		flags.PassThrough = action.ActionData{Type: action.StringType, Data: chat.Message.Fragments[1].Mention.UserID}
		flags.AddActions = action.Flag{
			Active: true,
			Actions: []action.Action{
				SendAnnouncement.Make([]action.ActionData{{Type: action.StringType, Data: fmt.Sprintf("HAY EVERYONE! go checkout @%s channel :3", user)}}),
				SendShoutOut.Make([]action.ActionData{{Type: action.StringType, Data: chat.Message.Fragments[1].Mention.UserID}}),
			},
		}
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "ShoutOut",
		Description: "Shout out a user on Twitch.",
		RunData: []action.ActionData{
			{
				Type:        action.ChatMessageType,
				Description: "Chat message containing the shoutout mention.",
				Required:    true,
			},
		},
	},
}

var SendShoutOut = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		targetUserID, err := parseStringParam(passThrough, actionData, 0)
		if err != nil {
			flags.Error = err
			return flags
		}

		_, err = twitchapi.AsBot().Post(
			fmt.Sprintf(
				"/chat/shoutouts?from_broadcaster_id=%s&to_broadcaster_id=%s&moderator_id=%s",
				context.Get().GetBroadcaster().Id,
				targetUserID,
				context.Get().GetBot().Id,
			),
			"",
		)
		if err != nil {
			flags.Error = err
		}
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "SendShoutOut",
		Description: "Send the Twitch shoutout API request.",
		RunData: []action.ActionData{
			{
				Type:        action.StringType,
				Description: "Target broadcaster user ID.",
				Required:    true,
			},
		},
	},
}

func init() {
	action.ActionMap[ShoutOut.MetaData.Name] = ShoutOut
	action.ActionMap[SendShoutOut.MetaData.Name] = SendShoutOut
}
