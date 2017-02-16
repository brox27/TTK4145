package fsm

import (
	."./driver"
	"time"
	"os"
	)

const (
	INITIALIZE = 0
	IDLE = 1
	RUNNING = 2
	ARRIVEDFLOOR = 3
	SHOULDSTOP = 4
	// trenger vi flere states?
)

type FsmEvent struct {
	Id int
	Dir int
	Floor int
}

type Elevator struct {	
	State int
	Dir int
	Floor int
	// interne ordre?
}

func (e Elevator) Run(to_main FsmEvent, from_main chan int) {
	// sends elev to nearest floor
	e.Dir = DIR_DOWN
	Set_motor_direction(DIR_DOWN)
	if Get_floor_sensor_signal() != 3 {
		Set_motor_direction(DIR_UP)
		e.Dir = DIR_UP
	}
	for Get_floor_sensor_signal() != -1 {
		time.Sleep(15 * time.Millisecond)
	}
	Set_motor_direction(- e.Dir)

	for Get_floor_sensor_signal() == -1 {
		time.Sleep(10 * time.Millisecond)
	}
	Set_motor_direction(DIR_STOP)

	to_main <- FsmEvent(Id: INITIALIZE)

	e.Floor = Get_floor_sensor_signal()
	Set_floor_indicator(Get_floor_sensor_signal())

	to_main <- FsmEvent{Id: ARRIVEDFLOOR, Dir: e.Dir, Floor: e.Floor}

	if SHOULDSTOP := <- from_main; SHOULDSTOP == 0 || SHOULDSTOP != 0 {
		e.State = IDLE
		Set_motor_direction(DIR_STOP)
	}

	for {
		switch(e.State) {
		case IDLE:
			to_main <- FsmEvent{Id: IDLE, Dir: e.Dir, Floor: e.Floor}

			e.Dir <- from_main
			Set_motor_direction(e.Dir)
			e.State = RUNNING
			ch.TimerReset <- true

		case RUNNING:
			if ArrivedFloor := Get_floor_sensor_signal(); ArrivedFloor != -1 {
				e.Floor = ArrivedFloor
				Set_floor_indicator(ArrivedFloor)

				to_main <- FsmEvent{Id: ARRIVEDFLOOR, Dir: e.Dir, Floor: e.Floor}

				if <- from_main == SHOULDSTOP {
					ch.TimerReset <- true
					e.State = IDLE
					Set_motor_direction(DIR_STOP)
					Set_door_open_light(true)
					time.Sleep(3 * time.Second)
					Set_door_open_light(false)
				}
			}
			time.Sleep(time.Millisecond)
		}
	}
}	