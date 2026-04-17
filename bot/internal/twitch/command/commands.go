package command

import (
	"encoding/json"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action/actions"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/tts"
)

var TestCommand = Command{
	Name:        "Test",
	Triggers:    []string{"!test"},
	Description: "A test command.",
	Usage:       "!test",
	Actions:     []action.Action{&actions.ReplyToMessage{Message: "test passed! :3"}},
	AlowedUsers: []string{"everyone"},
	QueueName:   "default",
	Blocking:    false,
	Active:      true,
}

var DiscordLink = Command{
	Name:        "Discord",
	Triggers:    []string{"!discord", "!dc"},
	Description: "Get the discord invite link.",
	Usage:       "!discord",
	Actions:     []action.Action{&actions.CreateDiscordInvite{}},
	AlowedUsers: []string{"everyone"},
	Active:      true,
}
var ShoutOutCommand = Command{
	Name:        "Shout out",
	Triggers:    []string{"!so"},
	Actions:     []action.Action{&actions.ShoutOut{}},
	AlowedUsers: []string{"moderator"},
	Active:      true,
}

var Lurk = Command{
	Name:        "Lurk",
	Triggers:    []string{"!lurk", "!lurking"},
	Description: "Thanks the viewer for lurking.",
	Usage:       "!lurk",
	Actions:     []action.Action{&actions.CreateLurkText{}, &actions.SendMessage{}},
	AlowedUsers: []string{"everyone"},
	Active:      true,
}

var CommandsCommand = Command{
	Name:        "Commands",
	Triggers:    []string{"!commands"},
	Description: "Gives a link to the webpage with all the commands",
	Usage:       "!commands",
	Actions:     []action.Action{&actions.SendMessage{Message: "Check out all my wondefull commands :3 (lazyturtle33.live/commands)"}},
	AlowedUsers: []string{"everyone"},
	Active:      true,
}

var KofiCommand = Command{
	Name:        "Ko-fi",
	Triggers:    []string{"!kofi", "!ko-fi"},
	Description: "Get the Ko-fi link. Used for TTS as well.",
	Usage:       "!kofi",
	Actions:     []action.Action{&actions.ReplyToMessage{Message: "If you'd like to support me (or send a TTS) chu can do that here ^w^ : https://ko-fi.com/thelazyturtle33"}},
	AlowedUsers: []string{"everyone"},
	Active:      true,
}

var TTSCommand = Command{
	Name:     "TTS",
	Triggers: []string{"!tts"},
	Actions:  []action.Action{&tts.TTS{Message: ""}},
	// Actions:  []action.Action{&actions.ReplyToMessage{Message: "still making TTS X.X try again when it ready :3 (or use !kofi for old tts)"}},
	Active: true,
}

var BRBCommand = Command{
	Name:        "Be Right Back",
	Triggers:    []string{"!brb"},
	Actions:     []action.Action{&actions.SetScene{Scene: "BRB"}, &actions.RunAD{}},
	Active:      true,
	AlowedUsers: []string{"moderator"},
}

var YappyChatCommad = Command{
	Name:     "Yappy Chat",
	Triggers: []string{"!yappychat", "!yc"},
	Actions: []action.Action{
		&actions.SetYappyChat{Toggle: true},
		&actions.IfBool{
			TrueActoins:  []action.Action{&actions.ReplyToMessage{Message: "Yappy Chat is now on!"}},
			FalseActoins: []action.Action{&actions.ReplyToMessage{Message: "Yappy Chat is now off :c"}},
		},
	},
	Active:      true,
	AlowedUsers: []string{"moderator"},
}

var CommandTriggers = map[string]*Command{}
var Commands = []*Command{
	&TestCommand,
	&DiscordLink,
	&Lurk,
	&CommandsCommand,
	&KofiCommand,
	&TTSCommand,
	&ShoutOutCommand,
	&BRBCommand,
	&YappyChatCommad,
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
