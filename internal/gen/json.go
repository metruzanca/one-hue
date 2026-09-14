package gen

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Value is a printable JSON value. Accepted dynamic types:
//
//	nil          → null
//	string, bool
//	int, float64 → number
//	[]Field      → object (empty slice renders "{}")
//	[]Value      → array  (empty slice renders "[]")
type Value any

// Field is one ordered key/value entry of a JSON object.
type Field struct {
	Key string
	Val Value
}

// Encoder writes JSON with stable key order and tab indentation, mirroring
// the pretty output of the original TypeScript generator.
type Encoder struct {
	b bytes.Buffer
}

// Encode serializes v to compact, tab-indented JSON bytes (no trailing
// newline, matching JSON.stringify(_, null, "\t")).
func Encode(v Value) []byte {
	e := &Encoder{}
	e.write(v, 0)
	return e.b.Bytes()
}

func indent(depth int) string {
	return strings.Repeat("\t", depth)
}

func (e *Encoder) write(v Value, depth int) {
	switch t := v.(type) {
	case nil:
		e.b.WriteString("null")
	case string:
		e.writeString(t)
	case bool:
		if t {
			e.b.WriteString("true")
		} else {
			e.b.WriteString("false")
		}
	case int:
		fmt.Fprintf(&e.b, "%d", t)
	case float64:
		e.b.WriteString(trimFloat(t))
	case []Field:
		if len(t) == 0 {
			e.b.WriteString("{}")
			return
		}
		e.b.WriteByte('{')
		for i, f := range t {
			if i > 0 {
				e.b.WriteByte(',')
			}
			e.b.WriteByte('\n')
			e.b.WriteString(indent(depth + 1))
			e.writeString(f.Key)
			e.b.WriteString(": ")
			e.write(f.Val, depth+1)
		}
		e.b.WriteByte('\n')
		e.b.WriteString(indent(depth))
		e.b.WriteByte('}')
	case []Value:
		if len(t) == 0 {
			e.b.WriteString("[]")
			return
		}
		e.b.WriteByte('[')
		for i, el := range t {
			if i > 0 {
				e.b.WriteByte(',')
			}
			e.b.WriteByte('\n')
			e.b.WriteString(indent(depth + 1))
			e.write(el, depth+1)
		}
		e.b.WriteByte('\n')
		e.b.WriteString(indent(depth))
		e.b.WriteByte(']')
	default:
		panic(fmt.Sprintf("gen: unsupported JSON value %T", v))
	}
}

func trimFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func (e *Encoder) writeString(s string) {
	e.b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			e.b.WriteString(`\"`)
		case '\\':
			e.b.WriteString(`\\`)
		case '\n':
			e.b.WriteString(`\n`)
		case '\t':
			e.b.WriteString(`\t`)
		case '\r':
			e.b.WriteString(`\r`)
		case '\b':
			e.b.WriteString(`\b`)
		case '\f':
			e.b.WriteString(`\f`)
		default:
			if r < 0x20 {
				fmt.Fprintf(&e.b, `\u%04x`, r)
			} else {
				e.b.WriteRune(r)
			}
		}
	}
	e.b.WriteByte('"')
}
