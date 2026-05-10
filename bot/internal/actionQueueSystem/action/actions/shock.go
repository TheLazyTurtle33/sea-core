package actions

import (
	"fmt"
	"strings"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
)

const shockPortName = "/dev/ttyACM0"

var Shock = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		if len(actionData) == 0 || actionData[0].Type != action.ChatMessageType {
			flags.Error = fmt.Errorf("Shock: expected chat message")
			return flags
		}

		chat, ok := actionData[0].Data.(datatypes.ChatMessageData)
		if !ok {
			flags.Error = fmt.Errorf("Shock: expected chat message, got %T", actionData[0].Data)
			return flags
		}

		for _, badge := range chat.Badges {
			if (badge.SetID == "moderator" || badge.SetID == "broadcaster") && strings.TrimSpace(strings.ToLower(chat.Message.Text)) == "!shock" {
				return setFlagsForToggle(flags)
			}
		}

		if !context.Get().CanShock {
			return setFlagsForReply(flags, `Sorryyy the "shock" collar isnt active rn >.<`, chat)
		}

		hasBits := false
		firstBits := 0
		secondBits := 0
		for _, fragment := range chat.Message.Fragments {
			if fragment.Text == "cheer" {
				hasBits = true
				if firstBits == 0 {
					firstBits = fragment.Cheer.Bits
				} else {
					secondBits += fragment.Cheer.Bits
				}
			}
		}

		if !hasBits || strings.TrimSpace(strings.ToLower(chat.Message.Text)) == "!shock -h" {
			return setFlagsForReply(flags, `Sorry it costs bits to "shock" me x.x use !shock {20-70bits(increments of 10)} or split it between the level and time vals, !shock {10|20bits} {10-50bits (increments of 10)}`, chat)
		}

		if secondBits == 0 {
			if firstBits < 20 {
				return setFlagsForReply(flags, `Sorry it costs a minimum of 20 bits to "shock" :c use !shock -h for more info`, chat)
			}
			steps := []int{70, 60, 50, 40, 30, 20}
			for _, step := range steps {
				if firstBits >= step {
					firstBits = step / 10
					break
				}
			}
			if firstBits == 2 {
				flags = setFlagsForLevel(flags, 1, time.Second)
				flags.AddActions.Actions = append(flags.AddActions.Actions, ReplyToMessage.Make([]action.ActionData{
					{Type: action.StringType, Data: `"shocking" at level 1 for 1 sec >:3`},
				}))
				return flags
			}
			firstBits -= 2
			flags = setFlagsForLevel(flags, 2, time.Duration(firstBits)*time.Second)
			flags.AddActions.Actions = append(flags.AddActions.Actions, ReplyToMessage.Make([]action.ActionData{
				{Type: action.StringType, Data: fmt.Sprintf(`"shocking" at level 2 for %d secs >:3`, firstBits)},
			}))
			return flags
		}

		if firstBits < 10 || secondBits < 10 {
			return setFlagsForReply(flags, `Sorry it costs a minimum of 200 bits to "shock" :c use !shock -h for more info`, chat)
		}

		steps := []int{20, 10}
		for _, step := range steps {
			if firstBits >= step {
				firstBits = step / 10
				break
			}
		}
		steps = []int{50, 40, 30, 20, 10}
		for _, step := range steps {
			if secondBits >= step {
				secondBits = step / 10
				break
			}
		}

		flags = setFlagsForLevel(flags, firstBits, time.Duration(secondBits)*time.Second)
		flags.AddActions.Actions = append(flags.AddActions.Actions, ReplyToMessage.Make([]action.ActionData{
			{Type: action.StringType, Data: fmt.Sprintf(`"shocking" at level %d for %d secs >:3`, firstBits, secondBits)},
		}))
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "Shock",
		Description: "Trigger the shock hardware sequence.",
		RunData: []action.ActionData{
			{
				Type:        action.ChatMessageType,
				Description: "Chat message that triggered the shock.",
				Required:    true,
			},
		},
	},
}

func setFlagsForReply(flags action.Flags, message string, chat datatypes.ChatMessageData) action.Flags {
	flags.AddActions = action.Flag{
		Active: true,
		Actions: []action.Action{
			ReplyToMessage.Make([]action.ActionData{
				{Type: action.ChatMessageType, Data: chat},
				{Type: action.StringType, Data: message},
			}),
		},
	}
	return flags
}

func setFlagsForToggle(flags action.Flags) action.Flags {
	context.Get().CanShock = !context.Get().CanShock
	message := "the shock collar is now off :c"
	if context.Get().CanShock {
		message = "THE (not so shock) SHOCK COLLAR IS NOW ON!!"
	}
	flags.AddActions = action.Flag{
		Active: true,
		Actions: []action.Action{
			SendMessage.Make([]action.ActionData{{Type: action.StringType, Data: message}}),
		},
	}
	return flags
}

func setFlagsForLevel(flags action.Flags, level int, duration time.Duration) action.Flags {
	flags.AddActions = action.Flag{
		Active: true,
		Actions: []action.Action{
			SendSerialData.Make([]action.ActionData{
				{Type: "int", Data: 115200},
				{Type: action.StringType, Data: shockPortName},
				{Type: action.StringType, Data: fmt.Sprintf("ON %d", level)},
			}),
			Delay.Make([]action.ActionData{{Type: "int", Data: int(duration.Seconds())}}),
			SendSerialData.Make([]action.ActionData{
				{Type: "int", Data: 115200},
				{Type: action.StringType, Data: shockPortName},
				{Type: action.StringType, Data: "OFF"},
			}),
		},
	}
	return flags
}

func init() {
	action.ActionMap[Shock.MetaData.Name] = Shock
}
