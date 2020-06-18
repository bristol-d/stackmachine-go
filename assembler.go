package main

import (
	"regexp"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var TABLE = map[string] (struct {opcode word; argument bool}) {
	// format is (opcode, takes_argument)
	"PUSH": {0x0001, true},
	"POP" : {0x0002, false},
	"ADD" : {0x0010, false},
}

type linedata struct {
	label *string
	len int
	opcode word
	operand word
	reference *string
} 

var LINE_RE = regexp.MustCompile("(([a-zA-Z_][a-zA-Z0-9_]+):)?[ \t]*([A-Z]+)?[ \t]*([0-9a-zA-Z_#]+)?")
var CONST_RE = regexp.MustCompile("#([0-9A-Fa-f]+)")
var REF_RE = regexp.MustCompile("([A-Za-z_][A-Za-z0-9_]+)")

func assemble_lines(lines []string) ([]word, error) {
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
	comment := strings.Index(line, "//")
	if comment != -1 {
		line = line[0:comment]
	}
	line = strings.Trim(line, " \t\n\r")
	return line
}

func parse_line(line string) (linedata, error) {
	// match returns: _, _, label, opcode, arg
	m := LINE_RE.FindStringSubmatch(line)
	if m == nil {
		return linedata{}, errors.New("Incorrect line format.")
	}
	var data = linedata{}
	if m[2] != "" {
		data.label = &m[2]
		if m[3] == "" {
			// pure label, no opcode
			return data, nil
		}
	}
	if m[3] == "" {
		return linedata{}, errors.New("Incorrect line.")
	}

	var keyword = m[3]
	s, ok := TABLE[keyword]
	if !ok {
		return linedata{}, fmt.Errorf("Illegal command: %s.", keyword)
	}
	data.opcode = s.opcode
	if s.argument {
		data.len = 2
		operand := m[4]
		if operand == "" {
			return linedata{}, fmt.Errorf("Instruction %s requires an argument.", keyword)
		}
		mm := CONST_RE.FindStringSubmatch(operand)
		if mm == nil {
			// not a number, then
			refm := REF_RE.FindStringSubmatch(operand)
			if refm == nil {
				return linedata{}, fmt.Errorf("Operand is neither a number nor a valid label: [%s].", operand)
			}
			str := refm[1]
			data.reference = &str
		} else {
			val := mm[1]
			num, err := strconv.ParseUint(val, 16, 16)
			if err != nil {
				return linedata{}, errors.New("Illegal constant operand.")
			}
			data.operand = word(num)
		}
	} else {
		data.len = 1
		if m[4] != "" {
			return linedata{}, fmt.Errorf("Instruction %s does not take an argument.", keyword)
		}
	}

	return data, nil
}
