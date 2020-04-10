// +build gofuzz

package ansi

func Fuzz(data []byte) int {
	a := New(&InMemory{})
	a.Parse(data)
	return 0
}
