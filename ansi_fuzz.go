//+build gofuzz

package ansi

func Fuzz(data []byte) int {
	a := NewWriter(new(Lines))
	a.Write(data)
	return 0
}
