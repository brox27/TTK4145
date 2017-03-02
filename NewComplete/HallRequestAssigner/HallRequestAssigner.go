package HallRequestAssigner

import (
	"../ConfigFile"
	"encoding/json"
//	"reflect"
	"fmt"
)


type TheOneAnders3Likes struct {
	HallRequests [2][2]bool		`json:"hallRequests"`
	States map[string]SomeSort	`json:"states"`
}

type SomeSort struct{
	Behaviour behave	`json:"behaviour"`
	Floor int 			`json:"floor"`
	Direction dirs		`json:"direction"`
	CabRequests [2]bool	`json:"cabRequests"`
}

type behave string
const ( 
	MOVING behave = "moving"
	IDLE behave = "idle"
	DOOROPEN behave = "doorOpen"
)

type dirs string
const ( 
	UP dirs = "up"
	DOWN dirs = "down"
	STOP dirs = "stop"
)


type ReturnFromAnders3 struct{
	LF map[string] [ConfigFile.Num_floors][2]bool
}



func main(){
}

type AsignerCompatibleInput struct{
	HallRequests [2][ConfigFile.Num_floors]bool 	`json:"hallRequests"`  //[Num_floors][Num_buttons - 1]OrderStatus
	States map[string]*ConfigFile.Elev 				`json:"states"`
}


func HallReq(
	ConsensusHallChan chan ConfigFile.ConsensusHall, 
	ConsensusCabChan chan map[string]*ConfigFile.ConsensusCab, 
	ElevatorStatesChan chan map[string]*ConfigFile.Elev, 
	LocalOrdersChan chan [ConfigFile.Num_floors][ConfigFile.Num_buttons]bool){
	// Får inn fra ConsHall, ConsCab, ElevStates
	// bruker dette til å "oppdatere via Anders3"

	// LocalID := "123"

	localCopy := AsignerCompatibleInput{}
	localCopy.States = make(map[string]*ConfigFile.Elev)
	var LastSentToFSM [ConfigFile.Num_floors][2]bool
	send := false

	for{	
	select{
		case newConsensusHall := <- ConsensusHallChan:
			// får inn (update) fra ConsensusHall
			for button := 0; button <2; button ++{
				for floor := 0; floor < ConfigFile.Num_floors; floor++{
					if newConsensusHall.HallButtons[floor][button].OrderState == ConfigFile.Active{
						localCopy.HallRequests[button][floor] = true
					}else{
						localCopy.HallRequests[button][floor] = false
					}
				}
			}
	
		case newConsensusCab := <- ConsensusCabChan:
			// får inn (update) fra ConsensusHall
			for elevID := range newConsensusCab {
				for floor := 0; floor < ConfigFile.Num_floors; floor++{
					if newConsensusCab[elevID].CabButtons[floor].OrderState == ConfigFile.Active{
						localCopy.States[elevID].CabOrders[floor] = true
					}else{
						localCopy.States[elevID].CabOrders[floor] = false
					}
				}
			}
	
		case newElevatorStates := <- ElevatorStatesChan:
			// får inn (update) fra ElevatorStates
			for elevID := range newElevatorStates {
				localCopy.States[elevID] = newElevatorStates[elevID]
			}
		default:
			if (2<1){		// Timer e.l. for å sjekke periodsik?
				buf, _ := json.Marshal(localCopy)
				fmt.Printf("sender dette til Anders %s\n", buf)
				var ReturnFromAnders [ConfigFile.Num_floors][2]bool
				for button := 0; button <2; button ++{
					for floor := 0; floor < ConfigFile.Num_floors; floor++{
						if LastSentToFSM[floor][button] != ReturnFromAnders[floor][button]{
							send = true
						}
					}
				}
				if send{
					// send til FSM
				}

			}
		}
	}
}
