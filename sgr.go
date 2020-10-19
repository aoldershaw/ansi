package ansi

// https://en.wikipedia.org/wiki/ANSI_escape_code#SGR

const maxCode = 128

var sgrParamToAction = [maxCode]Action{
	0:  Reset{},
	1:  SetBold(true),
	2:  SetFaint(true),
	3:  SetItalic(true),
	4:  SetUnderline(true),
	5:  SetBlink(true),
	7:  SetInverted(true),
	20: SetFraktur(true),

	30: SetForeground(Black),
	31: SetForeground(Red),
	32: SetForeground(Green),
	33: SetForeground(Yellow),
	34: SetForeground(Blue),
	35: SetForeground(Magenta),
	36: SetForeground(Cyan),
	37: SetForeground(White),
	39: SetForeground(DefaultColor),

	40: SetBackground(Black),
	41: SetBackground(Red),
	42: SetBackground(Green),
	43: SetBackground(Yellow),
	44: SetBackground(Blue),
	45: SetBackground(Magenta),
	46: SetBackground(Cyan),
	47: SetBackground(White),
	49: SetBackground(DefaultColor),

	90: SetForeground(BrightBlack),
	91: SetForeground(BrightRed),
	92: SetForeground(BrightGreen),
	93: SetForeground(BrightYellow),
	94: SetForeground(BrightBlue),
	95: SetForeground(BrightMagenta),
	96: SetForeground(BrightCyan),
	97: SetForeground(BrightWhite),

	100: SetBackground(BrightBlack),
	101: SetBackground(BrightRed),
	102: SetBackground(BrightGreen),
	103: SetBackground(BrightYellow),
	104: SetBackground(BrightBlue),
	105: SetBackground(BrightMagenta),
	106: SetBackground(BrightCyan),
	107: SetBackground(BrightWhite),
}

func sgrLookup(code int) (Action, bool) {
	if code >= maxCode || code < 0 {
		return nil, false
	}
	act := sgrParamToAction[code]
	if act == nil {
		return nil, false
	}
	return act, true
}
