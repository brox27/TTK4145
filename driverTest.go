package driverTest

import (
	. "./driver"
	. "./globals"
	. "fmt"
	"runtime"
	//. "time"
)

func driverTest() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	InitElev()

	/*
		for{
			floor := GetFloorSensorSignal()
			if (floor==-1){
				Println("no floor")
				SetDoorOpenLamp(0)
				SetMotorDirection(UP)
			}else{
				Println("%d", floor)
				varr := floor -1
				SetFloorLight(varr)
				SetDoorOpenLamp(1)
				SetMotorDirection(IDLE)
			}
			if (GetButtonSignal(2, floor)==1){
				Println("button pushed")
				SetButtonLamp(Button_Order_Command, floor-1, 1)
			}

		}
	*/

	/*
		buttonChannel := make(chan [2]int, N_BUTTONS*N_FLOORS)
		go ButtonPoll(buttonChannel)
		for {
			data := <-buttonChannel //data = {button, floor}
			Println(data)
		}
	*/

	/*
		for {
			floor := GetFloorSensorSignal()
			if GetButtonSignal(Button_Order_Up, 1) == 1 {
				Println("up")
				SetMotorDirection(UP)
			}
			if GetButtonSignal(Button_Order_Down, 4) == 1 {
				Println("down")
				SetMotorDirection(DOWN)
			}
			if GetFloorSensorSignal() != -1 {
				SetFloorLight(floor - 1)
				SetMotorDirection(IDLE)
			}
			for j := 1; j < 5; j++ {
				if GetButtonSignal(Button_Order_Down, j) == 1 {
					SetButtonLamp(Button_Order_Down, j, 1)
				}
				if GetButtonSignal(Button_Order_Up, j) == 1 {
					SetButtonLamp(Button_Order_Up, j, 1)
				}
				if GetButtonSignal(Button_Order_Command, j) == 1 {
					SetButtonLamp(Button_Order_Command, j, 1)
				}
			}
			//SetButtonLamp(Button_Order_Down,4,0)
			//SetButtonLamp(Button_Order_Command,4,0)
			varrr := floor
			varrr += 1
		}
	*/

	eventChan := make(chan map[Event]interface{})
	go EventHandler(eventChan)
	for {
		Println(<-eventChan)
	}

	defer SetMotorDirection(NEUTRAL)
}
