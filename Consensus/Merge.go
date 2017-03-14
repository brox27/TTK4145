package Consensus

import (
	"../ConfigFile"
	"sort"
)

func Merge(local *ConfigFile.OrderStatus, remote ConfigFile.OrderStatus, RemoteID string, LivingPeers []string, onActive func(), onInactive func()) {

	switch local.OrderState {

	case ConfigFile.Default:
		switch remote.OrderState {
		case ConfigFile.Default:
			break
		case ConfigFile.Inactive:
			local.OrderState = ConfigFile.Inactive
			local.AckdBy = local.AckdBy[:0]
			break
		case ConfigFile.PendingAck:
			local.OrderState = ConfigFile.PendingAck
			local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID)
			break
		case ConfigFile.Active:
			local.OrderState = ConfigFile.Active
			local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID)
			onActive()
		}

	case ConfigFile.Inactive:
		switch remote.OrderState {
		case ConfigFile.Default:
			break
		case ConfigFile.Inactive:
			break
		case ConfigFile.PendingAck:
			local.OrderState = remote.OrderState
			local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID)
			break
		case ConfigFile.Active:
			break
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
			if checkAcks(local.AckdBy, LivingPeers) {
				local.OrderState = ConfigFile.Active
				onActive()
			}
			break

		case ConfigFile.Active:
			local.OrderState = ConfigFile.Active
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
			local.AckdBy = local.AckdBy[:0]
			onInactive()
			break

		case ConfigFile.PendingAck:
			break

		case ConfigFile.Active:
			local.AckdBy = append(remote.AckdBy, ConfigFile.LocalID)
			break

		}
	}

	sort.Strings(local.AckdBy)
	local.AckdBy = removeDuplicates(local.AckdBy)

}

func Activate(local *ConfigFile.OrderStatus) {
	switch local.OrderState {
	case ConfigFile.Default:
		fallthrough
	case ConfigFile.Inactive:
		local.OrderState = ConfigFile.PendingAck
		local.AckdBy = append(local.AckdBy, ConfigFile.LocalID)
		break
	}
}

func Deactivate(local *ConfigFile.OrderStatus, LivingPeers []string) {
	if len(LivingPeers) == 0 || (len(LivingPeers) == 1 && LivingPeers[0] == ConfigFile.LocalID) {
		local.OrderState = ConfigFile.Default
	} else {
		local.OrderState = ConfigFile.Inactive
	}
	local.AckdBy = local.AckdBy[:0]
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func checkAcks(acks []string, LivingPeers []string) bool {
	if (len(acks) >= len(LivingPeers)) && (LivingPeers != nil) {
		globalFlag := true
		for _, i := range LivingPeers {
			flag := false
			for _, j := range acks {
				if i == j {
					flag = true
				}
			}
			if !flag {
				globalFlag = false
			}
		}
		return globalFlag
	}
	return false
}
