package redeems

import (
	"fmt"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action/actions"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
)

// default actions and queue for redeems defined outside of the bot.
var DefaultExternalRedeem = Redeem{
	Actions:   []action.Action{&actions.Log{Message: "External redeem triggerd", Formatter: redeemLogFormatter}},
	QueueName: "redeems",
}

func redeemLogFormatter(data []any) ([]any, error) {
	if len(data) == 0 {
		return data, fmt.Errorf("no data provided")
	}
	redeem, ok := data[0].(datatypes.RedemptionData)
	if !ok {
		return data, fmt.Errorf("expected RedemptionData, got %T", data[0])
	}
	data = data[1:] // remove redeem data from data passed to log action
	return []any{
		"ID", redeem.ID,
		"User", redeem.UserLogin,
		"Reward", redeem.Reward.Title,
		"Reward_ID", redeem.Reward.ID,
	}, nil
}

// example of an external redeem override. This is used to added extar logic to a redeem defined outside of the bot.
var HydraRedeem = Redeem{
	Actions:   []action.Action{&actions.SendMessage{Message: "WATER BITCH BOI >:3"}, &actions.Log{Message: "Override redeem triggerd", Formatter: redeemLogFormatter}},
	QueueName: "redeems",
	ID:        "0f58303d-e5d8-4c24-87c5-b2558017a98f", // key making this an override is that an ID is provided
}

var TestRedeem = Redeem{
	Actions:         []action.Action{&actions.SendMessage{Message: "This is a test redeem!"}, &CompleteRedeemAction{}, &actions.Log{Message: "Internal redeem triggerd", Formatter: redeemLogFormatter}},
	QueueName:       "redeems",
	Title:           "Test Redeem",
	Cost:            10,
	BackgroundColor: "#00BFFF",
	IsEnabled:       true,
	IsPaused:        false,
}

var definedRedeems = map[string]Redeem{
	TestRedeem.Title: TestRedeem,
}

var definedOverrides = map[string]Redeem{
	HydraRedeem.ID: HydraRedeem,
}

type CompleteRedeemAction struct {
	action.Action
}

func (a *CompleteRedeemAction) Run(passThough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThough}
	if len(v) == 0 {
		flags.Error = fmt.Errorf("no data provided")
		return flags
	}
	data, ok := v[0].(datatypes.RedemptionData)
	if !ok {
		flags.Error = fmt.Errorf("expected RedemptionData, got %T", v[0])
		return flags
	}
	_, err := twitchapi.AsUser().Patch("/channel_points/custom_rewards/redemptions?broadcaster_id="+context.Get().GetBroadcaster().Id+"&reward_id="+data.Reward.ID+"&id="+data.ID, `{"status":"FULFILLED"}`)
	if err != nil {
		flags.Error = err
		return flags
	}
	logger.Log("Redeem completed", "ID", data.Reward.ID) // , "Response", string(resp)

	return flags
}

func (a *CompleteRedeemAction) OnAdd(v ...any) action.Flags {
	flags := action.Flags{}
	return flags
}
