package FSM

import (
	"../driver"
	"../ConfigFile"
	"fmt"
	"time"
)


func RUN(
	FloorChan chan int, StateChan chan ConfigFile.Elev,
	LocalOrdersChan chan [][]bool,
	ClearHallOrdersChan chan [2]int, ClearCabOrderChan chan int, TransmitEnable chan bool) {

	LocalElev := ConfigFile.NewElev()
	var doorTimerChan <-chan time.Time
	var orderTimerChan <- chan time.Time

	{
		floor := driver.GetFloorSensorSignal()
		if floor == -1 {
			driver.SetMotorDirection(ConfigFile.DOWN)
			LocalElev.State = ConfigFile.INITIALIZE
		}
	}

	for {
		select {
		case newFloor := <-FloorChan:
			fmt.Printf("New floor: %+v\n", newFloor)
			LocalElev.Floor = newFloor
			//fmt.Printf("above StateChan \n ")
			//StateChan <- LocalElev
			//fmt.Printf("below StateChan \n ")
			driver.SetFloorLight(newFloor)

			switch LocalElev.State {

			case ConfigFile.INITIALIZE:
				driver.SetMotorDirection(ConfigFile.NEUTRAL)
				LocalElev.State = ConfigFile.IDLE
				LocalElev.Direction = ConfigFile.NEUTRAL
				StateChan <- LocalElev
				break

			case ConfigFile.IDLE:
				break

			case ConfigFile.MOVING:
				if ordersAbove(LocalElev) || ordersBelow(LocalElev){
					orderTimerChan = time.After(15*time.Second)
				}
				if shouldStop(LocalElev) { // se over, kan ha noen mangler, eks. når heisen allerede står i etg hvor det bestilles
					for button := 0; button < ConfigFile.Num_buttons; button++ {
						if LocalElev.Orders[LocalElev.Floor][button] {
							if button < ConfigFile.Num_buttons-1 {
								fmt.Printf("*FSM above ClearHallOrdersChan NewFLOOR\n")
								ClearHallOrdersChan <- [2]int{LocalElev.Floor, button}
								fmt.Printf("*FSM below ClearHallOrdersChan NewFLOOR\n")
							} else {
								fmt.Printf("*FSM above ClearCabOrderChan NewFLOOR\n")
								ClearCabOrderChan <- LocalElev.Floor
								fmt.Printf("*FSM below ClearCabOrderChan NewFLOOR\n")
							}
						}
					}

					driver.SetMotorDirection(ConfigFile.NEUTRAL)
					doorTimerChan = time.After(3*time.Second)
					fmt.Printf("Door open\n")
					driver.SetDoorOpenLamp(1)
					LocalElev.State = ConfigFile.DOORSOPEN
					fmt.Printf("above StateChan \n ")
					StateChan <- LocalElev
					fmt.Printf("below StateChan \n ")
					break
				}

			case ConfigFile.DOORSOPEN:
				break
			}


		case newOrders := <-LocalOrdersChan:
			fmt.Printf("*FSM got new orders \n")
			switch LocalElev.State {
			case ConfigFile.INITIALIZE:
				break

			case ConfigFile.IDLE:
				fmt.Printf("*FSM  was IDLE \n")
				if hasNewOrders(newOrders, LocalElev){
					orderTimerChan = time.After(15*time.Second)
				}
				LocalElev.Orders = newOrders

				if nextDirection(LocalElev) != ConfigFile.NEUTRAL {
					LocalElev.State = ConfigFile.MOVING
					LocalElev.Direction = nextDirection(LocalElev)
					driver.SetMotorDirection(LocalElev.Direction)
					fmt.Printf("*FSM above StateChan")
					StateChan <- LocalElev
					fmt.Printf("*FSM below StateChan")
				}else{
					for button := 0; button < ConfigFile.Num_buttons; button++ {
						if LocalElev.Orders[LocalElev.Floor][button] {
							doorTimerChan = time.After(3*time.Second)
							driver.SetDoorOpenLamp(1)
							LocalElev.State = ConfigFile.DOORSOPEN
							if button < ConfigFile.Num_buttons-1 {
								fmt.Printf("*FSM above ClearHallOrdersChan")
								ClearHallOrdersChan <- [2]int{LocalElev.Floor, button}
								fmt.Printf("*FSM below ClearHallOrdersChan")
							} else {
								fmt.Printf("*FSM above ClearCabOrderChan")
								ClearCabOrderChan <- LocalElev.Floor
								fmt.Printf("*FSM below ClearCabOrderChan")
							}
						}
					}
				}
				break

			case ConfigFile.MOVING:
				fmt.Printf("*FSM  was MOVING \n")
				if hasNewOrders(newOrders, LocalElev){
					orderTimerChan = time.After(15*time.Second)
				}
				LocalElev.Orders = newOrders
				break

			case ConfigFile.DOORSOPEN:
				fmt.Printf("*FSM  was DOORSOPEN \n")
				if hasNewOrders(newOrders, LocalElev){
					orderTimerChan = time.After(15*time.Second)
				}
				LocalElev.Orders = newOrders
				break
			}

		case <-doorTimerChan:
			fmt.Printf("Door close\n")
			switch LocalElev.State {

			case ConfigFile.INITIALIZE:
				break

			case ConfigFile.IDLE:
				break

			case ConfigFile.MOVING:
				break

			case ConfigFile.DOORSOPEN:
				driver.SetDoorOpenLamp(0)
				LocalElev.Direction = nextDirection(LocalElev)

				if LocalElev.Direction != ConfigFile.NEUTRAL {
					LocalElev.State = ConfigFile.MOVING
					driver.SetMotorDirection(LocalElev.Direction)
					StateChan <- LocalElev
				} else {
					LocalElev.State = ConfigFile.IDLE
					StateChan<-LocalElev
				}
				break
			}

		case <- orderTimerChan:
			if(LocalElev.State != ConfigFile.IDLE){
				TransmitEnable <- false
				driver.SetMotorDirection(ConfigFile.NEUTRAL)
				time.Sleep(20* time.Second)
				TransmitEnable <- true
				driver.SetMotorDirection(LocalElev.Direction)
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
	for f := floor; f < ConfigFile.Num_floors; f++ {
		for b := 0; b < ConfigFile.Num_buttons; b++ {
			if LocalElev.Orders[f][b] != false{
				return true
			}
		}
	}
	return false
}

func ordersBelow(LocalElev ConfigFile.Elev) bool {
	floor := LocalElev.Floor-1
	for f := floor; f >= 0; f-- {
		for b := 0; b < ConfigFile.Num_buttons; b++ {
			if LocalElev.Orders[f][b] != false {
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

func hasNewOrders(newOrders [][]bool, LocalElev ConfigFile.Elev) bool{
	for f := 0; f < ConfigFile.Num_floors; f++{
		for b := 0; b < ConfigFile.Num_buttons; b++{
			if (LocalElev.Orders[f][b] != newOrders[f][b]){
				return true
			}
		}
	}
	return false
}