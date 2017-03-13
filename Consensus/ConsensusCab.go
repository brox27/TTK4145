package Consensus

import (
	"../ConfigFile"
	. "../Network"
	"../driver"
	"time"
	"fmt"
	"reflect"
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
		select {

		case remoteCabConsensus := <-cabOrdersRx:
			for elevID := range remoteCabConsensus {
				_, exists := AllCabOrders[elevID];
                if !exists {
                    fmt.Printf(ConfigFile.ColorCC+"[CC]:  New elevator: %v\n"+ConfigFile.ColorNone, remoteCabConsensus[elevID])
					AllCabOrders[elevID] = remoteCabConsensus[elevID]
				}

				if !reflect.DeepEqual(remoteCabConsensus[elevID], AllCabOrders[elevID]) {
                    fmt.Printf(ConfigFile.ColorCC+"[CC]:  New worldview: %v\n"+ConfigFile.ColorNone, remoteCabConsensus[elevID])
                }
											
				for floor := 0; floor < ConfigFile.Num_floors; floor++ {


					remote := remoteCabConsensus[elevID].CabButtons[floor]
					fmt.Printf("på f: %+v : %+v merges med %+v, Min status: %+v remote: %+v \n", floor, ConfigFile.LocalID, elevID,AllCabOrders[elevID].CabButtons[floor], remote)
					Merge(&AllCabOrders[elevID].CabButtons[floor], remote, elevID, LivingPeers, 
						func() {
							if elevID == ConfigFile.LocalID {
								fmt.Println("Cab order light ON at floor", floor)
								driver.SetButtonLamp(ConfigFile.BUTTON_ORDER_COMMAND, floor, 1)
								ConsensusCabChan <- AllCabOrders
							}
						}, 
						func() {
							fmt.Printf("%+v completed a CAB order at floor %+v \n", elevID, floor)
							ConsensusCabChan <- AllCabOrders
						})
					fmt.Printf("Merge resultat; %+v\n", AllCabOrders[elevID].CabButtons[floor])
				}

				if !reflect.DeepEqual(remoteCabConsensus[elevID], AllCabOrders[elevID]) {
                    fmt.Printf(ConfigFile.ColorCC+"[CC]:  Worldview updated: \n   From: %v\n   To:   %v\n"+ConfigFile.ColorNone, remoteCabConsensus[elevID], AllCabOrders[elevID])
                }
                //fmt.Printf("*CC end of remoteCabConsensus channel \n")
			}

		case ClearedCabOrder := <- ClearCabOrderChan:
			fmt.Printf(ConfigFile.ColorCC+"[CC]:  Cleared cab order: %+v\n"+ConfigFile.ColorNone, ClearedCabOrder)
			Deactivate(&AllCabOrders[ConfigFile.LocalID].CabButtons[ClearedCabOrder], LivingPeers)
			driver.SetButtonLamp(ConfigFile.BUTTON_ORDER_COMMAND, ClearedCabOrder, 0)
			fmt.Printf("*CC over ConsensusCabChan \n")
			ConsensusCabChan <- AllCabOrders
			fmt.Printf("*CC under ConsensusCabChan \n")

		case NewCabButton := <-CabButtonChan:
			fmt.Printf(ConfigFile.ColorCC+"[CC]:  New cab button: %+v\n"+ConfigFile.ColorNone, NewCabButton)
			Activate(&AllCabOrders[ConfigFile.LocalID].CabButtons[NewCabButton])

		case <- transmittTimer:
			cabOrdersTx <- AllCabOrders

		case PeerUpdate := <-PeerUpdateChan:
			fmt.Printf("Peer update:\n  %+v\n", PeerUpdate)
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