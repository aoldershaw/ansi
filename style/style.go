package style

type Style struct {
	Foreground Color `json:"fg,omitempty"`
	Background Color `json:"bg,omitempty"`
	Bold       bool  `json:"bold,omitempty"`
	Faint      bool  `json:"faint,omitempty"`
	Italic     bool  `json:"italic,omitempty"`
	Underline  bool  `json:"underline,omitempty"`
	Blink      bool  `json:"blink,omitempty"`
	Inverted   bool  `json:"inverted,omitempty"`
	Fraktur    bool  `json:"fraktur,omitempty"`
	Framed     bool  `json:"framed,omitempty"`
}
