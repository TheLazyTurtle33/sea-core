package actions

import "github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"

type SendTTS struct {
	action.Action
}

func (a SendTTS) Run(v ...any) action.Flags {
	flags := action.Flags{}

	return flags
}

func (a SendTTS) OnAdd(v ...any) action.Flags {
	flags := action.Flags{}
	return flags
}
