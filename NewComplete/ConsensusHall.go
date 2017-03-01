package ConsensusHall

import (
	"../ConfigFile"
	//	"fmt"
)

var localHallConsensus ConfigFile.ConsensusHall

func ConsensusHall(ClearHallOrderChan chan [2]int, HallButtonChan chan [2]int, PeerUpdateChan chan ConfigFile.PeerUpdate) {
	// TIL FRA NETTVERK ********
	hallordersRx := make(chan ConfigFile.ConsensusHall)
	// BEM EM IN SCUTTY
	// hallordersTx := make(chan ConfigFile.AllHallOrders)						// REMEMBER TO USE
	// BEM EM UP SCUTTY
	// ********
	RemoteID = "123"
	LocalID = "321"
	for {
		select {
		case remoteHallConsensus := <-hallordersRx:
			for floor := 0; floor < ConfigFile.Num_floors; floor++ {
				for button := 0; button < ConfigFile.Num_buttons-1; button++ {

					local := localHallConsensus.HallOrders[floor][button]
					remote := remoteHallConsensus.HallButtons[floor][button]
					merge(&local, remote, LocalID, RemoteID) // sette til å oppdatere local også!!!
				}
			}
		default:
			break
		}
	}
}

// --> Ander$ Code <--
/*
func doostuff(onActive func(), onInactive func()) {
	f()
}
*/

func merge(local *ConfigFile.OrderStatus, remote ConfigFile.OrderStatus, LocalID string, RemoteID string) {

	switch local.OrderState {

	case ConfigFile.Default:
		if RemoteID != LocalID {
			switch remote.OrderState {
			case ConfigFile.Default:
				break
			case ConfigFile.Inactive:
				local.OrderState = ConfigFile.Inactive
				local.AckdBy = local.AckdBy[:0] // destroy
				break
			case ConfigFile.PendingAck:
				local.OrderState = ConfigFile.PendingAck
				local.AckdBy = append(local.AckdBy, LocalID)
				break
			case ConfigFile.Active:
				local.OrderState = ConfigFile.Active
				local.AckdBy = append(remote.AckdBy, LocalID)
				//onActive() // usikker på hvordan dette kan oversettes til golang
			}
		}

	case ConfigFile.Inactive:
		switch remote.OrderState {
		case ConfigFile.Inactive:
			break // if state remains the same, do nothing

		case ConfigFile.PendingAck:
			local.OrderState = remote.OrderState
			local.AckdBy = append(remote.AckdBy, LocalID) // adder egen ACK
			break

		case ConfigFile.Active:
			break // cannot skip PendingAck-state

		}

	case ConfigFile.PendingAck:
		switch remote.OrderState {
		case ConfigFile.Inactive:
			break

		case ConfigFile.PendingAck:
			// ADD all others ACKed to local
			if len(local.AckdBy) == 3 { // denne må selvsagt byttes til en "dynamisk" sak, og ikke bare sjekker antall!
				local.OrderState = ConfigFile.Active
				// onActive
			}
			break

		case ConfigFile.Active:
			local.OrderState = ConfigFile.Active
			local.AckdBy = append(local.AckdBy, remote.AckdBy, LocalID)
			//local.AckdBy
			break

		}

	case ConfigFile.Active:
		switch remote.OrderState {
		case ConfigFile.Inactive:
			local.OrderState = remote.OrderState
			local.AckdBy = local.AckdBy[:0] // destroy
			// onInnactive
			break

		case ConfigFile.PendingAck:
			break

		case ConfigFile.Active:
			local.AckdBy = append(local.AckdBy, remote.AckdBy, LocalID)
			break

		}
	}
}
