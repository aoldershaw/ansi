package parser

import (
	"github.com/aoldershaw/ansi/action"
	"github.com/aoldershaw/ansi/style"
)

// https://en.wikipedia.org/wiki/ANSI_escape_code#SGR

const maxCode = 128

var sgrParamToAction [maxCode]action.Action

func init() {
	codeActionsMap := map[int]action.Action{
		0:  action.Reset{},
		1:  action.SetBold(true),
		2:  action.SetFaint(true),
		3:  action.SetItalic(true),
		4:  action.SetUnderline(true),
		5:  action.SetBlink(true),
		7:  action.SetInverted(true),
		20: action.SetFraktur(true),

		30: action.SetForeground(style.Black),
		31: action.SetForeground(style.Red),
		32: action.SetForeground(style.Green),
		33: action.SetForeground(style.Yellow),
		34: action.SetForeground(style.Blue),
		35: action.SetForeground(style.Magenta),
		36: action.SetForeground(style.Cyan),
		37: action.SetForeground(style.White),
		39: action.SetForeground(style.DefaultColor),

		40: action.SetBackground(style.Black),
		41: action.SetBackground(style.Red),
		42: action.SetBackground(style.Green),
		43: action.SetBackground(style.Yellow),
		44: action.SetBackground(style.Blue),
		45: action.SetBackground(style.Magenta),
		46: action.SetBackground(style.Cyan),
		47: action.SetBackground(style.White),
		49: action.SetBackground(style.DefaultColor),

		90: action.SetForeground(style.BrightBlack),
		91: action.SetForeground(style.BrightRed),
		92: action.SetForeground(style.BrightGreen),
		93: action.SetForeground(style.BrightYellow),
		94: action.SetForeground(style.BrightBlue),
		95: action.SetForeground(style.BrightMagenta),
		96: action.SetForeground(style.BrightCyan),
		97: action.SetForeground(style.BrightWhite),

		100: action.SetBackground(style.BrightBlack),
		101: action.SetBackground(style.BrightRed),
		102: action.SetBackground(style.BrightGreen),
		103: action.SetBackground(style.BrightYellow),
		104: action.SetBackground(style.BrightBlue),
		105: action.SetBackground(style.BrightMagenta),
		106: action.SetBackground(style.BrightCyan),
		107: action.SetBackground(style.BrightWhite),
	}
	for code, act := range codeActionsMap {
		sgrParamToAction[code] = act
	}
}

func sgrLookup(code int) (action.Action, bool) {
	if code >= maxCode || code < 0 {
		return nil, false
	}
	act := sgrParamToAction[code]
	if act == nil {
		return nil, false
	}
	return act, true
}
