package main

import (
	. "../driver"
	. "../globals"
	. "fmt"
)

type State int

const (
	INITIALIZE State = iota
	IDLE
	RUNNING
	ARRIVEDFLOOR
	SHOULDSTOP
)

type Elevator struct {
	dir    Direction
	floor  int
	orders [][]bool
}

func (e *Elevator) ElevRun() {
	eventChan := make(chan map[Event]interface{})

	for {
		eventMap := <-eventChan
		for i, j := range eventMap {
			Println(i, j)
		}
	}
}
