package actions

import (
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/shared/obs"
)

var SetScene = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		scene, err := parseStringParam(passThrough, actionData, 0)
		if err != nil {
			flags.Error = err
			return flags
		}

		client, err := obs.Get()
		if err != nil {
			flags.Error = err
			return flags
		}

		response, err := client.SetScene(scene)
		if err != nil {
			flags.Error = err
			return flags
		}

		flags.PassThrough = action.ActionData{Type: action.StringType, Data: response}
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "SetScene",
		Description: "Switch the OBS scene.",
		RunData: []action.ActionData{
			{
				Type:        action.StringType,
				Description: "Scene name.",
				Required:    true,
			},
		},
	},
}

func init() {
	action.ActionMap[SetScene.MetaData.Name] = SetScene
}
