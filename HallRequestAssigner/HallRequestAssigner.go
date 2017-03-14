package hallRequestAssigner

import (
	"../ConfigFile"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type AssignerCompatibleElev struct {
	sync.RWMutex `json:"-"`
	Behaviour    string  `json:"behaviour"`
	Floor        int     `json:"floor"`
	Direction    string  `json:"direction"`
	CabRequests  [4]bool `json:"cabRequests"`
}

type AssignerCompatibleInput struct {
	sync.RWMutex `json:"-"`
	HallRequests [4][2]bool                         `json:"hallRequests"`
	States       map[string]*AssignerCompatibleElev `json:"states"`
}

func toAssignerCompatible(elev ConfigFile.Elev) AssignerCompatibleElev {
	temp := AssignerCompatibleElev{}
	switch elev.State {
	case ConfigFile.INITIALIZE:
		fallthrough
	case ConfigFile.IDLE:
		temp.Behaviour = "idle"
	case ConfigFile.MOVING:
		temp.Behaviour = "moving"
	case ConfigFile.DOORSOPEN:
		temp.Behaviour = "doorOpen"
	}
	temp.Floor = elev.Floor
	switch elev.Direction {
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
	ElevatorStatesChan chan ConfigFile.AllStates,
	LocalOrdersChan chan [][]bool,
	FromPeersToHallReqAss chan ConfigFile.PeerUpdate) {

	var LostPeers []string
	var LivingPeers []string

	localCopy := AssignerCompatibleInput{}
	localCopy.States = make(map[string]*AssignerCompatibleElev)
	localCopy.States[ConfigFile.LocalID] = &AssignerCompatibleElev{}

	for {
		select {
		case newConsensusHall := <-ConsensusHallChan:
			for button := 0; button < 2; button++ {
				for floor := 0; floor < ConfigFile.Num_floors; floor++ {
					if newConsensusHall.HallButtons[floor][button].OrderState == ConfigFile.Active {
						localCopy.HallRequests[floor][button] = true
					} else {
						localCopy.HallRequests[floor][button] = false
					}
				}
			}

		case newConsensusCab := <-ConsensusCabChan:
			for elevID := range newConsensusCab {
				if _, ok := localCopy.States[elevID]; ok {
					for floor := 0; floor < ConfigFile.Num_floors; floor++ {
						if newConsensusCab[elevID].CabButtons[floor].OrderState == ConfigFile.Active {
							localCopy.States[elevID].CabRequests[floor] = true
						} else {
							localCopy.States[elevID].CabRequests[floor] = false
						}
					}
				}
			}

		case newElevatorStates := <-ElevatorStatesChan:
			//	Var1 := ConfigFile.AllStates{}
			localCopy.Lock()
			//	Var1.StateMap = make(map[string]*ConfigFile.Elev)
			for _, peer := range LivingPeers {
				if peer != "" {
					if _, ok := newElevatorStates.StateMap[peer]; ok {
						if _, ok := localCopy.States[peer]; ok {

							//*newElevatorStates.StateMap[peer]

							temp := toAssignerCompatible(*newElevatorStates.StateMap[peer])
							temp.Lock()
							localCopy.States[peer].Behaviour = temp.Behaviour
							localCopy.States[peer].Floor = temp.Floor
							localCopy.States[peer].Direction = temp.Direction
							temp.Unlock()
						}
					}
				}
				time.Sleep(3 * time.Millisecond)
			}
			localCopy.Unlock()

			/* FROM:
			newElevatorStates.Lock()
			for elevID := range newElevatorStates.StateMap {
				if elevID != "" {
					localCopy.Lock()
					//annotherCopy := localCopy.States
					if _, ok := newElevatorStates.StateMap[elevID]; ok {
						if _, ok := localCopy.States[elevID]; ok {
							temp := toAssignerCompatible(*newElevatorStates.StateMap[elevID])
							temp.Lock()
							localCopy.States[elevID].Behaviour = temp.Behaviour
							localCopy.States[elevID].Floor = temp.Floor
							localCopy.States[elevID].Direction = temp.Direction
							temp.Unlock()
						}
					}
					localCopy.Unlock()
				}
			}
			newElevatorStates.Unlock()
			*/

		case PeerUpdate := <-FromPeersToHallReqAss:
			fmt.Printf("Peer status %+v \n", PeerUpdate)
			LostPeers = PeerUpdate.Lost
			LivingPeers = PeerUpdate.Peers
			if _, ok := localCopy.States[PeerUpdate.New]; !ok {
				localCopy.States[PeerUpdate.New] = &AssignerCompatibleElev{}

			}
			if LostPeers != nil {
				temp := make(map[string]*AssignerCompatibleElev)
				temp[ConfigFile.LocalID] = localCopy.States[ConfigFile.LocalID]
				for _, elevID := range PeerUpdate.Peers {
					temp[elevID] = localCopy.States[elevID]
				}
				localCopy.States = temp
			}
		}

		arg, _ := json.Marshal(localCopy)
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		result, err := exec.Command("sh", "-c", dir+"/HALL --input '"+string(arg)+"'").Output()

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
			LocalOrdersChan <- assignedOrders
		} else {
			fmt.Printf("error : %+v : %+v\n", err, result)
		}
	}
}
