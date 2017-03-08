package ConfigFile

const NUM_ELEVATORS = 3

var Elevators_Alive = 3
var LocalID = ""

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

var ELEVATOR_IPS = [NUM_ELEVATORS]string{"123.123.123", "321.321.321", "asd.asd.asd"}

type CabOrders struct {
	CabOrders [Num_floors]OrderStatus
	Direction Direction
	Floor     int
}

var AllOrders map[string]*Elev

var AllCabOrders map[string]*CabOrders

// med ny fsm

/////// ANDERS FIXES!! /////// /////// /////// /////// /////// /////// /////// /////// /////// ///////
const Port = 15647

type States int

const (
	INITIALIZE States = iota
	IDLE
	RUNNING
	DOORSOPEN
)

type Elev struct {
	State     States           `json:"hallRequests"`
	Floor     int              `json:"floor"`
	Direction Direction        `json:"direction"`
	CabOrders [Num_floors]bool `json:"cabRequests"`
	Orders    [Num_floors][Num_buttons]bool
}

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

type OrderState int

const ( //OrderState
	Default OrderState = iota + 1
	Inactive
	PendingAck
	Active
)

type OrderStatus struct {
	OrderState OrderState
	AckdBy     []string
}

type ConsensusCab struct {
	CabButtons [Num_floors]OrderStatus
}

type ConsensusHall struct {
	HallButtons [Num_floors][Num_buttons - 1]OrderStatus
}
