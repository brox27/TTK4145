package hallRequestAssigner

import (
	"../ConfigFile"
	"encoding/json"
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
	HallRequests [4][2]bool 								`json:"hallRequests"` 
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
	for f := range elev.Orders {
        temp.CabRequests[f] = elev.Orders[f][ConfigFile.BUTTON_ORDER_COMMAND]
    }
	return temp
}


func HallRequestAssigner(
	ConsensusHallChan chan ConfigFile.ConsensusHall,
	ConsensusCabChan chan map[string]*ConfigFile.ConsensusCab,
	ElevatorStatesChan chan map[string]*ConfigFile.Elev,
	LocalOrdersChan chan [][]bool,
	FromPeersToHallReqAss chan ConfigFile.PeerUpdate) {

	var LostPeers []string

	localCopy := AssignerCompatibleInput{}
	localCopy.States = make(map[string]*AssignerCompatibleElev)
	localCopy.States[ConfigFile.LocalID] = &AssignerCompatibleElev{}

	for {
		select {
		case newConsensusHall := <-ConsensusHallChan:
			fmt.Printf("*HRA newConsensusHall\n")
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
			fmt.Printf("*HRA newConsensusCab\n")
			for elevID := range newConsensusCab {
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
			fmt.Printf("*HRA newElevatorStates\n")
			for elevID := range newElevatorStates {
				if (elevID!=""){
					if _, ok := localCopy.States[elevID]; ok {
						temp := toAssignerCompatible(newElevatorStates[elevID])
						localCopy.States[elevID].Behaviour = temp.Behaviour
						localCopy.States[elevID].Floor = temp.Floor
						localCopy.States[elevID].Direction = temp.Direction

						/*
						localCopy.States[elevID].Behaviour
						 = toAssignerCompatible(newElevatorStates[elevID]).Behaviour
						localCopy.States[elevID].Floor = toAssignerCompatible(newElevatorStates[elevID]).Floor
						localCopy.States[elevID].Direction = toAssignerCompatible(newElevatorStates[elevID]).Direction
						*/
					}
				}
			}

		case PeerUpdate := <- FromPeersToHallReqAss:
			fmt.Printf("Peer status %+v \n",PeerUpdate )
			LostPeers = PeerUpdate.Lost
			if _, ok := localCopy.States[PeerUpdate.New]; !ok {
				localCopy.States[PeerUpdate.New] = &AssignerCompatibleElev{}

			}
			if (LostPeers != nil){
				temp := make (map[string]*AssignerCompatibleElev)
				temp[ConfigFile.LocalID]=localCopy.States[ConfigFile.LocalID]
				for _, elevID := range PeerUpdate.Peers{
					temp[elevID]=localCopy.States[elevID]
				}
				localCopy.States=temp
        	}
		}


		fmt.Printf(ConfigFile.ColorHRA+"[HRA]: local copy:\n"+ConfigFile.ColorNone)
        fmt.Printf(ConfigFile.ColorHRA+"   HallRequests : %+v\n"+ConfigFile.ColorNone, localCopy.HallRequests)
        fmt.Printf(ConfigFile.ColorHRA+"   States : \n"+ConfigFile.ColorNone)

        for e := range localCopy.States {
            fmt.Printf(ConfigFile.ColorHRA+"     %v : %+v\n"+ConfigFile.ColorNone, e, localCopy.States[e])
        }

        arg, _ := json.Marshal(localCopy)
        dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
        result, err := exec.Command("sh","-c", dir+"/HALL --input '" + string(arg) + "'").Output()

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
            fmt.Printf("*HRA aboce LocalOrdersChan\n")
            LocalOrdersChan <- assignedOrders
            fmt.Printf("*HRA below LocalOrdersChan\n")
        }else{
            fmt.Printf("err : %+v : %+v\n", err, result)
        }
        fmt.Printf("*HRA end of for \n")        
    }
}
