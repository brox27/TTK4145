package driver

import (
	. "../globals"
	//. "fmt"
	"time"
)

type ButtonState int

const (
	PRESSED ButtonState = iota
	RELEASED
)

type Event int

const (
	BUTTON_PRESSED Event = iota
	NEW_FLOOR
	//TimedOut
)

func ButtonPoll(buttonChan chan [2]int) {
	ButtonMap := make(map[[2]int]ButtonState) //Ensures that each
	for i := 0; i < N_BUTTONS; i++ {          //buttonpress is registered only once
		for j := 1; j <= N_FLOORS; j++ {
			ButtonMap[[2]int{i, j}] = RELEASED
		}
	}

	for {
		for i := 0; i < N_BUTTONS; i++ {
			for j := 1; j <= N_FLOORS; j++ {
				if (GetButtonSignal(ButtonType(i), j) == 1) && (ButtonMap[[2]int{i, j}] != Pressed) {
					buttonChan <- [2]int{i, j}
					ButtonMap[[2]int{i, j}] = PRESSED
				}
			}
		}

		for i := 0; i < N_BUTTONS; i++ {
			for j := 1; j <= N_FLOORS; j++ {
				if GetButtonSignal(ButtonType(i), j) == 0 {
					ButtonMap[[2]int{i, j}] = RELEASED
				}
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func floorPoll(floorChan chan int) {
	for {
		floorChan <- GetFloorSensorSignal()
	}
}

func EventHandler(eventChan chan map[Event]interface{}) {

	eventMap := make(map[Event]interface{}) //maps necessary data (value) to an event (key)

	buttonChan := make(chan [2]int)
	go ButtonPoll(buttonChan)
	buttonData := <-buttonChan

	for {

		data, ok := eventMap[ButtonPressed].([2]int)

		if !ok {
			eventMap[BUTTON_PRESSED] = buttonData
			eventChan <- eventMap

		} else {

			if &buttonData != nil {

				if !((buttonData[0] == 0 && buttonData[1] == 4) || (buttonData[0] == 1 && buttonData[1] == 1)) {

					if !(buttonData[0] == data[0] && buttonData[1] == data[1]) {
						eventMap[BUTTON_PRESSED] = buttonData
						eventChan <- eventMap
					}
				}
			}
		}

		floor := GetFloorSensorSignal()

		data2, ok2 := eventMap[NewFloor].(int)

		if !ok2 {
			eventMap[NEW_FLOOR] = floor
			eventChan <- eventMap

		} else {

			if floor != data2 {
				eventMap[NEW_FLOOR] = floor
				eventChan <- eventMap
			}
		}

		//what else
	}

}
