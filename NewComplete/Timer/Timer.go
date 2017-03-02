package main

import (
	"fmt"
	"time"
)

func main() {
	// DOORS OPEN
	fmt.Printf("Lolz \n")
	timer := time.NewTimer(3*time.Second)
	go func (){
		<- timer.C
		println("OUT!")
	}()

	// Spam på nettverk
	ticker := time.NewTicker(time.Millisecond * 20)
	go func (){
		for t:= range ticker.C{
			fmt.Printf("lala LAND %d\n", t)
		}
	}()
	for{
		// så ikke programmet stopper og kræsjer GO-rutninen
	}
}