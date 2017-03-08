package main

import (
	"./ConfigFile"
	"./Consensuscab"
	"./Consensushall"
	"./ElevatorStates"
	"./FSM"
	"./HallRequestAssigner"
	"./Peers"
	"./driver"
	"flag"
	"runtime"
)

func main() {

	flag.StringVar(&ConfigFile.LocalID, "id", "", "id of this peer")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())
	// ** setter opp channels ** \\
	FloorChan := make(chan int)                                        // Fra driver.go til FSM.go
	HallButtonChan := make(chan [2]int)                                // fra driver.go til ConsensusHall
	CabButtonChan := make(chan int)                                    // Fra driver.go til ConsensusCab
	ClearHallOrderChan := make(chan [2]int)                            // Fra FSM.go til ConsensusHall
	ClearCabOrderChan := make(chan int)                                // Fra FSM.go til ConsensusCab
	StateChan := make(chan ConfigFile.Elev)                            // Fra FSM.go til StateFileNotYetMade
	PeerUpdateChan := make(chan ConfigFile.PeerUpdate)                 // Fra "egen NW modul fra Anders" Til ConsHall/ConsCab/HallReqAss
	ConsensusCabChan := make(chan map[string]*ConfigFile.ConsensusCab) // Fra ConsensusCab til HallReqAss
	ConsensusHallChan := make(chan ConfigFile.ConsensusHall)           // Fra ConsensusHall til HallReqAss
	ElevatorStatesChan := make(chan map[string]*ConfigFile.Elev)       // Fra ElevatorStates til HallReqAss
	TransmitEnable := make(chan bool)                                  // Fra Peers til consensusCab/consensusHaøll/HallReqAss
	//NewOrderConsensusChan := make(chan ConfigFile.OrderMsg)			// not sure...

	LocalOrdersChan := make(chan [ConfigFile.Num_floors][ConfigFile.Num_buttons]bool) // Fra HallReqAss til FSM.go

	driver.InitElev()

	// ** starte GO routines ** \\
	go driver.ButtonPoll(HallButtonChan, CabButtonChan)
	go driver.FloorPoll(FloorChan)
	go FSM.RUN(FloorChan, StateChan, LocalOrdersChan, ClearHallOrderChan, ClearCabOrderChan)
	go Consensushall.ConsensusHall(ClearHallOrderChan, HallButtonChan, PeerUpdateChan)
	go Consensuscab.ConsensusCab(ClearCabOrderChan, ConsensusCabChan, CabButtonChan, PeerUpdateChan)
	go ElevatorStates.DoSomethingSmart(StateChan, ElevatorStatesChan)
	go HallRequestAssigner.HallReq(ConsensusHallChan, ConsensusCabChan, ElevatorStatesChan, LocalOrdersChan)
	go Peers.Transmitter(ConfigFile.Port, "123", TransmitEnable)
	go Peers.Receiver(ConfigFile.Port, PeerUpdateChan)

	println("Hallooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo1")
	for {

	}
}
