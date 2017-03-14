package main

import (
	"./ConfigFile"
	"./Consensus"
	"./Elevator"
	"./ElevatorStates"
	"./OrderAssigner"
	"./Peers"
	"./driver"
	"flag"
	"os"
	"os/signal"
	"runtime"
)

func main() {

	flag.StringVar(&ConfigFile.LocalID, "id", "", "id of this peer")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	FloorChan := make(chan int)
	HallButtonChan := make(chan [2]int)
	CabButtonChan := make(chan int)
	ClearHallOrderChan := make(chan [2]int, 3)
	ClearCabOrderChan := make(chan int, 3)
	StateChan := make(chan ConfigFile.Elev, 3)
	ConsensusCabChan := make(chan map[string]*ConfigFile.ConsensusCab, 3)
	ConsensusHallChan := make(chan ConfigFile.ConsensusHall, 3)
	ElevatorStatesChan := make(chan ConfigFile.AllStates, 3)
	LocalOrdersChan := make(chan [][]bool)
	TransmitEnableChan := make(chan bool)
	PeerUpdateChan := make(chan ConfigFile.PeerUpdate)
	PeersToConsensusHallChan := make(chan ConfigFile.PeerUpdate)
	PeersToConsensusCabChan := make(chan ConfigFile.PeerUpdate)
	PeersToOrderAssignerChan := make(chan ConfigFile.PeerUpdate)

	driver.InitElev()

	go driver.ButtonPoll(HallButtonChan, CabButtonChan)
	go driver.FloorPoll(FloorChan)
	go Elevator.ElevatorController(FloorChan, StateChan, LocalOrdersChan, ClearHallOrderChan, ClearCabOrderChan, TransmitEnableChan)
	go Consensus.ConsensusHall(ClearHallOrderChan, ConsensusHallChan, HallButtonChan, PeersToConsensusHallChan)
	go Consensus.ConsensusCab(ClearCabOrderChan, ConsensusCabChan, CabButtonChan, PeersToConsensusCabChan)
	go ElevatorStates.ElevatorStatesCoordinator(StateChan, ElevatorStatesChan)
	go orderAssigner.OrderAssigner(ConsensusHallChan, ConsensusCabChan, ElevatorStatesChan, LocalOrdersChan, PeersToOrderAssignerChan)
	go Peers.Transmitter(ConfigFile.PeersPort, ConfigFile.LocalID, TransmitEnableChan)
	go Peers.Receiver(ConfigFile.PeersPort, PeerUpdateChan)
	go Peers.Repeater(PeerUpdateChan, PeersToConsensusHallChan, PeersToConsensusCabChan, PeersToOrderAssignerChan)

	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, os.Interrupt)

	println("System is running, all goroutines are active")
	select {
	case <-osSignal:
		println("Interrupt detected, shutting down motor")
		driver.SetMotorDirection(ConfigFile.NEUTRAL)
	}
}
