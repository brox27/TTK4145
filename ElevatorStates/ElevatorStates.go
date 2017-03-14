package ElevatorStates

import (
	"../ConfigFile"
	. "../Network"
	//. "fmt"
	"time"
)

func ElevatorStatesCoordinator(StateChan chan ConfigFile.Elev, ElevatorStatesChan chan map[string]*ConfigFile.Elev) {

	StateNetworkRx := make(chan map[string]*ConfigFile.Elev)
	StateNetworkTx := make(chan map[string]*ConfigFile.Elev)
	go Transmitter(ConfigFile.ElevatorStatesPort, StateNetworkTx)
	go Receiver(ConfigFile.ElevatorStatesPort, StateNetworkRx)
	transmittTimer := time.NewTicker(time.Millisecond * 50).C

	AllStates := map[string]*ConfigFile.Elev{}

	for {
		select {
		case newLocalState := <-StateChan:
			AllStates[ConfigFile.LocalID] = &newLocalState
			ElevatorStatesChan <- AllStates

		case newRemoteStates := <-StateNetworkRx:
			for elevID := range newRemoteStates {
				if elevID != ConfigFile.LocalID  &&  AllStates[elevID] != newRemoteStates[elevID] {
 					AllStates[elevID] = newRemoteStates[elevID]
                    ElevatorStatesChan <- AllStates
				}
			}

		case <- transmittTimer:
			StateNetworkTx <- AllStates
		}
		
	}
}
