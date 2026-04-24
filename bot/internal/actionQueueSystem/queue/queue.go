package queue

import (
	"slices"
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
	actionData     [][]any
	actionIndex    int
	actionDataNext any
	running        bool
	puased         time.Duration
}

func (q *Queue) AddAction(a action.Action, data []any) {
	if q.locked {
		logger.Log("Queue Locked!", "queue", q.name)
		return
	}
	q.actions = append(q.actions, a)
	q.actionData = append(q.actionData, data)
	q.OnAddAction(a, data)
}

func (q *Queue) AddActionAtIndex(a action.Action, index int, data ...any) {
	if q.locked {
		logger.Log("Queue Locked!", "queue", q.name)
		return
	}
	logger.Debug("num of actoins", "num", len(q.actions), "queue", q.name, "action index", q.actionIndex, "inser index", index)
	q.actions = slices.Insert(q.actions, index, a)
	q.actionData = slices.Insert(q.actionData, index, data)
	logger.Debug("num of actoins post insert", "num", len(q.actions), "queue", q.name, "action index", q.actionIndex)
	q.OnAddAction(a, data)
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

		q.actionIndex++
		logger.Debug("running actoin", "queue", q.name, "index", q.actionIndex)
		if !q.repeating {
			q.RunAction(q.actions[0], q.actionData[0]...)
			q.actions = q.actions[1:]
			q.actionData = q.actionData[1:]
			q.actionIndex--
			logger.Debug("index decremented", "queue", q.name, "index", q.actionIndex)
		} else {
			data := []any{}
			if len(q.actionData) > 0 {
				data = q.actionData[q.actionIndex-1]
			}
			q.RunAction(q.actions[q.actionIndex-1], data...)
		}

		if len(q.actions) == q.actionIndex {
			time.Sleep(q.repeatDelay)
			q.actionIndex = 0
		}
	}
}

func (q *Queue) RunAction(act action.Action, data ...any) {
	q.RunActionFunc(act.Run, data...)
}

func (q *Queue) OnAddAction(act action.Action, data ...any) {
	q.RunActionFunc(act.OnAdd, data...)
}

func (q *Queue) RunActionFunc(fn func(passThrough any, v ...any) action.Flags, data ...any) {
	flags := fn(q.actionDataNext, data...)
	if flags.Error != nil {
		logger.Debug("actoin had flag Error")
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
		logger.Debug("actoin had flag Add Actoin")

		ActionData := flags.AddActions.ActionData
		if dif := len(flags.AddActions.Actions) - len(ActionData); dif > 0 {
			for range dif {
				ActionData = append(ActionData, nil)
			}
		}

		queue := q
		if flags.AddActions.StringData != "" {
			queue = GetQueue(flags.AddActions.StringData)
		}
		index := queue.actionIndex
		index += flags.AddActions.IntData // use releteve pos

		if flags.AddActions.BoolData { // use absolute pos
			index = flags.AddActions.IntData
		}

		for i, a := range flags.AddActions.Actions {
			queue.AddActionAtIndex(a, index+i, ActionData[i]...)
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
	q.actionDataNext = flags.PassThrough
}

func (q *Queue) Skip(skip int) {
	if skip == 0 {
		skip = 1
	}
	for range skip {
		if q.repeating {
			q.actionIndex++
		} else {
			q.actions = q.actions[1:]
		}
	}
}
