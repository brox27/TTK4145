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
	NO := NewOrder{}
	CO := CompleteOrder{}
	H := Heartbeat{}
	A := Acknowledge{}

//	NO := make(chan NewOrder)	// MÅ en spessifisere hvor mange en vil ha..??
	completeOrder := make(chan CompleteOrder, 5)
	acknowledge := make(chan Acknowledge, 5)
	heartbeat := make(chan Heartbeat, 5)
	newOrder := make(chan NewOrder, 5)// --> må vel ha en til her for "ALL-orders"
	//allOrders := make(chan map[string]*Elev, 5)
	



	// drit i alt under
	hesten := Num_floors
	hesten ++

	newOrder <- NO
	acknowledge <- A
	heartbeat <- H
	completeOrder <- CO

// skal dette i main/init eller 
}