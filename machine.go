package main

type word = uint16

type machine struct {
	stack [256] word
	nstack uint16
	retstack [256] word
	nret uint16
	code [4096] word
	data [4096] word
	pc word
	err uint8
}

const (
	OK = iota
	ILLEGAL = iota
	UNDERFLOW = iota
	RET_UNDERFLOW = iota
	OVERFLOW = iota
	RET_OVERFLOW = iota
)

func reset(m *machine) {
	m.pc = 0
	m.err = OK
}

// pre: stack not empty
func pop (m *machine) word {
	if m.nstack == 0 {
		panic("pop on empty stack")
	}
	var x word = m.stack[m.nstack+1]
	m.nstack--
	return x
}

func push(m *machine, x word) uint8 {
	if m.nstack > 255 {
		return OVERFLOW
	}
	m.nstack++
	m.stack[m.nstack] = x
	return OK
}

func add (m *machine) uint8 {
	if m.nstack < 2 {
		m.err = UNDERFLOW
		return m.err
	}
	var y word = pop(m)
	var x word = pop(m)
	x = x + y
	return push(m, x)
}
