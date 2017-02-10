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

type NewOrder struct {
	MsgId  int
	Floor  int
	Button int
}

type CompleteOrder struct {
	MsgId  int
	Floor  int
	Button int
}

type Acknowledge struct {
	MsgId int
}

type Heartbeat struct {
	MsgId int
}

type Elev struct {
	LastFloor int
	Direction Direction
	Orders    [Num_floors][Num_buttons]int
	//id 			string
}

var AllOrders map[string]*Elev
