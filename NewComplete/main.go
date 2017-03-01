package main

import (
	"./ConfigFile"
	"runtime"
	//	"./ConsensusCab"
	"./Consensushall"
	"./FSM"
	"./driver"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// sette opp channels
	FloorChan := make(chan int)
	HallButtonChan := make(chan [2]int)
	CabButtonChan := make(chan int)
	ClearHallOrderChan := make(chan [2]int)
	//	ClearCabOrderChan := make(chan int)
	StateChan := make(chan ConfigFile.Elev)
	PeerUpdateChan := make(chan ConfigFile.PeerUpdate)
	//	ButtonLightsChan := make(chan [2]int) // evt. struct -> se hva som blir "lettest" å lese
	//	ConsensusCabChan := make(chan map[string]*ConfigFile.ConsensusCab)
	//	ConsensusHallChan := make(chan ConfigFile.ConsensusHall)
	//	ElevatorStatesChan := make(chan map[string]*ConfigFile.Elev)
	LocalOrdersChan := make(chan [ConfigFile.Num_floors][ConfigFile.Num_buttons]bool)
	// CHANNEL TIL HALL REQ. ASS

	// ** starte GO routines ** \\
	go driver.ButtonPoll(HallButtonChan, CabButtonChan)
	go driver.FloorPoll(FloorChan)
	go FSM.RUN(FloorChan, StateChan, LocalOrdersChan)
	go Consensushall.ConsensusHall(ClearHallOrderChan, HallButtonChan, PeerUpdateChan) // Chan til her
}
