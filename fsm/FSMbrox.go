package main

import (
	. "../driver2"
	//	. "../globals"
	. "../ConfigFile"
	. "../OrderMAnager"
	//. "fmt"
	"runtime"
	//"time"
)

func main() {
	InitElev()
	println(" ░░░░░░░░░░░░░░░░░░░░░ \n ░░░░░░░░░░░░▄▀▀▀▀▄░░░ \n ░░░░░░░░░░▄▀░░▄░▄░█░░ \n ░▄▄░░░░░▄▀░░░░▄▄▄▄█░░ \n █░░▀▄░▄▀░░░░░░░░░░█░░ \n ░▀▄░░▀▄░░░░█░░░░░░█░░ \n ░░░▀▄░░▀░░░█░░░░░░█░░ \n ░░░▄▀░░░░░░█░░░░▄▀░░░ \n ░░░▀▄▀▄▄▀░░█▀░▄▀░░░░░ \n ░░░░░░░░█▀▀█▀▀░░░░░░░ \n ░░░░░░░░▀▀░▀▀░░░░░░░░ \n")
	runtime.GOMAXPROCS(runtime.NumCPU())
	localId := "123.123"

	RUN(localId)
}

func RUN(localId string) { // mulig localID kan/bør ligge i Config?, eventChan chan map[Event]interface{}
	State := INITIALIZE
	eventChan := make(chan Event)
	go EventHandler(eventChan)
	hest := GetFloorSensorSignal()
	println(hest)

	EventFloor := false
	for {
		select {
		case LOLL := <-eventChan:
			//update based on what arrived
			println("lolz noe skjedde!")
			if LOLL.EventType == NEWFLOOR {
				println("new floor motherfucker!")
				// OPPDATERE egen variable på FLOOR!
				EventFloor = true
			} else {
				println("new button motherfucker")
			}

		default:
			switch State {

			case INITIALIZE:
				if GetFloorSensorSignal() != -1 {
					SetMotorDirection(NEUTRAL)
					SetFloorLight(GetFloorSensorSignal())
					// oppdatere egen dir/floor
					// synkronisere ordre mot andre på nettverket!

					// WaitGroup e.l.? for å være sikker på klar
					// er det en ide å starte eventHandler her..? så vi ikke tar inn knappetrykk eller no shit før klar?
					State = IDLE

					break
				} else {
					SetMotorDirection(DOWN) // sjekke noe så den ikke "settes" så ofte?
				}
				break

			case IDLE:
				nextDir := DOWN // HARDCODED !!___________________________________________________________________________________ <-- SE!
				if nextDir != NEUTRAL {
					SetMotorDirection(nextDir)
					State = RUNNING
				}
				break

			case RUNNING:
				if EventFloor {
					SetFloorLight(GetFloorSensorSignal() - 1)
					EventFloor = false

					//HVIS SKAL STOPPE
					if OrderManager.ShouldStop {
						SetMotorDirection(NEUTRAL)
						State = DOORSOPEN
					}
					// oppdatere ordreliste
					// si OrderComplete til andre
					break
				}

				//event arrived at floor
				//	Oppdatere ny etg!
				//	Oppdatere lys for etg.
				//	if OrderManager.Shouldstop(){
				//		SetMotorDirection(NEUTRAL)
				//oppdatere direction?
				//sette lys for dør åpen
				// starte timer
				//oppdatere orders lista
				// si ordre utført! til andre

			case DOORSOPEN:

				//if TIMEOUT
				//	Skru av timer?
				//	Sett lys dør av
				//	nextDir = OrderManager.GetNextOrder()
				//	if nextDir == NEUTRAL{
				//		elev.State=IDLE
				//	}else{
				//		setMotorDir(nextDir)
				//		elev.State=RUNNING
				//	}
			}
		}
	}

	defer SetMotorDirection(NEUTRAL)
}
