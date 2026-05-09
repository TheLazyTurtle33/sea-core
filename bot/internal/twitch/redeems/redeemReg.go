package redeems

import (
	"encoding/json"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
	"github.com/TheLazyTurtle33/sea-core/shared/logger"
)

type TwitchRedeem struct {
	BroadcasterName  string `json:"broadcaster_name"`
	BroadcasterLogin string `json:"broadcaster_login"`
	BroadcasterID    string `json:"broadcaster_id"`
	ID               string `json:"id"`
	Image            struct {
		URL1X string `json:"url_1x"`
		URL2X string `json:"url_2x"`
		URL4X string `json:"url_4x"`
	} `json:"image"`
	BackgroundColor     string `json:"background_color"`
	IsEnabled           bool   `json:"is_enabled"`
	Cost                int    `json:"cost"`
	Title               string `json:"title"`
	Prompt              string `json:"prompt"`
	IsUserInputRequired bool   `json:"is_user_input_required"`
	MaxPerStreamSetting struct {
		IsEnabled    bool `json:"is_enabled"`
		MaxPerStream int  `json:"max_per_stream"`
	} `json:"max_per_stream_setting"`
	MaxPerUserPerStreamSetting struct {
		IsEnabled           bool `json:"is_enabled"`
		MaxPerUserPerStream int  `json:"max_per_user_per_stream"`
	} `json:"max_per_user_per_stream_setting"`
	GlobalCooldownSetting struct {
		IsEnabled             bool `json:"is_enabled"`
		GlobalCooldownSeconds int  `json:"global_cooldown_seconds"`
	} `json:"global_cooldown_setting"`
	IsPaused     bool `json:"is_paused"`
	IsInStock    bool `json:"is_in_stock"`
	DefaultImage struct {
		URL1X string `json:"url_1x"`
		URL2X string `json:"url_2x"`
		URL4X string `json:"url_4x"`
	} `json:"default_image"`
	ShouldRedemptionsSkipRequestQueue bool   `json:"should_redemptions_skip_request_queue"`
	RedemptionsRedeemedCurrentStream  int    `json:"redemptions_redeemed_current_stream"`
	CooldownExpiresAt                 string `json:"cooldown_expires_at"`
}

var iDRedeemsMap = map[string]Redeem{}

func RegisterRedeems() {
	logger.Log("Registering redeems")

	RequestRedeems()

	// create redeems defined in the bot that don't exist on Twitch yet.
	if len(definedRedeems) > 0 {
		for _, redeem := range definedRedeems {
			id := CreateRedeem(redeem)
			if id != "" {
				redeem.ID = id
				iDRedeemsMap[id] = redeem
			}
		}
	}

	logger.Debug("registered redeems", "count", len(iDRedeemsMap), "redeemsExport", string(ExportRedeems()))
}

func RequestRedeems() {
	resp, err := twitchapi.AsUser().Get("/channel_points/custom_rewards?broadcaster_id=" + context.Get().GetBroadcaster().Id)
	if err != nil {
		logger.Error("failed to request redeems", err)
		return
	}
	var data struct {
		Data []TwitchRedeem `json:"data"`
	}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		logger.Error("failed to unmarshal redeems", err)
		return
	}

	for _, redeem := range data.Data {
		if fromDefinedRedeems, ok := definedRedeems[redeem.Title]; ok {
			fromDefinedRedeems.ID = redeem.ID
			iDRedeemsMap[redeem.ID] = fromDefinedRedeems
			// UpdateRedeem(redeem.ID, fromDefinedRedeems)
			delete(definedRedeems, redeem.Title)
		} else if fromDefinedOverrides, ok := definedOverrides[redeem.ID]; ok {
			iDRedeemsMap[redeem.ID] = MakeRedeem(redeem, fromDefinedOverrides)
		} else {
			// iDRedeemsMap[redeem.ID] = MakeRedeem(redeem, DefaultExternalRedeem)
		}
	}
}

