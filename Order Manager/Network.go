package main
import (
//	."fmt"
)

// Variables that should come from somewhere else!!
var NUM_ELEVATORS = 3;
var Elevators_Alive = 3;
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
   		orders [][]int
	}
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
	//Cheapest_Elev += 12;
}


func main() {
	Assign_Order();
}