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

const HallConesnsusPort = 15647
const CabConsensusPort = 15657
const ElevatorStatesPort = 15667
const PeersPort = 15677

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
	CabOrders []bool 		   `json:"cabRequests"`
	Orders    [][]bool
}

// Ander lagde, ikke i bruk \\
func NewElev() Elev {
	var e Elev
	e.CabOrders = make([]bool, Num_floors)
	e.Orders = make([][]bool, Num_floors)
	for i := range e.Orders {
		e.Orders[i] = make([]bool, Num_buttons)
	}
	return e
}

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

type OrderState int
const (
	Default OrderState = iota
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
	ID string
}

type ConsensusHall struct {
	HallButtons [Num_floors][Num_buttons - 1]OrderStatus
	ID string
}

const ColorCC = "\x1b[38;5;177m"
const ColorCH = "\x1b[38;5;175m"
const ColorFSM = "\x1b[38;5;208m"
const ColorHRA = "\x1b[38;5;79m"
const ColorNone = "\x1b[0m" 