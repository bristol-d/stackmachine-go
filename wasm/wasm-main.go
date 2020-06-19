package main

import (
	"fmt"
	"../machine"
)

var m *machine.Machine

func main() {
	var M = machine.Machine {}
	m = &M
	machine.Reset(m)
	fmt.Println("Loaded stack machine.")
}