package queue

var DefaultQueue = Queue{
	name:        "default",
	locked:      false,
	repeating:   false,
	persistent:  true,
	repeatDelay: 0,
}

var RedeemsQueue = Queue{
	name:       "redeems",
	persistent: true,
}

// var DiscordQueue = Queue{
// 	name:        "Discord",
// 	locked:      true,
// 	repeating:   true,
// 	persistent:  true,
// 	repeatDelay: 45 * time.Minute,
// 	actions:     []action.Action{&actions.CreateDiscordInvite{}, &actions.SendMessage{}},
// }

var TTSQueue = Queue{
	name:       "tts",
	locked:     false,
	repeating:  false,
	persistent: true,
}

var Queues = []*Queue{
	&DefaultQueue,
	// &DiscordQueue,
	&RedeemsQueue,
	&TTSQueue,
}

func GetQueue(name string) *Queue {
	for _, q := range Queues {
		if q.name == name {
			return q
		}
	}
	return nil
}

func StartUp() {
	DefaultQueue.Start()
	// DiscordQueue.puased += 30 * time.Minute
	// DiscordQueue.Start()
}
