package style

import "github.com/aoldershaw/ansi/action"

type Style struct {
	Foreground action.Color `json:"fg,omitempty"`
	Background action.Color `json:"bg,omitempty"`
	Bold       bool         `json:"bold,omitempty"`
	Faint      bool         `json:"faint,omitempty"`
	Italic     bool         `json:"italic,omitempty"`
	Underline  bool         `json:"underline,omitempty"`
	Blink      bool         `json:"blink,omitempty"`
	Inverted   bool         `json:"inverted,omitempty"`
	Fraktur    bool         `json:"fraktur,omitempty"`
	Framed     bool         `json:"framed,omitempty"`
}
