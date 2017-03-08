package Consensus

import (
	"../ConfigFile"
	. "../Network"
	"../driver"
	//	"fmt"
)

func ConsensusHall(ClearHallOrderChan chan [2]int, HallButtonChan chan [2]int, PeerUpdateChan chan ConfigFile.PeerUpdate) {
	hallordersRx := make(chan ConfigFile.ConsensusHall)
	hallordersTx := make(chan ConfigFile.ConsensusHall)
	go Transmitter(ConfigFile.Port, hallordersTx)
	go Receiver(ConfigFile.Port, hallordersRx)
	var localHallConsensus ConfigFile.ConsensusHall

	var RemoteID string
	var LivingPeers []string
	for {
		select {
		case remoteHallConsensus := <-hallordersRx:
			for floor := 0; floor < ConfigFile.Num_floors; floor++ {
				for button := 0; button < ConfigFile.Num_buttons-1; button++ {
					local := localHallConsensus.HallButtons[floor][button]
					remote := remoteHallConsensus.HallButtons[floor][button]
					Merge(&local, remote, RemoteID, LivingPeers, onActive func() { driver.SetButtonLamp(ConfigFile.ButtonType(button), floor, 1) }(), onInactive())
					localHallConsensus.HallButtons[floor][button] = local
				}
			}

		case PeerUpdate := <-PeerUpdateChan:
			LivingPeers = PeerUpdate.Peers
			if (len(LivingPeers) == 1) && (LivingPeers[0] == ConfigFile.LocalID) || (len(LivingPeers) < 1) {
				for floor := 0; floor < ConfigFile.Num_floors; floor++ {
					for button := 0; button < ConfigFile.Num_buttons-1; button++ {
						if localHallConsensus.HallButtons[floor][button].OrderState == ConfigFile.Innactive {
							localHallConsensus.HallButtons[floor][button].OrderState = ConfigFile.Default
						}
					}
				}
			}

		case NewHallButton := <-HallButtonChan:
			if (localHallConsensus.HallButtons[NewHallButton[0]][NewHallButton[1]].OrderState == ConfigFile.Default) || (localHallConsensus.HallButtons[NewHallButton[0]][NewHallButton[1]].OrderState == ConfigFile.Inactive) {
				localHallConsensus.HallButtons[NewHallButton[0]][NewHallButton[1]].OrderState = ConfigFile.PendingAck
			}
		default:
			break
		}
	}
}

func onActive(int button, int floor) {
	driver.SetButtonLamp(button, floor, 1)
	// oppdatere noe?
}

func onInactive() {
	driver.SetButtonLamp(button, floor, 0)
	// oppdatere noe?
}
