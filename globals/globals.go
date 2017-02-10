package globals
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "channels.h"
*/

import(
	"C"
	//. "time"
)

const N_ELEVATORS = 3
const N_FLOORS = 4
const N_BUTTONS = 3

const MOTOR_SPEED = 2800

type Direction int
const (
	UP Direction = iota
	DOWN
	IDLE
)

type ButtonType int
const(
	Button_Order_Up ButtonType = iota
	Button_Order_Down
	Button_Order_Command
)


type Elevator struct {
	dir Direction
	floor int
	orders [][]bool
}