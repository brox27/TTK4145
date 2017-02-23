package driver

import (
	. "../ConfigFile"
	"time"
)

type ButtonState int
const (
	PRESSED ButtonState = iota
	RELEASED
)
/*
func ButtonPoll(buttonChan chan [2]int) {
	LastStatus := [Num_floors][Num_buttons]bool{}
	for floor := 0; floor <Num_floors; floor++{
		for button := 0; button<Num_buttons; button++{

		}
	}
}

func floorPoll(floorChan chan int) {
	for {
		//floorChan <- GetFloorSensorSignal()
		println("lal")
	}
}

func EventHandler(eventChan chan Event) {
	thisEvent := Event{}

	buttonChan := make(chan [2]int)
	go ButtonPoll(buttonChan)
	floorChan := make(chan int)
	go floorPoll(floorChan)

	for{
		select {
		case bc := <-buttonChan:
			println("noe på btn!")
		case fl := <-floorChan:
			println("noe på etg")

		default:
			println("lllaaall")
		// sjekke buttons 

		// sjekke etg.
		}
	}
	

}
*/

func Lala(){
	println("KUN TEST!")
	for{
		hesten := GetFloorSensorSignal()
		if hesten != -1{
			println(hesten)
		}
	}
	time.Sleep(Num_floors*time.Second)
}