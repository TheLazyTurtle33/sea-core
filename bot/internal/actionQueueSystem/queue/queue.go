package queue

import (
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
)

type Queue struct {
	name           string
	locked         bool
	repeating      bool
	repeatDelay    time.Duration // how long to sleep before looping
	persistent     bool          // keep worker alive when empty, or spin down
	actions        []action.Action
	actionData     []any
	actionIndex    int
	actionDataNext any
	running        bool
	puased         time.Duration
}

func (q *Queue) AddActions(a action.Action, data any) {
	if q.locked {
		return
	}
	q.actions = append(q.actions, a)
	q.actionData = append(q.actionData, data)
}

func (q *Queue) Lock() {
	q.locked = true
}

func (q *Queue) Unlock() {
	q.locked = false
}

func (q *Queue) Start() {
	if q.running {
		return
	}
	q.running = true
	logger.Log("Queue started", "queue", q.name)
	go q.worker()
}

func (q *Queue) Stop() {
	q.running = false
}

func (q *Queue) worker() {
	for {
		if !q.running {
			return
		}

		if q.puased > 0 {
			time.Sleep(q.puased)
			q.puased = 0
		}

		if len(q.actions) == 0 {
			if q.persistent {
				time.Sleep(1 * time.Second)
				continue
			}
			return
		}

		if !q.repeating {
			q.runActions(q.actions[0], q.actionData[0])
			q.actions = q.actions[1:]
			q.actionData = q.actionData[1:]
		} else {
			var data any
			if len(q.actionData) > 0 {
				data = q.actionData[q.actionIndex]
			}
			q.runActions(q.actions[q.actionIndex], data)
			q.actionIndex++
		}

		if len(q.actions) == q.actionIndex {
			time.Sleep(q.repeatDelay)
			q.actionIndex = 0
		}
	}
}

func (q *Queue) runActions(act action.Action, data any) {

	flags := act.Run(q.actionDataNext, data)
	if flags.Error != nil {
		logger.Error("error running action", flags.Error)
		return
	}
	if flags.Lock.Active {
		if flags.Lock.StringData != "" {
			GetQueue(flags.Lock.StringData).Lock()
		} else {
			q.Lock()
		}
	}
	if flags.Unlock.Active {
		if flags.Unlock.StringData != "" {
			GetQueue(flags.Unlock.StringData).Unlock()
		} else {
			q.Unlock()
		}
	}
	if flags.StartQueue.Active {
		GetQueue(flags.StartQueue.StringData).Start()
	}
	if flags.StopQueue.Active {
		if flags.StopQueue.StringData == "" {
			q.Stop()
		} else {
			GetQueue(flags.StopQueue.StringData).Stop()
		}
	}
	if flags.AddActions.Active {
		if flags.AddActions.StringData == "" {
			for i, a := range flags.AddActions.Actions {
				q.AddActions(a, flags.AddActions.ActionData[i])
			}
		} else {
			for i, a := range flags.AddActions.Actions {
				GetQueue(flags.AddActions.StringData).AddActions(a, flags.AddActions.ActionData[i])
			}
		}
	}
	if flags.Pause.Active {
		if flags.Pause.StringData == "" {
			q.puased += time.Duration(flags.Pause.IntData) * time.Second
		} else {
			GetQueue(flags.Pause.StringData).puased += time.Duration(flags.Pause.IntData) * time.Second
		}
	}
	if flags.Skip.Active {
		if flags.Skip.StringData == "" {
			q.Skip(flags.Skip.IntData)
		} else {
			GetQueue(flags.Skip.StringData).Skip(flags.Skip.IntData)
		}
	}
	q.actionDataNext = flags.DataForNextAction
}

func (q *Queue) Skip(skip int) {
	if skip == 0 {
		return
	}
	for range skip {
		if q.repeating {
			q.actionIndex++
		} else {
			q.actions = q.actions[1:]
		}
	}
}
