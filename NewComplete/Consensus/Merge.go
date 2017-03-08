package Consensus

import (
	"../ConfigFile"
	. "../Network"
	//	"fmt"
)

func Merge(local *ConfigFile.OrderStatus, remote ConfigFile.OrderStatus, RemoteID string, LivingPeers []string, onActive func(), onInactive func()) {
	// LOCAL liste/var
	// Global Liste/var
	// ConfigFile.LocalID - ConfigFile
	// RemoteID
	// Perr liste
	// onActive
	// onInactive

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
			temp := append(remote.AckdBy, ConfigFile.LocalId)
			local.AckdBy = append(local.AckdBy, temp)  // legger til alle andre..? ANDERS:  local.ackdBy ~= remote.ackdBy ~ [localId];
			if len(local.AckdBy) >= len(LivingPeers) { // denne må selvsagt byttes til en "dynamisk" sak, og ikke bare sjekker antall!
				local.OrderState = ConfigFile.Active
				onActive()
			}
			break

		case ConfigFile.Active:
			local.OrderState = ConfigFile.Active
			temp := append(remote.AckdBy, ConfigFile.LocalId)
			local.AckdBy = append(local.AckdBy, temp)
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
			onInnactive()
			break

		case ConfigFile.PendingAck:
			break

		case ConfigFile.Active:
			temp := append(remote.AckdBy, ConfigFile.LocalId)
			local.AckdBy = append(local.AckdBy, temp)
			break

		}
	}
}
