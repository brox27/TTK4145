package main

import (
	. "./driver"
	. "./globals"
	. "./fsm"

)

func main(){
	Elev1, Elev2, Elev3 := new(Elevator{})
	Elev1.Id, Elev2.Id, ELev3.Id = ELEVATOR_IPS[0], ELEVATOR_IPS[1], ELEVATOR_IPS[2]
	AllOrder[Elev1.Id]=Elev1
	AllOrder[Elev2.Id]=Elev2
	AllOrder[Elev3.Id]=Elev3

	go Elev1.RUN()
	go Elev2.RUN()
	go Elev3.RUN()
}