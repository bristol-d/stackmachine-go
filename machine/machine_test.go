package machine

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func load_program(m *Machine, p []word) {
    for i := 0; i < len(p); i++ {
        m.code[i] = p[i]
    }
}

func TestReset(t *testing.T) {
    var M Machine
    m := &M
    Reset(m)
}

func TestPushPop(t *testing.T) {
    var M Machine
    m := &M
    Reset(m)
    if m.nstack != 0 {
        t.Errorf("stack not empty")
    }
    // push 0x2002; push 0x3003; pop
    load_program(m, []word{0x0001, 0x2002, 0x0001, 0x3003, 0x0002})
	
	step(m)
	assert.Equal(t, word(0x0001), m.nstack)
	assert.Equal(t, word(0x2002), m.stack[m.nstack-1])
	assert.Equal(t, word(2), m.pc)

	step(m)
	assert.Equal(t, word(0x0002), m.nstack)
	assert.Equal(t, word(0x3003), m.stack[m.nstack-1])
	assert.Equal(t, word(4), m.pc)

	step(m)
	assert.Equal(t, word(0x0001), m.nstack)
	assert.Equal(t, word(5), m.pc)
}

func TestAdd(t *testing.T) {
	var M Machine
    m := &M
	Reset(m)
	// push 2; push 3; add
	load_program(m, []word{0x0001, 0x0002, 0x0001, 0x0003, 0x0101})
	step(m)
	step(m)
	step(m)
	assert.Equal(t, word(1), m.nstack)
	assert.Equal(t, word(5), peek(m))
}

func TestShifts(t *testing.T) {
	var M Machine
	m := &M
	Reset(m)

	// 0x0100 >> 0x0001 == 0x0080
	// PUSH 0x0100, PUSH 1, SHR
	load_program(m, []word{0x0001, 0x0100, 0x0001, 0x0001, 0x0206})
	step(m)
	step(m)
	step(m)
	assert.Equal(t, word(1), m.nstack)
	assert.Equal(t, word(0x0080), peek(m))

	// 0xFFFF >> 0x0002 == 0x3FFF
	// PUSH 0xFFFF, PUSH 0x0002, SHR
	Reset(m)
	load_program(m, []word{0x0001, 0xFFFF, 0x0001, 0x0002, 0x0206})
	step(m)
	step(m)
	step(m)
	assert.Equal(t, word(1), m.nstack)
	assert.Equal(t, word(0x3FFF), peek(m))

	// 0xFFFF >>> 0x0002 == 0xFFFF
	// PUSH 0xFFFF, PUSH 0x0002, SSR
	Reset(m)
	load_program(m, []word{0x0001, 0xFFFF, 0x0001, 0x0002, 0x0207})
	step(m)
	step(m)
	step(m)
	assert.Equal(t, word(1), m.nstack)
	assert.Equal(t, word(0xFFFF), peek(m))

	// 0xF000 >>> 0x0004 == 0xFF00
	// PUSH 0xF000, PUSH 0x0004, SSR
	Reset(m)
	load_program(m, []word{0x0001, 0xF000, 0x0001, 0x0004, 0x0207})
	step(m)
	step(m)
	step(m)
	assert.Equal(t, word(1), m.nstack)
	assert.Equal(t, word(0xFF00), peek(m))
}
