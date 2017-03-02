package main

import (
	"../ConfigFile"
	"./network/bcast"
	"./network/localip"
	"./network/peers"
	"flag"
	"fmt"
	"os"
	"time"
)

type channels struct {
	NewOrder      chan ConfigFile.NewOrder
	CompleteOrder chan ConfigFile.CompleteOrder
	Acknowledge   chan ConfigFile.Acknowledge
	Heartbeat     chan ConfigFile.Heartbeat
}

var sendChan channels
var recieveChan channels

func main() {

	// setter opp -- aka Anders work
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid()) // fjerne os.Getpid for å ikke prosessor ID -> stor forldel om restarte prosessen!
	}

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// Setter opp EGNE recieve channels
	recieveChan.NewOrder = make(chan ConfigFile.NewOrder, 5)
	recieveChan.CompleteOrder = make(chan ConfigFile.CompleteOrder, 5)
	recieveChan.Acknowledge = make(chan ConfigFile.Acknowledge, 5)
	recieveChan.Heartbeat = make(chan ConfigFile.Heartbeat, 5)

	// < ALL ORDERS > \\

	// Setter opp EGNE Send channels
	sendChan.NewOrder = make(chan ConfigFile.NewOrder, 5)
	sendChan.CompleteOrder = make(chan ConfigFile.CompleteOrder, 5)
	sendChan.Acknowledge = make(chan ConfigFile.Acknowledge, 5)
	sendChan.Heartbeat = make(chan ConfigFile.Heartbeat, 5)
	// < ALL ORDERS > \\
	go bcast.Transmitter(16569, sendChan.NewOrder, sendChan.CompleteOrder, sendChan.Acknowledge, sendChan.Heartbeat)
	go bcast.Receiver(16569, recieveChan.NewOrder, recieveChan.CompleteOrder, recieveChan.Acknowledge, recieveChan.Heartbeat)

	// IGNORA ALL BELLOW -- test variabler ++
	p := <-peerUpdateCh
	println("lol<")
	fmt.Println(p.Peers)
	println("lol>")

	time.Sleep(1 * time.Second)

	hest := ConfigFile.NewOrder{}
	Sender(hest)
}

func Sender(data interface{}) bool {
	ThisMsgId := 12 // sett ID på melding som sendes ut --> denne MÅ være dynamisk

	switch t := data.(type) {
	case ConfigFile.NewOrder:
		fmt.Println("New Order")
		sendChan.NewOrder <- t

	case ConfigFile.CompleteOrder:
		fmt.Println("CompleteOrder")
		sendChan.CompleteOrder <- t

	case ConfigFile.Acknowledge:
		fmt.Println("Acknowledge")
		sendChan.Acknowledge <- t

	case ConfigFile.Heartbeat:
		fmt.Println("Heartbeat")
		sendChan.Heartbeat <- t

	default:
		fmt.Printf("ERROR: Unknown type: %T", t)
		return false
	}

	ThisMsgId++ // update msg Id somehow
	numAcks := 0
	numAcks++

	// sjekke for Acks somehow
	for {

		// TIMEPOUT
		// BREAK
	}

	return false
}
