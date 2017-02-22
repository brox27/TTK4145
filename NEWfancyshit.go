package main

import (
	"fmt"
	"./ConfigFile"
)

var lolz ConfigFile.OrderStatus

func main() {
	lolz.OrderState = ConfigFile.Active
	fmt.Printf("lolz \n")
	switchen(lolz)
}

func switchen(local ConfigFile.OrderStatus) {
	rec :=ConfigFile.OrderStatus{}

	rec.OrderState = ConfigFile.Inactive

	switch local.OrderState{

	case ConfigFile.Inactive:
		switch rec.OrderState{
			case ConfigFile.Inactive:
				break

			case ConfigFile.PendingAck:
				local.OrderState = rec.OrderState
				local.AckdBy = append(local.AckdBy, "123.123.123.123")	//sett inn IP/ID
				break

			case ConfigFile.Active:
				break

		}

	case ConfigFile.PendingAck:
		switch rec.OrderState{
			case ConfigFile.Inactive:
				break

			case ConfigFile.PendingAck:
				// ADD all others ACKed to local

				if len(local.AckdBy) == 3{			// denne må selvsagt byttes til en "dynamisk" sak, og ikke bare sjekker antall!
					local.OrderState = ConfigFile.Active
				}
				break

			case ConfigFile.Active:
				break

		}
			
	case ConfigFile.Active:
		switch rec.OrderState{
			case ConfigFile.Inactive:
				local.OrderState = rec.OrderState
				local.AckdBy = local.AckdBy[:0]
				break

			case ConfigFile.PendingAck:
				break

			case ConfigFile.Active:
				// Anders har appendet IP til ACked by her også... why? eller ser Why, men er det nødvendig?
				break

		}
	}
}