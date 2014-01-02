package sxed

import (
	"os"
	"fmt"
	"bytes"
)

func Eval(text []byte, chain Chain) []byte {
	var indicies [][]int
	if len(chain) == 1 && chain[0].Name == 'p' {
		// A lone 'p' is treated specially
		os.Stdout.Write(text)
		return text
	}
	for _, cmd := range chain {
		switch cmd.Name {
		case 'x':
			// TODO FindAllSubmatchIndex?
			indicies = cmd.Pattern.FindAllIndex(text, -1)
		case 'y':
			match_indicies := cmd.Pattern.FindAllIndex(text, -1)
			n := len(match_indicies)
			// TODO Eliminate zero-width slices
			indicies = make([][]int, n + 1)
			indicies[0][0] = 0
			indicies[n][1] = len(text)
			for i := 0; i < n; i++ {
				indicies[i][1] = match_indicies[i][0]
				indicies[i+1][0] = match_indicies[i][1]
			}
		case 'g':
			old_indicies := indicies
			indicies = make([][]int, 0, len(old_indicies))
			for _, idx := range old_indicies {
				if cmd.Pattern.Match(text[idx[0]:idx[1]]) {
					indicies = append(indicies, idx)
				}
			}
		case 'v':
			// Duplicate with 'g'
			old_indicies := indicies
			indicies = make([][]int, 0, len(old_indicies))
			for _, idx := range old_indicies {
				if !cmd.Pattern.Match(text[idx[0]:idx[1]]) {
					indicies = append(indicies, idx)
				}
			}
		case 'p':
			for _, idx := range indicies {
				_, err := os.Stdout.Write(text[idx[0]:idx[1]])
				// TODO Maybe Eval should return an error
				if err != nil {
					fmt.Println(os.Stderr, err)
				}
			}
		case 'd', 'c': // 'd' is equivalent to 'c//'
			if len(indicies) > 0 {
				// TODO Would it be beneficial to eliminate this conversion?
				// TODO Support $n patterns in substitute
				substitute := []byte(cmd.Operand)
				// TODO Implement a quick path where all matches and the
				// substitute are of the same length
				// TODO Benchmark whether it will faster to compute needed
				// buffer size exactly when possible
				buf := bytes.NewBuffer(make([]byte, 0, len(text)))
				buf.Write(text[:indicies[0][0]])
				for i := 0; i < len(indicies) - 1; i++ {
					buf.Write(substitute)
					buf.Write(text[indicies[i][1]:indicies[i+1][0]])
				}
				buf.Write(substitute)
				buf.Write(text[indicies[len(indicies)-1][1]:])
				text = buf.Bytes()
				// TODO Update indicies to reflect new text?
				indicies = nil
			}
		default:
			panic(fmt.Errorf("Unknown command %c", cmd.Name))
		}
	}
	return text
}
