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

type newOrder struct {
	MsgId  int
	Floor  int
	Button int
}

type completeOrder struct {
	MsgId  int
	Floor  int
	Button int
}

type acknowledge struct {
	MsgId int
}

type Heartbeat struct {
	MsgId int
}

type AllOrders struct {
	Elev1 struct {
		LastFloor int
		Direction Direction
		Orders    [Num_floors][Num_buttons]int
	}

	Elev2 struct {
		LastFloor int
		Direction Direction
		Orders    [Num_floors][Num_buttons]int
	}

	Elev3 struct {
		LastFloor int
		Direction Direction
		Orders    [Num_floors][Num_buttons]int
	}
	Hest int
}
