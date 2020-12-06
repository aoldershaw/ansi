package ansi

const (
	defaultLines = 48
	defaultCols  = 80
)

type LineDiscipline int

const (
	Raw LineDiscipline = iota
	Cooked
)

type State struct {
	Style          Style
	LineDiscipline LineDiscipline
	Position       Pos
	SavedPosition  *Pos

	MaxLine int
	MaxCol  int
}

type Writer struct {
	State
	Parser *Parser
	Output Output
}

func NewWriter(output Output, opts ...WriterOption) *Writer {
	w := &Writer{
		State: State{
			MaxLine: defaultLines,
			MaxCol:  defaultCols,

			LineDiscipline: Cooked,
		},
		Parser: NewParser(),
		Output: output,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *Writer) Write(input []byte) (int64, error) {
	n := len(input)
	for {
		action, ok, newInput := w.Parser.Parse(input)
		if !ok {
			break
		}
		if err := w.Action(action); err != nil {
			return int64(n - len(input)), err
		}
		input = newInput
	}
	return int64(n), nil
}

func (w *Writer) Action(act Action) error {
	switch v := act.(type) {
	case Print:
		if err := w.Output.Print(v, w.Style, w.Position); err != nil {
			return err
		}
		endCol := w.Position.Col + len(v)
		if endCol > w.MaxCol {
			w.MaxCol = endCol
		}
		w.Position.Col = endCol
	case Reset:
		w.Style = Style{}
	case SetForeground:
		w.Style.Foreground = Color(v)
	case SetBackground:
		w.Style.Background = Color(v)
	case SetBold:
		w.Style.Modifier.applyBit(bool(v), Bold)
	case SetFaint:
		w.Style.Modifier.applyBit(bool(v), Faint)
	case SetItalic:
		w.Style.Modifier.applyBit(bool(v), Italic)
	case SetUnderline:
		w.Style.Modifier.applyBit(bool(v), Underline)
	case SetBlink:
		w.Style.Modifier.applyBit(bool(v), Blink)
	case SetInverted:
		w.Style.Modifier.applyBit(bool(v), Inverted)
	case SetFraktur:
		w.Style.Modifier.applyBit(bool(v), Fraktur)
	case SetFramed:
		w.Style.Modifier.applyBit(bool(v), Framed)
	case CursorPosition:
		w.moveCursorTo(v.Line, v.Col)
	case CursorUp:
		w.moveCursor(-int(v), 0)
	case CursorDown:
		w.moveCursor(int(v), 0)
	case CursorForward:
		w.moveCursor(0, int(v))
	case CursorBack:
		w.moveCursor(0, -int(v))
	case CursorColumn:
		w.moveCursorTo(w.Position.Line, int(v))
	case Linebreak:
		switch w.LineDiscipline {
		case Raw:
			w.Position.Line++
		case Cooked:
			w.Position.Line++
			w.Position.Col = 0
		}
		if w.Position.Line > w.MaxLine {
			w.MaxLine = w.Position.Line
		}
	case CarriageReturn:
		w.Position.Col = 0
	case SaveCursorPosition:
		pos := w.Position
		w.SavedPosition = &pos
	case RestoreCursorPosition:
		if w.SavedPosition != nil {
			w.Position = *w.SavedPosition
		}
	case EraseLine:
		startOfLine := w.Position
		startOfLine.Col = 0
		switch EraseMode(v) {
		case EraseToBeginning:
			if w.Position.Col == 0 {
				return nil
			}
			empty := spacer(w.Position.Col)
			w.Output.Print(empty, Style{}, startOfLine)
		case EraseToEnd:
			pos := w.Position
			pos.Col++
			if err := w.Output.ClearRight(pos); err != nil {
				return err
			}
		case EraseAll:
			if err := w.Output.ClearRight(startOfLine); err != nil {
				return err
			}
		}

	case EraseDisplay:
		// unsupported
	}

	return nil
}

func (w *Writer) moveCursorTo(l, c int) {
	w.Position.Line = l
	w.Position.Col = c
	if w.Position.Line < 0 {
		w.Position.Line = 0
	}
	if w.Position.Col < 0 {
		w.Position.Col = 0
	}
	if w.Position.Line > w.MaxLine {
		w.Position.Line = w.MaxLine
	}
	if w.Position.Col > w.MaxCol {
		w.Position.Col = w.MaxCol
	}
}

func (w *Writer) moveCursor(dl, dc int) {
	w.moveCursorTo(w.Position.Line+dl, w.Position.Col+dc)
}

type WriterOption func(*Writer)

func WithLineDiscipline(d LineDiscipline) WriterOption {
	return func(w *Writer) {
		w.State.LineDiscipline = d
	}
}

func WithInitialScreenSize(lines, cols int) WriterOption {
	return func(w *Writer) {
		if lines > 0 {
			w.State.MaxLine = lines
		}
		if cols > 0 {
			w.State.MaxCol = cols
		}
	}
}
