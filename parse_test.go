package sxed

import (
	"testing"
)

func TestParse(t *testing.T) {
	inputs := []string{
		`x/[a-zA-Z0-9]+/ g/n/ v/../ c/num/; p`,
		`x|/usr/lib| c|/usr/local/lib|`,
	}
	wanted := []Program{
		Program{Chain{Command{'x', `[a-zA-Z0-9]+`, nil},
		              Command{'g', `n`, nil},
		              Command{'v', `..`, nil},
		              Command{'c', `num`, nil}},
		        Chain{Command{'p', ``, nil}}},
		Program{Chain{Command{'x', `/usr/lib`, nil},
		              Command{'c', `/usr/local/lib`, nil}}}}
	for i, in := range inputs {
		out, err := Parse(in)
		if err != nil {
			t.Errorf("Parse(%v) got error %v", in, err)
		} else if !out.Equals(wanted[i]) {
			t.Errorf("Parse(%v) = %v, want %v", in, out, wanted[i])
		}
	}
}
