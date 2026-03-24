package chat

import (
	"log"
)

type EventData struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	ChatterUserID        string `json:"chatter_user_id"`
	ChatterUserLogin     string `json:"chatter_user_login"`
	ChatterUserName      string `json:"chatter_user_name"`
	MessageID            string `json:"message_id"`
	Message              struct {
		Text      string `json:"text"`
		Fragments []struct {
			Type      string      `json:"type"`
			Text      string      `json:"text"`
			Cheermote interface{} `json:"cheermote"`
			Emote     interface{} `json:"emote"`
			Mention   interface{} `json:"mention"`
		} `json:"fragments"`
	} `json:"message"`
	Color  string `json:"color"`
	Badges []struct {
		SetID string `json:"set_id"`
		ID    string `json:"id"`
		Info  string `json:"info"`
	} `json:"badges"`
	MessageType                 string      `json:"message_type"`
	Cheer                       interface{} `json:"cheer"`
	Reply                       interface{} `json:"reply"`
	ChannelPointsCustomRewardID interface{} `json:"channel_points_custom_reward_id"`
	SourceBroadcasterUserID     string      `json:"source_broadcaster_user_id"`
	SourceBroadcasterUserLogin  string      `json:"source_broadcaster_user_login"`
	SourceBroadcasterUserName   string      `json:"source_broadcaster_user_name"`
	SourceMessageID             string      `json:"source_message_id"`
	SourceBadges                []struct {
		SetID string `json:"set_id"`
		ID    string `json:"id"`
		Info  string `json:"info"`
	} `json:"source_badges"`
	IsSourceOnly bool `json:"is_source_only"`
}

func HandleMessage(data EventData) {
	log.Println(data.Message.Text)
}
