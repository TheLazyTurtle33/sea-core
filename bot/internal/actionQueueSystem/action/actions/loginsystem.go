package actions

import (
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
)

var LogIn = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		logger.Log("hi from login")
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		flags.AddActions = action.Flag{
			Active: true,
			Actions: []action.Action{
				GetLogIns.Make(nil),
			},
		}
		return flags
	},
	MetaData: action.ActionMetaData{
		Name:        "LogIn",
		Description: "Execute a login workflow.",
	},
}

var GetLogIns = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		logger.Log("hi from get logins")
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "GetLogIns",
		Description: "Fetch login records.",
	},
}

func init() {
	action.ActionMap[LogIn.MetaData.Name] = LogIn
	action.ActionMap[GetLogIns.MetaData.Name] = GetLogIns
}
