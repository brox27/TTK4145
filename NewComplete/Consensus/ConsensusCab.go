package Consensuscab

import (
	"../ConfigFile"
	. "../Network"
	. "fmt"
	"time"
)

/*
func main() {
	fmt.Printf("we programmers now? ¯\\_(ツ)_/¯ \n")
	ConfigFile.AllCabOrders = make(map[string]*ConfigFile.CabOrders)
	switchen()
}
*/

func ConsensusCab(ClearCabOrderChan chan int, ConsensusCabChan chan map[string]*ConfigFile.ConsensusCab, CabButtonChan chan int, PeerUpdateChan chan ConfigFile.PeerUpdate) {

	RemoteID := "321"
	cabOrdersRx := make(chan map[string]*ConfigFile.ConsensusCab)
	cabOrdersTx := make(chan map[string]*ConfigFile.ConsensusCab)
	go Transmitter(ConfigFile.Port, cabOrdersTx)
	go Receiver(ConfigFile.Port, cabOrdersRx)
	var LivingPeers []string

	timerChan := make(chan int)

	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			timerChan <- 1
		}
	}()

	AllCabOrders := make(map[string]*ConfigFile.ConsensusCab)
	thisCab := ConfigFile.ConsensusCab{}
	AllCabOrders[ConfigFile.LocalID] = &thisCab
	for floor := 0; floor < ConfigFile.Num_floors; floor++ {
		thisCab.CabButtons[floor].OrderState = ConfigFile.Default
	}

	for {
		select {
		case remoteCabConsensus := <-cabOrdersRx:
			for elevID := range remoteCabConsensus { // kan det bli "kræsj" med lengden på "local" og "remote"...?
				for floor := 0; floor < ConfigFile.Num_floors; floor++ {

					remote := remoteCabConsensus[elevID].CabButtons[floor]
					local := AllCabOrders[elevID].CabButtons[floor]
					Println("hallo")
					merge(&local, remote, RemoteID, floor, LivingPeers) // usikker på om denne trenger å vite Floor

				}
			}
			// burde sende her?
			ConsensusCabChan <- AllCabOrders
		case NewCabButton := <-CabButtonChan:
			if (AllCabOrders[ConfigFile.LocalID].CabButtons[NewCabButton].OrderState == ConfigFile.Default) || (AllCabOrders[ConfigFile.LocalID].CabButtons[NewCabButton].OrderState == ConfigFile.Inactive) {
				AllCabOrders[ConfigFile.LocalID].CabButtons[NewCabButton].OrderState = ConfigFile.PendingAck
			}
		case timeOut := <-timerChan:
			_ = timeOut
			cabOrdersTx <- AllCabOrders
		default:
			break
		}
	}
}

func merge(local *ConfigFile.OrderStatus, remote ConfigFile.OrderStatus, RemoteID string, floor int, LivingPeers []string) {
	switch local.OrderState {
	case ConfigFile.Default:
		if RemoteID != ConfigFile.LocalID {
			switch remote.OrderState {
			case ConfigFile.Default:
				break
			case ConfigFile.Inactive:
				local.OrderState = ConfigFile.Inactive
				local.AckdBy = local.AckdBy[:0] // destroy
				break
			case ConfigFile.PendingAck:
				local.OrderState = ConfigFile.PendingAck
				local.AckdBy = append(local.AckdBy, ConfigFile.LocalID) // bør det være local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID) her?
				break
			case ConfigFile.Active:
				local.OrderState = ConfigFile.Active
				local.AckdBy = append(local.AckdBy, ConfigFile.LocalID) // bør det være local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID) her?
				//onActive() // usikker på hvordan dette kan oversettes til golang ** onActive -> i Hall
			}
		}

	case ConfigFile.Inactive:
		switch remote.OrderState {
		case ConfigFile.Inactive:
			break // if state remains the same, do nothing

		case ConfigFile.PendingAck:
			local.OrderState = remote.OrderState
			local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID) //sett inn IP/ID
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
			if len(local.AckdBy) >= len(LivingPeers) { // denne må selvsagt byttes til en "dynamisk" sak, og ikke bare sjekker antall!
				local.OrderState = ConfigFile.Active
				// onACTIVE!!
			}
			break
		case ConfigFile.Active:
			local.OrderState = ConfigFile.Active
			local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID)
			break

		}

	case ConfigFile.Active:
		switch remote.OrderState {
		case ConfigFile.Inactive:
			local.OrderState = remote.OrderState
			local.AckdBy = local.AckdBy[:0] // destroy
			break

		case ConfigFile.PendingAck:
			break

		case ConfigFile.Active:
			local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID)
			break

		}
	}

}
