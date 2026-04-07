package actions

import (
	"fmt"
	"reflect"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
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

type Log struct {
	action.Action
	Message string
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
	fmt.Println(a.Message)
	return flags
}

func (a *Log) OnAdd(v ...any) action.Flags {
	flags := action.Flags{}
	return flags
}
