package actions

import (
	"fmt"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
	"github.com/TheLazyTurtle33/sea-core/shared/obs"
)

var JumpScareAction = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		client, err := obs.Get()
		if err != nil {
			flags.Error = fmt.Errorf("JumpScareAction: failed to get OBS client: %w", err)
			flags.Skip.Active = true
			return flags
		}

		if _, err := client.SetSourceVisability("JumpScare", true); err != nil {
			flags.Error = fmt.Errorf("JumpScareAction: failed to set visibility true: %w", err)
			flags.Skip.Active = true
			return flags
		}

		go func() {
			time.Sleep(4 * time.Second)
			if _, err := client.SetSourceVisability("JumpScare", false); err != nil {
				logger.Error("JumpScareAction: failed to set visibility false", err)
			}
		}()

		flags.AddActions = action.Flag{
			Active: true,
			Actions: []action.Action{
				CompleteRedeemAction.Make(nil),
			},
		}
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "JumpScareAction",
		Description: "Trigger the OBS jump scare scene and complete the redeem.",
		RunData:     []action.ActionData{},
	},
}

func init() {
	action.ActionMap[JumpScareAction.MetaData.Name] = JumpScareAction
}
