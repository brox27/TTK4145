package main

import (
	. "./driver"
	. "./globals"
	. "fmt"
	. "time"
)

func main() {
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
	buttonChannel := make(chan ButtonType, N_BUTTONS*N_FLOORS)
	go ButtonCheck(buttonChannel)
	for {
		button := <-buttonChannel
		Println(button)
		Sleep(Second)
	}
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

	defer SetMotorDirection(IDLE)
}
