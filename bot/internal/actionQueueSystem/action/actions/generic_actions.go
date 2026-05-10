package actions

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
	"go.bug.st/serial"
)

const (
	defaultADLength = 60
)

var ExampleAction = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "ExampleAction",
		Description: "A template action that does nothing.",
	},
}

var ReplyToMessage = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		if len(actionData) < 1 {
			flags.Error = fmt.Errorf("no data provided")
			return flags
		}

		if actionData[0].Type != action.ChatMessageType {
			flags.Error = fmt.Errorf("expected ChatMessageData, got %T", actionData[0].Data)
			return flags
		}
		data, ok := actionData[0].Data.(datatypes.ChatMessageData)
		if !ok {
			flags.Error = fmt.Errorf("expected ChatMessageData, got %T", actionData[0].Data)
			return flags
		}

		message, err := parseStringParam(passThrough, actionData, 1)
		if err != nil {
			flags.Error = err
			return flags
		}

		_, err = twitchapi.AsBot().SendReply(message, data.MessageID)
		if err != nil {
			flags.Error = err
		}
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "ReplyToMessage",
		Description: "Reply to a chat message.",
		RunData: []action.ActionData{
			{
				Type:        action.ChatMessageType,
				Description: "The chat message to reply to.",
				Required:    true,
			},
			{
				Type:        action.StringType,
				Description: "The reply text. If omitted, passThrough must contain a string.",
				Required:    false,
			},
		},
	},
}

var SendMessage = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		message, err := parseStringParam(passThrough, actionData, 0)
		if err != nil {
			flags.Error = err
			return flags
		}
		_, err = twitchapi.AsBot().SendMessage(message)
		if err != nil {
			flags.Error = err
		}
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "SendMessage",
		Description: "Send a standard chat message.",
		RunData: []action.ActionData{
			{
				Type:        action.StringType,
				Description: "The message text to send.",
				Required:    true,
			},
		},
	},
}

var SendAnnouncement = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		message, err := parseStringParam(passThrough, actionData, 0)
		if err != nil {
			flags.Error = err
			return flags
		}
		color, _ := parseStringParam(passThrough, actionData, 1)
		if color == "" {
			color = "primary"
		}
		_, err = twitchapi.AsBot().SendAnnouncement(message, color)
		if err != nil {
			flags.Error = err
		}
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "SendAnnouncement",
		Description: "Send a channel announcement.",
		RunData: []action.ActionData{
			{
				Type:        action.StringType,
				Description: "Announcement text.",
				Required:    true,
			},
			{
				Type:        action.StringType,
				Description: "Announcement color.",
				Required:    false,
			},
		},
	},
}

var Delay = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		durationSeconds, err := parseIntParam(passThrough, actionData, 0)
		if err != nil {
			flags.Error = err
			return flags
		}
		time.Sleep(time.Duration(durationSeconds) * time.Second)
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "Delay",
		Description: "Pause execution for a number of seconds.",
		RunData: []action.ActionData{
			{
				Type:        "int",
				Description: "Number of seconds to pause.",
				Required:    true,
			},
		},
	},
}

var Log = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		message, err := parseStringParam(passThrough, actionData, 0)
		if err != nil {
			flags.Error = err
			return flags
		}
		logger.Log(message)
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "Log",
		Description: "Log a message to the bot logger.",
		RunData: []action.ActionData{
			{
				Type:        action.StringType,
				Description: "Log message.",
				Required:    true,
			},
		},
	},
}

var RunAD = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		duration, err := parseIntParam(passThrough, actionData, 0)
		if err != nil {
			duration = resolveADTime(passThrough, actionData)
		}
		steps := []int{180, 150, 120, 90, 60, 30}
		for _, step := range steps {
			if duration >= step {
				duration = step
				break
			}
		}
		if _, err := twitchapi.AsUser().Post("/channels/commercial", fmt.Sprintf(`{"broadcaster_id": %s, "length": %d}`, context.Get().GetBroadcaster().Id, duration)); err != nil {
			flags.Error = err
		}
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "RunAD",
		Description: "Start a commercial for a broadcaster.",
		RunData: []action.ActionData{
			{
				Type:        "int",
				Description: "Commercial length in seconds.",
				Required:    false,
			},
		},
	},
}

