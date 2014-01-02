package main

import (
	"io"
	"os"
	"fmt"
	"github.com/xiaq/sxed"
)

const (
	READ_BLOCK = 4 * 1024
)

var usage = `Usage: sxed PROGRAM`

func slurp(f *os.File) ([]byte, error) {
	bs := make([]byte, 0, READ_BLOCK)
	for {
		b := make([]byte, READ_BLOCK)
		_, err := f.Read(b)
		switch err {
		case nil:
			bs = append(bs, b...)
		case io.EOF:
			return bs, nil
		default:
			return nil, err
		}
	}
}

func main() {
	// Parse args.
	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}
	progText := os.Args[1]
	program, err := sxed.Parse(progText)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// Read input.
	text, err := slurp(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

	// Evaluate.
	for _, chain := range program {
		text = sxed.Eval(text, chain)
	}
}
