package actions

import (
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
)

type ExampleActoin struct {
	action.Action
}

func (a *ExampleActoin) Run(v ...any) action.Flags {
	flags := action.Flags{}
	return flags
}

func (a *ExampleActoin) OnAdd(v ...any) action.Flags {
	flags := action.Flags{}
	return flags
}

// generic actions
type ReplyToMessage struct {
	action.Action
	Message string
}

func (a *ReplyToMessage) Run(v ...any) action.Flags { // v[0] is the message, v[1] is the CommandData or message is defined at creation, then v[0] is the CommandData
	flags := action.Flags{}
	if len(v) == 0 {
		flags.Error = fmt.Errorf("no data provided")
		return flags
	}
	if a.Message == "" {
		if reflect.TypeOf(v[0]).Kind() != reflect.String {
			flags.Error = fmt.Errorf("no message provided")
			return flags
		}
		a.Message = v[0].(string)
	}
	v = v[1:] // remove message from v if there, else remove nil
	data, ok := v[0].(datatypes.ChatMessageData)
	if !ok {
		flags.Error = fmt.Errorf("expected ChatMessageData or nil, got %T", v[0])
		return flags
	}
	_, err := twitchapi.As("bot").SendReply(a.Message, data.MessageID)
	if err != nil {
		flags.Error = err
	}
	return flags
}

func (a *ReplyToMessage) OnAdd(v ...any) action.Flags {
	flags := action.Flags{}
	return flags
}

type SendMessage struct {
	action.Action
	Message string
}

func (a *SendMessage) Run(v ...any) action.Flags {
	flags := action.Flags{}
	if len(v) == 0 {
		flags.Error = fmt.Errorf("no data provided")
		return flags
	}
	if a.Message == "" {
		if reflect.TypeOf(v[0]).Kind() != reflect.String {
			flags.Error = fmt.Errorf("no message provided")
			return flags
		}
		a.Message = v[0].(string)
	}
	_, err := twitchapi.As("bot").SendMessage(a.Message)
	if err != nil {
		flags.Error = err
	}
	return flags
}

func (a *SendMessage) OnAdd(v ...any) action.Flags {
	flags := action.Flags{}
	return flags
}

type Delay struct {
	action.Action
	Duration time.Duration
}

func (a *Delay) Run(v ...any) action.Flags {
	flags := action.Flags{}
	if len(v) == 0 {
		flags.Error = fmt.Errorf("no data provided")
		return flags
	}
	if a.Duration == 0 {
		if reflect.TypeOf(v[0]).Kind() != reflect.String {
			flags.Error = fmt.Errorf("no duration provided")
			return flags
		}
		a.Duration = time.Duration(v[0].(int))
	}
	time.Sleep(a.Duration)
	return flags
}

func (a *Delay) OnAdd(v ...any) action.Flags {
	flags := action.Flags{}
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

func (a *Log) Run(v ...any) action.Flags {
	flags := action.Flags{}
	if len(v) == 0 {
		flags.Error = fmt.Errorf("no data provided")
		return flags
	}

	if a.Message == "" {
		if reflect.TypeOf(v[0]).Kind() != reflect.String {
			flags.Error = fmt.Errorf("no message provided")
			return flags
		}
		a.Message = v[0].(string)
	}

	if len(a.Data) == 0 {
		a.Data = v[1:]
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

func (a *Log) OnAdd(v ...any) action.Flags {
	flags := action.Flags{}
	return flags
}
