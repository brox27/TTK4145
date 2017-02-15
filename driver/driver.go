package driver

import (
	. "../globals"
	. "fmt"
)

type Event int

const (
	ButtonPressed Event = iota
)

func ButtonCheck(ch chan ButtonType) {
	for {
		pressedButtons := []ButtonType{}
		for i := 0; i < N_BUTTONS; i++ {
			for j := 1; j <= N_FLOORS; j++ {
				if GetButtonSignal(ButtonType(i), j) == 1 {
					Println(ButtonType(i))

					pressedButtons = append(pressedButtons, ButtonType(i))
				}
			}
		}
		for i := range pressedButtons {
			ch <- pressedButtons[i]
		}
	}
}
