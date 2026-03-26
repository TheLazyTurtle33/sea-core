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
	actionDataNext any
	running        bool
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
			for i, a := range q.actions {
				q.runActions(a, q.actionData[i])
			}
			time.Sleep(q.repeatDelay)
		}

	}
}

func (q *Queue) runActions(act action.Action, data any) {

	flags := act.Run(q.actionDataNext, data)
	if flags.Error != nil {
		logger.Error(flags.Error, "error running action")
		return
	}
	if flags.Lock.Active {
		if flags.Lock.QueueName != "" {
			GetQueue(flags.Lock.QueueName).Lock()
		} else {
			q.Lock()
		}
	}
	if flags.Unlock.Active {
		if flags.Unlock.QueueName != "" {
			GetQueue(flags.Unlock.QueueName).Unlock()
		} else {
			q.Unlock()
		}
	}
	if flags.StartQueue.Active {
		GetQueue(flags.StartQueue.QueueName).Start()
	}
	if flags.StopQueue.Active {
		if flags.StopQueue.QueueName != "" {
			GetQueue(flags.StopQueue.QueueName).Stop()
		} else {
			q.Stop()
		}
	}
	if flags.AddActions.Active {
		if flags.AddActions.QueueName != "" {
			for i, a := range flags.AddActions.Actions {
				q.AddActions(a, flags.AddActions.ActionData[i])
			}
		} else {
			for i, a := range flags.AddActions.Actions {
				GetQueue(flags.AddActions.QueueName).AddActions(a, flags.AddActions.ActionData[i])
			}
		}
	}
	q.actionDataNext = flags.DataForNextAction
}
