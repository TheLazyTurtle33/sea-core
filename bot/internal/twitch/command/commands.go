package command

import (
	"encoding/json"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action/actions"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
)

var TestCommand = Command{
	Name:        "Test",
	Triggers:    []string{"!test"},
	Description: "A test command.",
	Usage:       "!test",
	Actions:     []action.Action{&actions.ReplyToMessage{Message: "test passed! :3"}},
	QueueName:   "default",
	Blocking:    false,
}

var DiscordLink = Command{
	Name:        "Discord",
	Triggers:    []string{"!discord", "!dc"},
	Description: "Get the discord invite link.",
	Usage:       "!discord",
	Actions:     []action.Action{&actions.CreateDiscordInvite{}, &actions.ReplyToMessage{}},
	QueueName:   "default",
	Blocking:    false,
}

var Lurk = Command{
	Name:        "Lurk",
	Triggers:    []string{"!lurk"},
	Description: "Thanks the viewer for lurking.",
	Usage:       "!lurk",
	Actions:     []action.Action{&actions.CreateLurkText{}, &actions.SendMessage{}},
	QueueName:   "default",
	Blocking:    false,
}

var CommandsCommand = Command{
	Name:        "Commands",
	Triggers:    []string{"!commands"},
	Description: "Gives a link to the webpage with all the commands",
	Usage:       "!commands",
	Actions:     []action.Action{&actions.SendMessage{Message: "Check out all my wondefull commands :3 (lazyturtle33.live/commands)"}},
	QueueName:   "default",
	Blocking:    false,
}

var CommandTriggers = map[string]*Command{}
var Commands = []*Command{
	&TestCommand,
	&DiscordLink,
	&Lurk,
	&CommandsCommand,
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
		logger.Error(err, "failed to marshal commands")
		return nil
	}
	return out

}
