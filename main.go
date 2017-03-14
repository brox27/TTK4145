package main

import (
	"./ConfigFile"
	"./Consensus"
	"./ElevatorStates"
	"./FSM"
	"./HallRequestAssigner"
	"./Peers"
	"./driver"
	"flag"
	"runtime"
	"os"
	"os/signal"
)

func main() {

	flag.StringVar(&ConfigFile.LocalID, "id", "", "id of this peer")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	// ** Setup of channels ** \\
	FloorChan := make(chan int)                                        		// Fra driver.go til FSM.go
	HallButtonChan := make(chan [2]int)                                		// fra driver.go til ConsensusHall
	CabButtonChan := make(chan int)                                    		// Fra driver.go til ConsensusCab
	ClearHallOrderChan := make(chan [2]int, 3)                            	// Fra FSM.go til ConsensusHall
	ClearCabOrderChan := make(chan int, 3)                                	// Fra FSM.go til ConsensusCab
	StateChan := make(chan ConfigFile.Elev,3)                           	// Fra FSM.go til ElevatorStatesCoord
	PeerUpdateChan := make(chan ConfigFile.PeerUpdate)                 		// Fra Peers til Repeater
	ConsensusCabChan := make(chan map[string]*ConfigFile.ConsensusCab, 3) 	// Fra ConsensusCab til HallReqAss
	ConsensusHallChan := make(chan ConfigFile.ConsensusHall, 3)           	// Fra ConsensusHall til HallReqAss
	ElevatorStatesChan := make(chan map[string]*ConfigFile.Elev,3)      	// Fra ElevatorStates til HallReqAss
	TransmitEnable := make(chan bool)                                  		// Fra FSM til peers
	LocalOrdersChan := make(chan [][]bool)							   		// Fra HallReqAss til FSM.go
	FromPeersToConsensusHall := make(chan ConfigFile.PeerUpdate)			// Fra Repeater til ConsensusHall
	FromPeersToConsensusCab := make(chan ConfigFile.PeerUpdate)				// Fra Repeater til ConsensusCab
	FromPeersToHallReqAss := make(chan ConfigFile.PeerUpdate)				// Fra Repeater til HallReqAss

	driver.InitElev()

	// ** start GO routines ** \\
	go driver.ButtonPoll(HallButtonChan, CabButtonChan)
	go driver.FloorPoll(FloorChan)
	go FSM.RUN(FloorChan, StateChan, LocalOrdersChan, ClearHallOrderChan, ClearCabOrderChan, TransmitEnable)
	go Consensus.ConsensusHall(ClearHallOrderChan, ConsensusHallChan, HallButtonChan, FromPeersToConsensusHall)
	go Consensus.ConsensusCab(ClearCabOrderChan, ConsensusCabChan, CabButtonChan, FromPeersToConsensusCab)
	go ElevatorStates.ElevatorStatesCoordinator(StateChan, ElevatorStatesChan)
	go hallRequestAssigner.HallRequestAssigner(ConsensusHallChan, ConsensusCabChan, ElevatorStatesChan, LocalOrdersChan, FromPeersToHallReqAss)
	go Peers.Transmitter(ConfigFile.PeersPort, ConfigFile.LocalID, TransmitEnable)
	go Peers.Receiver(ConfigFile.PeersPort, PeerUpdateChan)
	go Peers.Repeater(PeerUpdateChan, FromPeersToConsensusHall, FromPeersToConsensusCab, FromPeersToHallReqAss)

	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, os.Interrupt)

	println("System is running, all goroutines are active")
	select {
		case <- osSignal:
			println("Interrupt detected, shutting down motor")
			println("Interrupt detected, shutting down motor")
			println("Interrupt detected, shutting down motor")
			println("Interrupt detected, shutting down motor")
			println("Interrupt detected, shutting down motor")
			println("Interrupt detected, shutting down motor")
			println("In****************************************************************************************************************************************************or")
			println("Interrupt detected, shutting down motor")
			println("Interrupt detected, shutting down motor")
			println("Interrupt detected, shutting down motor")
			println("Interrupt detected, shutting down motor")
			println("Interrupt detected, shutting down motor")
			println("Interrupt detected, shutting down motor")
			println("Interrupt detected, shutting down motor")
			driver.SetMotorDirection(ConfigFile.NEUTRAL)
	}
}
