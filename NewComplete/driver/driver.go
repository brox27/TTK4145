package driver

import (
	. "../ConfigFile"
	"time"
)

func ButtonPoll(HallButtonChan chan [2]int, CabButtonChan chan int) {
	LastStatus := [Num_floors][Num_buttons]int{}
	for {
		for floor := 0; floor < Num_floors; floor++ {
			for button := 0; button < Num_buttons; button++ {

				newStatus := GetButtonSignal(floor, button)
				if (LastStatus[floor][button] != newStatus) && (newStatus == 1) {
					if button == 2 {
						CabButtonChan <- floor
					} else {
						var hallOrder [2]int
						hallOrder[0] = floor
						hallOrder[1] = button
						HallButtonChan <- hallOrder
					}
				}
				LastStatus[floor][button] = newStatus
			}
		}
		time.Sleep(1 * time.Millisecond) // vi tror dette er god kode... Anders???
	}
}

func FloorPoll(FloorChan chan int) {
	LastStatus := -2
	for {
		floor := GetFloorSensorSignal()
		if (floor != -1) && (floor != LastStatus) {
			FloorChan <- floor
		}
		LastStatus = floor
	}
}
