package actions

import "github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"

var SendTTS = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "SendTTS",
		Description: "Placeholder TTS action.",
	},
}

func init() {
	action.ActionMap[SendTTS.MetaData.Name] = SendTTS
}
