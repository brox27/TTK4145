package globals

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "channels.h"
*/

import (
)

const N_ELEVATORS = 3
const N_FLOORS = 4
const N_BUTTONS = 3
const MOTOR_SPEED = 2800

type Direction int

const (
	UP Direction = iota
	DOWN
	NEUTRAL
)

type ButtonType int

const (
	BUTTON_ORDER_UP ButtonType = iota
	BUTTON_ORDER_DOWN
	BUTTON_ORDER_COMMAND
)

type NewOrder struct {
	MsgId  int
	Floor  int
	Button int
}

type CompleteOrder struct {
	MsgId  int
	Floor  int
	Button int
}



var ELEVATOR_IPS=[N_ELEVATORS]string{"123.123.123","321.321.321","asd.asd.asd"}