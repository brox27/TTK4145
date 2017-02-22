package main

import (
	"../ConfigFile"
	"./network/bcast"
	"./network/localip"
	"./network/peers"
	"fmt"
	"os"
	"time"
)

struct sendChan{
	chan ConfigFile.NewOrder
	chan ConfigFile.CompleteOrder
	chan ConfigFile.Acknowledge
	chan ConfigFile.Heartbeat
}

struct recieveChan{
	chan ConfigFile.NewOrder
	chan ConfigFile.CompleteOrder
	chan ConfigFile.Acknowledge
	chan ConfigFile.Heartbeat
}

func main() {

	// setter opp -- aka Anders work
	var id string
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// Setter opp EGNE recieve channels
	NewOrderRx := make(chan ConfigFile.NewOrder, 5)
	CompleteOrderRx := make(chan ConfigFile.CompleteOrder, 5)
	AcknowledgeRx := make(chan ConfigFile.Acknowledge, 5)
	HeartbeatRx := make(chan ConfigFile.Heartbeat, 5)

	// < ALL ORDERS > \\

	// Setter opp EGNE Send channels
	NewOrderTx := make(chan ConfigFile.NewOrder, 5)
	CompleteOrderTx := make(chan ConfigFile.CompleteOrder, 5)
	AcknowledgeTx := make(chan ConfigFile.Acknowledge, 5)
	HeartbeatTx := make(chan ConfigFile.Heartbeat, 5)
	// < ALL ORDERS > \\

	// setter opp alle rutinene for å sende..?
	go bcast.Transmitter(16569, NewOrderTx, CompleteOrderTx, AcknowledgeTx, HeartbeatTx)
	go bcast.Receiver(16569, NewOrderRx, CompleteOrderRx, AcknowledgeRx, HeartbeatRx)

	time.Sleep(1 * time.Second)

	hest := ConfigFile.NewOrder{}
	Sender(hest)
}

func Sender(data interface{}) bool {
	ThisMsgId := 12 // sett ID på melding som sendes ut --> denne MÅ være dynamisk

	switch t := data.(type) {
	case ConfigFile.NewOrder:
		fmt.Println("New Order")

	case ConfigFile.CompleteOrder:
		fmt.Println("CompleteOrder")

	case ConfigFile.Acknowledge:
		fmt.Println("Acknowledge")

	case ConfigFile.Heartbeat:
		fmt.Println("Heartbeat")

	default:
		fmt.Printf("ERROR: Unknown type: %T", t)
	}
	
	ThisMsgId++
	numAcks := 0

	for i := 0; i < 100; i++ {
		inncomming := ConfigFile.Acknowledge{}
		if inncomming.MsgId == ThisMsgId {
			numAcks++
		} else {
			// put back on channeL???
		}
		fmt.Println("lolz kjorer")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("lolz ute m.", numAcks)
	return numAcks
}
