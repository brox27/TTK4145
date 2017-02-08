package main
import (
//	."fmt"
)

// Variables that should come from somewhere else!!
const NUM_ELEVATORS = 3;
var Elevators_Alive = 3;
const num_floors = 4;
type Direction int;
const (
 DOWN Direction = -1 + iota
 NEUTRAL
 UP
)
// end of var list


type All_Orders struct {
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
	hesten int
}



func Assign_Order(){
//	var cost_table [3]int;
	lowest := 1;
	cost := -1;
	Cheapest_Elev := -2;
	Cheapest_Elev +=1;
	for elev := 0; elev < Elevators_Alive; elev++{
		if (cost > lowest){
			println("kjorer cont..");
			//continue;		
		}

		if (cost == lowest){
			//CHECK IP addr.!
			println("LIK!")
			//continue;
		}

		if (cost < lowest) {
			println("cost less -> new best");
			lowest = cost;
			Cheapest_Elev = elev;
		}
	}
}

func (f *All_Orders) Update_Orders(floor int, elevator int){
	for i := 0; i<3;i++{
		f.elev3.orders[floor][i]=1
	}
}

// minor change
func main() {
	All_Orders := All_Orders{}
	All_Orders.Update_Orders(1,3)
	println(All_Orders.elev3.orders[1][1])
}