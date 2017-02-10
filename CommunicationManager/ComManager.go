// Com manager

package main

import (
	. "../ConfigFile"
//S	. "../Network"
	//	"../CommunicationManager"
	//	"flag"
	//	"fmt"
	//	"os"
	//	"time"
)


func main() {
	println("hei")
//	Test := All_Information{}
//	Test.function()
	hestefaen := NewOrder{}

	NO := make(chan NewOrder)	// MÅ en spessifisere hvor mange en vil ha..??
//	completeOrder := make(chan completeOrder)
//	acknowledge := make(chan acknowledge)
//	Heartbeat := make(chan Heartbeat)
	// newOrder := make(chan newOrder) --> må vel ha en til her for "ALL-orders"

	NO <- hestefaen

	hesten := Num_floors
	hesten ++
// skal dette i main/init eller 
}