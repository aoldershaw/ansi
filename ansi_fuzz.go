// +build gofuzz

package ansi

import "github.com/aoldershaw/ansi/output"

func Fuzz(data []byte) int {
	a := New(&output.InMemory{})
	a.Parse(data)
	return 0
}
