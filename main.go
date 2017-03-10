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
)

func main() {

	flag.StringVar(&ConfigFile.LocalID, "id", "", "id of this peer")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	// ** setter opp channels ** \\
	FloorChan := make(chan int)                                        	// Fra driver.go til FSM.go
	HallButtonChan := make(chan [2]int)                                	// fra driver.go til ConsensusHall
	CabButtonChan := make(chan int)                                    	// Fra driver.go til ConsensusCab
	ClearHallOrderChan := make(chan [2]int)                            	// Fra FSM.go til ConsensusHall
	ClearCabOrderChan := make(chan int)                                	// Fra FSM.go til ConsensusCab
	StateChan := make(chan ConfigFile.Elev)                            	// Fra FSM.go til StateFileNotYetMade
	PeerUpdateChan := make(chan ConfigFile.PeerUpdate)                 	// Fra "egen NW modul fra Anders" Til ConsHall/ConsCab/HallReqAss
	ConsensusCabChan := make(chan map[string]*ConfigFile.ConsensusCab) 	// Fra ConsensusCab til HallReqAss
	ConsensusHallChan := make(chan ConfigFile.ConsensusHall)           	// Fra ConsensusHall til HallReqAss
	ElevatorStatesChan := make(chan map[string]*ConfigFile.Elev)       	// Fra ElevatorStates til HallReqAss
	TransmitEnable := make(chan bool)                                  	// Fra Peers til consensusCab/consensusHall/HallReqAss
	LocalOrdersChan := make(chan [][]bool)							   	// Fra HallReqAss til FSM.go
	FromPeersToConsensusHall := make(chan ConfigFile.PeerUpdate)
	FromPeersToConsensusCab := make(chan ConfigFile.PeerUpdate)
	FromPeersToHallReqAss := make(chan ConfigFile.PeerUpdate)

	driver.InitElev(driver.ET_Comedi)

	// ** starte GO routines ** \\
	go driver.ButtonPoll(HallButtonChan, CabButtonChan)
	go driver.FloorPoll(FloorChan)
	go FSM.RUN(FloorChan, StateChan, LocalOrdersChan, ClearHallOrderChan, ClearCabOrderChan)
	go Consensus.ConsensusHall(ClearHallOrderChan, ConsensusHallChan, HallButtonChan, FromPeersToConsensusHall)
	go Consensus.ConsensusCab(ClearCabOrderChan, ConsensusCabChan, CabButtonChan, FromPeersToConsensusCab)
	go ElevatorStates.DoSomethingSmart(StateChan, ElevatorStatesChan)
	go HallRequestAssigner.HallReq(ConsensusHallChan, ConsensusCabChan, ElevatorStatesChan, LocalOrdersChan, FromPeersToHallReqAss)
	go Peers.Transmitter(ConfigFile.PeersPort, ConfigFile.LocalID, TransmitEnable)
	go Peers.Receiver(ConfigFile.PeersPort, PeerUpdateChan)
	go Peers.Repeater(PeerUpdateChan, FromPeersToConsensusHall, FromPeersToConsensusCab, FromPeersToHallReqAss)//{FromPeersToConsensusHall, FromPeersToConsensusCab, FromPeersToHallReqAss})	//3](chan ConfigFile.PeerUpdate){

	println("Main ferdig")
	select {}
}
