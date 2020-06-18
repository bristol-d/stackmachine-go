package main

import "fmt"

func main() {
	var M machine
	var m *machine = &M
	reset(m)

	fmt.Println("main")
}