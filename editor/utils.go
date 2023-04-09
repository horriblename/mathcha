package editor

import (
	"github.com/horriblename/mathcha/latex"
)

// Returns the currently focused text field (RunesContainer) or nil
func (e *Editor) FocusedTextField() latex.RunesContainer {
	if text, ok := e.getParent().(latex.RunesContainer); ok {
		return text
	}

	return nil
}