func CreateRedeem(redeem Redeem) string {
	jsonData, err := json.Marshal(redeem)
	if err != nil {
		logger.Error("failed to marshal redeem", err)
		return ""
	}
	resp, err := twitchapi.AsUser().Post("/channel_points/custom_rewards?broadcaster_id="+context.Get().GetBroadcaster().Id, string(jsonData))
	if err != nil {
		logger.Error("failed to create redeem", err)
		return ""
	}
	var data struct {
		Data []TwitchRedeem `json:"data"`
	}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		logger.Error("failed to unmarshal redeem", err)
		return ""
	}
	if len(data.Data) == 0 {
		logger.Error("failed to create redeem", err)
		return ""

	}
	return data.Data[0].ID
}

func UpdateRedeem(id string, redeem Redeem) {
	logger.Log("Updating redeem", "ID", id)
	jsonData, err := json.Marshal(redeem)
	if err != nil {
		logger.Error("failed to marshal redeem", err)
		return
	}
	_, err = twitchapi.AsUser().Patch("/channel_points/custom_rewards?broadcaster_id="+context.Get().GetBroadcaster().Id+"&id="+id, string(jsonData))
	if err != nil {
		logger.Error("failed to update redeem", err)
		return
	}
}

func MakeRedeem(data TwitchRedeem, override Redeem) Redeem {
	return Redeem{
		Actions:                           override.Actions,
		QueueName:                         override.QueueName,
		ID:                                data.ID,
		Title:                             data.Title,
		Cost:                              data.Cost,
		IsEnabled:                         data.IsEnabled,
		Prompt:                            data.Prompt,
		IsUserInputRequired:               data.IsUserInputRequired,
		IsMaxPerStreamEnabled:             data.MaxPerStreamSetting.IsEnabled,
		MaxPerStream:                      data.MaxPerStreamSetting.MaxPerStream,
		IsMaxPerUserPerStreamEnabled:      data.MaxPerUserPerStreamSetting.IsEnabled,
		MaxPerUserPerStream:               data.MaxPerUserPerStreamSetting.MaxPerUserPerStream,
		IsGlobalCooldownEnabled:           data.GlobalCooldownSetting.IsEnabled,
		GlobalCooldownSeconds:             data.GlobalCooldownSetting.GlobalCooldownSeconds,
		IsPaused:                          data.IsPaused,
		ShouldRedemptionsSkipRequestQueue: data.ShouldRedemptionsSkipRequestQueue,
	}
}

func GetRedeem(id string) Redeem {
	return iDRedeemsMap[id]
}

func ExportRedeems() []byte {
	redeems := []Redeem{}
	for _, redeem := range iDRedeemsMap {
		redeems = append(redeems, redeem)
	}
	out, err := json.Marshal(redeems)
	if err != nil {
		logger.Error("failed to marshal redeems", err)
		return nil
	}
	return out
}

func PauseRedeem(id string) {
	redeem := GetRedeem(id)
	redeem.IsPaused = true
	iDRedeemsMap[id] = redeem
	UpdateRedeem(id, redeem)
}

func UnpauseRedeem(id string) {
	redeem := GetRedeem(id)
	redeem.IsPaused = false
	iDRedeemsMap[id] = redeem
	UpdateRedeem(id, redeem)
}

func EnableRedeem(id string) {
	redeem := GetRedeem(id)
	redeem.IsEnabled = true
	iDRedeemsMap[id] = redeem
	UpdateRedeem(id, redeem)
}

func DisableRedeem(id string) {
	redeem := GetRedeem(id)
	redeem.IsEnabled = false
	iDRedeemsMap[id] = redeem
	UpdateRedeem(id, redeem)
}

func DeleteRedeem(id string) {
	body, err := twitchapi.AsUser().Delete("/channel_points/custom_rewards?broadcaster_id=" + context.Get().GetBroadcaster().Id + "&id=" + id)
	if err != nil {
		logger.Error("failed to delete redeem", err)
		return
	}
	logger.Log("Delete redeem", "ID", id, "Response", string(body))
	delete(iDRedeemsMap, id)
}
