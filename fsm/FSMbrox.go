package main

import (
	. "../driver2"
//	. "../globals"
//	"../OrderManager"
	. "../ConfigFile"
	//. "fmt"
//	"runtime"
	//"time"
)

func main() {
//	println(" ░░░░░░░░░░░░░░░░░░░░░ \n ░░░░░░░░░░░░▄▀▀▀▀▄░░░ \n ░░░░░░░░░░▄▀░░▄░▄░█░░ \n ░▄▄░░░░░▄▀░░░░▄▄▄▄█░░ \n █░░▀▄░▄▀░░░░░░░░░░█░░ \n ░▀▄░░▀▄░░░░█░░░░░░█░░ \n ░░░▀▄░░▀░░░█░░░░░░█░░ \n ░░░▄▀░░░░░░█░░░░▄▀░░░ \n ░░░▀▄▀▄▄▀░░█▀░▄▀░░░░░ \n ░░░░░░░░█▀▀█▀▀░░░░░░░ \n ░░░░░░░░▀▀░▀▀░░░░░░░░ \n")
//	runtime.GOMAXPROCS(runtime.NumCPU())
//	localId := "123.123"
//	InitElev()
//	eventChan := make(chan Event)
//	go EventHandler(eventChan)
//	RUN(localId, eventChan)
	Lala()
}


func RUN(localId string, eventChan chan map[Event]interface{}) {	// mulig localID kan/bør ligge i Config?
	State := INITIALIZE
	for {
		select{
		//case eventMap := <- eventChan:
			//update based on what arrived
		//	println("lolz")

		default:
			switch(State){

			case INITIALIZE:
				if(GetFloorSensorSignal()!=-1){
					SetMotorDirection(NEUTRAL)
					SetFloorLight(GetFloorSensorSignal())
					// oppdatere egen dir/floor
					// synkronisere ordre mot andre på nettverket!

					// WaitGroup e.l.? for å være sikker på klar
					State=IDLE

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
	}

	defer SetMotorDirection(NEUTRAL)
}
