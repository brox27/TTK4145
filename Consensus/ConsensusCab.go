package Consensus

import (
	"../ConfigFile"
	. "../Network"
	"../driver"
	"time"
	"fmt"
)


func ConsensusCab(ClearCabOrderChan chan int, ConsensusCabChan chan map[string]*ConfigFile.ConsensusCab, CabButtonChan chan int, PeerUpdateChan chan ConfigFile.PeerUpdate) {
	cabOrdersRx := make(chan map[string]*ConfigFile.ConsensusCab)
	cabOrdersTx := make(chan map[string]*ConfigFile.ConsensusCab)
	go Transmitter(ConfigFile.CabConsensusPort, cabOrdersTx)
	go Receiver(ConfigFile.CabConsensusPort, cabOrdersRx)
	transmittTimer := time.NewTicker(time.Millisecond * 50).C

	var LivingPeers []string

	AllCabOrders := make(map[string]*ConfigFile.ConsensusCab)
	thisCab := ConfigFile.ConsensusCab{}
    thisCab.ID = ConfigFile.LocalID
	AllCabOrders[ConfigFile.LocalID] = &thisCab
	for floor := 0; floor < ConfigFile.Num_floors; floor++ {
		thisCab.CabButtons[floor].OrderState = ConfigFile.Default
	}
	for {
		//fmt.Printf("\n \n")
		select {

		case remoteCabConsensus := <-cabOrdersRx:
			for elevID := range remoteCabConsensus {
            
				// IF "ny" add to map!
				_, exists := AllCabOrders[elevID];
                if !exists {
					thisCab := ConfigFile.ConsensusCab{}
					AllCabOrders[elevID] = &thisCab
				}
                
				//fmt.Printf("merger meg: %+v med: %+v \n", ConfigFile.LocalID, elevID)													// SE HER!
				for floor := 0; floor < ConfigFile.Num_floors; floor++ {

					remote := remoteCabConsensus[elevID].CabButtons[floor]
					Merge(&AllCabOrders[elevID].CabButtons[floor], remote, elevID, LivingPeers, 
						func() {
							if elevID == ConfigFile.LocalID {
								fmt.Printf(ConfigFile.ColorCC+"Cab order light ON at floor %v\n"+ConfigFile.ColorNone, floor)
								driver.SetButtonLamp(ConfigFile.BUTTON_ORDER_COMMAND, floor, 1)
                                ConsensusCabChan <- AllCabOrders
							}
						}, 
						func() {
							fmt.Printf(ConfigFile.ColorCC+"%+v completed a cab order at floor %+v\n"+ConfigFile.ColorNone, elevID, floor)
                            ConsensusCabChan <- AllCabOrders
						})
				}
				//fmt.Printf("%+v has the following CAB statuses: %+v\n", elevID, AllCabOrders[elevID])									// SE HER
			}

		case ClearedCabOrder := <- ClearCabOrderChan:
            fmt.Printf(ConfigFile.ColorCC+"[CC]:  Cleared cab order: %+v\n"+ConfigFile.ColorNone, ClearedCabOrder)
			Deactivate(&AllCabOrders[ConfigFile.LocalID].CabButtons[ClearedCabOrder], LivingPeers)
			fmt.Println(ConfigFile.ColorCC+"Cab order light OFF at floor %v\n"+ConfigFile.ColorNone, ClearedCabOrder)
			driver.SetButtonLamp(ConfigFile.BUTTON_ORDER_COMMAND, ClearedCabOrder, 0)

		case NewCabButton := <-CabButtonChan:
            fmt.Printf(ConfigFile.ColorCC+"[CC]:  New cab button: %+v\n"+ConfigFile.ColorNone, NewCabButton)
			Activate(&AllCabOrders[ConfigFile.LocalID].CabButtons[NewCabButton])

		case <- transmittTimer:
			cabOrdersTx <- AllCabOrders

		case PeerUpdate := <-PeerUpdateChan:
            fmt.Printf(ConfigFile.ColorCC+"[CC]:  Peer update: %+v\n"+ConfigFile.ColorNone, PeerUpdate)
			LivingPeers = PeerUpdate.Peers
			if len(PeerUpdate.Lost)!=0{
				fmt.Printf("Lost: %+v\n", PeerUpdate.Lost)
				fmt.Printf("\ntest\n%+v\n", AllCabOrders)
				for _, lostID := range PeerUpdate.Lost {
					for floor := 0; floor < ConfigFile.Num_floors; floor++{
						if (AllCabOrders[lostID].CabButtons[floor].OrderState == ConfigFile.Inactive){
							AllCabOrders[lostID].CabButtons[floor].OrderState = ConfigFile.Default
						}
					}
				}
			}
		}
	}
}