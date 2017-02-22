package fsm

func (e *someprocessname) AddElev(message Message, SendToNet chan Message){
	e.elevators[message.Source] = new(Elevator)
	e.elevators[message.Source].Dir = DIR_DOWN
	fmt.Println("Elevator ", message.Source, " added to system")
	syncOrders() 		// sync external and internal orders
						// how to differentiate between dead elev/lost net/new elev?
}

func (e *someprocessname) DeleteElev(id int, to_network chan Message) { // deletes elev if dead
	delete(e.elevators, id)
	fmt.Println("Elevator ", id, " is removed from the system")
						// put some sort of iteration through orderarray here?
	RecalculateOrders() // have to redistribute orders from dead elev to remaining elevs
}