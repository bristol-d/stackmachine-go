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
	"HALT": {0x0000, false},

	"PUSH": {0x0001, true},
	"POP" : {0x0002, false},
	"DUP" : {0x0003, false},
	"SWAP": {0x0004, false},

	"ADD" : {0x0101, false},
	"SUB" : {0x0103, false},
	"MUL" : {0x0104, false},
	"MULC": {0x0105, false},
	"DIV" : {0x0106, false},
	"MOD" : {0x0107, false},

	"AND" : {0x0201, false},
	"OR"  : {0x0202, false},
	"XOR" : {0x0203, false},
	"NAND": {0x0204, false},
	"NOT" : {0x0205, false},
	"SHR" : {0x0206, false},
	"SSR" : {0x0207, false},
	"SHL" : {0x0208, false},
	"SWE" : {0x0209, false},

	"CEQ" : {0x0301, false},
	"CNE" : {0x0302, false},
	"CGT" : {0x0303, false},
	"CGE" : {0x0304, false},
	"CLT" : {0x0305, false},
	"CLE" : {0x0306, false},
	"TZ"  : {0x0310, false},
	"TN"  : {0x0311, false},
	"TM"  : {0x0312, false},
	"TL"  : {0x0313, false},

	"J"   : {0x0401, true},
	"JS"  : {0x0402, false},
	"JT"  : {0x0403, true},
	"JTS" : {0x0404, false},
	"JF"  : {0x0405, true},
	"JFS" : {0x0406, false},

	"LOAD": {0x0501, true},
	"STOR": {0x0502, true},
	"LODS": {0x0503, false},
	"STRS": {0x0504, false},

	"INT" : {0xffff, false},
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

func is_digit(r rune) bool {
	return r >= '0' && r <= '9'
}

func is_xdigit(r rune) bool {
	return is_digit(r) || (r >= 'A' && r <= 'F')
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
	return source[0:count], s
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

/// Consume an int in decimal notation.
func consume_int(r []rune) (*uint16, []rune) {
	var x uint16 = 0
	found := false
	for len(r) > 0 && r[0] >= '0' && r[0] <= '9' {
		found = true
		x = 10*x + uint16(r[0] - '0')
		r = r[1:]
	}
	if found {
		return &x, r
	} else {
		return nil, r
	}
}

/// Consume a one-character token, if it exists.
/// Returns true/false to indicate if the token was found,
/// and the rest of the stream.
func consume_token1(token rune, r []rune) (bool, []rune) {
	if len(r) > 0 && r[0] == token {
		return true, r[1:]
	} else {
		return false, r
	}
}

func parse_xdigit(r rune) word {
	if is_digit(r) {
		return word(r - '0')
	} else {
		return word(r - 'A' + 10)
	}
}

/// Consume a hexadecimal constant in the format [0-9A-F]{1,4}
func consume_constant(r []rune) (*word, []rune) {
	if len(r) == 0 || !is_xdigit(r[0]) {
		return nil, r
	}
	var value word = parse_xdigit(r[0])
	r = r[1:]
	if len(r) > 0 && is_xdigit(r[0]) {
		value *= 16
		value += parse_xdigit(r[0])
		r = r[1:]
	}
	if len(r) > 0 && is_xdigit(r[0]) {
		value *= 16
		value += parse_xdigit(r[0])
		r = r[1:]
	}
	if len(r) > 0 && is_xdigit(r[0]) {
		value *= 16
		value += parse_xdigit(r[0])
		r = r[1:]
	}
	return &value, r
}

/// Parse a line in data mode.
/// Format is LABEL: [length] = [value [value]...]
/// Default length is 1
func parse_data_line(line string) (string, uint16, []word) {
	runes := []rune(line)
	var data []word
	label, rest := consume_label(runes)
	if label == nil {
		return "", 0, data
	}
	_, rest = consume_spaces(rest)
	length, rest := consume_int(rest)
	var l uint16 = 0
	if length == nil {
		l = 1
	} else {
		l = *length
	}

	_, rest = consume_spaces(rest)
	has_data, rest := consume_token1('=', rest)
	if !has_data {
		return string(label), l, data
	}

	_, rest = consume_spaces(rest)
	var c *word
	var change bool = false

	for i := word(0); i < l; i++ {
		c, rest = consume_constant(rest)
		if c != nil {
			data = append(data, *c)
		} else {
			break
		}
		change, rest = consume_spaces(rest)
		if !change && len(rest) > 0 {
			// something after a constant that's not a separator - error
			return "", 0, data
		}
	}

	return string(label), l, data
}

func Assemble_lines(lines []string) ([]word, []word, error) {
	offset := word(0)
	code := []word {}
	data := []word {}
	var labels = map[string] word {}
	var refs = map[word] string {}
	// first pass
	data_mode := false
	var data_index uint16 = 0
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		line = strip_line(line)
		if line == "" { continue }
		if data_mode {
			label, length, d := parse_data_line(line)
			if length > 0 {
				labels[label] = data_index
				data_index += length
				for i := 0; i < len(d); i++ {
					data = append(data, d[i])
				}
				if len(d) < int(length) {
					for i := 0; i < int(length) - len(d); i++ {
						data = append(data, word(0))
					}
				}
			} else {
				return nil, nil, fmt.Errorf("Line %d: invalid line in data section", i+1)
			}
		} else {
			if strings.HasPrefix(line, ".DATA") {
				data_mode = true
				continue
			}
			d, e := parse_line(line)
			if e != nil {
				return nil, nil, fmt.Errorf("Line %d: %s", i+1, e.Error())
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
	}
	// second pass
	for k, v := range refs {
		dest, exists := labels[v]
		if !exists {
			return nil, nil, fmt.Errorf("Label %s referenced but never defined.", v)
		}
		code[k] = dest 
	}
	return code, data, nil
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

	if opcode == nil && label == nil {
		return linedata{}, errors.New("Incorrect line format or no opcode found.")
	}

	var data = linedata{}
	if label == nil {
		data.label = nil
	} else {
		str := string(label)
		data.label = &str
	}

	if opcode == nil {
		return data, nil
	}

	opstr := string(opcode)
	op, ok := TABLE[opstr]
	if !ok {
		return linedata{}, fmt.Errorf("Illegal command: %s.", opstr)
	}
	data.opcode = op.opcode
	data.len = 1
	if operand == nil {
		if op.argument {
			return linedata{}, fmt.Errorf("Instruction %s requires an argument.", opstr)
		} 
	} else {
		data.len = 2
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

// numbering==false disassembles only a single instruction
func Disassemble(code []word, numbering bool) []string {
	text := []string {}

	for count := 0; count < len(code); count++ {
		m := "??"
		arg := false
		for name, op := range TABLE {
			if code[count] == op.opcode {
				m = name
				for len(m) < 4 {
					m += " "
				}
				if op.argument {
					arg = true
					count++
					m = m + " #" + fmt.Sprintf("%x", code[count])
				}
				break
			}
		}
		if arg {
			m = fmt.Sprintf("%04x %04x %s", code[count-1], code[count], m)
		} else {
			m = fmt.Sprintf("%04x      %s", code[count], m)
		}
		if (numbering) {
			pos := count
			if arg {
				pos--
			}
			m = fmt.Sprintf("0x%04x: %s", pos, m)
		}
		text = append(text, m)
		if !numbering {
			return text
		}
	}
	return text
}