package main

import (
	. "./driver"
	. "./globals"
	. "fmt"
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
	state  States
	dir    Direction
	floor  int
	orders [][]bool
}

func main() {
	InitElev()

	eventChan := make(chan map[Event]interface{})
	go EventHandler(eventChan)

	for {
		Println(<-eventChan)

	}

	defer SetMotorDirection(NEUTRAL)
}