var SetYappyChat = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		mode, err := parseBoolParam(passThrough, actionData, 0)
		if err == nil {
			context.Get().TTSContext.YappyChat = mode
		} else if len(actionData) > 0 && actionData[0].Type == action.StringType {
			command, _ := actionData[0].Data.(string)
			if strings.EqualFold(command, "toggle") {
				context.Get().TTSContext.YappyChat = !context.Get().TTSContext.YappyChat
			} else {
				mode = strings.EqualFold(command, "true") || strings.EqualFold(command, "on")
				context.Get().TTSContext.YappyChat = mode
			}
		} else {
			flags.Error = fmt.Errorf("no bool or toggle value provided")
			return flags
		}

		if context.Get().TTSContext.YappyChat {
			context.Get().TTSContext.Delay = 0
		} else {
			context.Get().TTSContext.Delay = context.TTSDelayDefualt
		}
		flags.PassThrough = action.ActionData{Type: "bool", Data: context.Get().TTSContext.YappyChat}
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "SetYappyChat",
		Description: "Enable or disable yappy chat mode.",
		RunData: []action.ActionData{
			{
				Type:        action.StringType,
				Description: "Toggle or set the mode. use 'toggle', 'true', or 'false'.",
				Required:    false,
			},
		},
	},
}

var IfBool = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		if passThrough.Type != "bool" {
			flags.Error = fmt.Errorf("expected bool passThrough for IfBool action")
			return flags
		}
		condition, ok := passThrough.Data.(bool)
		if !ok {
			flags.Error = fmt.Errorf("expected bool passThrough for IfBool action")
			return flags
		}
		if len(actionData) < 2 {
			flags.Error = fmt.Errorf("IfBool requires two Actions values: true and false action lists")
			return flags
		}

		trueActions, ok := actionData[0].Data.([]action.Action)
		if !ok {
			flags.Error = fmt.Errorf("expected []action.Action for true actions, got %T", actionData[0].Data)
			return flags
		}
		falseActions, ok := actionData[1].Data.([]action.Action)
		if !ok {
			flags.Error = fmt.Errorf("expected []action.Action for false actions, got %T", actionData[1].Data)
			return flags
		}

		if condition {
			flags.AddActions = action.Flag{Active: true, Actions: trueActions}
		} else {
			flags.AddActions = action.Flag{Active: true, Actions: falseActions}
		}
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "IfBool",
		Description: "Run either a true or false action list based on a boolean passThrough.",
		RunData: []action.ActionData{
			{
				Type:        action.ActionsType,
				Description: "Actions to run when passThrough is true.",
				Required:    true,
			},
			{
				Type:        action.ActionsType,
				Description: "Actions to run when passThrough is false.",
				Required:    true,
			},
		},
	},
}

var SendSerialData = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		if len(actionData) < 3 {
			flags.Error = fmt.Errorf("SendSerialData requires baud rate, port, and message")
			return flags
		}

		baudRate, err := parseIntParam(passThrough, actionData, 0)
		if err != nil {
			flags.Error = err
			return flags
		}
		portName, err := parseStringParam(passThrough, actionData, 1)
		if err != nil {
			flags.Error = err
			return flags
		}
		message, err := parseStringParam(passThrough, actionData, 2)
		if err != nil {
			flags.Error = err
			return flags
		}

		mode := &serial.Mode{BaudRate: baudRate}
		port, err := serial.Open(portName, mode)
		if err != nil {
			flags.Error = fmt.Errorf("SendSerialData: failed to open port: %w", err)
			return flags
		}
		defer port.Close()

		n, err := port.Write([]byte(message))
		if err != nil {
			flags.Error = fmt.Errorf("SendSerialData: failed to write port: %w", err)
			return flags
		}
		logger.Log("SendSerialData: sent bytes", "port", portName, "bytes", n)
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "SendSerialData",
		Description: "Write raw data to a serial port.",
		RunData: []action.ActionData{
			{Type: "int", Description: "Baud rate.", Required: true},
			{Type: action.StringType, Description: "Serial port path.", Required: true},
			{Type: action.StringType, Description: "Message to send.", Required: true},
		},
	},
}

var Function = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		if len(actionData) == 0 {
			flags.Error = fmt.Errorf("Function action requires a function in actionData[0].Data")
			return flags
		}
		fn, ok := actionData[0].Data.(func(action.ActionData, []action.ActionData) action.Flags)
		if !ok {
			flags.Error = fmt.Errorf("Function action data was not executable")
			return flags
		}
		return fn(passThrough, actionData[1:])
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "Function",
		Description: "Run a custom function defined in actionData[0].Data.",
	},
}

