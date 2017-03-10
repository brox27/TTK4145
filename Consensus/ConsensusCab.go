package Consensus

import (
	"../ConfigFile"
	. "../Network"
	"../driver"
	"time"
	. "fmt"
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
	AllCabOrders[ConfigFile.LocalID] = &thisCab
	for floor := 0; floor < ConfigFile.Num_floors; floor++ {
		thisCab.CabButtons[floor].OrderState = ConfigFile.Default
	}
	for {
		select {

		case remoteCabConsensus := <-cabOrdersRx:
			isThere := false
			for elevID := range remoteCabConsensus {
				// IF "ny" add to map!
				if _, ok := AllCabOrders[elevID]; ok {
					isThere = true
				}
				if !isThere{
					thisCab := ConfigFile.ConsensusCab{}
					AllCabOrders[elevID] = &thisCab
				}

				for floor := 0; floor < ConfigFile.Num_floors; floor++ {

					remote := remoteCabConsensus[elevID].CabButtons[floor]
					Merge(&AllCabOrders[ConfigFile.LocalID].CabButtons[floor], remote, elevID, LivingPeers, 
						func() {
							if elevID == ConfigFile.LocalID {
								Println("Cab order light ON at floor", floor)
								driver.SetButtonLamp(ConfigFile.BUTTON_ORDER_COMMAND, floor, 1)
							}
						}, 
						func() {
							Println("%+v completed a cab order at floor %+v\n", elevID, floor)
						})
				}
			}
			ConsensusCabChan <- AllCabOrders

		case ClearedCabOrder := <- ClearCabOrderChan:
			Deactivate(&AllCabOrders[ConfigFile.LocalID].CabButtons[ClearedCabOrder], LivingPeers)
			Println("Cab order light OFF at floor", ClearedCabOrder)
			driver.SetButtonLamp(ConfigFile.BUTTON_ORDER_COMMAND, ClearedCabOrder, 0)

		case NewCabButton := <-CabButtonChan:
			Activate(&AllCabOrders[ConfigFile.LocalID].CabButtons[NewCabButton])

		case <- transmittTimer:
			cabOrdersTx <- AllCabOrders

		case PeerUpdate := <-PeerUpdateChan:
			Printf("Peer update:\n  %+v\n", PeerUpdate)
			LivingPeers = PeerUpdate.Peers
			if len(PeerUpdate.Lost)!=0{
				Printf("Lost: %+v\n", PeerUpdate.Lost)
				Printf("\ntest\n%+v\n", AllCabOrders)
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