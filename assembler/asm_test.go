package assembler

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

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
		0x0010,
		0x0002,
	}

	code, err := Assemble_lines(source)
	assert.Nil(t, err)
	assert.Equal(t, len(bin), len(code))

	for i := 0; i < len(bin); i++ {
		assert.Equal(t, bin[i], code[i], "Word #%u", i)
	}
}