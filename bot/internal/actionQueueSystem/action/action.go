package action

type Flag struct {
	Active     bool     // if a flag is set
	StringData string   // if a flag needs to pass string data
	IntData    int      // if a flag needs to pass int data
	BoolData   bool     // if a flag needs to pass bool data
	Actions    []Action // the actoins for if a flag needs it
	ActionData []any    // the data for the added actions
}

type Flags struct {
	Error error

	Lock   Flag // locks the given queue. if QueueName is "" use curent queue
	Unlock Flag // unlocks the given queue. if QueueName is "" use curent queue

	AddActions  Flag // add an action at a given queue
	PassThrough any  // data for the next action in the queue

	StartQueue Flag // tell a given queue to spin up, usefull for non-persistent repeting queues
	StopQueue  Flag // tell a given queue to spin down, usefull for non-persistent repeting queues. if QueueName is "" use curent queue

	Pause Flag // tells a given queue to pause for a given number of seconds
	Skip  Flag // skips/culls the next IntData of actoins.
}

type Action interface {
	Run(passThrough any, v ...any) Flags
	OnAdd(passThrough any, v ...any) Flags // ran when this action is added to a queue
}
