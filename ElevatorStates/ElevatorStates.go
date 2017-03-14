package ElevatorStates

import (
	"../ConfigFile"
	"../Network"
	"time"
)

func ElevatorStatesCoordinator(StateChan chan ConfigFile.Elev, ElevatorStatesChan chan ConfigFile.AllStates) {

	StateNetworkRx := make(chan ConfigFile.AllStates)
	StateNetworkTx := make(chan ConfigFile.AllStates)
	go Network.Transmitter(ConfigFile.ElevatorStatesPort, StateNetworkTx)
	go Network.Receiver(ConfigFile.ElevatorStatesPort, StateNetworkRx)
	transmittTimer := time.NewTicker(time.Millisecond * 50).C

	States := ConfigFile.AllStates{}
	States.StateMap = make(map[string]*ConfigFile.Elev)
	updateFlag := true

	for {
		select {
		case newLocalState := <-StateChan:
			States.StateMap[ConfigFile.LocalID] = &newLocalState
			ElevatorStatesChan <- States

		case newRemoteStates := <-StateNetworkRx:
			States.Lock()
			newRemoteStates.Lock()
			for elevID := range newRemoteStates.StateMap {
				if elevID != ConfigFile.LocalID && States.StateMap[elevID] != newRemoteStates.StateMap[elevID] {
					States.StateMap[elevID] = newRemoteStates.StateMap[elevID]
					updateFlag = true
				}
			}
			if updateFlag {
				updateFlag = false
				ElevatorStatesChan <- States
			}
			newRemoteStates.Unlock()
			States.Unlock()

		case <-transmittTimer:
			StateNetworkTx <- States
		}

	}
}
