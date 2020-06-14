package ansi

import "encoding/json"

type Color uint8

const (
	DefaultColor Color = iota
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

var colourNames = [17]string{
	"",
	"black",
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"white",
	"bright-black",
	"bright-red",
	"bright-green",
	"bright-yellow",
	"bright-blue",
	"bright-magenta",
	"bright-cyan",
	"bright-white",
}

func (c Color) String() string {
	if int(c) >= len(colourNames) {
		return ""
	}
	return colourNames[c]
}
func (c Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}
