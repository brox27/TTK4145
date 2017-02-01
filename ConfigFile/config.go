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

type Signal int;
const (
 New_Order Signal = 1 + iota
 Order_Complete
 Heartbeat
 Acknowledge
 All_Orders
 Request_All // veldig trolig ikke nødvendig...
)


type All_Information struct {
	id int
	signal Signal
	elev_id int

	
	elev1 struct{
		Last_Floor int
    		Direction Direction

   		orders [num_floors][num_buttons]int
	}

   	elev2 struct{
		Last_Floor int
    		Direction Direction

   		orders [num_floors][num_buttons]int
	}

	elev3 struct{
		Last_Floor int
    		Direction Direction

   		orders [num_floors][num_buttons]int
	}
}


