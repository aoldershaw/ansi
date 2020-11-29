// +build gofuzz

package ansi

func Fuzz(data []byte) int {
	a := NewWriter(&InMemory{})
	a.Write(data)
	return 0
}
