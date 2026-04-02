package redeems

import (
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action/actions"
)

var HydrateRedeems = Redeem{
	Name:        "Hydrate",
	RewardID:    "0f58303d-e5d8-4c24-87c5-b2558017a98f",
	Description: "Remind the streamer to hydrate.",
	Actions:     []action.Action{&actions.SendMessage{Message: "HYDRATED! >:3"}},
	QueueName:   "default",
	Blocking:    false,
	Active:      true,
}

var redeems = []Redeem{
	HydrateRedeems,
}

func GetRedemptions() []Redeem {
	return redeems
}
