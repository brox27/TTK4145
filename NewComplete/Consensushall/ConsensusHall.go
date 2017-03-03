package Consensushall

import (
	"../ConfigFile"
	. "../Network"
	//	"fmt"
)

var localHallConsensus ConfigFile.ConsensusHall

func ConsensusHall(ClearHallOrderChan chan [2]int, HallButtonChan chan [2]int, PeerUpdateChan chan ConfigFile.PeerUpdate) {
	hallordersRx := make(chan ConfigFile.ConsensusHall)
	hallordersTx := make(chan ConfigFile.AllHallOrders)
	go Transmitter(ConfigFile.Port, hallordersTx)
	go Receiver(ConfigFile.Port, hallordersRx)

	RemoteID := "123"
	LocalID := "321"
	var LivingPeers []string
	for {
		select {
		case remoteHallConsensus := <-hallordersRx:
			for floor := 0; floor < ConfigFile.Num_floors; floor++ {
				for button := 0; button < ConfigFile.Num_buttons-1; button++ {

					local := localHallConsensus.HallButtons[floor][button]
					remote := remoteHallConsensus.HallButtons[floor][button]
					merge(&local, remote, LocalID, RemoteID, LivingPeers) // sette til å oppdatere local også!!!
				}
			}
		case PeerUpdate := <- PeerUpdateChan:
			LivingPeers = PeerUpdate.Peers
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

func merge(local *ConfigFile.OrderStatus, remote ConfigFile.OrderStatus, LocalID string, RemoteID string, LivingPeers []string ) {

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
				local.AckdBy = append(remote.AckdBy, LocalID)			// bør det være local.AckdBy = append(remote.AckdBy, LocalID) her?
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
			local.AckdBy = append(remote.AckdBy, LocalID) // legger til alle andre..?
			if len(local.AckdBy) >= len(LivingPeers) { // denne må selvsagt byttes til en "dynamisk" sak, og ikke bare sjekker antall!
				local.OrderState = ConfigFile.Active
				// onActive
			}
			break

		case ConfigFile.Active:
			local.OrderState = ConfigFile.Active
			local.AckdBy = append(remote.AckdBy, LocalID)
			// onActive
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
			local.AckdBy = append(remote.AckdBy, LocalID)
			break

		}
	}
	// local.ackdBy = local.ackdBy.sort().uniq.array;  -> hva i huleste faen gjør detta?? -> anders har den her
}

/*
func Activate(local *OrderStatus, LocalID string){
	switch (local.OrderState){
	}
}

func DeActivate(){
	if (peers == LocalID){
		local.OrderState = ConfigFile.Default
	}else{
		local.OrderState = ConfigFile.Inactive
	}
	local.AckdBy = local.AckdBy[:0]	// Destroy!
}
*/

// GLOBAL QUESTION? 			local.AckdBy = append(remote.AckdBy, LocalID) er nå alle remote og egen lagt til i lokal?