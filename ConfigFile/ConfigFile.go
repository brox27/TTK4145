package ConfigFile

const NUM_ELEVATORS = 3

var Elevators_Alive = 3

const Num_floors = 4
const Num_buttons = 3

const MOTOR_SPEED = 2800

type Direction int

const (
	UP Direction = iota
	DOWN
	NEUTRAL
)

type ButtonType int

const (
	BUTTON_ORDER_UP ButtonType = iota
	BUTTON_ORDER_DOWN
	BUTTON_ORDER_COMMAND
)




type MsgType int

const (
	NEW MsgType = iota + 1
	ACKNOWLEDGE
	COMPLETE
)

type OrderMsg struct {
	Floor   int
	Button  int
	MsgType int
}


var ELEVATOR_IPS=[NUM_ELEVATORS]string{"123.123.123","321.321.321","asd.asd.asd"}

type OrderState int

/*interface merge{
	LocalID string
	RemoteID string
	Peers[] string
}
*/

const (	//OrderState 
	Default OrderState = iota +1
	Inactive
	PendingAck
	Active
)

type OrderStatus struct {
	OrderState  OrderState
	AckdBy 		[]string
}

type AllHallOrders struct{
	HallOrders [Num_floors][(Num_buttons-1)] OrderStatus
}

type CabOrders struct{
	CabOrders [Num_floors] OrderStatus
	Direction Direction
	Floor int
}


type Elev struct {
	LastFloor int
	Direction Direction
	Orders    [Num_floors][Num_buttons]int
	// bytte ut over med caborders
	//id 			string
}

var AllOrders map[string]*Elev
var AllCabOrders map[string]*CabOrders


// med ny fsm

type States int
const (
	INITIALIZE States = iota
	IDLE
	RUNNING
	DOORSOPEN
)

type EventType int
const (
	BUTTONPRESSED = iota
	NEWFLOOR
)

type Event struct{
	EventType EventType
	Floor int
	Button int
}
