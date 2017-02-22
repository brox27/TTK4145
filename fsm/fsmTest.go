package fsm

import (
	. "../driver"
	. "../globals"
	. "fmt"
	"runtime"
	//"time"
)

type States int

const (
	INITIALIZE States = iota
	IDLE
	RUNNING
	ARRIVEDFLOOR
	SHOULDSTOP
)

type Elevator struct {
	State States
	LastFloor int
	Direction Direction
	Orders    [N_FLOORS][N_BUTTONS]int
	id 			string
}

var ALL_ELEVATORS map[string]*Elevator

func (elev *Elevator) RUN() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	elev.State=INITIALIZE
	InitElev()

	eventChan := make(chan map[Event]interface{})


	go EventHandler(eventChan)
	eventMap:= <-eventChan

	for {
		switch(elev.State){

		case INITIALIZE:

			for{
				if(GetFloorSignal()!=-1){
					elev.State=IDLE
					elev.dir=IDLE
					break
				}
			}

		case IDLE:

			if eventMap[BUTTON_PRESSED]!=nil{
				j:=0
				data := eventMap[BUTTON_PRESSED].([2]int)

				for i := range ELEVATOR_IPS{
					if(ALL_ORDERS[i].Orders[data[1]][data[0]]){
						break
					}
					j++
				}

				if j==3{
					order := NewOrder{}
					order.Button=data[0]
					order.Floor=data[1]

					go Sender(order)

					CalculateCost(, elev.Id)
					elev.State=RUNNING
				} else {
					break
				}
			}

		case RUNNING:

		}
	}

	defer SetMotorDirection(NEUTRAL)
}
