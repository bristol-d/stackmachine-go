package main

import (
	"fmt"
	"os"
	// "io"
	"bufio"
	"strings"
)

func usage() {
	fmt.Printf("Usage: %s [assemble]\n", os.Args[0])
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
		return
	}

	if args[0] == "assemble" {
		assemble()
	} else {
		fmt.Printf("Unknown command %s.", args[0])
		usage()
		return
	}

	var M machine
	var m *machine = &M
	reset(m)
}

func assemble() {
	read := bufio.NewReader(os.Stdin)
	fmt.Println("Enter an assembly program, end with a period (.) alone on a line")
	done := false
	lines := []string {}
	for !done {
		line, _ := read.ReadString('\n')
		line = strings.Trim(line, "\r\n")
		if line == "." {
			done = true
		} else {
			lines = append(lines, line)
		}
	}
	program, err := assemble_lines(lines)
	if err != nil {
		fmt.Printf("Error assembling program: %s\n", err.Error())
		return
	}
	fmt.Print("[")
	for i := 0; i < len(program); i++ {
		fmt.Printf("0x%04x ", program[i])
	}
	fmt.Println("]")
}
