package optimalHallReq

import(
	. "../ConfigFile"
	"encoding/json"
	"fmt"
)

func optimalHallReq(LocalID, HallRequests [][2]bool, CabRequests, AllStates map[string]*Elev, Peers []string)[][2]bool {
	if val, ok := ALlStates[LocalID] ; !ok{
		return HallRequests
	}
	
}