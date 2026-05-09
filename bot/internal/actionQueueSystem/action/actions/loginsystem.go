package actions

// import (
// 	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
// 	"github.com/TheLazyTurtle33/sea-core/shared/logger"
// )

// type LogIn struct {
// 	action.Action
// }

// func (a LogIn) Run(passThrough any, v ...any) action.Flags {
// 	flags := action.Flags{PassThrough: passThrough}
// 	// if len(v) == 0 {
// 	// 	flags.Error = fmt.Errorf("Expected ether chat data or redeem data. got nothing")
// 	// 	return flags
// 	// }

// 	// var user datatypes.User

// 	// if redeemData, ok := v[0].(datatypes.RedemptionData); ok {
// 	// 	user.Id = redeemData.UserID
// 	// 	user.Login = redeemData.UserLogin
// 	// 	user.DisplayName = redeemData.UserName
// 	// } else if chatData, ok := v[0].(datatypes.ChatMessageData); ok {
// 	// 	user.Id = chatData.ChatterUserID
// 	// 	user.Login = chatData.ChatterUserLogin
// 	// 	user.DisplayName = chatData.ChatterUserName
// 	// } else {
// 	// 	flags.Error = fmt.Errorf("Expected ether chat data or redeem data")
// 	// 	return flags
// 	// }
// 	logger.Log("hi from login")

// 	return flags
// }

// func (a LogIn) OnAdd(passThrough any, v ...any) action.Flags {
// 	flags := action.Flags{
// 		AddActions: action.Flag{
// 			Active: true,
// 			Actions: []action.Action{
// 				&GetLogIns{},
// 			},
// 			ActionData: [][]any{},
// 		},
// 	}
// 	return flags
// }

// type GetLogIns struct {
// 	action.Action
// }

// func (a GetLogIns) Run(passThrough any, v ...any) action.Flags {
// 	flags := action.Flags{PassThrough: passThrough}
// 	logger.Log("hi from get logins")
// 	return flags
// }

// func (a GetLogIns) OnAdd(passThrough any, v ...any) action.Flags {
// 	flags := action.Flags{PassThrough: passThrough}
// 	return flags
// }
