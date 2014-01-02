package sxed

import (
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"
)

const (
	NO_OP int = iota
	LITERAL_OP
	REGEXP_OP
)

const (
	EOF rune = -1
)

var opOfCommand = map[rune]int{
	'x': REGEXP_OP,
	'y': REGEXP_OP,
	'g': REGEXP_OP,
	'v': REGEXP_OP,
	'p': NO_OP,
	'd': NO_OP,
	'c': LITERAL_OP,
	'a': LITERAL_OP, // TODO Implement 'a'
}

type Command struct {
	Name rune
	Operand string
	Pattern *regexp.Regexp // compiled operand
}

func (cmd Command) Equals(cmd2 Command) bool {
	return cmd.Name == cmd2.Name && cmd.Operand == cmd2.Operand
}

type Chain []Command

func (ch Chain) Equals(ch2 Chain) bool {
	if len(ch) != len(ch2) {
		return false
	}
	for i, cmd := range ch {
		if !cmd.Equals(ch2[i]) {
			return false
		}
	}
	return true
}

type Program []Chain

func (p Program) Equals(p2 Program) bool {
	if len(p) != len(p2) {
		return false
	}
	for i, ch := range p {
		if !ch.Equals(p2[i]) {
			return false
		}
	}
	return true
}

type scanner struct {
	text string
	i int
}

func newScanner(text string) *scanner {
	return &scanner{text, 0}
}

func (sc *scanner) next() rune {
	if sc.i == len(sc.text) {
		return EOF
	}
	r, size := utf8.DecodeRuneInString(sc.text[sc.i:])
	sc.i += size
	return r
}

func (sc *scanner) nextNonSpace() rune {
	for {
		r := sc.next()
		if r == EOF || !unicode.IsSpace(r) {
			return r
		}
	}
}

func Parse(text string) (Program, error) {
	sc := newScanner(text)
	prog := make(Program, 0, 1)

	for {
		// Parse a chain
		chain := make(Chain, 0, 1)
		for {
			// Parse a command
			name := sc.nextNonSpace()
			if name == ';' || name == EOF {
				break
			}
			op, ok := opOfCommand[name]
			if !ok {
				return nil, fmt.Errorf("bad command name: %c", name)
			}
			cmd := Command{Name: name}
			if op != NO_OP {
				// Parse an operand
				var delim rune
				for {
					r := sc.next()
					if r == EOF {
						return nil, fmt.Errorf("missing operand for command %c", name)
					} else if !unicode.IsSpace(r) {
						delim = r
						break
					}
				}
				istart := sc.i
				op_loop: for {
					switch sc.next() {
					case EOF:
						return nil, fmt.Errorf("incomplete operand")
					case delim:
						break op_loop
					case '\\':
						sc.next()
					}
				}
				cmd.Operand = sc.text[istart:sc.i-1]
				if op == REGEXP_OP {
					var e error
					cmd.Pattern, e = regexp.Compile(cmd.Operand)
					if e != nil {
						return nil, e
					}
				}
			}
			chain = append(chain, cmd)
		}
		if len(chain) > 0 {
			prog = append(prog, chain)
		} else {
			break
		}
	}
	return prog, nil
}
