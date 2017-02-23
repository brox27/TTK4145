package driver

import (
	. "../ConfigFile"
	"time"
//	"fmt"
)

type ButtonState int
const (
	PRESSED ButtonState = iota
	RELEASED
)

func ButtonPoll(buttonChan chan [2]int) {
	LastStatus := [Num_floors][Num_buttons]int{}
	for{
	for floor := 0; floor <Num_floors; floor++{
		for button := 0; button<Num_buttons; button++{
			newStatus := GetButtonSignal(floor, button)
			//fmt.Printf("floor %d, button %d, status %d\n",floor, button, hest)
			if (LastStatus[floor][button] != newStatus) && (newStatus == 1){
				var temp [2]int
				temp[0] = floor
				temp[1] = button
				buttonChan <- temp
			}
			LastStatus[floor][button] = newStatus
		}
	}
	time.Sleep(1*time.Millisecond)
	}
}

func floorPoll(floorChan chan int) {
	LastStatus := -2
	for {
		floor := GetFloorSensorSignal()
		if (floor != -1) && (floor != LastStatus){
			floorChan <- floor
		}
		LastStatus = floor
	}
}



func EventHandler(eventChan chan Event) {
//	thisEvent := Event{}

	buttonChan := make(chan [2]int)
	go ButtonPoll(buttonChan)

	floorChan := make(chan int)
	go floorPoll(floorChan)
	for{
		select {
		case bc := <-buttonChan:
			//println("noe på btn!")
			temp := Event{}
			temp.EventType = BUTTONPRESSED
			temp.Button = bc[1]
			temp.Floor = bc[0]
			eventChan <- temp
		case fl := <-floorChan:
			//println("noe på etg")
			temp := Event{}
			temp.EventType = NEWFLOOR
			temp.Floor = fl
			eventChan <- temp

		default:
		}
	}

}

