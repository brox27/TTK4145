package main

import (
	. "../ConfigFile"
)

func ordersAbove(f map[string]*Elev, ip string) bool {
	floor := f[ip].LastFloor 
	for i := floor; i < Num_floors; i++ {
		for j := 0; j < Num_buttons; j++ {
			if f[ip].Orders[i][j] == 1 { 
				return true
			}
		}
	}
	return false
}

func ordersBelow(f map[string]*Elev, ip string) bool{
	// trekker fra 2 for å "0" indeksere OG ikke sjekke etg den er i
	floor := f[ip].LastFloor - 2 
	for i := floor; i >= 0; i-- {
		for j := 0; j < Num_buttons; j++ {
			if f[ip].Orders[i][j] == 1 { 
				return true
			}
		}
	}
	return false
}

func nextDirection(f map[string]*Elev, ip string) Direction {
	
	if f[ip].Direction == UP {
		if ordersAbove(f, ip) {
			println("OPPOVER")
			return UP
		}

		if ordersBelow(f, ip) {
			println("NEDOVER")
			return DOWN
		} else {
			println("STANDA STILLE")
			return NEUTRAL
		}
	} else {
		if ordersBelow(f, ip) {
			println("NEDOVER")
			return DOWN
		}

		if ordersAbove(f, ip) {
			println("OPPOVER")
			return UP
		} else {
			println("STANDA STILLE")
			return NEUTRAL
		}
	}
}

func ShouldStop(f map[string]*Elev, ip string, curFloor int) bool {			// NB! sender inn curr floor! Se når struckt oppdateres om vi må det
	for i := 0; i < Num_buttons; i++ {
		if f[ip].Orders[curFloor-1][i] == 1 { 
			println("SHOULD STOP!")
			return true
		}
	}
	println("I DONT GIVE A FUCK!!! i just drive")
	return false
}

func calculateCost(f map[string]*Elev, ip string) int {
	return 12

}

func ipOfElevThatShouldTakeOrder(f map[string]*Elev, ip string, floor int, button int) string {			//Nytt navn!! + array med alle IPer??
	winnerIP := ""
	lowestCost := 99999999

	// her må en loop settes inn for alle..
	thisCost := calculateCost(f, ip)
	if (thisCost < lowestCost){
		winnerIP = ip
		lowestCost = thisCost	
	}

	if ( thisCost == lowestCost){
		if (ip < winnerIP){		// V. LIKHET ER DET DEN MED LAVEST IP SOM VINNER!!
			winnerIP = ip
		}
	}
	// slutt loopen

	return winnerIP
}


func recalculateOrders(f map[string]*Elev, ipOfDead string, ip string){	// bytte ut IP med array? over de som er levende
	for floor := 0; floor <Num_floors; floor++{
		for button := 0; button < (Num_buttons-1); button++{
			if f[ipOfDead].Orders[floor][button] == 1{
				ipWinner := ipOfElevThatShouldTakeOrder(f, ip, floor, button)
				f[ipWinner].Orders[floor][button] = 1	// bør en kanskje heller kalle en funksjon for dette..? mtp god kode?
			}
		}
	}
}

func main() {
	elev1 := Elev{}
	AllOrders = make(map[string]*Elev)
	AllOrders["123.123.13.123"] = &elev1


//

//
// UNDER ER TEST VARIABLER/FUNCTIONER
	TEST(AllOrders)
	ShouldStop(AllOrders, "123.123.13.123", 2)
	nextDirection(AllOrders, "123.123.13.123")
	println(ipOfElevThatShouldTakeOrder(AllOrders, "123.123.13.123",1,1))
}

//
//
//
//
//
//
//
//
//
//
//
//
//

func TEST(f map[string]*Elev) {
	


	// 1 etg
	for i := 0; i < 3; i++ {
		f["123.123.13.123"].Orders[0][i] = 1
	}

	// 2 etg
	for i := 0; i < 3; i++ {
		f["123.123.13.123"].Orders[1][i] = 0
	}

	//3 etg
	for i := 0; i < 3; i++ {
		f["123.123.13.123"].Orders[2][i] = 1
	}

	// 4etg
	for i := 0; i < 3; i++ {
		f["123.123.13.123"].Orders[3][i] = 0
	}
}