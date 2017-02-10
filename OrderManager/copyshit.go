package main

import (
	. "../ConfigFile"
)

func assignOrder() {
	//	var cost_table [3]int;
	lowest := 1
	cost := -1
	Cheapest_Elev := -2
	Cheapest_Elev += 1
	for elev := 0; elev < 2; elev++ {
		if cost > lowest {
			println("kjorer cont..")
			//continue;
		}

		if cost == lowest {
			//CHECK IP addr.!
			println("LIK!")
		}

		if cost < lowest {
			println("cost less -> new best")
			lowest = cost
			Cheapest_Elev = elev
		}
	}
}

func (f *Orders) UpdateOrders(floor int, elevator int) {
	f.Hest = 145
}

func (f *Orders) ShouldStop(floor int, elevator int) bool {
	for i := 0; i < Num_buttons; i++ {
		if f.Elev1.Orders[floor][i] == 1 { //NB!! HARDKODET!!! PÅ ELEV1 $$$$ DETTE MÅ ENDRES!\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ <-- se der er du grei
			return true
		}
	}
	return false
}

func (f *Orders) ordersAbove() bool {
	floor := f.Elev1.LastFloor //NB!! HARDKODET!!! PÅ ELEV1 $$$$ DETTE MÅ ENDRES!\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ <-- se der er du grei
	for i := floor; i < Num_floors; i++ {
		for j := 0; j < Num_buttons; j++ {
			if f.Elev1.Orders[i][j] == 1 { //NB!! HARDKODET!!! PÅ ELEV1 $$$$ DETTE MÅ ENDRES!\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ <-- se der er du grei
				return true
			}
		}
	}
	return false
}

func ordersBelow() bool { // tar inn den MAP saken
	// trekker fra 2 for å "0" indeksere OG ikke sjekke etg den er i
	floor := f.Elev1.LastFloor - 2 //NB!! HARDKODET!!! PÅ ELEV1 $$$$ DETTE MÅ ENDRES!\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ <-- se der er du grei
	for i := floor; i >= 0; i-- {
		for j := 0; j < Num_buttons; j++ {
			if f.Elev1.Orders[i][j] == 1 { //NB!! HARDKODET!!! PÅ ELEV1 $$$$ DETTE MÅ ENDRES!\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ <-- se der er du grei
				return true
			}
		}
	}
	return false
}

func (f *Orders) nextDirection(e Elev) Dirn {
	direction := f.Elev1.Direction //NB!! HARDKODET!!! PÅ ELEV1 $$$$ DETTE MÅ ENDRES!\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ <-- se der er du grei
	if direction == UP {
		if f.ordersAbove() {
			println("OPPOVER")
		}

		if f.ordersBelow() {
			println("NEDOVER")
		} else {
			println("STANDA STILLE")
		}
	} else {
		if f.ordersBelow() {
			println("NEDOVER")
		}

		if f.ordersAbove() {
			println("OPPOVER")
		} else {
			println("STANDA STILLE")
		}
	}
}

func main() {
	elev1 := Elev{}
	allorder["minIP"] = elev1
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

func (f *Orders) TEST() {
	// 1 etg
	f.Elev1.LastFloor = 2
	f.Elev1.Direction = DOWN

	for i := 0; i < 3; i++ {
		f.Elev1.Orders[0][i] = 0
	}

	// 2etg
	for i := 0; i < 3; i++ {
		f.Elev1.Orders[1][i] = 1
	}

	//3 etg
	for i := 0; i < 3; i++ {
		f.Elev1.Orders[2][i] = 0
	}

	// 4etg
	for i := 0; i < 3; i++ {
		f.Elev1.Orders[3][i] = 0
	}
}
