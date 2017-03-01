package FSM

import (
	. "../driver"
	//	. "../globals"
	"../ConfigFile"
	//	. "../OrderManager"
	//. "fmt"
	//	"runtime"
	//"time"
)

/*
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
*/

// 	ClearHallOrderChan chan []int, ClearCabOrderChan chan int)
func RUN(FloorChan chan int, StateChan chan ConfigFile.Elev, LocalOrdersChan chan [ConfigFile.Num_floors][ConfigFile.Num_buttons]bool) {
	//	State := INITIALIZE
	//	go EventHandler(eventChan)
	//	hest := GetFloorSensorSignal()
	//	println(hest)
	LocalElev := ConfigFile.Elev{}

	for {
		select {
		case newFloor := <-FloorChan:
			LocalElev.Floor = newFloor
			StateChan <- LocalElev
			SetFloorLight(newFloor - 1) // oppdatere mtp 1 indeksering **********************************************************************************
		case newOrders := <-LocalOrdersChan:
			LocalElev.Orders = newOrders

		default:
			switch LocalElev.State {
			case ConfigFile.INITIALIZE:
				if GetFloorSensorSignal() != -1 {
					LocalElev.Floor = GetFloorSensorSignal()
					SetFloorLight(GetFloorSensorSignal())
					SetMotorDirection(ConfigFile.NEUTRAL)
					LocalElev.State = ConfigFile.IDLE
					LocalElev.Direction = ConfigFile.NEUTRAL
					StateChan <- LocalElev

					break
				} else {
					SetMotorDirection(ConfigFile.DOWN) // sjekke noe så den ikke "settes" så ofte?
				}
				break

			case ConfigFile.IDLE:
				if nextDirection(LocalElev) != ConfigFile.NEUTRAL {
					SetMotorDirection(LocalElev.Direction)
					LocalElev.State = ConfigFile.RUNNING
					LocalElev.Direction = nextDirection(LocalElev)
					StateChan <- LocalElev
				}
				if shouldStop(LocalElev) {
					LocalElev.State = ConfigFile.DOORSOPEN
					// oppdatere ORDER COMPLETE!
				}
				break

			case ConfigFile.RUNNING:
				if shouldStop(LocalElev) {
					SetMotorDirection(ConfigFile.NEUTRAL)
					// Starte timer
					SetDoorOpenLamp(1)
					// Oppdatere mtp ordre ferdig
					LocalElev.State = ConfigFile.DOORSOPEN
					StateChan <- LocalElev
					break
				}

			case ConfigFile.DOORSOPEN:
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

	defer SetMotorDirection(ConfigFile.NEUTRAL)
}

func nextDirection(LocalElev ConfigFile.Elev) ConfigFile.Direction {
	if LocalElev.Direction == ConfigFile.UP {
		if ordersAbove(LocalElev) {
			println("OPPOVER")
			return ConfigFile.UP
		}

		if ordersBelow(LocalElev) {
			println("NEDOVER")
			return ConfigFile.DOWN
		} else {
			println("STANDA STILLE")
			return ConfigFile.NEUTRAL
		}
	} else {
		if ordersBelow(LocalElev) {
			println("NEDOVER")
			return ConfigFile.DOWN
		}

		if ordersAbove(LocalElev) {
			println("OPPOVER")
			return ConfigFile.UP
		} else {
			println("STANDA STILLE")
			return ConfigFile.NEUTRAL
		}
	}
}

func ordersAbove(LocalElev ConfigFile.Elev) bool {
	floor := LocalElev.Floor
	for i := floor; i < ConfigFile.Num_floors; i++ {
		for j := 0; j < ConfigFile.Num_buttons; j++ {
			if LocalElev.Orders[i][j] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(LocalElev ConfigFile.Elev) bool {
	// trekker fra 2 for å "0" indeksere OG ikke sjekke etg den er i
	floor := LocalElev.Floor - 2
	for i := floor; i >= 0; i-- {
		for j := 0; j < ConfigFile.Num_buttons; j++ {
			if LocalElev.Orders[i][j] {
				return true
			}
		}
	}
	return false
}

func shouldStop(LocalElev ConfigFile.Elev) bool {
	if LocalElev.Orders[LocalElev.Floor][2] {
		return true
	} else if LocalElev.Direction == ConfigFile.UP {
		if LocalElev.Orders[LocalElev.Floor][1] {
			return true
		} else {
			return (!ordersAbove(LocalElev))
		}
	} else if LocalElev.Direction == ConfigFile.DOWN {
		if LocalElev.Orders[LocalElev.Floor][0] {
			return true
		} else {
			return (!ordersBelow(LocalElev))
		}
	}
	return false
}
