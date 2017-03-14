package Consensus

import (
	"../ConfigFile"
	"../Network"
	"../driver"
	"fmt"
	"time"
)

func ConsensusCab(ClearCabOrderChan chan int, ConsensusCabChan chan map[string]*ConfigFile.ConsensusCab, CabButtonChan chan int, PeerUpdateChan chan ConfigFile.PeerUpdate) {
	cabOrdersRx := make(chan map[string]*ConfigFile.ConsensusCab)
	cabOrdersTx := make(chan map[string]*ConfigFile.ConsensusCab)
	localcabOrdersTx := make(chan map[string]*ConfigFile.ConsensusCab)

	go Network.Transmitter(ConfigFile.CabConsensusPort, cabOrdersTx)
	go Network.LocalTransmitter(ConfigFile.CabConsensusPort, localcabOrdersTx)
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
			fmt.Printf("ny remote: %+v \n", remoteCabConsensus)

			for elevID := range remoteCabConsensus {
				_, exists := AllCabOrders[elevID]
				if !exists {
					AllCabOrders[elevID] = remoteCabConsensus[elevID]
				}

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
			}

		case ClearedCabOrder := <-ClearCabOrderChan:
			Deactivate(&AllCabOrders[ConfigFile.LocalID].CabButtons[ClearedCabOrder], LivingPeers)
			driver.SetButtonLamp(ConfigFile.BUTTON_ORDER_COMMAND, ClearedCabOrder, 0)
			ConsensusCabChan <- AllCabOrders

		case NewCabButton := <-CabButtonChan:
			Activate(&AllCabOrders[ConfigFile.LocalID].CabButtons[NewCabButton])

		case <-transmittTimer:
			fmt.Printf("nå sender til remote \n")
			cabOrdersTx <- AllCabOrders
			fmt.Printf("nå sender til local \n")
			localcabOrdersTx <- AllCabOrders

		case PeerUpdate := <-PeerUpdateChan:
			LivingPeers = PeerUpdate.Peers
			if len(PeerUpdate.Lost) != 0 {
				for _, lostID := range PeerUpdate.Lost {
					for floor := 0; floor < ConfigFile.Num_floors; floor++ {
						if AllCabOrders[lostID].CabButtons[floor].OrderState == ConfigFile.Inactive {
							AllCabOrders[lostID].CabButtons[floor].OrderState = ConfigFile.Default
						}
					}
				}
			}
		}
	}
}
