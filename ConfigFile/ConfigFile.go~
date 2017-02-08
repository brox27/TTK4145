package ConfigFile


const NUM_ELEVATORS = 3;
var Elevators_Alive = 3;
const num_floors = 4;
const num_buttons = 3;

type Direction int;
const (
 DOWN Direction = -1 + iota
 NEUTRAL
 UP
)

type newOrder struct {
	msgId int
	floor int
	button int

}

type completeOrder struct{
	msgId int
	floor int
	button int
}

type acknowledge struct {
	msgId int
}

type heartbeat struct{
	msgId int
}

type allOrders struct {

	elev1 struct{
		Last_Floor int
    		Direction Direction
   		orders [][]int
	}

	elev2 struct{
		Last_Floor int
    		Direction Direction
   		orders [][]int
	}

	elev3 struct{
		Last_Floor int
    		Direction Direction
   		orders [4][3]int
	}
}
