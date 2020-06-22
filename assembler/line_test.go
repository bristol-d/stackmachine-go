package assembler

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConsumeLabel(t *testing.T) {
	s1 := "start:"
	label, rest := consume_label([]rune(s1))
	assert.Equal(t, "start", string(label))
	assert.Equal(t, "", string(rest))

	s2 := "start: rest"
	label, rest = consume_label([]rune(s2))
	assert.Equal(t, "start", string(label))
	assert.Equal(t, " rest", string(rest))

	s3 := "nolabel"
	label, rest = consume_label([]rune(s3))
	assert.Nil(t, label)
	assert.Equal(t, "nolabel", string(rest))

}

func TestConsumeSpaces(t *testing.T) {
	c, s := consume_spaces([]rune("word"))
	assert.False(t, c)
	assert.Equal(t, "word", string(s))

	c, s = consume_spaces([]rune("  word"))
	assert.True(t, c)
	assert.Equal(t, "word", string(s))

	c, s = consume_spaces([]rune("\tword"))
	assert.True(t, c)
	assert.Equal(t, "word", string(s))
}

func TestConsumeOpcode(t *testing.T) {
	opcode, rest := consume_opcode([]rune("PUSH"))
	assert.Equal(t, "PUSH", string(opcode))
	assert.Equal(t, "", string(rest))
	
	opcode, rest = consume_opcode([]rune("PUSH #20"))
	assert.Equal(t, "PUSH", string(opcode))
	assert.Equal(t, " #20", string(rest))

	opcode, rest = consume_opcode([]rune(" PUSH"))
	assert.Nil(t, opcode)
	assert.Equal(t, " PUSH", string(rest))
}

func TestConsumeOperand(t *testing.T) {
	operand, rest := consume_operand([]rune ("#1"))
	assert.Equal(t, []rune("#1"), operand)
	assert.Equal(t, []rune(""), rest)

	operand, rest = consume_operand([]rune ("#1234"))
	assert.Equal(t, []rune("#1234"), operand)
	assert.Equal(t, []rune(""), rest)

	operand, rest = consume_operand([]rune ("#12 "))
	assert.Equal(t, []rune("#12"), operand)
	assert.Equal(t, []rune(" "), rest)
}

func TestParseLabel(t *testing.T) {
	d, e := parse_line("start:")
	assert.Nil(t, e)
	assert.Equal(t, "start", *d.label)
	assert.Equal(t, 0, d.len)
}

func TestParsePop(t *testing.T) {
	d, e := parse_line("POP")
	assert.Nil(t, e)
	assert.Nil(t, d.label)
	assert.Equal(t, 1, d.len)
	assert.Equal(t, word(0x0002), d.opcode)
}

func TestParsePushConstant(t *testing.T) {
	d, e := parse_line("PUSH #20a2")
	assert.Nil(t, e)
	assert.Nil(t, d.label)
	assert.Equal(t, 2, d.len)
	assert.Equal(t, word(0x0001), d.opcode)
	assert.Equal(t, word(0x20A2), d.operand)
	assert.Nil(t, d.reference)
}

func TestParsePushReference(t *testing.T) {
	d, e := parse_line("PUSH start")
	assert.Nil(t, e)
	assert.Nil(t, d.label)
	assert.Equal(t, 2, d.len)
	assert.Equal(t, word(0x0001), d.opcode)
	assert.Equal(t, word(0x0000), d.operand)
	assert.Equal(t, "start", *d.reference)
}

func TestPopWithLabel(t *testing.T) {
	d, e := parse_line("start: POP")
	assert.Nil(t, e)
	assert.Equal(t, "start", *d.label)
	assert.Equal(t, 1, d.len)
	assert.Equal(t, word(0x0002), d.opcode)
}

func TestParsePushConstantLabel(t *testing.T) {
	d, e := parse_line("start: PUSH #20a2")
	assert.Nil(t, e)
	assert.Equal(t, "start", *d.label)
	assert.Equal(t, 2, d.len)
	assert.Equal(t, word(0x0001), d.opcode)
	assert.Equal(t, word(0x20A2), d.operand)
	assert.Nil(t, d.reference)
}

func TestParsePushReferenceLabel(t *testing.T) {
	d, e := parse_line("start: PUSH start")
	assert.Nil(t, e)
	assert.Equal(t, "start", *d.label)
	assert.Equal(t, 2, d.len)
	assert.Equal(t, word(0x0001), d.opcode)
	assert.Equal(t, word(0x0000), d.operand)
	assert.Equal(t, "start", *d.reference)
}
