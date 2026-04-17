package actions

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
)

type ExampleActoin struct {
	action.Action
}

func (a ExampleActoin) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}

func (a ExampleActoin) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}

// generic actions
type ReplyToMessage struct {
	action.Action
	Message string
}

func (a ReplyToMessage) Run(passThrough any, v ...any) action.Flags { // v[0] is the message, v[1] is the CommandData or message is defined at creation, then v[0] is the CommandData
	flags := action.Flags{PassThrough: passThrough}
	if len(v) == 0 {
		flags.Error = fmt.Errorf("no data provided")
		return flags
	}
	if a.Message == "" {
		message, ok := passThrough.(string)
		if !ok {
			flags.Error = fmt.Errorf("no message provided")
			return flags
		}
		a.Message = message
	}
	data, ok := v[0].(datatypes.ChatMessageData)
	if !ok {
		flags.Error = fmt.Errorf("expected ChatMessageData or nil, got %T", v[0])
		return flags
	}
	_, err := twitchapi.AsBot().SendReply(a.Message, data.MessageID)
	if err != nil {
		flags.Error = err
	}
	return flags
}

func (a ReplyToMessage) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}

type SendMessage struct {
	action.Action
	Message string
}

func (a SendMessage) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	if a.Message == "" {
		message, ok := passThrough.(string)
		if !ok {
			flags.Error = fmt.Errorf("no message provided")
			return flags
		}
		a.Message = message
	}
	_, err := twitchapi.AsBot().SendMessage(a.Message)
	if err != nil {
		flags.Error = err
	}
	return flags
}

func (a SendMessage) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}

type SendAnnouncement struct {
	action.Action
	Message string
	Color   string
}

func (a SendAnnouncement) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	if a.Message == "" {
		message, ok := passThrough.(string)
		if ok {
			a.Message = message
		} else {
			announcement, ok := passThrough.([]string)
			if !ok || len(announcement) != 2 {
				flags.Error = fmt.Errorf("no message or color provided")
				return flags
			}
			a.Message = announcement[0]
			a.Color = announcement[1]
		}
	}
	if a.Color == "" {
		a.Color = "primary"
	}
	logger.Debug("Sending Announcement", "message", a.Message, "color", a.Color)
	body, err := twitchapi.AsBot().SendAnnouncement(a.Message, a.Color)
	if err != nil {
		flags.Error = err
		return flags
	}

	logger.Debug("response", "body", string(body))

	return flags
}

func (a SendAnnouncement) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}

type Delay struct {
	action.Action
	Duration time.Duration
}

func (a Delay) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	if a.Duration == 0 {
		duration, ok := passThrough.(int)
		if !ok {
			flags.Error = fmt.Errorf("no duration provided")
			return flags
		}
		a.Duration = time.Duration(duration * int(time.Second))
	}
	time.Sleep(a.Duration)
	return flags
}

func (a Delay) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}

const (
	// Log Sceeme format:
	// reserved (bit 5-15 ) | Log level (bit 2-4) | Loggers (bit 0-1)
	// (0 0 0 0 0 0 0 0 0 0 0) (0 0 0) (0 0)

	LogInfo    = 0b00
	LogWarning = 0b01
	LogError   = 0b10
	LogDebug   = 0b11

	// logger logic inverted, when 0 print to that logger, when 1 do not print to that logger. this allows for a default of printing to all loggers when sceeme is 0, and allows for easy combination of loggers by setting bits to 1 for loggers you do not want to print to.
	LogAllLogger    = 0b00000
	LogLoggerStdout = 0b11000 // skips other loggers, not intended to be staked.
	LogLoggerFile   = 0b10100 // skips other loggers, not intended to be staked.
	LogLoggerWeb    = 0b01100 // skips other loggers, not intended to be staked. not implemented yet, would log to web dashboard
	LogSkipStdout   = 0b00100
	LogSkipFile     = 0b01000
	LogSkipWeb      = 0b10000 // not implemented yet, would skip logging to web dashboard

)

