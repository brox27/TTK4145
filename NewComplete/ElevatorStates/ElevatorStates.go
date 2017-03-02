package ElevatorStates

import (
	"../ConfigFile"
)
//StateChan := make(chan ConfigFile.Elev)									// Fra FSM.go til StateFileNotYetMade
//ElevatorStatesChan := make(chan map[string]*ConfigFile.Elev)		// Fra ElevatorStates til HallReqAss

func DoSomethingSmart(StateChan chan ConfigFile.Elev, ElevatorStatesChan chan map[string]*ConfigFile.Elev){
	LocalId := "123.123"
	// tar inn ELEV struct fra StateChan -> sender "noe" på NW
	//						-> sender "noe" til HallReqAss på ElevatorStatesChan
	StateNetworkRx := make(chan map[string]*ConfigFile.Elev)
//	StateNetworkTx := make(chan map[string]*ConfigFile.Elev)
	
	AllStates := map[string]*ConfigFile.Elev{}

	for{
		select{
		case newLocalState := <- StateChan:
			AllStates[LocalId] = &newLocalState
			ElevatorStatesChan <- AllStates

		case newRemoteStates := <- StateNetworkRx:
			for elevID := range newRemoteStates {
				if (elevID == LocalId){
					// IGNORE!, we know our own state "best"
				}else{
					AllStates[elevID] = newRemoteStates[elevID]
				}
			}
			// sende hver gang får inn eller, bare hvis sikker på oppdatering?
			ElevatorStatesChan <- AllStates
		default: 
			// timer som "spammer" på nettverket? altså sender AllStates
			break
		}
	}
}