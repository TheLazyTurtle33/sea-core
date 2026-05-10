package actions

import (
	"fmt"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
)

const InviteLink = "https://discord.gg/GJpybRAEq5"

var CreateDiscordInvite = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		if len(actionData) == 0 {
			return scheduledDiscordMessage(flags)
		}
		if actionData[0].Type == action.ChatMessageType {
			if chat, ok := actionData[0].Data.(datatypes.ChatMessageData); ok {
				return discordCommand(flags, chat)
			}
		}
		return scheduledDiscordMessage(flags)
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "CreateDiscordInvite",
		Description: "Send or schedule a Discord invite.",
		RunData: []action.ActionData{
			{
				Type:        action.ChatMessageType,
				Description: "Chat message that triggered the invite.",
				Required:    false,
			},
		},
	},
}

func discordCommand(flags action.Flags, chat datatypes.ChatMessageData) action.Flags {
	message := fmt.Sprintf("YESSS join the server >:3 (%s)", InviteLink)
	if chat.SourceBroadcasterUserID != "" && chat.SourceBroadcasterUserID != context.Get().GetBroadcaster().Id {
		message = fmt.Sprintf("Join my Discord as well while you're at it ;3 (%s)", InviteLink)
	}
	flags.AddActions = action.Flag{
		Active: true,
		Actions: []action.Action{
			ReplyToMessage.Make([]action.ActionData{
				{Type: action.ChatMessageType, Data: chat},
				{Type: action.StringType, Data: message},
			}),
		},
	}
	return flags
}

func scheduledDiscordMessage(flags action.Flags) action.Flags {
	logger.Debug("scheduled Discord message", "last chatter id", context.Get().LastChat.ChatterUserID, "bot id", context.Get().GetBot().Id)
	if context.Get().LastChat.ChatterUserID != context.Get().GetBot().Id {
		flags.PassThrough = action.ActionData{Type: action.StringType, Data: fmt.Sprintf("Come check out the discord if ya wanna hang after stream ^w^ (%s)", InviteLink)}
		return flags
	}

	logger.Log("Last Message was from bot, sleeping")
	flags.AddActions = action.Flag{
		Active: true,
		Actions: []action.Action{
			Delay.Make([]action.ActionData{{Type: "int", Data: 5 * 60}}),
		},
	}
	flags.Skip.Active = true
	return flags
}

func init() {
	action.ActionMap[CreateDiscordInvite.MetaData.Name] = CreateDiscordInvite
}