var CompleteRedeemAction = action.Action{
	Run: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		flags := action.Flags{PassThrough: passThrough}
		if len(actionData) == 0 {
			flags.Error = fmt.Errorf("CompleteRedeemAction requires redemption data")
			return flags
		}
		data, ok := actionData[0].Data.(datatypes.RedemptionData)
		if !ok {
			flags.Error = fmt.Errorf("expected RedemptionData, got %T", actionData[0].Data)
			return flags
		}
		_, err := twitchapi.AsUser().Patch(
			"/channel_points/custom_rewards/redemptions?broadcaster_id="+context.Get().GetBroadcaster().Id+"&reward_id="+data.Reward.ID+"&id="+data.ID,
			`{"status":"FULFILLED"}`,
		)
		if err != nil {
			flags.Error = err
		}
		return flags
	},
	OnAdd: func(passThrough action.ActionData, actionData []action.ActionData) action.Flags {
		return action.Flags{PassThrough: passThrough}
	},
	MetaData: action.ActionMetaData{
		Name:        "CompleteRedeemAction",
		Description: "Mark a channel points redemption as fulfilled.",
		RunData: []action.ActionData{
			{
				Type:        action.RedeemType,
				Description: "The redemption data to complete.",
				Required:    true,
			},
		},
	},
}

func parseStringParam(passThrough action.ActionData, actionData []action.ActionData, idx int) (string, error) {
	if len(actionData) > idx {
		if actionData[idx].Type == action.StringType {
			if value, ok := actionData[idx].Data.(string); ok {
				return value, nil
			}
			if actionData[idx].Data == nil {
				return "", fmt.Errorf("expected string data for param %d, got nil", idx)
			}
			return "", fmt.Errorf("expected string data for param %d, got %T", idx, actionData[idx].Data)
		}
	}
	if passThrough.Type == action.StringType {
		if value, ok := passThrough.Data.(string); ok {
			return value, nil
		}
		return "", fmt.Errorf("expected passThrough string data, got %T", passThrough.Data)
	}
	return "", fmt.Errorf("no string data provided")
}

func parseIntParam(passThrough action.ActionData, actionData []action.ActionData, idx int) (int, error) {
	if len(actionData) > idx {
		switch v := actionData[idx].Data.(type) {
		case int:
			return v, nil
		case int64:
			return int(v), nil
		case float64:
			return int(v), nil
		case string:
			i, err := strconv.Atoi(v)
			if err != nil {
				return 0, err
			}
			return i, nil
		case nil:
			return 0, fmt.Errorf("expected int data for param %d, got nil", idx)
		default:
			return 0, fmt.Errorf("expected int data for param %d, got %T", idx, v)
		}
	}
	switch v := passThrough.Data.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, fmt.Errorf("no int data provided")
}

func parseBoolParam(passThrough action.ActionData, actionData []action.ActionData, idx int) (bool, error) {
	if len(actionData) > idx {
		switch v := actionData[idx].Data.(type) {
		case bool:
			return v, nil
		case string:
			return strings.EqualFold(v, "true") || strings.EqualFold(v, "on"), nil
		case nil:
			return false, fmt.Errorf("expected bool data for param %d, got nil", idx)
		default:
			return false, fmt.Errorf("expected bool data for param %d, got %T", idx, v)
		}
	}
	if passThrough.Type == "bool" {
		if value, ok := passThrough.Data.(bool); ok {
			return value, nil
		}
		return false, fmt.Errorf("expected passThrough bool data, got %T", passThrough.Data)
	}
	return false, fmt.Errorf("no bool data provided")
}

func resolveADTime(passThrough action.ActionData, actionData []action.ActionData) int {
	if passThrough.Type == "int" {
		if v, ok := passThrough.Data.(int); ok {
			return v
		}
	}
	if len(actionData) > 0 && actionData[0].Type == action.ChatMessageType {
		if chat, ok := actionData[0].Data.(datatypes.ChatMessageData); ok {
			words := strings.Split(chat.Message.Text, " ")
			if len(words) > 1 {
				if t, err := strconv.Atoi(words[1]); err == nil {
					return t
				}
			}
		}
	}
	return defaultADLength
}

func init() {
	action.ActionMap[ExampleAction.MetaData.Name] = ExampleAction
	action.ActionMap[ReplyToMessage.MetaData.Name] = ReplyToMessage
	action.ActionMap[SendMessage.MetaData.Name] = SendMessage
	action.ActionMap[SendAnnouncement.MetaData.Name] = SendAnnouncement
	action.ActionMap[Delay.MetaData.Name] = Delay
	action.ActionMap[Log.MetaData.Name] = Log
	action.ActionMap[RunAD.MetaData.Name] = RunAD
	action.ActionMap[SetYappyChat.MetaData.Name] = SetYappyChat
	action.ActionMap[IfBool.MetaData.Name] = IfBool
	action.ActionMap[SendSerialData.MetaData.Name] = SendSerialData
	action.ActionMap[Function.MetaData.Name] = Function
	action.ActionMap[CompleteRedeemAction.MetaData.Name] = CompleteRedeemAction
}