// This is a generic log action that can be used to log messages to the console, file, or web dashboard.
// The message can be defined at creation or passed in as data.
// The sceeme can be used to define the format of the log message and which loggers to use.
type Log struct {
	action.Action
	Message   string                          // optional. if not provided, will use v[0] as the message if it is a string
	Sceeme    int                             // optional. if not provided, will output in text format to all loggers. see Log Sceeme constants for more details.
	Data      []any                           // optional. additional data to pass to log. sceemes can use this data to output in different formats or include additional information in the log message. if empty, will use v[1:] as data. if nil will not pass any additional data to log.
	Formatter func(data []any) ([]any, error) // optional. if sceeme is LogSceemeCustom, this formatter will be used to format the log message. it takes in the data and returns the formatted data to pass to the logger. if not provided, will pass data as is.
}

func (a Log) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	if a.Message == "" {
		message, ok := passThrough.(string)
		if !ok {
			flags.Error = fmt.Errorf("no message provided")
			return flags
		}
		a.Message = message
	}

	if len(a.Data) == 0 {
		a.Data = v
	} else if a.Data == nil {
		a.Data = []any{}
	}

	outData := []any{}
	loggers := []*slog.Logger{}

	if a.Formatter != nil {
		var err error
		outData, err = a.Formatter(a.Data)
		if err != nil {
			flags.Error = err
			return flags
		}
	} else {
		outData = a.Data
	}

	if a.Sceeme&LogSkipStdout == 0 {
		loggers = append(loggers, logger.StanderOutLogger)
	}
	if a.Sceeme&LogSkipFile == 0 {
		loggers = append(loggers, logger.FileLogger)
	}
	// if a.Sceeme&LogSkipWeb == 0 {
	// 	loggers = append(loggers, logger.WebLogger)
	// }

	for _, l := range loggers {
		switch a.Sceeme & 0b11 {
		case LogInfo:
			l.Info(a.Message, outData...)
		case LogWarning:
			l.Warn(a.Message, outData...)
		case LogError:
			l.Error(a.Message, outData...)
		case LogDebug:
			logger.DebugToLogger(l, a.Message, outData...)
		default:
			flags.Error = fmt.Errorf("invalid log level")
			return flags
		}
	}

	return flags
}

func (a Log) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}

const defaultADLeanth = 60

type RunAD struct {
	action.Action
	Time int
}

func (a RunAD) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}

	if a.Time >= 0 {
		if time, ok := passThrough.(int); ok {
			a.Time = time
		} else {
			time = defaultADLeanth
		}
	}

	if a.Time > 180 {
		a.Time = 180
	} else if a.Time < 30 {
		a.Time = 30
	}

	if _, err := twitchapi.AsUser().Post("/channels/commercial", fmt.Sprintf(`
		{
			"broadcaster_id": %s,
  			"length": %d
		}
		`,
		context.Get().GetBroadcaster().Id,
		a.Time,
	)); err != nil {
		flags.Error = err
	}

	return flags
}

func (a RunAD) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}

type SetYappyChat struct {
	action.Action
	Mode   bool
	Toggle bool
}

func (a SetYappyChat) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	if a.Toggle {
		context.Get().YappyChat = !context.Get().YappyChat
	} else {
		context.Get().YappyChat = a.Mode
	}
	flags.PassThrough = context.Get().YappyChat
	return flags
}

func (a SetYappyChat) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}

type IfBool struct {
	action.Action
	TrueActoins  []action.Action
	FalseActoins []action.Action
	ActoinData   []any
}

func (a IfBool) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	b, ok := passThrough.(bool)
	if !ok {
		flags.Error = fmt.Errorf("Expexctd bool")
		return flags
	}

	if b {
		flags.AddActions = action.Flag{
			Active:     true,
			Actions:    a.TrueActoins,
			ActionData: a.ActoinData,
		}
	} else {
		flags.AddActions = action.Flag{
			Active:     true,
			Actions:    a.FalseActoins,
			ActionData: a.ActoinData,
		}
	}

	if len(a.ActoinData) == 0 {
		flags.AddActions.ActionData = v
	}
	return flags
}

func (a IfBool) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}
