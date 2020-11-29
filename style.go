package ansi

import "encoding/json"

type StyleModifier uint8

const (
	Bold StyleModifier = 1 << iota
	Faint
	Italic
	Underline
	Blink
	Inverted
	Fraktur
	Framed
)

type Style struct {
	Foreground Color
	Background Color
	Modifier   StyleModifier
}

func (s Style) MarshalJSON() ([]byte, error) {
	return json.Marshal(style{
		Foreground: s.Foreground,
		Background: s.Background,
		Bold:       s.Modifier&Bold != 0,
		Faint:      s.Modifier&Faint != 0,
		Italic:     s.Modifier&Italic != 0,
		Underline:  s.Modifier&Underline != 0,
		Blink:      s.Modifier&Blink != 0,
		Inverted:   s.Modifier&Inverted != 0,
		Fraktur:    s.Modifier&Fraktur != 0,
		Framed:     s.Modifier&Framed != 0,
	})
}

func (s *Style) UnmarshalJSON(data []byte) error {
	var ss style
	if err := json.Unmarshal(data, &ss); err != nil {
		return err
	}
	s.Foreground = ss.Foreground
	s.Background = ss.Background
	s.Modifier.applyBit(ss.Bold, Bold)
	s.Modifier.applyBit(ss.Faint, Faint)
	s.Modifier.applyBit(ss.Italic, Italic)
	s.Modifier.applyBit(ss.Underline, Underline)
	s.Modifier.applyBit(ss.Blink, Blink)
	s.Modifier.applyBit(ss.Inverted, Inverted)
	s.Modifier.applyBit(ss.Fraktur, Fraktur)
	s.Modifier.applyBit(ss.Framed, Framed)
	return nil
}

func (s *StyleModifier) applyBit(b bool, bit StyleModifier) {
	if b {
		*s |= bit
	} else {
		*s &= ^bit
	}
}

type style struct {
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
