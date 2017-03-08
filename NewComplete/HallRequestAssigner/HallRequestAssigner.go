package HallRequestAssigner

import (
	"../ConfigFile"
	"encoding/json"
	"time"
	//	"reflect"
	"fmt"
	"os/exec"
)

type TheOneAnders3Likes struct {
	HallRequests [2][2]bool          `json:"hallRequests"`
	States       map[string]SomeSort `json:"states"`
}

type SomeSort struct {
	Behaviour   behave  `json:"behaviour"`
	Floor       int     `json:"floor"`
	Direction   dirs    `json:"direction"`
	CabRequests [2]bool `json:"cabRequests"`
}

type behave string

const (
	MOVING   behave = "moving"
	IDLE     behave = "idle"
	DOOROPEN behave = "doorOpen"
)

type dirs string

const (
	UP   dirs = "up"
	DOWN dirs = "down"
	STOP dirs = "stop"
)

type ReturnFromAnders3 struct {
	LF map[string][ConfigFile.Num_floors][2]bool
}

func main() {
}

type AssignerCompatibleElev struct {
	Behaviour string `json:behaviour`
	Floor int `json:floor`
	Direction string `json:direction`
	CabRequests []bool `json:cabRequests`
}

func toAssignerCompatible(ConfigFile.Elev) AssignerCompatibleElev {

}

type AsignerCompatibleInput struct {
	HallRequests [2][ConfigFile.Num_floors]bool `json:"hallRequests"` //[Num_floors][Num_buttons - 1]OrderStatus
	States       map[string]*AssignerCompatibleElev    `json:"states"`
}

func HallReq(
	ConsensusHallChan chan ConfigFile.ConsensusHall,
	ConsensusCabChan chan map[string]*ConfigFile.ConsensusCab,
	ElevatorStatesChan chan map[string]*ConfigFile.Elev,
	LocalOrdersChan chan [ConfigFile.Num_floors][ConfigFile.Num_buttons]bool) {
	// Får inn fra ConsHall, ConsCab, ElevStates
	//// bruker dette til å "oppdatere via Anders3"
	//timerChan := make(chan int)
	//go timer(timerChan)
	//var tester [ConfigFile.Num_floors][3]bool

	localCopy := AsignerCompatibleInput{}
	localCopy.States = make(map[string]*ConfigFile.Elev)
	localCopy.States[ConfigFile.LocalID] = &ConfigFile.Elev{}
	var LastSentToFSM [ConfigFile.Num_floors][2]bool // hvorfor skrev vi 2???

	send := false

	for {
		select {
		case newConsensusHall := <-ConsensusHallChan:
			// får inn (update) fra ConsensusHall
			for button := 0; button < 2; button++ {
				for floor := 0; floor < ConfigFile.Num_floors; floor++ {
					if newConsensusHall.HallButtons[floor][button].OrderState == ConfigFile.Active {
						localCopy.HallRequests[button][floor] = true
					} else {
						localCopy.HallRequests[button][floor] = false
					}
				}
			}

		case newConsensusCab := <-ConsensusCabChan:
			// får inn (update) fra ConsensusHall
			for elevID := range newConsensusCab {
				for floor := 0; floor < ConfigFile.Num_floors; floor++ {
					if newConsensusCab[elevID].CabButtons[floor].OrderState == ConfigFile.Active {
						localCopy.States[elevID].CabOrders[floor] = true
						// 																																	** DEBUG!
						//tester[floor][2] = true
						//																																		** Sluitt
					} else {
						fmt.Printf("%+v\n%+v\n", elevID, localCopy)
						localCopy.States[elevID].CabOrders[floor] = false
					}
				}
			}

		case newElevatorStates := <-ElevatorStatesChan:
			// får inn (update) fra ElevatorStates
			for elevID := range newElevatorStates {
				localCopy.States[elevID] = newElevatorStates[elevID]
			}
		// DEBUG!!!!
		//case lol := <-timerChan:
		//_ = lol
		//LocalOrdersChan <- tester
		// DEBUG SLUTT
			/*
			if 2 < 1 { // Timer e.l. for å sjekke periodsik?
				buf, _ := json.Marshal(localCopy)
				fmt.Printf("sender dette til Anders %s\n", buf)
				var ReturnFromAnders [ConfigFile.Num_floors][2]bool
				for button := 0; button < 2; button++ {
					for floor := 0; floor < ConfigFile.Num_floors; floor++ {
						if LastSentToFSM[floor][button] != ReturnFromAnders[floor][button] {
							send = true
						}
					}
				}
				if send {
					// send til FSM
					send = false
				}

			}
			*/
		}


		arg, _ ::= json.Marshal(localCopy)
		result, err := exec.Command("hall_request_assigner", "--input " + string(arg)).Output()
		if err == nil {
			var a map[string][]bool
			json.Unmarshal(result, &a)
			//someChan <- a[ConfigFile.LocalID]
			fmt.Printf("%+v\n", a)
		}


	}
}

func timer(timerChan chan int) {
	time.Sleep(3 * time.Second)
	timerChan <- 1
}
