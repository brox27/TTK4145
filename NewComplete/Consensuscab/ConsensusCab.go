package Consensuscab

import (
	"../ConfigFile"
	"fmt"
)

/*
func main() {
	fmt.Printf("we programmers now? ¯\\_(ツ)_/¯ \n")
	ConfigFile.AllCabOrders = make(map[string]*ConfigFile.CabOrders)
	switchen()
}
*/

func ConsensusCab() {
	// DEBUG SHIT!
	/*
		newOrderConsensus := make(chan ConfigFile.OrderMsg, 9)      // denne går TIL (FSM?) og bør lagres der og tatt inn som argument
		completeOrderConsensus := make(chan ConfigFile.OrderMsg, 9) // denne går TIL (FSM?) og bør lagres der og tatt inn som argument
		thisElev := ConfigFile.CabOrders{}
		ConfigFile.AllCabOrders[LocalID] = &thisElev
		newElev := ConfigFile.CabOrders{}
		//END DEBUG SHIT
		ConfigFile.AllCabOrders[LocalID] = &newElev
		ConfigFile.AllCabOrders[LocalID].Floor = 3
	*/
	cabOrdersRx := make(chan map[string]*ConfigFile.CabOrders)
	// BEM EM IN SCUTTY
	// hallordersTx := make(chan ConfigFile.AllHallOrders)						// REMEMBER TO USE
	// BEM EM UP SCUTTY

	var AllCabOrders map[string]*ConfigFile.ConsensusCab
	thisCab := ConfigFile.ConsensusCab
	allCabOrders[LocalID] = &thisCab
	LocalID := "123"
	RemoteID := "321"
	for {
		select {
		case remoteCabConsensus := <-cabOrdersRx:
			for elevator := range remoteCabConsensus { // kan det bli "kræsj" med lengden på "local" og "remote"...?
				for floor := 0; floor < ConfigFile.Num_floors; floor++ {

					remote := remoteCabConsensus[elevator].ConsensusCab[floor]
					local := AllCabOrders[elevator].ConsensusCab[floor]
					merge(&local, remote, LocalID, RemoteID)

				}
			}
		default:
			break
		}
	}
}

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
				local.AckdBy = append(local.AckdBy, LocalID)
				//onActive() // usikker på hvordan dette kan oversettes til golang
				varible := ConfigFile.OrderMsg{}
				varible.Floor = floor
				newOrderConsensus <- varible
			}
		}

	case ConfigFile.Inactive:
		switch remote.OrderState {
		case ConfigFile.Inactive:
			break // if state remains the same, do nothing

		case ConfigFile.PendingAck:
			local.OrderState = remote.OrderState
			local.AckdBy = append(local.AckdBy, LocalID) //sett inn IP/ID
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
				varible := ConfigFile.OrderMsg{}
				varible.Floor = floor
				newOrderConsensus <- varible

			}
			break
		case ConfigFile.Active:
			local.OrderState = ConfigFile.Active
			//local.AckdBy
			break

		}

	case ConfigFile.Active:
		switch remote.OrderState {
		case ConfigFile.Inactive:
			local.OrderState = remote.OrderState
			local.AckdBy = local.AckdBy[:0] // destroy
			varible := ConfigFile.OrderMsg{}
			varible.Floor = floor
			completeOrderConsensus <- varible
			break

		case ConfigFile.PendingAck:
			break

		case ConfigFile.Active:
			// Anders har appendet IP til ACked by her også... why? eller ser Why, men er det nødvendig?
			break

		}
	}

}
