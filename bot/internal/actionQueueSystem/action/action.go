package action

import (
	"fmt"
)

var actionMap = make(map[string]Action)

type Flag struct {
	Active     bool     // if a flag is set
	StringData string   // if a flag needs to pass string data
	IntData    int      // if a flag needs to pass int data
	BoolData   bool     // if a flag needs to pass bool data
	Actions    []Action // the actoins for if a flag needs it
}

type Flags struct {
	Error error

	Lock   Flag // locks the given queue. if QueueName is "" use curent queue
	Unlock Flag // unlocks the given queue. if QueueName is "" use curent queue

	AddActions  Flag       // add an action at a given queue
	PassThrough ActionData // data for the next action in the queue

	StartQueue Flag // tell a given queue to spin up, usefull for non-persistent repeting queues
	StopQueue  Flag // tell a given queue to spin down, usefull for non-persistent repeting queues. if QueueName is "" use curent queue

	Pause Flag // tells a given queue to pause for a given number of seconds
	Skip  Flag // skips/culls the next IntData of actoins.
}

type ActionMetaData struct {
	// basic info
	Name        string `json:"name"`
	Description string `json:"description"`

	// Run func info
	RunDescription string       `json:"runDescription"`
	RunData        []ActionData `json:"runData"`
	RunPassThrough struct {
		ModefiesPassThrough bool       `json:"modefiesPassThrough"`
		Out                 ActionData `json:"out"`
	} `json:"runPassThrough"`

	// OnAdd func info
	OnAdd            bool         `json:"hasOnAdd"`
	OnAddDescription string       `json:"onAddDescription"`
	OnAddData        []ActionData `json:"onAddData"`
	OnAddPassThrough struct {
		ModefiesPassThrough bool       `json:"modefiesPassThrough"`
		Out                 ActionData `json:"out"`
	} `json:"onAddPassThrough"`
}

// type Action interface {
// 	Run(passThrough any, actionData []ActionData) Flags
// 	OnAdd(passThrough any, actionData []ActionData) Flags // ran when this action is added to a queue
// 	GetMetaData() ActionMetaData
// }

const ChatMessageType = "ChatMessageData" // use as ActionData to pass the chat message data to the action, if data is nil the command will replace it with ChatMessageData
const ActionsType = "Actions"             // use as ActionData to pass the action data to the next action
const StringType = "string"               // use as ActionData to pass string data
type ActionData struct {
	// for actions
	Type string `json:"type"`
	Data any    `json:"-"`
	// for matadate
	Required    bool   `json:"required"`
	Description string `json:"description"`
}
type Action struct {
	Run   func(passThrough ActionData, actionData []ActionData) Flags `json:"-"`
	OnAdd func(passThrough ActionData, actionData []ActionData) Flags `json:"-"`

	MetaData   ActionMetaData
	ActionData []ActionData `json:"-"`
}

func (a *Action) Make(actionData []ActionData) Action {
	return Action{
		Run:        a.Run,
		OnAdd:      a.OnAdd,
		MetaData:   a.MetaData,
		ActionData: actionData,
	}
}

func (a *Action) UnmarshalJSON(data []byte) error {
	return nil
}

func (a *Action) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func Errorf(a Action, format string, v ...any) Flags {
	formatFull := fmt.Sprintf("action name: %s, error from action: %s", a.MetaData.Name, format)
	return Flags{Error: fmt.Errorf(formatFull, v...)}
}
