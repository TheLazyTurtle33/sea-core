package eventsub

import (
	"fmt"
	"time"

	"github.com/TheLazyTurtle33/sea-core/bot/internal/context"
	"github.com/TheLazyTurtle33/sea-core/bot/internal/logger"
	twitchapi "github.com/TheLazyTurtle33/sea-core/bot/internal/twitch/api"
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
	if data.FromBroadcasterUserID == context.Get().GetBroadcaster().Id {
		if _, err := twitchapi.AsBot().SendMessage(
			fmt.Sprintf("BYYEEE EVERYONE :3 have fun in %s's channel ^w^, remember to behave >:c, raid messages coming in!",
				data.ToBroadcasterUserName,
			),
		); err != nil {
			logger.Error("Raid: Error sending message", err)
		}
		if _, err := twitchapi.AsBot().SendMessage("TURTLE RAID! TURTLE RAID! TURTLE RAID! TURTLE RAID! TURTLE RAID! TURTLE RAID! "); err != nil {
			logger.Error("Raid: Error sending message", err)
		}
		// sub raid, still need emots lol
		// if _, err := twitchapi.AsBot().SendMessage(); err != nil {
		// 	logger.Error("Raid: Error sending message", err)
		// }
		return
	}
	if data.ToBroadcasterUserID == context.Get().GetBroadcaster().Id {
		client, err := obs.Get()
		if err != nil {
			logger.Error("Raid: Error getting OBS client", err)
			return
		}

		raidText := fmt.Sprintf("%s Washed in \n with %d castaways", data.FromBroadcasterUserName, data.Viewers)

		if _, err = client.SetSourceVisability("RaidText", false); err != nil {
			logger.Error("Raid: Error hiding RaidText", err)
			return
		}

		if _, err = client.SetSourceVisability("Raid", true); err != nil {
			logger.Error("Raid: Error showing Raid group", err)
			return
		}

		if _, err = client.SetTextSourceText("RaidText", raidText); err != nil {
			logger.Error("Raid: Error uppdating RaidText", err)
			return
		}

		time.Sleep(TextShowWait)

		if _, err = client.SetSourceVisability("RaidText", true); err != nil {
			logger.Error("Raid: Error showing RaidText", err)
			return
		}

		time.Sleep(GroupHideWait)

		if _, err = client.SetSourceVisability("Raid", false); err != nil {
			logger.Error("Raid: Error showing Raid group", err)
			return
		}

		time.Sleep(TextHideWait)

		if _, err = client.SetSourceVisability("RaidText", false); err != nil {
			logger.Error("Raid: Error hiding RaidText", err)
			return
		}

		return
	}

	logger.Error("raid: unknown raid event format", nil)
}
