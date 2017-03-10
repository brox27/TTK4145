package FSM

import (
	. "../driver"
	"../ConfigFile"
	"fmt"
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

func RUN(
	FloorChan chan int, StateChan chan ConfigFile.Elev,
	LocalOrdersChan chan [][]bool,
	ClearHallOrdersChan chan [2]int, ClearCabOrderChan chan int) {


	defer SetMotorDirection(ConfigFile.NEUTRAL)
	LocalElev := ConfigFile.NewElev()
	var doorTimerChan <-chan time.Time

	{
		f := GetFloorSensorSignal()
		if f == -1 {
			SetMotorDirection(ConfigFile.DOWN)
			LocalElev.State = ConfigFile.INITIALIZE
		}
	}

	for {
		select {
		case newFloor := <-FloorChan:
			fmt.Printf("New floor: %+v\n", newFloor)
			LocalElev.Floor = newFloor
			StateChan <- LocalElev
			SetFloorLight(newFloor) // oppdatere mtp 1 indeksering **********************code quality************************************************************

			switch LocalElev.State {

			case ConfigFile.INITIALIZE:
				SetMotorDirection(ConfigFile.NEUTRAL)
				LocalElev.State = ConfigFile.IDLE
				LocalElev.Direction = ConfigFile.NEUTRAL
				break

			case ConfigFile.IDLE:
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
					//LocalElev.Direction=ConfigFile.NEUTRAL
					doorTimerChan = time.After(3*time.Second)
					fmt.Printf("Door open\n")
					SetDoorOpenLamp(1)
					LocalElev.State = ConfigFile.DOORSOPEN
					StateChan <- LocalElev
					break
				}

			case ConfigFile.DOORSOPEN:
				break
			}


		case newOrders := <-LocalOrdersChan:

			switch LocalElev.State {

			case ConfigFile.INITIALIZE:
				break

			case ConfigFile.IDLE:
				LocalElev.Orders = newOrders
				if nextDirection(LocalElev) != ConfigFile.NEUTRAL {
					LocalElev.State = ConfigFile.RUNNING
					LocalElev.Direction = nextDirection(LocalElev)
					SetMotorDirection(LocalElev.Direction)
					StateChan <- LocalElev
				}else{
					for button := 0; button < ConfigFile.Num_buttons; button++ {
						if LocalElev.Orders[LocalElev.Floor][button] {
							doorTimerChan = time.After(3*time.Second)
							SetDoorOpenLamp(1)
							LocalElev.State = ConfigFile.DOORSOPEN
							if button < ConfigFile.Num_buttons-1 {
								ClearHallOrdersChan <- [2]int{LocalElev.Floor, button}
							} else {
								ClearCabOrderChan <- LocalElev.Floor
							}
						}
					}

				}
				break

			case ConfigFile.RUNNING:

				LocalElev.Orders = newOrders
				break

			case ConfigFile.DOORSOPEN:
				LocalElev.Orders = newOrders
				// if order at this floor, keep door open longer?
				break
			}

		case <-doorTimerChan:
			fmt.Printf("Door close\n")
			switch LocalElev.State {

			case ConfigFile.INITIALIZE:
				break

			case ConfigFile.IDLE:
				break

			case ConfigFile.RUNNING:
				break

			case ConfigFile.DOORSOPEN:
				SetDoorOpenLamp(0)
				LocalElev.Direction = nextDirection(LocalElev)

				if LocalElev.Direction != ConfigFile.NEUTRAL {
					LocalElev.State = ConfigFile.RUNNING
					SetMotorDirection(LocalElev.Direction)
					StateChan <- LocalElev
				} else {
					LocalElev.State = ConfigFile.IDLE
					StateChan<-LocalElev
				}
				break
			}
		}
	}
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
	floor := LocalElev.Floor+1
	for i := floor; i < ConfigFile.Num_floors; i++ {
		for j := 0; j < ConfigFile.Num_buttons; j++ {
			if LocalElev.Orders[i][j] != false{
				return true
			}
		}
	}
	return false
}

func ordersBelow(LocalElev ConfigFile.Elev) bool {
	floor := LocalElev.Floor-1
	for i := floor; i >= 0; i-- {
		for j := 0; j < ConfigFile.Num_buttons; j++ {
			if LocalElev.Orders[i][j] != false {
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
		if LocalElev.Orders[LocalElev.Floor][0] {
			return true
		} else {
			return (!ordersAbove(LocalElev))
		}
	} else if LocalElev.Direction == ConfigFile.DOWN {
		if LocalElev.Orders[LocalElev.Floor][1] {
			return true
		} else {
			return (!ordersBelow(LocalElev))
		}
	}
	return false
}
