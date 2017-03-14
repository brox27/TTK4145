package Peers

import (
	"./conn"
	"../ConfigFile"
	"fmt"
	"net"
	"sort"
	"time"
	"reflect"
)

// ** This is the "peers module" from the "ofical" Network Module provided to us i GOlang - NOT something we wrote ourselves ** \\

const interval = 10 * time.Millisecond
const timeout = 200 * time.Millisecond

func Transmitter(port int, id string, transmitEnable <-chan bool) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			conn.WriteTo([]byte(id), addr)
		}
	}
}

func Receiver(port int, peerUpdateCh chan ConfigFile.PeerUpdate) {

	var buf [1024]byte
	var p ConfigFile.PeerUpdate
	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])

		id := string(buf[:n])

		// ** Adding new connection ** \\
		p.New = ""
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}

		// ** Removing dead connection ** \\
		p.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}
		// ** Sending update ** \\
		if updated {
			p.Peers = make([]string, 0, len(lastSeen))

			for k, _ := range lastSeen {
				p.Peers = append(p.Peers, k)
			}

			sort.Strings(p.Peers)
			sort.Strings(p.Lost)
			peerUpdateCh <- p
		}
	}
}


func Repeater(ch_in interface{}, chs_out ...interface{}) {
	for {
		v, _ := reflect.ValueOf(ch_in).Recv()
		for _, c := range chs_out {
			reflect.ValueOf(c).Send(v)
		}
	}
}