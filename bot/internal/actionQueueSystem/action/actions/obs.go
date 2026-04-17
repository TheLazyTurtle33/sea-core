package actions

import (
	"fmt"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/shared/obs"
)

type SetScene struct {
	action.Action
	Scene string
}

func (a SetScene) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}

	if a.Scene == "" {
		val, ok := passThrough.(string)
		if !ok {
			flags.Error = fmt.Errorf("expected pass through value to be sting")
			return flags
		}
		a.Scene = val
	}

	client, err := obs.Get()
	if err != nil {
		flags.Error = err
		return flags
	}

	responce, err := client.SetScene(a.Scene)

	flags.PassThrough = responce
	flags.Error = err

	return flags
}
func (a SetScene) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}
