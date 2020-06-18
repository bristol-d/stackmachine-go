package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

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
