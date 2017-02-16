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
	go bcast.Transmitter(16569, NewOrderTx)
	go bcast.Receiver(16569, NewOrderRx)

	go bcast.Transmitter(16569, CompleteOrderTx)
	go bcast.Receiver(16569, CompleteOrderRx)

	go bcast.Transmitter(16569, AcknowledgeTx)
	go bcast.Receiver(16569, AcknowledgeRx)

	go bcast.Transmitter(16569, HeartbeatTx)
	go bcast.Receiver(16569, HeartbeatRx)

	// < ALL ORDERS > \\
	// < ALL ORDERS > \\

	//

	// %% her (tror) jeg vi må ha ett eller annet mer, men faen ikke sikker på hva

	//

	//

	//
	time.Sleep(1 * time.Second)
	//
	//	go DoSomethingSmart(CompleteOrderTx, CompleteOrderRx)
	//

	hesten := ConfigFile.NewOrder{}
	hesten.MsgId = 12
	Send(hesten)

}

// Den som "snakker" med FSM o.l. -> input arg er en struct basert på hva som vil sendes (Heartbeat/NewOrder e.l.)
func Send(name interface{}) {
	if c, ok := name.(ConfigFile.NewOrder); ok {
		fmt.Println(c.MsgId)
	}
	println(name)
}

/*
func DoSomethingSmart(channel chan ConfigFile.CompleteOrder, inn chan ConfigFile.CompleteOrder) {

	// f.eks sende signal og vente på ack/telle ack
	varr := ConfigFile.CompleteOrder{}
	varr.MsgId = 12
	varr.Button = 3
	varr.Floor = 2

	channel <- varr
	time.Sleep(1 * time.Second)
	recieved := <-inn
	fmt.Printf("Received: %#v\n", recieved)

}

*/
