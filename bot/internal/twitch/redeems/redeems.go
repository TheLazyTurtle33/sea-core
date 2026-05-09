package redeems

// default actions and queue for redeems defined outside of the bot.
// var DefaultExternalRedeem = Redeem{
// 	Actions:   []action.Action{&actions.Log{Message: "External redeem triggerd", Formatter: redeemLogFormatter}},
// 	QueueName: "redeems",
// }

// func redeemLogFormatter(data []any) ([]any, error) {
// 	if len(data) == 0 {
// 		return data, fmt.Errorf("no data provided")
// 	}
// 	redeem, ok := data[0].(datatypes.RedemptionData)
// 	if !ok {
// 		return data, fmt.Errorf("expected RedemptionData, got %T", data[0])
// 	}
// 	data = data[1:] // remove redeem data from data passed to log action
// 	return []any{
// 		"ID", redeem.ID,
// 		"User", redeem.UserLogin,
// 		"Reward", redeem.Reward.Title,
// 		"Reward_ID", redeem.Reward.ID,
// 	}, nil
// }

// example of an external redeem override. This is used to added extar logic to a redeem defined outside of the bot.
// var HydraRedeem = Redeem{
// 	Actions:   []action.Action{&actions.SendMessage{Message: "WATER BITCH BOI >:3"}, &actions.Log{Message: "Override redeem triggerd", Formatter: redeemLogFormatter}},
// 	QueueName: "redeems",
// 	ID:        "0f58303d-e5d8-4c24-87c5-b2558017a98f", // key making this an override is that an ID is provided
// }

// var TestRedeem = Redeem{
// 	Actions:         []action.Action{&actions.SendMessage{Message: "This is a test redeem!"}, &actions.CompleteRedeemAction{}, &actions.Log{Message: "Internal redeem triggerd", Formatter: redeemLogFormatter}},
// 	QueueName:       "redeems",
// 	Title:           "Test Redeem",
// 	Cost:            10,
// 	BackgroundColor: "#00BFFF",
// 	IsEnabled:       true,
// 	IsPaused:        false,
// }

// var LogInRedeem = Redeem{
// 	Title:           "temp name",
// 	Cost:            1000,
// 	BackgroundColor: "#00BFFF",
// 	IsEnabled:       true,
// 	IsPaused:        false,
// 	QueueName:       "redeems",
// 	Actions:         []action.Action{&actions.LogIn{}, &actions.CompleteRedeemAction{}},
// }

// var JumpScare = Redeem{
// 	Title:           "Jump Scare",
// 	Cost:            1000,
// 	BackgroundColor: "#FF0000",
// 	IsEnabled:       true,
// 	IsPaused:        false,
// 	QueueName:       "redeems",
// 	Actions:         []action.Action{&actions.JumpScareAction{}},
// }

var definedRedeems = map[string]Redeem{
	// TestRedeem.Title:  TestRedeem,
	// LogInRedeem.Title: LogInRedeem,
	// JumpScare.Title:   JumpScare,
}

var definedOverrides = map[string]Redeem{
	// HydraRedeem.ID: HydraRedeem,
}
