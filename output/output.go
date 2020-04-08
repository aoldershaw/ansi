package output

import (
	"github.com/aoldershaw/ansi/action"
	"github.com/aoldershaw/ansi/style"
)

type Output interface {
	Print(data []byte, style style.Style, pos action.Pos)
	ClearRight(pos action.Pos)
}
