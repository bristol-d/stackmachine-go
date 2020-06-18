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
	m.nstack = 0
	m.nret = 0
	m.err = OK
}

// pre: stack not empty
// this is the internal helper function
func _pop (m *machine) word {
	if m.nstack == 0 {
		panic("pop on empty stack")
	}
	var x word = m.stack[m.nstack+1]
	m.nstack--
	return x
}

// this is an internal helper, not the instruction
func _push(m *machine, x word) {
	if m.nstack > 255 {
		panic("stack overflow")
	}
	m.stack[m.nstack] = x
	m.nstack++
}

func peek(m *machine) word {
	if m.nstack == 0 {
		panic("peek on empty stack")
	}
	return m.stack[m.nstack - 1]
}

func step (m *machine) {
	// fetch
	var instruction word = m.code[m.pc]
	m.pc++
	// decode
	var f, err = decode(instruction)
	if err != OK {
		m.err = ILLEGAL
		return
	}
	// execute
	result := f(m)
	if result != OK {
		m.err = result
	}
}

// instructions //

func add (m *machine) uint8 {
	if m.nstack < 2 {
		m.err = UNDERFLOW
		return m.err
	}
	var y word = _pop(m)
	var x word = _pop(m)
	x = x + y
	// cannot fail, as we just popped two
	_push(m, x)
	return OK
}

// this is the push instruction
func push (m *machine) uint8 {
	if m.nstack > 255 {
		m.err = OVERFLOW
		return m.err
	}
	var x word = m.code[m.pc]
	m.pc++
	// this is safe, as checked above
	_push(m, x)
	return OK
}

// this is the instruction
func pop (m *machine) uint8 {
	if m.nstack == 0 {
		m.err = UNDERFLOW
		return m.err
	}
	m.nstack--
	return OK
}

// the decoding table //

var INSTRUCTIONS = map[word] func(*machine) uint8 {
	0x0001: push,
	0x0002: pop,

	0x0010: add,
}

func decode(instruction word) (func(*machine) uint8, uint8) {
	f, exists := INSTRUCTIONS[instruction]
	if exists == false {
		return nil, ILLEGAL
	}
	return f, OK
}
