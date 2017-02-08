package main

import (
	. "../ConfigFile"
)

func main() {
	elev1 := Elev{}
	AllOrders = make(map[string]*Elev)
	AllOrders["123.123.13.123"] = &elev1
}
