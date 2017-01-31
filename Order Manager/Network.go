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


type Elev struct {
    Last_Floor int
    Direction Direction

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

    orders := [][]int{}
    row1 := []int{1, 2, 3}
    row2 := []int{4, 5, 6}
    orders = append(orders, row1)
    orders = append(orders, row2)
}