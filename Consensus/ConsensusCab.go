package Consensus

import (
	"../ConfigFile"
	"../Network"
	"../driver"
	"time"
//	"fmt"
//	"reflect"
)

func ConsensusCab(ClearCabOrderChan chan int, ConsensusCabChan chan map[string]*ConfigFile.ConsensusCab, CabButtonChan chan int, PeerUpdateChan chan ConfigFile.PeerUpdate) {
	cabOrdersRx := make(chan map[string]*ConfigFile.ConsensusCab)
	cabOrdersTx := make(chan map[string]*ConfigFile.ConsensusCab)
	go Network.Transmitter(ConfigFile.CabConsensusPort, cabOrdersTx)
	go Network.LocalTransmitter(ConfigFile.CabConsensusPort, cabOrdersTx)
	go Network.Receiver(ConfigFile.CabConsensusPort, cabOrdersRx)

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
		select {

		case remoteCabConsensus := <-cabOrdersRx:
			
			for elevID := range remoteCabConsensus {
				_, exists := AllCabOrders[elevID];
               if !exists {
					AllCabOrders[elevID] = remoteCabConsensus[elevID]
				}

	//			if !reflect.DeepEqual(remoteCabConsensus[elevID], AllCabOrders[elevID]) {
  	//          	fmt.Printf(ConfigFile.ColorCC+"[CC]:  New worldview: %v\n"+ConfigFile.ColorNone, remoteCabConsensus[elevID])
    //          }
				for floor := 0; floor < ConfigFile.Num_floors; floor++ {


					remote := remoteCabConsensus[elevID].CabButtons[floor]
					Merge(&AllCabOrders[elevID].CabButtons[floor], remote, elevID, LivingPeers, 
						func() {
							if elevID == ConfigFile.LocalID {
								driver.SetButtonLamp(ConfigFile.BUTTON_ORDER_COMMAND, floor, 1)
								ConsensusCabChan <- AllCabOrders
							}
						}, 
						func() {
							ConsensusCabChan <- AllCabOrders
						})
				}

	//			if !reflect.DeepEqual(remoteCabConsensus[elevID], AllCabOrders[elevID]) {
    //		      	fmt.Printf(ConfigFile.ColorCC+"[CC]:  Worldview updated: \n   From: %v\n   To:   %v\n"+ConfigFile.ColorNone, remoteCabConsensus[elevID], AllCabOrders[elevID])
    // 		    }
    //		fmt.Printf("*CC end of remoteCabConsensus channel \n")
			}

		case ClearedCabOrder := <- ClearCabOrderChan:
			Deactivate(&AllCabOrders[ConfigFile.LocalID].CabButtons[ClearedCabOrder], LivingPeers)
			driver.SetButtonLamp(ConfigFile.BUTTON_ORDER_COMMAND, ClearedCabOrder, 0)
			ConsensusCabChan <- AllCabOrders

		case NewCabButton := <-CabButtonChan:
			Activate(&AllCabOrders[ConfigFile.LocalID].CabButtons[NewCabButton])

		case <- transmittTimer:
			cabOrdersTx <- AllCabOrders


		case PeerUpdate := <-PeerUpdateChan:
			LivingPeers = PeerUpdate.Peers
			if len(PeerUpdate.Lost)!=0{
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
