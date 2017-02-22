package ConfigFile

const NUM_ELEVATORS = 3

var Elevators_Alive = 3

const Num_floors = 4
const Num_buttons = 3

type Direction int

const (
	DOWN Direction = -1 + iota
	NEUTRAL
	UP
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

const OrderState (
	inactive
	staged
	active
)
type OrderStatus struct {
	OrderState  int
	AckdBy 		[]string
}


type Elev struct {
	LastFloor int
	Direction Direction
	Orders    [Num_floors][Num_buttons]int
	// bytte ut over med caborders
	//id 			string
}

var AllOrders map[string]*Elev
