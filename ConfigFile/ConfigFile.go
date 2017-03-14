package ConfigFile

import (
	"sync"
)

const Num_floors = 4
const Num_buttons = 3

const MOTOR_SPEED = 2800

const HallConsensusPort = 15647
const CabConsensusPort = 15657
const ElevatorStatesPort = 15667
const PeersPort = 15677
const LocalHostPort = 15687

var LocalID = ""

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

type States int

const (
	INITIALIZE States = iota
	IDLE
	MOVING
	DOORSOPEN
)

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

type Elev struct {
	sync.RWMutex `json:"-"`
	State        States
	Floor        int
	Direction    Direction
	Orders       [][]bool
}

func NewElev() Elev {
	var e Elev
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

type ConsensusCab struct {
	CabButtons [Num_floors]OrderStatus
	ID         string
}

type ConsensusHall struct {
	HallButtons [Num_floors][Num_buttons - 1]OrderStatus
	ID          string
}

type AllStates struct {
	sync.RWMutex `json:"-"`
	StateMap     map[string]*Elev
}
