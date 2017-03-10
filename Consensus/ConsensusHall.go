package Consensus

import (
	"../ConfigFile"
	. "../Network"
	"../driver"
	"time"
//	"fmt"
)

func ConsensusHall(ClearHallOrderChan chan [2]int, ConsensusHallChan chan ConfigFile.ConsensusHall,  HallButtonChan chan [2]int, PeerUpdateChan chan ConfigFile.PeerUpdate) {
	hallordersRx := make(chan ConfigFile.ConsensusHall)
	hallordersTx := make(chan ConfigFile.ConsensusHall)
	go Transmitter(ConfigFile.HallConesnsusPort, hallordersTx)
	go Receiver(ConfigFile.HallConesnsusPort, hallordersRx)
	transmittTimer := time.NewTicker(time.Millisecond * 50).C

	var LivingPeers []string

	localHallConsensus := ConfigFile.ConsensusHall{}
	for floor := 0; floor < ConfigFile.Num_floors; floor++ {
		for button := 0; button < ConfigFile.Num_buttons-1; button ++{
			localHallConsensus.HallButtons[floor][button].OrderState = ConfigFile.Default
		}
	}

	for {
		select {
		case remoteHallConsensus := <-hallordersRx:
			RemoteID := remoteHallConsensus.ID
			for floor := 0; floor < ConfigFile.Num_floors; floor++ {
				for button := 0; button < ConfigFile.Num_buttons-1; button++ {

					remote := remoteHallConsensus.HallButtons[floor][button]
					Merge(&localHallConsensus.HallButtons[floor][button], remote, RemoteID, LivingPeers, 
						func() {driver.SetButtonLamp(ConfigFile.ButtonType(button), floor, 1)}, 
						func() {driver.SetButtonLamp(ConfigFile.ButtonType(button), floor, 0)})
				}
			}
			ConsensusHallChan <- localHallConsensus

		case ClearedHallOrder := <- ClearHallOrderChan:
			Deactivate(&localHallConsensus.HallButtons[ClearedHallOrder[0]][ClearedHallOrder[1]], LivingPeers)
			driver.SetButtonLamp(ConfigFile.ButtonType(ClearedHallOrder[1]), ClearedHallOrder[0], 0)
		
		case NewHallButton := <-HallButtonChan:
			Activate(&localHallConsensus.HallButtons[NewHallButton[0]][NewHallButton[1]])

		case <- transmittTimer:
            hallordersTx <- localHallConsensus

        case PeerUpdate := <-PeerUpdateChan:
			LivingPeers = PeerUpdate.Peers
		}
	}
}
