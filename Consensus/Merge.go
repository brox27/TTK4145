package Consensus

import (
	"../ConfigFile"
	//"fmt"
	"sort"
)

func Merge(local *ConfigFile.OrderStatus, remote ConfigFile.OrderStatus, RemoteID string, LivingPeers []string, onActive func(), onInactive func()) {

	switch local.OrderState {

	case ConfigFile.Default:
		if RemoteID != ConfigFile.LocalID {
			switch remote.OrderState {
			case ConfigFile.Default:
				break
			case ConfigFile.Inactive:
				local.OrderState = ConfigFile.Inactive
				local.AckdBy = local.AckdBy[:0] // destroy
				break
			case ConfigFile.PendingAck:
				local.OrderState = ConfigFile.PendingAck
				local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID) // bør det være local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID) her? -> anders: local.ackdBy = remote.ackdBy ~ [localId];
				break
			case ConfigFile.Active:
				local.OrderState = ConfigFile.Active
				local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID) // Anders  local.ackdBy = remote.ackdBy ~ [localId];
				onActive()
			}
		}

	case ConfigFile.Inactive:
		switch remote.OrderState {
		case ConfigFile.Default:
			break
		case ConfigFile.Inactive:
			break
		case ConfigFile.PendingAck:
			local.OrderState = remote.OrderState
			local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID) // adder egen ACK ANDERS: local.ackdBy = remote.ackdBy ~ [localId];
			break
		case ConfigFile.Active:
			break // cannot skip PendingAck-state
		}

	case ConfigFile.PendingAck:
		switch remote.OrderState {
		case ConfigFile.Default:
			break
		case ConfigFile.Inactive:
			break

		case ConfigFile.PendingAck:
			local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID)
			sort.Strings(local.AckdBy)
			local.AckdBy = removeDuplicates(local.AckdBy)
//			local.AckdBy = append(local.AckdBy, temp)  // legger til alle andre..? ANDERS:  local.ackdBy ~= remote.ackdBy ~ [localId];
			//if len(local.AckdBy) >= len(LivingPeers) { // denne må selvsagt byttes til en "dynamisk" sak, og ikke bare sjekker antall!
			if checkAcks(local.AckdBy, LivingPeers){
				local.OrderState = ConfigFile.Active
				onActive()
			}
			break

		case ConfigFile.Active:
			local.OrderState = ConfigFile.Active
//			temp := append(remote.AckdBy, ConfigFile.LocalID)
//			local.AckdBy = append(local.AckdBy, temp)
			local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID)
			onActive()
			break

		}

	case ConfigFile.Active:
		switch remote.OrderState {
		case ConfigFile.Default:
			break
		case ConfigFile.Inactive:
			local.OrderState = remote.OrderState
			local.AckdBy = local.AckdBy[:0] // destroy
			onInactive()
			break

		case ConfigFile.PendingAck:
			break

		case ConfigFile.Active:
			//temp := append(remote.AckdBy, ConfigFile.LocalId)
			//local.AckdBy = append(local.AckdBy, temp)
			local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID)
			break

		}
	}

	sort.Strings(local.AckdBy)
	local.AckdBy = removeDuplicates(local.AckdBy)

}

func Activate(local *ConfigFile.OrderStatus){
	switch local.OrderState {
	case ConfigFile.Default:
		fallthrough
	case ConfigFile.Inactive:
		local.OrderState = ConfigFile.PendingAck
		local.AckdBy = append(local.AckdBy, ConfigFile.LocalID)
		break
	}
}

func Deactivate(local *ConfigFile.OrderStatus, LivingPeers []string){
	if len(LivingPeers) == 0 || (len(LivingPeers) == 1  &&  LivingPeers[0] == ConfigFile.LocalID) {
		// Only us or "noone" on the network
		local.OrderState = ConfigFile.Inactive 																// SE MEG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! VVIKTIG(litt) endret til default
	} else {
		local.OrderState = ConfigFile.Inactive
	}
	local.AckdBy = local.AckdBy[:0]
}




func removeDuplicates(elements []string) []string {
    // Use map to record duplicates as we find them.
    encountered := map[string]bool{}
    result := []string{}

    for v := range elements {
        if encountered[elements[v]] == true {
            // Do not add duplicate.
        } else {
            // Record this element as an encountered element.
            encountered[elements[v]] = true
            // Append to result slice.
            result = append(result, elements[v])
        }
    }
    // Return the new slice.
    return result
}


func checkAcks(acks []string, LivingPeers []string) bool{
	if (len(acks) >= len(LivingPeers)) && (LivingPeers != nil){
		globalFlag := true
		for _, i := range LivingPeers{
			flag := false
			for _, j :=range acks{
				if(i == j){
					flag = true
				}
			}
			if(!flag){
				globalFlag = false
			}
		}
		return globalFlag
	}
	return false
}