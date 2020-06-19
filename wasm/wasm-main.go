package main

import (
	"fmt"
	"../machine"
	"../assembler"
	"strings"
)

var m *machine.Machine

func main() {
	var M = machine.Machine {}
	m = &M
	machine.Reset(m)
	fmt.Println("Loaded stack machine.")
}

// go:export assemble
func assemble (source string) string {
	lines := strings.Split(source, "\n")
	_, err := assembler.Assemble_lines(lines)
	if err != nil {
		return err.Error()
	}
	return "OK"
}