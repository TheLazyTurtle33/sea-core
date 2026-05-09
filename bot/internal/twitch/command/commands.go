package command

import (
	"encoding/json"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action/actions"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
)

const (
	moderator = "moderator"
	everyone  = "everyone"
)

var TestCommand = Command{
	Name:        "Test",
	Triggers:    []string{"!test"},
	Description: "A test command.",
	Usage:       "!test",
	Actions: []action.Action{
		actions.ReplyToMessage.Make(
			[]action.ActionData{
				{
					Type: action.ChatMessageType,
					Data: nil,
				},
				{
					Type: action.StringType,
					Data: "Test Passed! :3",
				},
			}),
	},

	AlowedUsers: []string{everyone},
	QueueName:   "default",
	Blocking:    false,
	Active:      true,
}

// var DiscordLink = Command{
// 	Name:        "Discord",
// 	Triggers:    []string{"!discord", "!dc"},
// 	Description: "Get the discord invite link.",
// 	Usage:       "!discord",
// 	Actions:     []action.Action{&actions.CreateDiscordInvite{}},
// 	Active:      true,
// }
// var ShoutOutCommand = Command{
// 	Name:        "Shout out",
// 	Triggers:    []string{"!so"},
// 	Actions:     []action.Action{&actions.ShoutOut{}},
// 	AlowedUsers: []string{moderator},
// 	Active:      true,
// }

// var Lurk = Command{
// 	Name:        "Lurk",
// 	Triggers:    []string{"!lurk", "!lurking"},
// 	Description: "Thanks the viewer for lurking.",
// 	Usage:       "!lurk",
// 	Actions:     []action.Action{&actions.CreateLurkText{}, &actions.SendMessage{}},
// 	Active:      true,
// }

// var CommandsCommand = Command{
// 	Name:        "Commands",
// 	Triggers:    []string{"!commands"},
// 	Description: "Gives a link to the webpage with all the commands",
// 	Usage:       "!commands",
// 	Actions:     []action.Action{&actions.SendMessage{Message: "Check out all my wondefull commands :3 (lazyturtle33.live/commands)"}},
// 	Active:      true,
// }

// var TTSCommand = Command{
// 	Name:      "TTS",
// 	Triggers:  []string{"!tts"},
// 	Actions:   []action.Action{&tts.TTS{}},
// 	QueueName: "tts",
// 	Active:    true,
// }

// var BRBCommand = Command{
// 	Name:        "Be Right Back",
// 	Triggers:    []string{"!brb"},
// 	Actions:     []action.Action{&actions.SetScene{Scene: "BRB"}, &actions.RunAD{}},
// 	Active:      true,
// 	AlowedUsers: []string{moderator},
// }

// var YappyChatCommad = Command{
// 	Name:     "Yappy Chat",
// 	Triggers: []string{"!yappychat", "!yc"},
// 	Actions: []action.Action{
// 		{
// 			Action:     &actions.SetYappyChat{Toggle: true},
// 			ActionData: nil,
// 		},
// 		{
// 			Action: &actions.IfBool{},
// 			ActionData: []action.ActionData{
// 				{
// 					Type: ActionsType,
// 					Data: []action.ActionFull{
// 						{
// 							Action: &actions.ReplyToMessage{Message: "Yappy Chat is now on :D"},
// 							ActionData: []action.ActionData{
// 								{
// 									Type: ChatMessageDataType,
// 									Data: nil,
// 								},
// 								{
// 									Type: "string",
// 									Data: "Yappy Chat is now on :D",
// 								},
// 							},
// 						},
// 					}, // true actions

// 				},
// 				{
// 					Type: ActionsType,
// 					Data: []action.Action{
// 						{
// 							Action: &actions.ReplyToMessage{Message: "Yappy Chat is now off :c"},
// 							ActionData: []action.ActionData{
// 								{
// 									Type: ChatMessageDataType,
// 									Data: nil,
// 								},
// 								{
// 									Type: "string",
// 									Data: "Yappy Chat is now off :c",
// 								},
// 							},
// 						},
// 					}, // false actions
// 				},
// 			},
// 		},
// 	},
// 	Active:      true,
// 	AlowedUsers: []string{moderator},
// }

// var ShockCommand = Command{
// 	Name:     "Shock",
// 	Triggers: []string{"!shock"},
// 	Actions:  []action.Action{&actions.Shock{}},
// 	Active:   true,
// }

// var SupportMeCommand = Command{
// 	Name:        "Support Me",
// 	Triggers:    []string{"!tip", "!support", "!kofi", "!kofi", "!donate"},
// 	Description: "Get the links you can support me at <3",
// 	Usage:       "!kofi",
// 	Actions:     []action.Action{&actions.ReplyToMessage{Message: "so many places to suport me :o pick your poison ig: https://ko-fi.com/thelazyturtle33, https://wish.ly/thelazyturtle33, https://fansly.com/LazyTurtle33"}},
// 	Active:      true,
// }

// var TimeCommand = Command{
// 	Name:        "Time",
// 	Triggers:    []string{"!time"},
// 	Description: "Get my current local time",
// 	Usage:       "!time",
// 	Actions: []action.Action{
// 		&actions.Functoin{
// 			Fn: func(passThrough any, v ...any) action.Flags {
// 				flags := action.Flags{}
// 				flags.PassThrough = fmt.Sprintf("It is currently %s for me :3", time.Now().Add(2*time.Hour).Format("15:04:05"))
// 				return flags
// 			},
// 		},
// 		&actions.ReplyToMessage{},
// 	},
// 	Active: true,
// }

// var SocalsCommand = Command{
// 	Name:        "Socals",
// 	Triggers:    []string{"!socals", "!socal", "!tiktok", "!youtube", "!yt"},
// 	Description: "Get the links to all my socals",
// 	Usage:       "!socals",
// 	Actions:     []action.Action{&actions.ReplyToMessage{Message: "Everywhere you can find me: https://www.youtube.com/@The_LazyTurtle33, https://www.tiktok.com/@thelazyturtle33, https://bsky.app/profile/lazyturtle33.bsky.social,  https://fansly.com/LazyTurtle33"}},
// 	Active:      true,
// }

// var ReminderCommand = Command{
// 	Name:        "Reminder",
// 	Triggers:    []string{"!remind", "!r"},
// 	Usage:       "!reminder Message",
// 	Actions:     []action.Action{&actions.AddReminder{}},
// 	AlowedUsers: []string{moderator},
// 	Active:      true,
// }

var CommandTriggers = map[string]*Command{}
var Commands = []*Command{
	&TestCommand,
	// &DiscordLink,
	// &Lurk,
	// &CommandsCommand,
	// &SupportMeCommand,
	// &TTSCommand,
	// &ShoutOutCommand,
	// &BRBCommand,
	// &YappyChatCommad,
	// &ShockCommand,
	// &TimeCommand,
	// &SocalsCommand,
	// &ReminderCommand,
}

func RegisterCommands() {
	if len(CommandTriggers) > 0 {
		return
	}
	for _, c := range Commands {
		for _, t := range c.Triggers {
			CommandTriggers[t] = c
		}
	}
}

func GetCommand(trigger string) *Command {
	return CommandTriggers[trigger]
}

func GetCommandsJson() []byte {
	out, err := json.Marshal(Commands)
	if err != nil {
		logger.Error("failed to marshal commands", err)
		return nil
	}
	return out

}
