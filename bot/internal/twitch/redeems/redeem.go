package redeems

import (
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/queue"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
)

type Redeem struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	RewardID    string          `json:"reward_id"`
	Actions     []action.Action `json:"-"`
	QueueName   string          `json:"queue"`
	Blocking    bool            `json:"blocking"`
	Active      bool            `json:"active"`
}

func HandleRedemption(data datatypes.RedemptionData) {
	for _, redeem := range GetRedemptions() {
		if redeem.RewardID == data.Reward.ID {
			redeem.AddActions(data)
		}
	}
}

func (r *Redeem) AddActions(data any) {
	if !r.Active {
		return
	}
	q := queue.GetQueue(r.QueueName)
	for _, a := range r.Actions {
		q.AddActions(a, data)
	}

	if r.Blocking {
		q.Lock()
	}

	q.Start()
}
