package command

import (
	"encoding/json"
	"fmt"
	"time"

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
	QueueName: "default",
	Blocking:  false,
	Active:    true,
}

var DiscordLink = Command{
	Name:        "Discord",
	Triggers:    []string{"!discord", "!dc"},
	Description: "Get the discord invite link.",
	Usage:       "!discord",
	Actions: []action.Action{
		actions.CreateDiscordInvite.Make([]action.ActionData{{Type: action.ChatMessageType, Data: nil}}),
	},
	Active: true,
}
var ShoutOutCommand = Command{
	Name:     "Shout out",
	Triggers: []string{"!so"},
	Actions: []action.Action{
		actions.ShoutOut.Make([]action.ActionData{{Type: action.ChatMessageType, Data: nil}}),
	},
	UsersList: []string{moderator},
	Active:    true,
}

var Lurk = Command{
	Name:        "Lurk",
	Triggers:    []string{"!lurk", "!lurking"},
	Description: "Thanks the viewer for lurking.",
	Usage:       "!lurk",
	Actions: []action.Action{
		actions.CreateLurkText.Make([]action.ActionData{{Type: action.ChatMessageType, Data: nil}}),
		actions.SendMessage.Make([]action.ActionData{{Type: action.StringType, Data: nil}}),
	},
	Active: true,
}

var CommandsCommand = Command{
	Name:        "Commands",
	Triggers:    []string{"!commands"},
	Description: "Gives a link to the webpage with all the commands",
	Usage:       "!commands",
	Actions: []action.Action{
		actions.SendMessage.Make([]action.ActionData{{Type: action.StringType, Data: "Check out all my wondefull commands :3 (lazyturtle33.live/commands)"}}),
	},
	Active: true,
}

var TTSCommand = Command{
	Name:     "TTS",
	Triggers: []string{"!tts"},
	Actions: []action.Action{
		actions.SendTTS.Make(nil),
	},
	QueueName: "tts",
	Active:    true,
}

var BRBCommand = Command{
	Name:     "Be Right Back",
	Triggers: []string{"!brb"},
	Actions: []action.Action{
		actions.SetScene.Make([]action.ActionData{{Type: action.StringType, Data: "BRB"}}),
		actions.RunAD.Make(nil),
	},
	Active:    true,
	UsersList: []string{moderator},
}

var YappyChatCommad = Command{
	Name:     "Yappy Chat",
	Triggers: []string{"!yappychat", "!yc"},
	Actions: []action.Action{
		actions.SetYappyChat.Make([]action.ActionData{{Type: action.StringType, Data: "toggle"}}),
		actions.IfBool.Make([]action.ActionData{
			{Type: action.ActionsType, Data: []action.Action{
				actions.ReplyToMessage.Make([]action.ActionData{{Type: action.ChatMessageType, Data: nil}, {Type: action.StringType, Data: "Yappy Chat is now on :D"}}),
			}},
			{Type: action.ActionsType, Data: []action.Action{
				actions.ReplyToMessage.Make([]action.ActionData{{Type: action.ChatMessageType, Data: nil}, {Type: action.StringType, Data: "Yappy Chat is now off :c"}}),
			}},
		}),
	},
	Active:    true,
	UsersList: []string{moderator},
}

var ShockCommand = Command{
	Name:     "Shock",
	Triggers: []string{"!shock"},
	Actions: []action.Action{
		actions.Shock.Make([]action.ActionData{{Type: action.ChatMessageType, Data: nil}}),
	},
	Active: true,
}

var SupportMeCommand = Command{
	Name:        "Support Me",
	Triggers:    []string{"!tip", "!support", "!kofi", "!kofi", "!donate"},
	Description: "Get the links you can support me at <3",
	Usage:       "!kofi",
	Actions: []action.Action{
		actions.ReplyToMessage.Make([]action.ActionData{{Type: action.ChatMessageType, Data: nil}, {Type: action.StringType, Data: "so many places to suport me :o pick your poison ig: https://ko-fi.com/thelazyturtle33, https://wish.ly/thelazyturtle33, https://fansly.com/LazyTurtle33"}}),
	},
	Active: true,
}

var TimeCommand = Command{
	Name:        "Time",
	Triggers:    []string{"!time"},
	Description: "Get my current local time",
	Usage:       "!time",
	Actions: []action.Action{
		actions.Function.Make([]action.ActionData{{Type: "function", Data: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
			flags := action.Flags{}
			flags.PassThrough = action.ActionData{Type: action.StringType, Data: fmt.Sprintf("It is currently %s for me :3", time.Now().Add(2*time.Hour).Format("15:04:05"))}
			return flags
		}}}),
		actions.ReplyToMessage.Make([]action.ActionData{{Type: action.ChatMessageType, Data: nil}}),
	},
	Active: true,
}

var SocalsCommand = Command{
	Name:        "Socals",
	Triggers:    []string{"!socals", "!socal", "!tiktok", "!youtube", "!yt"},
	Description: "Get the links to all my socals",
	Usage:       "!socals",
	Actions: []action.Action{
		actions.ReplyToMessage.Make([]action.ActionData{{Type: action.ChatMessageType, Data: nil}, {Type: action.StringType, Data: "Everywhere you can find me: https://www.youtube.com/@The_LazyTurtle33, https://www.tiktok.com/@thelazyturtle33, https://bsky.app/profile/lazyturtle33.bsky.social,  https://fansly.com/LazyTurtle33"}}),
	},
	Active: true,
}

var ReminderCommand = Command{
	Name:     "Reminder",
	Triggers: []string{"!remind", "!r"},
	Usage:    "!reminder Message",
	Actions: []action.Action{
		actions.AddReminder.Make([]action.ActionData{{Type: action.ChatMessageType, Data: nil}}),
	},
	UsersList: []string{moderator},
	Active:    true,
}

var CommandTriggers = map[string]*Command{}
var Commands = []*Command{
	&TestCommand,
	&DiscordLink,
	&Lurk,
	&CommandsCommand,
	&SupportMeCommand,
	&TTSCommand,
	&ShoutOutCommand,
	&BRBCommand,
	&YappyChatCommad,
	&ShockCommand,
	&TimeCommand,
	&SocalsCommand,
	&ReminderCommand,
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
