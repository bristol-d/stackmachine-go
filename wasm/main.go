package main

import (
	"fmt"
	"go-stackmachine/machine"
	"go-stackmachine/assembler"
	"strings"
	"syscall/js"
	"encoding/json"
)

var m *machine.Machine
var data []uint16 // global so it survives resets of the machine

func main() {
	var M = machine.Machine {}
	m = &M
	machine.Reset(m)

	js.Global().Set("assemble", js.FuncOf(assemble))
	js.Global().Set("reset_simulation", js.FuncOf(reset_simulation))
	js.Global().Set("step_simulation", js.FuncOf(step_simulation))
	js.Global().Set("dump_simulation", js.FuncOf(dump_simulation))

	fmt.Println("Loaded stack machine.")
	c := make(chan struct{}, 0)
	<- c
}

func dump_simulation(this js.Value, input []js.Value) interface{} {
	d := machine.Dump(m)
	j, _ := json.Marshal(d)
	return js.ValueOf(string(j))
}

func step_simulation(this js.Value, input []js.Value) interface{} {
	machine.Step(m)
	var dump = machine.Dump(m)
	j, err := json.Marshal(dump)
	s := string(j)
	if err != nil {
		return js.ValueOf("Error encoding JSON: " + err.Error())
	}
	return js.ValueOf(s)
}

func reset_simulation(this js.Value, input []js.Value) interface{} {
	machine.Reset(m)
	machine.LoadData(m, data)
	return js.ValueOf(true)
}

func assemble (this js.Value, input []js.Value) interface{} {
	source := input[0].String()
	lines := strings.Split(source, "\n")
	code, d, err := assembler.Assemble_lines(lines)
	if err != nil {
		return js.ValueOf("ERROR: " + err.Error())
	}

	text := assembler.Disassemble(code, true)
	machine.Load(m, code)
	data = d
	machine.LoadData(m, d)

	return js.ValueOf(strings.Join(text, "<br />"))
}