package renderer

import (
	parser "github.com/horriblename/mathcha/latex"
)

type Renderer struct {
	Color                bool
	UnicodeSuperscript   bool
	Buffer               string
	LatexTree            parser.FlexContainer
	FocusOn              parser.Container
	HasSelection         bool
	Focus                bool
}

func New(color bool, unicodeSuperscript bool) Renderer {
	root := &parser.UnboundCompExpr{}
	return Renderer{
		Color:              color,
		UnicodeSuperscript: unicodeSuperscript,
		Buffer:             "",
		LatexTree:          root,
		FocusOn:            root,
		HasSelection:       false,
		Focus:              false,
	}
}

func FromFormula(formula string, color bool, unicodeSuperscript bool) *Renderer {
	root := parser.Parse(formula)
	return &Renderer{
		Color:              color,
		UnicodeSuperscript: unicodeSuperscript,
		Buffer:             "",
		LatexTree:          root,
		FocusOn:            root,
		HasSelection:       false,
		Focus:              false,
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
