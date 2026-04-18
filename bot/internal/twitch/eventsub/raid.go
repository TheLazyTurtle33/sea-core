package eventsub

import (
	"fmt"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	"github.com/TheLazyTurtle33/sea-core/shared/obs"
)

type raidData struct {
	FromBroadcasterUserID    string `json:"from_broadcaster_user_id"`
	FromBroadcasterUserLogin string `json:"from_broadcaster_user_login"`
	FromBroadcasterUserName  string `json:"from_broadcaster_user_name"`
	ToBroadcasterUserID      string `json:"to_broadcaster_user_id"`
	ToBroadcasterUserLogin   string `json:"to_broadcaster_user_login"`
	ToBroadcasterUserName    string `json:"to_broadcaster_user_name"`
	Viewers                  int    `json:"viewers"`
}

const (
	TextShowWait  time.Duration = 1 * time.Second
	GroupHideWait time.Duration = 2 * time.Second
	TextHideWait  time.Duration = 1 * time.Second
)

func HandleRaid(data raidData) {
	if context.Get().GetBroadcaster().Id == data.FromBroadcasterUserID {
	} else {
		client, err := obs.Get()
		if err != nil {
			logger.Error("Error getting OBS client", err)
			return
		}

		raidText := fmt.Sprintf("%s Washed in \n with %d castaways", data.FromBroadcasterUserName, data.Viewers)

		_, err = client.SetSourceVisability("RaidText", false)
		if err != nil {
			logger.Error("Error hiding RaidText", err)
			return
		}
		_, err = client.SetSourceVisability("Raid", true)
		if err != nil {
			logger.Error("Error showing Raid group", err)
			return
		}

		_, err = client.SetTextSourceText("RaidText", raidText)
		if err != nil {
			logger.Error("Error uppdating RaidText", err)
			return
		}

		time.Sleep(TextShowWait)

		_, err = client.SetSourceVisability("RaidText", true)
		if err != nil {
			logger.Error("Error showing RaidText", err)
			return
		}

		time.Sleep(GroupHideWait)

		_, err = client.SetSourceVisability("Raid", false)
		if err != nil {
			logger.Error("Error showing Raid group", err)
			return
		}

		time.Sleep(TextHideWait)

		_, err = client.SetSourceVisability("RaidText", false)
		if err != nil {
			logger.Error("Error hiding RaidText", err)
			return
		}

	}
}
