package assembler

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAssembleNoOperand(t *testing.T) {
	var source = []string {"ADD"}
		var bin = []word {0x0101}
	code, err := Assemble_lines(source)
	assert.Nil(t, err)
	assert.Equal(t, bin, code)
}

func TestAssembleSimple(t *testing.T) {
	var source = []string  {
		"; test program",
		"  PUSH #1234",
		"  PUSH #face",
		"  ADD",
		"  POP",
	}

	var bin = []word {
		0x0001, 0x1234,
		0x0001, 0xface,
		0x0101,
		0x0002,
	}

	code, err := Assemble_lines(source)
	assert.Nil(t, err)
	assert.Equal(t, len(bin), len(code))

	for i := 0; i < len(bin); i++ {
		assert.Equal(t, bin[i], code[i], "Word #%u", i)
	}
}

func TestDisassemble (t *testing.T) {
	code := []word {0x0002}
	text := Disassemble(code, false)
	assert.Equal(t, []string{"0002      POP "}, text)
}

func TestDisassembleOperand (t *testing.T) {
	code := []word {0x0001, 0x0002}
	text := Disassemble(code, false)
	assert.Equal(t, []string{"0001 0002 PUSH #2"}, text)
}

func TestConsumeConstant(t *testing.T) {
	s := "3"
	r := []rune(s)
	w, rest := consume_constant(r)
	assert.Equal(t, word(3), *w)
	assert.Equal(t, "", string(rest))

	s = ""
	r = []rune(s)
	w, rest = consume_constant(r)
	assert.Nil(t, w)	
}

func TestParseDataLine(t *testing.T) {
	label, length, data := parse_data_line("main:")
	assert.Equal(t, "main", label)
	assert.Equal(t, word(1), length)
	assert.Equal(t, 0, len(data))

	label, length, data = parse_data_line("a: 2")
	assert.Equal(t, "a", label)
	assert.Equal(t, word(2), length)
	assert.Equal(t, 0, len(data))

	label, length, data = parse_data_line("b: = 4")
	assert.Equal(t, "b", label)
	assert.Equal(t, word(1), length)
	assert.Equal(t, 1, len(data))
	assert.Equal(t, word(4), data[0])

	label, length, data = parse_data_line("c: 2 = 3")
	assert.Equal(t, "c", label)
	assert.Equal(t, word(2), length)
	assert.Equal(t, 1, len(data))
	assert.Equal(t, word(3), data[0])

	label, length, data = parse_data_line("d: 2 = 4 5")
	assert.Equal(t, "d", label)
	assert.Equal(t, word(2), length)
	assert.Equal(t, 2, len(data))
	assert.Equal(t, word(4), data[0])
	assert.Equal(t, word(5), data[1])

	label, length, data = parse_data_line("ee: = 0ABC")
	assert.Equal(t, "ee", label)
	assert.Equal(t, word(1), length)
	assert.Equal(t, 1, len(data))
	assert.Equal(t, word(0x0abc), data[0])
}
