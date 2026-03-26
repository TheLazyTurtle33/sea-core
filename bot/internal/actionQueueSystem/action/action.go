package action

type Flag struct {
	Active     bool     // if a flag is set
	Actions    []Action // the actoins for if a flag needs it
	ActionData []any    // the data for the added actions
	QueueName  string   // the name of a queue if a flag needs it
}

type Flags struct {
	Error             error
	Lock              Flag // locks the given queue. if QueueName is "" use curent queue
	Unlock            Flag // unlocks the given queue. if QueueName is "" use curent queue
	AddActions        Flag // add an action at a given queue
	DataForNextAction any  // data for the next action in the queue
	StartQueue        Flag // tell a given queue to spin up, usefull for non-persistent repeting queues
	StopQueue         Flag // tell a given queue to spin down, usefull for non-persistent repeting queues. if QueueName is "" use curent queue
}

type Action interface {
	Run(v ...any) Flags   // v[0] is the data from the previous action
	OnAdd(v ...any) Flags // ran when this action is added to a queue
}
