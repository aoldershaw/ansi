package ansi

type Output interface {
	// Print must not modify the slice data, not even temporarily.
	// Implementations must not retain a reference to data.
	Print(data []byte, style Style, pos Pos) error
	ClearRight(pos Pos) error
}
