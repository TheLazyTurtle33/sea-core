package actions

import (
	"fmt"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
	"github.com/TheLazyTurtle33/sea-core/shared/obs"
)

type JumpScareAction struct {
	action.Action
}

func (a JumpScareAction) Run(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}

	client, err := obs.Get()
	if err != nil {
		flags.Error = fmt.Errorf("JumpScare: faild to get obs client")
		flags.Skip.Active = true
		return flags
	}

	if _, err := client.SetSourceVisability("JumpScare", true); err != nil {
		flags.Error = fmt.Errorf("JumpScare: fiald to set visability")
		flags.Skip.Active = true
		return flags
	}

	go func() {
		time.Sleep(4 * time.Second)
		if _, err := client.SetSourceVisability("JumpScare", false); err != nil {
			logger.Error("JumpScare: fiald to set visability", nil)
		}
	}()

	return flags
}

func (a JumpScareAction) OnAdd(passThrough any, v ...any) action.Flags {
	flags := action.Flags{PassThrough: passThrough}
	return flags
}
