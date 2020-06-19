package main

import (
	"fmt"
	"../machine"
	"../assembler"
	"strings"
	"syscall/js"
)

var m *machine.Machine

func main() {
	var M = machine.Machine {}
	m = &M
	machine.Reset(m)

	js.Global().Set("assemble", js.FuncOf(assemble))

	fmt.Println("Loaded stack machine.")
	c := make(chan struct{}, 0)
	<- c
}

func assemble (this js.Value, input []js.Value) interface{} {
	source := input[0].String()
	lines := strings.Split(source, "\n")
	_, err := assembler.Assemble_lines(lines)
	if err != nil {
		return js.ValueOf(err.Error())
	}
	return js.ValueOf("OK")
}