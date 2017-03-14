package ElevatorStates

import (
	"../ConfigFile"
	. "../Network"
	"time"
)

func ElevatorStatesCoordinator(StateChan chan ConfigFile.Elev, ElevatorStatesChan chan ConfigFile.AllStates) {

	StateNetworkRx := make(chan ConfigFile.AllStates)
	StateNetworkTx := make(chan ConfigFile.AllStates)
	go Transmitter(ConfigFile.ElevatorStatesPort, StateNetworkTx)
	go Receiver(ConfigFile.ElevatorStatesPort, StateNetworkRx)
	transmittTimer := time.NewTicker(time.Millisecond * 50).C

	States := ConfigFile.AllStates{}
	States.StateMap = make(map[string]*ConfigFile.Elev)
	//	AllStates := map[string]*ConfigFile.Elev{}

	for {
		select {
		case newLocalState := <-StateChan:
			States.StateMap[ConfigFile.LocalID] = &newLocalState
			ElevatorStatesChan <- States

		case newRemoteStates := <-StateNetworkRx:
			States.Lock()
			for elevID := range newRemoteStates.StateMap {
				if elevID != ConfigFile.LocalID && States.StateMap[elevID] != newRemoteStates.StateMap[elevID] {
					States.StateMap[elevID] = newRemoteStates.StateMap[elevID]
					ElevatorStatesChan <- States
				}
			}
			States.Unlock()

		case <-transmittTimer:
			StateNetworkTx <- States
		}

	}
}
