package parser

import "github.com/aoldershaw/ansi/action"

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

		30: action.SetForeground(action.Black),
		31: action.SetForeground(action.Red),
		32: action.SetForeground(action.Green),
		33: action.SetForeground(action.Yellow),
		34: action.SetForeground(action.Blue),
		35: action.SetForeground(action.Magenta),
		36: action.SetForeground(action.Cyan),
		37: action.SetForeground(action.White),
		39: action.SetForeground(action.DefaultColor),

		40: action.SetBackground(action.Black),
		41: action.SetBackground(action.Red),
		42: action.SetBackground(action.Green),
		43: action.SetBackground(action.Yellow),
		44: action.SetBackground(action.Blue),
		45: action.SetBackground(action.Magenta),
		46: action.SetBackground(action.Cyan),
		47: action.SetBackground(action.White),
		49: action.SetBackground(action.DefaultColor),

		90: action.SetForeground(action.BrightBlack),
		91: action.SetForeground(action.BrightRed),
		92: action.SetForeground(action.BrightGreen),
		93: action.SetForeground(action.BrightYellow),
		94: action.SetForeground(action.BrightBlue),
		95: action.SetForeground(action.BrightMagenta),
		96: action.SetForeground(action.BrightCyan),
		97: action.SetForeground(action.BrightWhite),

		100: action.SetBackground(action.BrightBlack),
		101: action.SetBackground(action.BrightRed),
		102: action.SetBackground(action.BrightGreen),
		103: action.SetBackground(action.BrightYellow),
		104: action.SetBackground(action.BrightBlue),
		105: action.SetBackground(action.BrightMagenta),
		106: action.SetBackground(action.BrightCyan),
		107: action.SetBackground(action.BrightWhite),
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
