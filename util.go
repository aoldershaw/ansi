package ansi

type interval struct {
	L int
	R int
}

func intervalWithWidth(start int, width int) interval {
	return interval{L: start, R: start + width - 1}
}

func (i interval) contains(i2 interval) bool {
	return i.L <= i2.L && i2.R <= i.R
}

func (i interval) overlaps(i2 interval) bool {
	return i.R >= i2.L && i.L <= i2.R
}

func (i interval) length() int {
	return i.R - i.L + 1
}
