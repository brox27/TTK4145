package driver

import (
	"fmt"
	"time"
)

type state int

const (
	initialize state = iota
	idle
	running
	doorOpen
	isDead
)

// legg inn eventhandler her?
// arrivedFloor := make(chan Order, 1)
//	prevOrder := make(chan Order, 1)

// legg til funksjon for å intialize heisen

func elevRun(arrivedFloor chan Order, prevOrder chan Order, delOrder chan Order, ..., ...,) { // legg til alle andre channels som er nødvendig
	Elev_set_speed(300 * nextOrder.Dir) // sender heisen i retning av den øverste ordren i ordrelisten dens
	state <- running
	var flag bool
	for {
		time.Sleep(10 * time.Milliesecond)
		prevFloor := Elev_get_floor_sensor_signal()
		if prevFloor != -1 {
			var current Order
			current.Floor = prevFloor
			current.Dir = nextOrder.dir // husk å endre navn til hva enn den er kalt i OM
			get_prev_floor_c <- current
			Elev_set_floor_indicator(prevFloor)
		}
		if prevFloor == nextOrder.Floor {
			Elevator_break(nextOrder.Dir)
			arrivedFloor <- nextOrder
			return
		}
		if Elev_get_stop_signal() {
			Elevator_break(nextOrder.Dir)
			stop <- true
			return
		}
	}
}

func newOrder(ch def.channels) {
	switch localElev.state {
	case def.running:
	case: def.idle
		// blabla hvordan vi handler dette i ordermanager?
	case def.doorOpen:
		if queue.shouldStop(def.localElev.floor, def.localElev.dir) {
			ch.timerReset <- true
		}
	}
}

func arrivedFloor(ch def.channels, newFloor int) {
	def.localElev.Floor = newFloor
	ch.newFloor <- def.localElev.floor
	switch def.localElev.state{
	case def.RUNNING:
		if orderManagerg.shouldStop(def.localElev.floor, def.localElev.dir) {
			def.localElev.state = def.doorOpen
			Elev_set_door_open_lamp(1)
			// ha all funksjonaliteten som skal kalles i doorOpen-staten i shouldStop?
		}
		else {
			ch.motorDir <- def.localElev.dir 
			def.localElev.state = def.running
		}
	case def.doorOpen:
		if orderManager.shouldStop(def.localElev.floor, def.localElev.dir) {
			ch.timerReset <- true //hvis dør har åpnet seg, reset timer på at elev er alive
		}
	default:
	}
}

func elevDoorOpen(nextOrder Order, delOrder chan Order, state chan state) {
	if Elev_get_floor_sensor_signal() != -1 {
		Elev_set_door_open_lamp(1)
		state <- doorOpen
		time.Sleep(3 * time.Second)
		// heller ha en egen doorTimer funksjon?
		nextDirection()
	}
	state <- idle	// ha en egen WAIT state?
	Elev_set_door_open_lamp(0)
	Elev_set_button_lamp(BUTTON_COMMAND, nextOrder.Floor, 0)
	delOrder <- nextOrder
}

func elevClearAllLights() {
	Elev_set_door_open_lamp(0)
	for i := 0; i < N_FLOORS; i++ {
		Elev_set_button_lamp(BUTTON_COMMAND, i, 0)
		if i > 0 {
			Elev_set_button_lamp(BUTTON_CALL_DOWN, i, 0)
		}
		if i < N_FLOORS-1 {
			Elev_set_button_lamp(BUTTON_CALL_UP, i, 0)
		}
	}
}