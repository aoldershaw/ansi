package ansi

const (
	defaultLines = 48
	defaultCols  = 80
)

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

type Writer struct {
	Parser *Parser
	State  *State
}

func NewWriter(output Output, opts ...WriterOption) *Writer {
	w := &Writer{
		Parser: NewParser(),
		State: &State{
			MaxLine: defaultLines,
			MaxCol:  defaultCols,

			LineDiscipline: Cooked,

			output: output,
		},
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *Writer) Write(input []byte) (int64, error) {
	n := len(input)
	for {
		var (
			action Action
			ok     bool
		)
		action, ok, input = w.Parser.Parse(input)
		if !ok {
			break
		}
		w.State.Action(action)
	}
	return int64(n), nil
}
