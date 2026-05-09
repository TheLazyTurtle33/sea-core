package actions

// import (
// 	"fmt"
// 	"strings"
// 	"time"

// 	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
// 	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
// 	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
// )

// type Shock struct {
// 	action.Action
// 	level string
// }

// const portName = "/dev/ttyACM0"

// func (a Shock) Run(passThrough any, v ...any) action.Flags {
// 	flags := action.Flags{PassThrough: passThrough}

// 	chat, ok := v[0].(datatypes.ChatMessageData)
// 	if !ok {
// 		flags.Error = fmt.Errorf("Shock: expectid v[0] to be ChatMessageData")
// 		return flags
// 	}

// 	for _, badge := range chat.Badges {
// 		if (badge.SetID == "moderator" || badge.SetID == "broadcaster") && strings.TrimSpace(strings.ToLower(chat.Message.Text)) == "!shock" {
// 			return setFlagsForToggle(flags)
// 		}
// 	}
// 	if !context.Get().CanShock {
// 		return setFlagsForReplay(flags, `Sorryyy the "shock" collar isnt active rn >.<`, chat)
// 	}

// 	hasBits := false

// 	firstBits := 0
// 	secondBits := 0

// 	for _, fragment := range chat.Message.Fragments {
// 		if fragment.Text == "cheer" {
// 			hasBits = true
// 			if firstBits == 0 {
// 				firstBits = fragment.Cheer.Bits
// 			} else {
// 				secondBits += fragment.Cheer.Bits
// 			}
// 		}
// 	}

// 	if !hasBits || chat.Message.Text == "!shock -h" {
// 		return setFlagsForReplay(flags, `Sorry it costs bits to "shock" me x.x use !shock {20-70bits(increments of 10)} or split it between the level and time vals, !shock {10|20bits} {10-50bits (increments of 10)}`, chat)
// 	}

// 	if secondBits == 0 {
// 		if firstBits < 20 {
// 			return setFlagsForReplay(flags, `Sorry it costs a minimum of 20 bits to "shock" :c use !shock -h for more info`, chat)
// 		}

// 		steps := []int{70, 60, 50, 40, 30, 20}
// 		for _, step := range steps {
// 			if firstBits >= step {
// 				firstBits = step / 10
// 				break
// 			}
// 		}
// 		if firstBits == 2 {
// 			flags = setFlagsForLevel(flags, 1, time.Second)
// 			flags.AddActions.Actions = append(flags.AddActions.Actions, &ReplyToMessage{Message: `"shocking" at level 1 for 1 sec >:3`})
// 			flags.AddActions.ActionData = append(flags.AddActions.ActionData, []any{chat})
// 			return flags
// 		}
// 		firstBits -= 2
// 		flags = setFlagsForLevel(flags, 2, time.Duration(firstBits)*time.Second)
// 		flags.AddActions.Actions = append(flags.AddActions.Actions, &ReplyToMessage{Message: fmt.Sprintf(`"shocking" at level 2 for %d secs >:3`, firstBits)})
// 		flags.AddActions.ActionData = append(flags.AddActions.ActionData, []any{chat})
// 		return flags

// 	}

// 	if firstBits < 10 || secondBits < 10 {
// 		return setFlagsForReplay(flags, `Sorry it costs a minimum of 200 bits to "shock" :c use !shock -h for more info`, chat)
// 	}

// 	steps := []int{20, 10}
// 	for _, step := range steps {
// 		if firstBits >= step {
// 			firstBits = step / 10
// 			break
// 		}
// 	}

// 	steps = []int{50, 40, 30, 20, 10}
// 	for _, step := range steps {
// 		if secondBits >= step {
// 			secondBits = step / 10
// 			break
// 		}
// 	}

// 	flags = setFlagsForLevel(flags, firstBits, time.Duration(secondBits)*time.Second)
// 	flags.AddActions.Actions = append(flags.AddActions.Actions, &ReplyToMessage{Message: fmt.Sprintf(`"shocking" at level %d for %d secs >:3`, firstBits, secondBits)})
// 	flags.AddActions.ActionData = append(flags.AddActions.ActionData, []any{chat})

// 	return flags
// }

// func setFlagsForReplay(flags action.Flags, message string, chat datatypes.ChatMessageData) action.Flags {
// 	flags.AddActions = action.Flag{
// 		Active: true,
// 		Actions: []action.Action{
// 			&ReplyToMessage{Message: message},
// 		},
// 		ActionData: [][]any{{chat}},
// 	}
// 	return flags
// }

// func setFlagsForToggle(flags action.Flags) action.Flags {

// 	context.Get().CanShock = !context.Get().CanShock

// 	var message string
// 	if context.Get().CanShock {
// 		message = "THE (not so shock) SHOCK COLLAR IS NOW ON!!"
// 	} else {
// 		message = "the shock collar is now off :c"
// 	}

// 	flags.AddActions = action.Flag{
// 		Active: true,
// 		Actions: []action.Action{
// 			&SendMessage{Message: message},
// 		},
// 	}
// 	return flags
// }

// func setFlagsForLevel(flags action.Flags, level int, duration time.Duration) action.Flags {
// 	flags.AddActions = action.Flag{
// 		Active: true,
// 		Actions: []action.Action{
// 			&SendSerialData{},
// 			&Delay{Duration: duration},
// 			&SendSerialData{},
// 		},
// 		ActionData: [][]any{
// 			{115200, portName, fmt.Sprintf("ON %d", level)},
// 			nil,
// 			{115200, portName, "OFF"},
// 		},
// 	}
// 	return flags
// }

// func (a Shock) OnAdd(passThrough any, v ...any) action.Flags {
// 	flags := action.Flags{PassThrough: passThrough}
// 	return flags
// }
