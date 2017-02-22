package fsm

import (
	."../driver"
	"time"
	"os"
	"fmt"
)

const (

	Init int iota
	Idle
	Running 
	ShouldStop
	DoorOpen
	// IsDead // Foreløpig ikke med pga usikker på implementering

	// FSM-EVENTS
	FSM_NewOrder
	FSM_NextDir
	FSM_ArrivedFloor
	FSM_ShouldStop
)
