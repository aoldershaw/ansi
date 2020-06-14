package ansi

import (
	"strconv"
)

func (a Print) ActionString() string         { return "Print(" + string(a) + ")" }
func (a Reset) ActionString() string         { return "Reset" }
func (a SetForeground) ActionString() string { return "SetForeground(" + Color(a).String() + ")" }
func (a SetBackground) ActionString() string { return "SetBackground(" + Color(a).String() + ")" }
func (a SetBold) ActionString() string       { return "SetBold(" + strconv.FormatBool(bool(a)) + ")" }
func (a SetFaint) ActionString() string      { return "SetFaint(" + strconv.FormatBool(bool(a)) + ")" }
func (a SetItalic) ActionString() string     { return "SetItalic(" + strconv.FormatBool(bool(a)) + ")" }
func (a SetUnderline) ActionString() string {
	return "SetUnderline(" + strconv.FormatBool(bool(a)) + ")"
}
func (a SetBlink) ActionString() string       { return "SetBlink(" + strconv.FormatBool(bool(a)) + ")" }
func (a SetInverted) ActionString() string    { return "SetInverted(" + strconv.FormatBool(bool(a)) + ")" }
func (a SetFraktur) ActionString() string     { return "SetFraktur(" + strconv.FormatBool(bool(a)) + ")" }
func (a SetFramed) ActionString() string      { return "SetFramed(" + strconv.FormatBool(bool(a)) + ")" }
func (a Linebreak) ActionString() string      { return "Linebreak" }
func (a CarriageReturn) ActionString() string { return "CarriageReturn" }
func (a CursorUp) ActionString() string       { return "CursorUp(" + strconv.FormatInt(int64(a), 10) + ")" }
func (a CursorDown) ActionString() string {
	return "CursorDown(" + strconv.FormatInt(int64(a), 10) + ")"
}
func (a CursorForward) ActionString() string {
	return "CursorForward(" + strconv.FormatInt(int64(a), 10) + ")"
}
func (a CursorBack) ActionString() string {
	return "CursorBack(" + strconv.FormatInt(int64(a), 10) + ")"
}
func (a CursorPosition) ActionString() string { return "CursorPosition(" + Pos(a).String() + ")" }
func (a CursorColumn) ActionString() string {
	return "CursorColumn(" + strconv.FormatInt(int64(a), 10) + ")"
}
func (a EraseDisplay) ActionString() string          { return "EraseDisplay(" + EraseMode(a).String() + ")" }
func (a EraseLine) ActionString() string             { return "EraseLine(" + EraseMode(a).String() + ")" }
func (a SaveCursorPosition) ActionString() string    { return "SaveCursorPosition" }
func (a RestoreCursorPosition) ActionString() string { return "RestoreCursorPosition" }

func (a Print) String() string                 { return a.ActionString() }
func (a Reset) String() string                 { return a.ActionString() }
func (a SetForeground) String() string         { return a.ActionString() }
func (a SetBackground) String() string         { return a.ActionString() }
func (a SetBold) String() string               { return a.ActionString() }
func (a SetFaint) String() string              { return a.ActionString() }
func (a SetItalic) String() string             { return a.ActionString() }
func (a SetUnderline) String() string          { return a.ActionString() }
func (a SetBlink) String() string              { return a.ActionString() }
func (a SetInverted) String() string           { return a.ActionString() }
func (a SetFraktur) String() string            { return a.ActionString() }
func (a SetFramed) String() string             { return a.ActionString() }
func (a Linebreak) String() string             { return a.ActionString() }
func (a CarriageReturn) String() string        { return a.ActionString() }
func (a CursorUp) String() string              { return a.ActionString() }
func (a CursorDown) String() string            { return a.ActionString() }
func (a CursorForward) String() string         { return a.ActionString() }
func (a CursorBack) String() string            { return a.ActionString() }
func (a CursorPosition) String() string        { return a.ActionString() }
func (a CursorColumn) String() string          { return a.ActionString() }
func (a EraseDisplay) String() string          { return a.ActionString() }
func (a EraseLine) String() string             { return a.ActionString() }
func (a SaveCursorPosition) String() string    { return a.ActionString() }
func (a RestoreCursorPosition) String() string { return a.ActionString() }

func (p Pos) String() string {
	return "L" + strconv.FormatInt(int64(p.Line), 10) + "C" + strconv.FormatInt(int64(p.Col), 10)
}

var eraseModeNames = [4]string{
	"undefined",
	"EraseToBeginning",
	"EraseToEnd",
	"EraseAll",
}

func (e EraseMode) String() string { return eraseModeNames[e] }
