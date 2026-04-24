package redeems

import (
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/action"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/actionQueueSystem/queue"
	datatypes "github.com/TheLazyTurtle33/sea-core/bot/internal/dataTypes"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
)

type Redeem struct {
	// internal fields
	Actions   []action.Action `json:"-"`     // required
	QueueName string          `json:"queue"` // required

	ID string `json:"id"` // only provide when overriding

	// fields from Twitch API
	Title                             string `json:"title"`                                 // required
	Cost                              int    `json:"cost"`                                  // required
	BackgroundColor                   string `json:"background_color"`                      // required
	IsEnabled                         bool   `json:"is_enabled"`                            // optional. default false
	Prompt                            string `json:"prompt"`                                // optional. only if IsUserInputRequired is true
	IsUserInputRequired               bool   `json:"is_user_input_required"`                // optional. default false
	IsMaxPerStreamEnabled             bool   `json:"is_max_per_stream_enabled"`             // optional. default false
	MaxPerStream                      int    `json:"max_per_stream"`                        // optional. only if IsMaxPerStreamEnabled is true
	IsMaxPerUserPerStreamEnabled      bool   `json:"is_max_per_user_per_stream_enabled"`    // optional. default false
	MaxPerUserPerStream               int    `json:"max_per_user_per_stream"`               // optional. only if IsMaxPerUserPerStreamEnabled is true
	IsGlobalCooldownEnabled           bool   `json:"is_global_cooldown_enabled"`            // optional. default false
	GlobalCooldownSeconds             int    `json:"global_cooldown_seconds"`               // optional. only if IsGlobalCooldownEnabled is true
	IsPaused                          bool   `json:"is_paused"`                             // optional. default false
	ShouldRedemptionsSkipRequestQueue bool   `json:"should_redemptions_skip_request_queue"` // optional. default false

}

func HandleRedemption(data datatypes.RedemptionData) {
	logger.Log("Redeem redeemed.", "ID", data.Reward.ID)
	if redeem, ok := iDRedeemsMap[data.Reward.ID]; ok {
		redeem.AddActions(data)
	}
}

func (r *Redeem) AddActions(data any) {
	if !r.IsEnabled {
		return
	}
	q := queue.GetQueue(r.QueueName)
	for _, a := range r.Actions {
		q.AddAction(a, []any{data})
	}

	q.Start()
}
