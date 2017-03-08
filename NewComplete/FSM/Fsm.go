package FSM

import (
	. "../driver"
	//	. "../globals"
	"../ConfigFile"
	//	. "../OrderManager"
	//. "fmt"
	//	"runtime"
	"time"
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
func RUN(
	FloorChan chan int, StateChan chan ConfigFile.Elev,
	LocalOrdersChan chan [ConfigFile.Num_floors][ConfigFile.Num_buttons]bool,
	ClearHallOrdersChan chan [2]int, ClearCabOrderChan chan int) {

	LocalElev := ConfigFile.Elev{}

	timerChan := make(chan int)
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
				break

			case ConfigFile.RUNNING:
				if shouldStop(LocalElev) { // se over, kan ha noen mangler, eks. når heisen allerede står i etg hvor det bestilles
					for button := 0; button < ConfigFile.Num_buttons; button++ {
						if LocalElev.Orders[LocalElev.Floor][button] {
							if button < ConfigFile.Num_buttons-1 {
								ClearHallOrdersChan <- [2]int{LocalElev.Floor, button}
							} else {
								ClearCabOrderChan <- LocalElev.Floor
							}
						}
					}
					SetMotorDirection(ConfigFile.NEUTRAL)
					go timer(timerChan)
					SetDoorOpenLamp(1)
					LocalElev.State = ConfigFile.DOORSOPEN
					StateChan <- LocalElev
					break
				}

			case ConfigFile.DOORSOPEN:
				temp := <-timerChan
				_ = temp
				SetDoorOpenLamp(0)
				LocalElev.State = ConfigFile.IDLE
				break
			}
		}
	}

	defer SetMotorDirection(ConfigFile.NEUTRAL)
}

func nextDirection(LocalElev ConfigFile.Elev) ConfigFile.Direction {
	if LocalElev.Direction == ConfigFile.UP {
		if ordersAbove(LocalElev) {
			return ConfigFile.UP
		}

		if ordersBelow(LocalElev) {
			return ConfigFile.DOWN
		} else {
			return ConfigFile.NEUTRAL
		}
	} else {
		if ordersBelow(LocalElev) {
			return ConfigFile.DOWN
		}

		if ordersAbove(LocalElev) {
			return ConfigFile.UP
		} else {
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

func timer(timerChan chan int) {
	time.Sleep(3 * time.Second)
	timerChan <- 1
}
