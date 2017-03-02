package main

import (
	"./ConfigFile"
	"runtime"
	"./Consensuscab"
	"./Consensushall"
	"./FSM"
	"./driver"
	"./ElevatorStates"
//	"./Peers"
	"./HallRequestAssigner"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// sette opp channels
	FloorChan := make(chan int)												// Fra driver.go til FSM.go
	HallButtonChan := make(chan [2]int)										// fra driver.go til ConsensusHall
	CabButtonChan := make(chan int)											// Fra driver.go til ConsensusCab
	ClearHallOrderChan := make(chan [2]int)									// Fra FSM.go til ConsensusHall
	ClearCabOrderChan := make(chan int)										// Fra FSM.go til ConsensusCab
	StateChan := make(chan ConfigFile.Elev)									// Fra FSM.go til StateFileNotYetMade
	PeerUpdateChan := make(chan ConfigFile.PeerUpdate)						// Fra "egen NW modul fra Anders" Til ConsHall/ConsCab/HallReqAss
	//	ButtonLightsChan := make(chan [2]int) // evt. struct -> se hva som blir "lettest" å lettest  FUNKSJONSKALL!!!!!
	ConsensusCabChan := make(chan map[string]*ConfigFile.ConsensusCab)		// Fra ConsensusCab til HallReqAss
	ConsensusHallChan := make(chan ConfigFile.ConsensusHall)			// Fra ConsensusHall til HallReqAss
	ElevatorStatesChan := make(chan map[string]*ConfigFile.Elev)		// Fra ElevatorStates til HallReqAss
	//NewOrderConsensusChan := make(chan ConfigFile.OrderMsg)				// Fnot sure...

	LocalOrdersChan := make(chan [ConfigFile.Num_floors][ConfigFile.Num_buttons]bool)	// Fra HallReqAss til FSM.go
	// CHANNEL TIL HALL REQ. ASS

	// ** starte GO routines ** \\
	go driver.ButtonPoll(HallButtonChan, CabButtonChan)
	go driver.FloorPoll(FloorChan)
	go FSM.RUN(FloorChan, StateChan, LocalOrdersChan)
	go Consensushall.ConsensusHall(ClearHallOrderChan, HallButtonChan, PeerUpdateChan)
	go Consensuscab.ConsensusCab(ClearCabOrderChan, ConsensusCabChan, CabButtonChan, PeerUpdateChan)
	go ElevatorStates.DoSomethingSmart(StateChan, ElevatorStatesChan)
	go HallRequestAssigner.HallReq(ConsensusHallChan, ConsensusCabChan, ElevatorStatesChan, LocalOrdersChan)
//	go Peers.Receiver(ConfigFile.Port, PeerUpdateChan)






	println("Hallooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo1")
	for{

	}
}
