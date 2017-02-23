package main

import (
	. "../driver"
//	. "../globals"
//	"../OrderManager"
	. "../ConfigFile"
	//. "fmt"
	"runtime"
	//"time"
)

func main() {
	println(" ░░░░░░░░░░░░░░░░░░░░░ \n ░░░░░░░░░░░░▄▀▀▀▀▄░░░ \n ░░░░░░░░░░▄▀░░▄░▄░█░░ \n ░▄▄░░░░░▄▀░░░░▄▄▄▄█░░ \n █░░▀▄░▄▀░░░░░░░░░░█░░ \n ░▀▄░░▀▄░░░░█░░░░░░█░░ \n ░░░▀▄░░▀░░░█░░░░░░█░░ \n ░░░▄▀░░░░░░█░░░░▄▀░░░ \n ░░░▀▄▀▄▄▀░░█▀░▄▀░░░░░ \n ░░░░░░░░█▀▀█▀▀░░░░░░░ \n ░░░░░░░░▀▀░▀▀░░░░░░░░ \n")
	localId = "123.123"
	RUN(LocalId)
}

type States int

const (
	INITIALIZE States = iota
	IDLE
	RUNNING
	ARRIVEDFLOOR
	SHOULDSTOP
	DOORSOPEN
)

type Elevator struct {
	State States
	LastFloor int
	Dir Direction
	Orders    [Num_floors][Num_buttons]int
	Id 			string
}

var ALL_ELEVATORS map[string]*Elevator

// OLD: func (elev *Elevator) RUN() {

func RUN(localId string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	elev.State=INITIALIZE
	InitElev()

	eventChan := make(chan map[Event]interface{})


	go EventHandler(eventChan)
	// eventMap:= <-eventChan

	for {

		// BROX FUCKER OPP (flyttet fra IDLE)
		/*
		if eventMap[BUTTON_PRESSED]!=nil{
			button:=0
			data := eventMap[BUTTON_PRESSED].([2]int)

			for floor := range ELEVATOR_IPS{		// IPS = IPs  / IP adresser  -> og sjekker om ordren finnes allerede
				if(ALL_ORDERS[floor].Orders[data[1]][data[0]]){
					break
				}
				button++
			}

			// dette er om det er en Cab order -> i så fall bare sender
			if button==ConfigFile.Num_buttons{
				order := OrderMsg{}
				order.Button=data[0]
				order.Floor=data[1]
				order.MsgType=NEW

				// denne må trolig endres etter Anders fix!!
				go Sender(order)

				// calculate cost -> Anders lager dette ?
				//CalculateCost(, elev.Id)
				// denne må fjernes nå som utenfor!!
				elev.State=RUNNING
			} else {
				break
			}
		}
		// BROX FUCKER OPP SLUTT --> over kan være en funksjon? bedre måte å løse dette på?
*/

		/*
		UAVHENIG AV HVILKEN CASE VI ER I!
		* sjekke innkommende events
		* sende på nettverk? -> egen og andre sine status?
		*/

		switch(elev.State){

		case INITIALIZE:
			if(GetFloorSensorSignal()!=-1){

				SetMotorDirection(NEUTRAL)		// endre navn? setter motor dir
				SetFloorLight(GetFloorSensorSignal())
				// oppdatere egen dir/floor
				// synkronisere ordre mot andre på nettverket!

				// WaitGroup e.l.? for å være sikker på klar
				elev.State=IDLE

				break
			}else{
				SetMotorDirection(DOWN)		// sjekke noe så den ikke "settes" så ofte?
			}
			break

		case IDLE:
			/*nextDir := OrderManager.NextDirection()
			if nextDir != IDLE{
				SetMotorDirection(nextDir)
				elev.State=RUNNING
			}*/
			break

		case RUNNING:
	/*
			event arrived at floor
				Oppdatere ny etg!
				Oppdatere lys for etg.
				if OrderManager.Shouldstop(){
					SetMotorDirection(NEUTRAL)
					//oppdatere direction?
					//sette lys for dør åpen
					// starte timer
					//oppdatere orders lista
					// si ordre utført! til andre
					elev.State=DOORSOPEN
			}
	*/

		case DOORSOPEN:
			/*
			if TIMEOUT
				Skru av timer?
				Sett lys dør av
				nextDir = OrderManager.GetNextOrder()
				if nextDir == NEUTRAL{
					elev.State=IDLE
				}else{
					setMotorDir(nextDir)
					elev.State=RUNNING
				}	
		*/

		}
	}

	defer SetMotorDirection(NEUTRAL)
}
