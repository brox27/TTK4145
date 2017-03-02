package main

import (
	. "../driver"
	//	. "../globals"
	. "../ConfigFile"
	//	. "../OrderManager"
	//. "fmt"
	"runtime"
	//"time"
)

func main() {
	InitElev()
	println(" 
		░░░░░░░░░░░░░░░░░░░░░ \n 
		░░░░░░░░░░░░▄▀▀▀▀▄░░░ \n 
		░░░░░░░░░░▄▀░░▄░▄░█░░ \n 
		░▄▄░░░░░▄▀░░░░▄▄▄▄█░░ \n 
		█░░▀▄░▄▀░░░░░░░░░░█░░ \n 
		░▀▄░░▀▄░░░░█░░░░░░█░░ \n 
		░░░▀▄░░▀░░░█░░░░░░█░░ \n 
		░░░▄▀░░░░░░█░░░░▄▀░░░ \n 
		░░░▀▄▀▄▄▀░░█▀░▄▀░░░░░ \n 
		░░░░░░░░█▀▀█▀▀░░░░░░░ \n
		░░░░░░░░▀▀░▀▀░░░░░░░░ \n")


	RUN(localId)
}

func RUN(FloorChan chan int, StateChan chan ConfigFile.State, LocalOrdersChan chan [ConfigFile.Num_floors][ConfigFile.Num_buttons]bool,
	ClearHallOrderChan chan []int, ClearCabOrderChan chan int)) {
//	State := INITIALIZE
	//	go EventHandler(eventChan)
	//	hest := GetFloorSensorSignal()
	//	println(hest)
	LocalElev = ConfigFile.Elev{}

	EventFloor := false
	for {
		select {
		case newFloor := <-FloorChan:
			LocalElev.Floor = newFloor
			StateChan <- LocalElev
			SetFloorLight(newFloor - 1) // oppdatere mtp 1 indeksering **********************************************************************************
		case newOrders := <- LocalOrdersChan:
			LocalElev.Orders = newOrders

		default:
			switch LocalElev.State {
			case INITIALIZE:
				if GetFloorSensorSignal() != -1 {
					LocalElev.Floor = GetFloorSensorSignal()
					SetFloorLight(GetFloorSensorSignal())
					SetMotorDirection(NEUTRAL)
					LocalElev.State = IDLE
					LocalElev.Direction = NEUTRAL
					StateChan <- LocalElev

					break
				} else {
					SetMotorDirection(DOWN) // sjekke noe så den ikke "settes" så ofte?
				}
				break

			case IDLE:
				if nextDirection(LocalElev) != NEUTRAL {
					SetMotorDirection(LocalElevState.Direction)
					LocalElev.State = RUNNING
					LocalElev.Direction = nextDirection(LocalElev)
					StateChan <- LocalElev
				}
				if shouldStop(LocalElev){
					LocalElev.State = DOORSOPEN
					// oppdatere ORDER COMPLETE!
				}
				break

			case RUNNING:
				if shouldStop(LocalElev) { 
					SetMotorDirection(NEUTRAL)
					// Starte timer
					elev.SetDoorOpenLamp(1)
					// Oppdatere mtp ordre ferdig
					LocalElev.State = DOORSOPEN
					StateChan <- LocalElev
					break
				}

			case DOORSOPEN:
	/*
				if TIMEOUT{											// LAG TIMEOUT
					Skru av timer?
					elev.SetDoorOpenLamp(0)
					nextDir = nextDirection(LocalElev)
					if nextDir == NEUTRAL{
						LocalElev.State=IDLE
					}else{
						elev.SetMotorDirection(nextDir)
						LocalElev.State=RUNNING
				}
	*/
			}
		}
	}

	defer SetMotorDirection(NEUTRAL)
}

func nextDirection(LocalElev ConfigFile.Elev) Direction {
		if LocalElev.Direction == UP {
		if ordersAbove(LocalElev) {
			println("OPPOVER")
			return UP
		}

		if ordersBelow(LocalElev) {
			println("NEDOVER")
			return DOWN
		} else {
			println("STANDA STILLE")
			return NEUTRAL
		}
	} else {
		if ordersBelow(LocalElev) {
			println("NEDOVER")
			return DOWN
		}

		if ordersAbove(LocalElev) {
			println("OPPOVER")
			return UP
		} else {
			println("STANDA STILLE")
			return NEUTRAL
		}
	}
}

func ordersAbove(LocalElev ConfigFile.Elev) bool {
	floor := LocalElev.Floor 
	for i := floor; i < Num_floors; i++ {
		for j := 0; j < Num_buttons; j++ {
			if LocalElev.Orders[i][j]{ 
				return true
			}
		}
	}
	return false
}

func ordersBelow(LocalElev ConfigFile.Elev) bool{
	// trekker fra 2 for å "0" indeksere OG ikke sjekke etg den er i
	floor := LocalElev.Floor - 2 
	for i := floor; i >= 0; i-- {
		for j := 0; j < Num_buttons; j++ {
			if LocalElev.Orders[i][j] { 
				return true
			}
		}
	}
	return false
}

func shouldStop(LocalElev ConfigFile.Elev) bool{
	if LocalElev.Orders[LocalElev.Floor][2]{
		return true
	}else if LocalElev.Direction == UP{
		if (LocalElev.Orders[LocalElev.Floor][1]){
			return true
		}else{
			return (!ordersAbove(LocalElev))
		}
	}else if LocalElev.Direction == DOWN{
		if LocalElev.Orders[LocalElev.Floor][0]{
			return true
		}else{
			return(!ordersBelow(LocalElev))
		}
	}
	return false
}