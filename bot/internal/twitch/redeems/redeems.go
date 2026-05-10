package redeems

import (
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action/actions"
)

// default actions and queue for redeems defined outside of the bot.
var DefaultExternalRedeem = Redeem{
	Actions: []action.Action{
		actions.Log.Make([]action.ActionData{{Type: action.StringType, Data: "External redeem triggered"}}),
	},
	QueueName: "redeems",
}

// example of an external redeem override. This is used to add extra logic to a redeem defined outside of the bot.
var HydraRedeem = Redeem{
	Actions: []action.Action{
		actions.SendMessage.Make([]action.ActionData{{Type: action.StringType, Data: "WATER BITCH BOI >:3"}}),
		actions.Log.Make([]action.ActionData{{Type: action.StringType, Data: "Override redeem triggered"}}),
	},
	QueueName: "redeems",
	ID:        "0f58303d-e5d8-4c24-87c5-b2558017a98f",
}

var TestRedeem = Redeem{
	Actions: []action.Action{
		actions.SendMessage.Make([]action.ActionData{{Type: action.StringType, Data: "This is a test redeem!"}}),
		actions.CompleteRedeemAction.Make([]action.ActionData{{Type: action.RedeemType, Data: nil}}),
		actions.Log.Make([]action.ActionData{{Type: action.StringType, Data: "Internal redeem triggered"}}),
	},
	QueueName:       "redeems",
	Title:           "Test Redeem",
	Cost:            10,
	BackgroundColor: "#00BFFF",
	IsEnabled:       true,
	IsPaused:        false,
}

var LogInRedeem = Redeem{
	Title:           "temp name",
	Cost:            1000,
	BackgroundColor: "#00BFFF",
	IsEnabled:       true,
	IsPaused:        false,
	QueueName:       "redeems",
	Actions: []action.Action{
		actions.LogIn.Make(nil),
		actions.CompleteRedeemAction.Make([]action.ActionData{{Type: action.RedeemType, Data: nil}}),
	},
}

var JumpScare = Redeem{
	Title:           "Jump Scare",
	Cost:            1000,
	BackgroundColor: "#FF0000",
	IsEnabled:       true,
	IsPaused:        false,
	QueueName:       "redeems",
	Actions: []action.Action{
		actions.JumpScareAction.Make(nil),
	},
}

var definedRedeems = map[string]Redeem{
	TestRedeem.Title:  TestRedeem,
	LogInRedeem.Title: LogInRedeem,
	JumpScare.Title:   JumpScare,
}

var definedOverrides = map[string]Redeem{
	HydraRedeem.ID: HydraRedeem,
}
