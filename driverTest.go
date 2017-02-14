package main

import(
	. "./globals"
	. "./driver"
	. "fmt"
)

func main(){
	InitElev()
	for{
		floor := GetFloorSensorSignal()
		/*
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
		*/
		if (GetButtonSignal(2, floor)==1){
			Println("button pushed")
			SetButtonLamp(Button_Order_Command, floor-1, 1)
		}

	}
	defer SetMotorDirection(IDLE)
}