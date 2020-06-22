package machine

import "fmt"

type word = uint16

type Machine struct {
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
	VALUE_ERROR = iota // e.g. division by 0
)

func Reset(m *Machine) {
	m.pc = 0
	m.nstack = 0
	m.nret = 0
	m.err = OK
}

// pre: stack not empty
// this is the internal helper function
func _pop (m *Machine) word {
	if m.nstack == 0 {
		panic("pop on empty stack")
	}
	m.nstack--
	var x word = m.stack[m.nstack]
	return x
}

// this is an internal helper, not the instruction
func _push(m *Machine, x word) {
	if m.nstack > 255 {
		panic("stack overflow")
	}
	m.stack[m.nstack] = x
	m.nstack++
}

func peek(m *Machine) word {
	if m.nstack == 0 {
		panic("peek on empty stack")
	}
	return m.stack[m.nstack - 1]
}

func dump (m *Machine) {
	fmt.Printf("M pc=%04x stack=%d\n", m.pc, m.nstack)
	fmt.Printf("  next=%04x\n", m.code[m.pc])
	if m.nstack > 0 {
		var min word = 1
		if m.nstack > 2 {
			min = m.nstack - 2
		}
		for i := m.nstack; i >= min; i-- {
			fmt.Printf("  |%04x|\n", m.stack[i-1])
		}
		fmt.Printf("  ------\n")
	}
}

func step (m *Machine) {
	// fetch
	var instruction word = m.code[m.pc]
	// decode
	var f, err = decode(instruction)
	if err != OK {
		m.err = ILLEGAL
		return
	}
	m.pc++
	// execute
	result := f(m)
	if result != OK {
		m.err = result
	}
}

// instructions //

func add (m *Machine) uint8 {
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
func push (m *Machine) uint8 {
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
func pop (m *Machine) uint8 {
	if m.nstack == 0 {
		m.err = UNDERFLOW
		return m.err
	}
	m.nstack--
	return OK
}

func dup (m *Machine) uint8 {
	if m.nstack == 0 {
		m.err = UNDERFLOW
		return m.err
	}
	if m.nstack > 255 {
		m.err = OVERFLOW
		return m.err
	}
	w := peek(m)
	_push(m, w)
	return OK
}

func swap (m *Machine) uint8 {
	if m.nstack < 2 {
		m.err = UNDERFLOW
		return m.err
	}
	x := _pop(m)
	y := _pop(m)
	_push(m, x)
	_push(m, y)
	return OK
}

func binary_operation(f func(word, word) word) func(*Machine) uint8 {
	return func(m *Machine) uint8 {
		if m.nstack < 2 {
			m.err = UNDERFLOW
			return m.err
		}
		var y word = _pop(m)
		var x word = _pop(m)
		var z word = f(x, y)
		// cannot fail, as we just popped two
		_push(m, z)
		return OK
	}
} 

func binary_not(m *Machine) uint8 {
	if m.nstack < 1 {
		m.err = UNDERFLOW
		return m.err
	}
	x := _pop(m)
	x = ^x
	_push(m, x)
	return OK
}

// the decoding table //

var INSTRUCTIONS = map[word] func(*Machine) uint8 {
	0x0001: push,
	0x0002: pop,
	0x0003: dup,
	0x0004: swap,

	0x0101: add,
	//0x0011: addc,
	0x0103: binary_operation(func(x, y word) word {return x - y}),
	0x0104: binary_operation(func(x, y word) word {return x * y}),
	//0x0014: muld,
	//0x0015: mod,
	//0x0016: div,
	0x0201: binary_operation(func(x, y word) word {return x & y}),
	0x0202: binary_operation(func(x, y word) word {return x | y}),
	0x0203: binary_operation(func(x, y word) word {return x ^ y}),
	0x0204: binary_operation(func(x, y word) word {return ^(x & y)}),
	0x0205: binary_not,
}

func decode(instruction word) (func(*Machine) uint8, uint8) {
	f, exists := INSTRUCTIONS[instruction]
	if exists == false {
		return nil, ILLEGAL
	}
	return f, OK
}
