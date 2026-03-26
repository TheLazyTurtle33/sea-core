package command

import (
	"encoding/json"
	"log"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action/actions"
)

var TestCommand = Command{
	Name:        "test",
	Triggers:    []string{"!test"},
	Description: "a test command",
	Usage:       "!test",
	Actions:     []action.Action{&actions.ReplyToMessage{Message: "test passed! :3"}},
	QueueName:   "default",
	blocking:    false,
}

var DiscordLink = Command{
	Name:        "discord",
	Triggers:    []string{"!discord", "!dc"},
	Description: "get the discord invite link",
	Usage:       "!discord",
	Actions:     []action.Action{&actions.CreateDiscordInvite{}, &actions.ReplyToMessage{}},
	QueueName:   "default",
	blocking:    false,
}

var CommandTriggers = map[string]*Command{}
var Commands = []*Command{
	&TestCommand,
	&DiscordLink,
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
		log.Println(err)
		return nil
	}
	return out

}
