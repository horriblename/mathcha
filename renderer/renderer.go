package renderer

import (
	parser "github.com/horriblename/mathcha/latex"
)

type Renderer struct {
	Buffer       string
	LatexTree    parser.FlexContainer
	FocusOn      parser.Container // the container in which the cursor is, a better implementation would be letting Render functions return a 'focused' flag when cursor is found
	HasSelection bool             // whether there is a selection in FocusOn
	Focus        bool             // whether the widget itself is focused
}

func New() Renderer {
	root := &parser.UnboundCompExpr{}
	return Renderer{
		Buffer:       "",
		LatexTree:    root,
		FocusOn:      root,
		HasSelection: false,
		Focus:        false,
	}
}

func (r *Renderer) Load(tree parser.FlexContainer) {
	r.LatexTree = tree
	r.Sync(nil, false)
}

// rerender the latex tree
func (r *Renderer) Sync(focus parser.Container, selected bool /*whether there is a selection*/) {
	r.FocusOn = focus
	r.HasSelection = selected
	r.DrawToBuffer(r.LatexTree)
}

func (r *Renderer) View() string {
	return r.Buffer
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(numbers ...int) int {
	if len(numbers) == 0 {
		panic("min was passed 0 parameters")
	} else if len(numbers) == 1 {
		return numbers[0]
	}

	m := numbers[0]
	for _, n := range numbers {
		if n < m {
			m = n
		}
	}
	return m
}
