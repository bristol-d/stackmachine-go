package assembler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type word = uint16

var TABLE = map[string] (struct {opcode word; argument bool}) {
	// format is (opcode, takes_argument)
	"PUSH": {0x0001, true},
	"POP" : {0x0002, false},
	"ADD" : {0x0010, false},
	"SUB" : {0x0012, false},
	"MUL" : {0x0013, false},
	"AND" : {0x0017, false},
	"OR"  : {0x0018, false},
	"XOR" : {0x0019, false},
	"NAND": {0x001A, false},
}

type linedata struct {
	label *string
	len int
	opcode word
	operand word
	reference *string
} 

func is_capital (r rune) bool {
	if r >= 'A' && r <= 'Z' { return true }
	return false
}

func is_alpha(r rune) bool {
	if r >= 'A' && r <= 'Z' { return true }
	if r >= 'a' && r <= 'z' { return true }
	if r == '_' { return true }
	return false
}

func is_alphanum(r rune) bool {
	if is_alpha(r) { return true }
	if r >= '0' && r <= '9' { return true }
	return false
}

// Consume a leading [A-Za-z_][A-Za-z_0-9]*[:]
// return (label, rest) where label may be nil
func consume_label(source []rune) ([]rune, []rune) {
	s := source
	if len(s) == 0 || !is_alpha(s[0]) { return nil, s }
	count := 1
	s = s[1:]
	for len(s) > 0 {
		if !is_alphanum(s[0]) { break }
		s = s[1:]
		count += 1
	}
	if len(s) > 0 && s[0] == ':' {
		return source[0:count], s[1:]
	} else {
		return nil, source
	}
}

func consume_operand(source []rune) ([]rune, []rune) {
	s := source
	if len(s) == 0 || (!is_alpha(s[0]) && s[0] != '#') { return nil, s }
	count := 1
	s = s[1:]
	for len(s) > 0 {
		if !is_alphanum(s[0]) { break }
		s = s[1:]
		count += 1
	}
	return source[0:count-1], s
}

func consume_spaces(source []rune) (bool, []rune) {
	change := false
	for len(source) > 0 {
		if source[0] == ' ' || source[0] == '\n' || source[0] == '\t' || source[0] == '\r' {
			source = source[1:]
			change = true
		} else {
			break
		}
	}
	return change, source
}

func consume_opcode(source []rune) ([]rune, []rune) {
	original := source
	if len(source) == 0 || !is_capital(source[0]) {
		return nil, source
	}
	count := 1
	for len(source) > 0 {
		if !is_capital(source[0]) { break }
		count += 1
		source = source[1:]
	}
	return original[0:count-1], source
}

func Assemble_lines(lines []string) ([]word, error) {
	offset := word(0)
	code := []word {}
	var labels = map[string] word {}
	var refs = map[word] string {}
	// first pass
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		line = strip_line(line)
		if line == "" { continue }
		d, e := parse_line(line)
		if e != nil {
			return nil, fmt.Errorf("Line %d: %s", i+1, e.Error())
		}
		if d.label != nil {
			labels[*d.label] = offset
		}
		if d.len > 0 {
			code = append(code, d.opcode)
			offset++
		}
		if d.len > 1 {
			code = append(code, d.operand)
			if d.reference != nil {
				// the operand is currently an unresolved reference
				// save it in the map for the second pass
				refs[offset] = *d.reference
			}
			offset++
		}
	}
	// second pass
	for k, v := range refs {
		dest, exists := labels[v]
		if !exists {
			return nil, fmt.Errorf("Label %s referenced but never defined.", v)
		}
		code[k] = dest 
	}
	return code, nil
}

func strip_line(line string) string {
	comment := strings.Index(line, ";")
	if comment != -1 {
		line = line[0:comment]
	}
	line = strings.Trim(line, " \t\n\r")
	return line
}

func parse_line(line string) (linedata, error) {
	runes := []rune(line)
	var label, opcode, operand []rune
	var space bool
	label, runes = consume_label(runes)
	_, runes = consume_spaces(runes)
	opcode, runes = consume_opcode(runes)
	space, runes = consume_spaces(runes)
	operand, runes = consume_operand(runes) 

	if opcode == nil {
		return linedata{}, errors.New("Incorrect line format or no opcode found.")
	}

	var data = linedata{}
	if label == nil {
		data.label = nil
	} else {
		str := string(label)
		data.label = &str
	}
	opstr := string(opcode)
	op, ok := TABLE[opstr]
	if !ok {
		return linedata{}, fmt.Errorf("Illegal command: %s.", opstr)
	}
	data.opcode = op.opcode
	if operand == nil {
		if op.argument {
			return linedata{}, fmt.Errorf("Instruction %s requires an argument.", opstr)
		} 
	} else {
		if !op.argument {
			return linedata{}, fmt.Errorf("Instruction %s does not take an argument.", opstr)
		}
		if!space {
			return linedata{}, errors.New("No space after opcode.")
		}
		if operand[0] == '#' {
			num, err := strconv.ParseUint(string(operand[1:]), 16, 16)
			if err != nil {
				return linedata{}, errors.New("Illegal constant operand.")
			}
			data.operand = word(num)
		} else {
			// it's a reference to a label
			ref := string(operand)
			data.reference = &ref
		}
	}

	return data, nil
}
