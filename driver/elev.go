package driver

/*
#cgo CFLAGS: -std=gnu11
#cgo LDFLAGS: -lpthread -lcomedi -lm
#include "elev.h"
*/
import "C"

import (
	. "../ConfigFile"
)

type ElevatorType int
const (
    ET_Comedi ElevatorType  = 0
    ET_Simulation           = 1
)

func InitElev(elevatorType ElevatorType) {
    C.elev_init(C.elev_type(elevatorType))

	for floor := 0; floor < Num_floors; floor++ {
		for button := BUTTON_ORDER_UP; button < Num_buttons; button++ {		// debug button=0
			SetButtonLamp(button, floor, 0)
		}
	}

	SetStopLamp(0)
	SetDoorOpenLamp(0)
	SetFloorLight(0)
	//fmt.Printf("init done \n")
}

func SetMotorDirection(dir Direction) {
	if dir == NEUTRAL {
        C.elev_set_motor_direction(C.elev_motor_direction_t(0))
	} else if dir == UP {
        C.elev_set_motor_direction(C.elev_motor_direction_t(1))
	} else if dir == DOWN {
        C.elev_set_motor_direction(C.elev_motor_direction_t(-1))
    }
}

func SetButtonLamp(button ButtonType, floor int, value int) {
    C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func SetFloorLight(floor int) {
    C.elev_set_floor_indicator(C.int(floor))
}

func SetDoorOpenLamp(value int) {
    C.elev_set_door_open_lamp(C.int(value))
}

func GetButtonSignal(floor int, button int) int {
    return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func GetFloorSensorSignal() int {
    return int(C.elev_get_floor_sensor_signal())
}

func SetStopLamp(value int) {
    C.elev_set_stop_lamp(C.int(value))
}
