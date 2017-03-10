package HallRequestAssigner

import (
	"../ConfigFile"
	"encoding/json"
	"time"
	//	"reflect"
	"fmt"
	"path/filepath"
	"os/exec"
	"os"
)


type AssignerCompatibleElev struct {
	Behaviour string 	`json:"behaviour"`
	Floor int 			`json:"floor"`
	Direction string 	`json:"direction"`
	CabRequests [4]bool `json:"cabRequests"`
}

type AssignerCompatibleInput struct {
	HallRequests [4][2]bool 								`json:"hallRequests"` 												// hvorfor kan ikke denne settes som "tom" slice e.g. [][] >> [4][2]
	States       map[string]*AssignerCompatibleElev			`json:"states"`
}

func toAssignerCompatible(elev *ConfigFile.Elev) AssignerCompatibleElev {
	temp := AssignerCompatibleElev{}
	switch elev.State{
		case ConfigFile.INITIALIZE:
			fallthrough
		case ConfigFile.IDLE:
			temp.Behaviour = "idle"
		case ConfigFile.RUNNING:
			temp.Behaviour = "moving"
		case ConfigFile.DOORSOPEN:
			temp.Behaviour = "doorOpen"
	}
	temp.Floor = elev.Floor
	switch elev.Direction{
		case ConfigFile.UP:
			temp.Direction = "up"
		case ConfigFile.DOWN:
			temp.Direction = "down"
		case ConfigFile.NEUTRAL:
			temp.Direction = "stop"
	}
	return temp
}


func HallReq(
	ConsensusHallChan chan ConfigFile.ConsensusHall,
	ConsensusCabChan chan map[string]*ConfigFile.ConsensusCab,
	ElevatorStatesChan chan map[string]*ConfigFile.Elev,
	LocalOrdersChan chan [][]bool,
	FromPeersToHallReqAss chan ConfigFile.PeerUpdate) {

	var LostPeers []string

	localCopy := AssignerCompatibleInput{}
	localCopy.States = make(map[string]*AssignerCompatibleElev)
	localCopy.States[ConfigFile.LocalID] = &AssignerCompatibleElev{}

//	localCope.States
	//var LastSentToFSM [ConfigFile.Num_floors][3]bool // hvorfor skrev vi 2???			TRENGER VI DENNE OM VI SPAMMER???
	//go tester(LocalOrdersChan)

	for {
	//	fmt.Printf("\n \n NEW ROUND! \n")
		select {
		case newConsensusHall := <-ConsensusHallChan:
            fmt.Printf(ConfigFile.ColorHRA+"[HRA]: new hall orders: %+v\n"+ConfigFile.ColorNone, newConsensusHall)
		//	fmt.Printf("ConsensusHall \n")
			// får inn (update) fra ConsensusHall
			for button := 0; button < 2; button++ {
				for floor := 0; floor < ConfigFile.Num_floors; floor++ {
					if newConsensusHall.HallButtons[floor][button].OrderState == ConfigFile.Active {
						localCopy.HallRequests[floor][button] = true
					}else {
						localCopy.HallRequests[floor][button] = false
					}
				}
			}

		case newConsensusCab := <-ConsensusCabChan:
            fmt.Printf(ConfigFile.ColorHRA+"[HRA]: new cab orders: %+v\n"+ConfigFile.ColorNone, newConsensusCab)
			for elevID := range newConsensusCab {
                fmt.Printf(ConfigFile.ColorHRA+"   %v : %+v\n"+ConfigFile.ColorNone, elevID, newConsensusCab[elevID])
				if _, ok := localCopy.States[elevID]; ok {
					for floor := 0; floor < ConfigFile.Num_floors; floor++ {
						if newConsensusCab[elevID].CabButtons[floor].OrderState == ConfigFile.Active {
							localCopy.States[elevID].CabRequests[floor] = true																																	
						}else {
							localCopy.States[elevID].CabRequests[floor] = false
						}
					}
				}
			}

		case newElevatorStates := <-ElevatorStatesChan:
            fmt.Printf(ConfigFile.ColorHRA+"[HRA]: new elevator states: %+v\n"+ConfigFile.ColorNone, newElevatorStates)

			for elevID := range newElevatorStates {
                fmt.Printf(ConfigFile.ColorHRA+"   %v : %+v\n"+ConfigFile.ColorNone, elevID, newElevatorStates[elevID])
				temp := toAssignerCompatible(newElevatorStates[elevID])
				localCopy.States[elevID] = &temp
			}

		case PeerUpdate := <- FromPeersToHallReqAss:
            fmt.Printf(ConfigFile.ColorHRA+"[HRA]: new peer list: %+v\n"+ConfigFile.ColorNone, PeerUpdate)
			LostPeers = PeerUpdate.Lost

        }

        
        fmt.Printf(ConfigFile.ColorHRA+"[HRA]: local copy:\n"+ConfigFile.ColorNone)
        fmt.Printf(ConfigFile.ColorHRA+"   HallRequests : %+v\n"+ConfigFile.ColorNone, localCopy.HallRequests)
        fmt.Printf(ConfigFile.ColorHRA+"   States : \n"+ConfigFile.ColorNone)
        for e := range localCopy.States {
            fmt.Printf(ConfigFile.ColorHRA+"     %v : %+v\n"+ConfigFile.ColorNone, e, localCopy.States[e])
        }

        // sjekke og evt. ta ut de som ikke lever \\
        if LostPeers != nil{
            for _, elevID := range LostPeers{
                delete(localCopy.States, elevID)
            }
        }

        arg, _ := json.Marshal(localCopy)
        dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
//		fmt.Printf("sender inn:     %+v\n", string(arg) )
        result, err := exec.Command("sh","-c", dir+"/HALL --input '" + string(arg) + "'").Output()
//		fmt.Printf("tilbake:   %+v\n%  +v\n\n", err, string(result) )
        if err == nil {
            var a map[string][][]bool
            json.Unmarshal(result, &a)

            assignedOrders := make([][]bool, ConfigFile.Num_floors)
            for i := range assignedOrders {
                assignedOrders[i] = make([]bool, 3)
            }

            for f := 0; f < ConfigFile.Num_floors; f++ {
                for b := 0; b < 2; b++ {
                    assignedOrders[f][b] = a[ConfigFile.LocalID][f][b]
                }
                assignedOrders[f][2] = localCopy.States[ConfigFile.LocalID].CabRequests[f]
            }
            //fmt.Println("assigned: %+v\n", assignedOrders)
//				fmt.Printf("min local id er: %+v",ConfigFile.LocalID)
//				fmt.Printf("sender orders: %+v\n",  assignedOrders)
            LocalOrdersChan <- assignedOrders
            //fmt.Printf("%+v\n", a)
        }else{
            fmt.Printf("err : %+v : %+v\n", err, result)
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
}

func tester(LocalOrdersChan chan [][]bool){
	time.Sleep(10*time.Second)
	println("*********************************************************************************")
	var temp [][]bool
	if temp!=nil{
		temp[2][2] = true
		LocalOrdersChan <- temp
	}
}