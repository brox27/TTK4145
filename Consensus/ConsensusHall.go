package Consensus

import (
	"../ConfigFile"
	. "../Network"
	"../driver"
	"time"
	"fmt"
    "reflect"
)



func ConsensusHall(ClearHallOrderChan chan [2]int, ConsensusHallChan chan ConfigFile.ConsensusHall,  HallButtonChan chan [2]int, PeerUpdateChan chan ConfigFile.PeerUpdate) {
	hallordersRx := make(chan ConfigFile.ConsensusHall)
	hallordersTx := make(chan ConfigFile.ConsensusHall)
	go Transmitter(ConfigFile.HallConesnsusPort, hallordersTx)
	go Receiver(ConfigFile.HallConesnsusPort, hallordersRx)
	transmittTimer := time.NewTicker(time.Millisecond * 50).C

	var LivingPeers []string

	localHallConsensus := ConfigFile.ConsensusHall{}
    localHallConsensus.ID = ConfigFile.LocalID
	for floor := 0; floor < ConfigFile.Num_floors; floor++ {
		for button := 0; button < ConfigFile.Num_buttons-1; button ++{
			localHallConsensus.HallButtons[floor][button].OrderState = ConfigFile.Default
		}
	}

	for {
		select {
		case remoteHallConsensus := <-hallordersRx:
        
            if !reflect.DeepEqual(remoteHallConsensus, localHallConsensus) {
                fmt.Printf(ConfigFile.ColorCH+"[CH]:  New worldview: %v\n"+ConfigFile.ColorNone, remoteHallConsensus)
            }
        
			RemoteID := remoteHallConsensus.ID
			for floor := 0; floor < ConfigFile.Num_floors; floor++ {
				for button := 0; button < ConfigFile.Num_buttons-1; button++ {

					remote := remoteHallConsensus.HallButtons[floor][button]
					Merge(&localHallConsensus.HallButtons[floor][button], remote, RemoteID, LivingPeers, 
						func() {
                            driver.SetButtonLamp(ConfigFile.ButtonType(button), floor, 1)
                            ConsensusHallChan <- localHallConsensus
                        }, 
						func() {
                            driver.SetButtonLamp(ConfigFile.ButtonType(button), floor, 0)
                            ConsensusHallChan <- localHallConsensus
                        })
				}
			}
            
            if !reflect.DeepEqual(remoteHallConsensus, localHallConsensus) {
                fmt.Printf(ConfigFile.ColorCH+"[CH]:  Worldview updated: \n   From: %v\n   To:   %v\n"+ConfigFile.ColorNone, remoteHallConsensus, localHallConsensus)
            }

		case ClearedHallOrder := <- ClearHallOrderChan:
            fmt.Printf(ConfigFile.ColorCH+"[CH]:  Cleared hall order: %+v\n"+ConfigFile.ColorNone, ClearedHallOrder)
            
			Deactivate(&localHallConsensus.HallButtons[ClearedHallOrder[0]][ClearedHallOrder[1]], LivingPeers)
			driver.SetButtonLamp(ConfigFile.ButtonType(ClearedHallOrder[1]), ClearedHallOrder[0], 0)
            ConsensusHallChan <- localHallConsensus
                        
            fmt.Printf(ConfigFile.ColorCH+"[CH]:  Our worldview: \n         %v\n"+ConfigFile.ColorNone, localHallConsensus)
		
		case NewHallButton := <-HallButtonChan:
            fmt.Printf(ConfigFile.ColorCH+"[CH]:  New hall button: %+v\n"+ConfigFile.ColorNone, NewHallButton)
            
			Activate(&localHallConsensus.HallButtons[NewHallButton[0]][NewHallButton[1]])

            fmt.Printf(ConfigFile.ColorCH+"[CH]:  Our worldview: \n         %v\n"+ConfigFile.ColorNone, localHallConsensus)
            
		case <- transmittTimer:
            hallordersTx <- localHallConsensus

        case PeerUpdate := <-PeerUpdateChan:
            fmt.Printf(ConfigFile.ColorCH+"[CH]:  Peer update: %+v\n"+ConfigFile.ColorNone, PeerUpdate)
			LivingPeers = PeerUpdate.Peers
		}
	}
}
