package ansi

// https://bluesock.org/~willkg/dev/ansi.html
// https://en.wikipedia.org/wiki/ANSI_escape_code#CSI_sequences

type Action interface {
	ActionString() string
}

type Print []byte
type Reset struct{}
type SetForeground Color
type SetBackground Color
type SetBold bool
type SetFaint bool
type SetItalic bool
type SetUnderline bool
type SetBlink bool
type SetInverted bool
type SetFraktur bool
type SetFramed bool
type Linebreak struct{}
type CarriageReturn struct{}
type CursorUp int
type CursorDown int
type CursorForward int
type CursorBack int
type CursorPosition Pos
type CursorColumn int
type EraseDisplay EraseMode
type EraseLine EraseMode
type SaveCursorPosition struct{}
type RestoreCursorPosition struct{}

type Pos struct {
	Line int
	Col  int
}

type EraseMode uint8

const (
	EraseToEnd EraseMode = iota
	EraseToBeginning
	EraseAll
)
