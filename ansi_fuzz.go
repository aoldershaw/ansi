// +build gofuzz

package ansi

import (
	"github.com/aoldershaw/ansi/output"
	"github.com/aoldershaw/ansi/parser"
)

func Fuzz(data []byte) int {
	out := &output.InMemory{}
	state := New(Cooked, out)
	parse := parser.New(state)

	parse.Parse(data)

	return 0
}