package main

import (
	"fmt"
	"./ConfigFile"
)



func main() {
	fmt.Printf("we programmers now? ¯\\_(ツ)_/¯ \n")
	ConfigFile.AllCabOrders = make(map[string]*ConfigFile.CabOrders)
	switchen()
}



func switchen() {
	// DEBUG SHIT!
	newOrderConsensus := make(chan ConfigFile.OrderMsg, 9)				// denne går TIL (FSM?) og bør lagres der og tatt inn som argument
	completeOrderConsensus := make(chan ConfigFile.OrderMsg, 9)		// denne går TIL (FSM?) og bør lagres der og tatt inn som argument
	RemoteID := "123.123"
	LocalID := "321.321"
	thisElev := ConfigFile.CabOrders{}
	ConfigFile.AllCabOrders[LocalID] = &thisElev
	newElev := ConfigFile.CabOrders{}
	//END DEBUG SHIT
	ConfigFile.AllCabOrders["LocalID"] = &newElev
	ConfigFile.AllCabOrders[LocalID].Floor = 3
	ConfigFile.AllCabOrders[LocalID].Direction = ConfigFile.UP
	cabOrdersRx := make(chan map[string]*ConfigFile.CabOrders)
	// BEM EM IN SCUTTY
	// hallordersTx := make(chan ConfigFile.AllHallOrders)						// REMEMBER TO USE 
	// BEM EM UP SCUTTY
	
	for{
		select{
		case p := <-cabOrdersRx:
			for elevator := range ConfigFile.AllCabOrders{								// kan det bli "kræsj" med lengden på "local" og "remote"...?
				fmt.Printf("Floor %d\n", ConfigFile.AllCabOrders[elevator].Floor)		// debug variable!
				for floor := 0; floor < ConfigFile.Num_floors; floor++{
					
					remote := p[elevator].CabOrders[floor]
					local := ConfigFile.AllCabOrders[elevator].CabOrders[floor]
					// nå går vi gjennom ALLE heiser, alle etg.
					switch local.OrderState{		// legge denne på egen linje eller noe? for en laaaaang jævel

					case ConfigFile.Default:
						if(RemoteID != LocalID){
							switch remote.OrderState{
							case ConfigFile.Default:
								break
							case ConfigFile.Inactive:
								local.OrderState = ConfigFile.Inactive
								local.AckdBy = local.AckdBy[:0] 		// destroy
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
						switch remote.OrderState{
							case ConfigFile.Inactive:
								break									// if state remains the same, do nothing

							case ConfigFile.PendingAck:
								local.OrderState = remote.OrderState
								local.AckdBy = append(local.AckdBy, LocalID)	//sett inn IP/ID
								break

							case ConfigFile.Active:
								break									// cannot skip PendingAck-state

						}

					case ConfigFile.PendingAck:
						switch remote.OrderState{
							case ConfigFile.Inactive:
								break

							case ConfigFile.PendingAck:
								// ADD all others ACKed to local

								if len(local.AckdBy) == 3{			// denne må selvsagt byttes til en "dynamisk" sak, og ikke bare sjekker antall!
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
						switch remote.OrderState{
							case ConfigFile.Inactive:
								local.OrderState = remote.OrderState
								local.AckdBy = local.AckdBy[:0]			// destroy
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

				} // her slutter vi å ittere over "floor"
				// må gå over og oppdatere direction/floor?
			}
		default:
			break
		}
	}
}