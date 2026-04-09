package renderer

import (
	"testing"

	parser "github.com/horriblename/mathcha/latex"
)

func TestLatexConsistent(t *testing.T) {
	testCases := []struct {
		desc  string
		input string
	}{
		{
			desc:  "NumberLit - single digit",
			input: "5",
		},
		{
			desc:  "NumberLit - multiple digits",
			input: "123",
		},
		{
			desc:  "VarLit - single letter",
			input: "x",
		},
		{
			desc:  "SimpleOpLit - plus sign",
			input: "+",
		},
		{
			desc:  "SimpleCmdLit - math symbol",
			input: "\\times",
		},
		{
			desc:  "SimpleCmdLit - greek letter",
			input: "\\pi",
		},
		{
			desc:  "CompositeExpr - simple braces",
			input: "{x}",
		},
		{
			desc:  "CompositeExpr - nested",
			input: "{a + b}",
		},
		{
			desc:  "SuperExpr - superscript",
			input: "x^2",
		},
		{
			desc:  "SubExpr - subscript",
			input: "x_1",
		},
		{
			desc:  "Cmd1ArgExpr - sqrt",
			input: "\\sqrt{x}",
		},
		{
			desc:  "Cmd1ArgExpr - underline",
			input: "\\underline{x}",
		},
		{
			desc:  "Cmd2ArgExpr - frac",
			input: "\\frac{1}{2}",
		},
		{
			desc:  "Cmd2ArgExpr - binom",
			input: "\\binom{a}{b}",
		},
		{
			desc:  "ParenCompExpr - left right parentheses",
			input: "\\left( x \\right)",
		},
		{
			desc:  "ParenCompExpr - brackets",
			input: "\\left[ x \\right]",
		},
		{
			desc:  "TextContainer - text command",
			input: "\\text{hello}",
		},
		{
			desc:  "EnvExpr - matrix environment",
			input: `\begin{matrix} a & b \\ c & d \end{matrix}`,
		},
		{
			desc:  "EnvExpr - single cell",
			input: `\begin{matrix} x \end{matrix}`,
		},
		{
			desc:  "Combined - simple expression",
			input: "x + 1",
		},
	}

	cfg := &LatexSourceConfig{}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tree := parser.Parse(tc.input)
			if tree == nil {
				t.Fatal("Parse returned nil")
			}

			rendered := cfg.ProduceLatex(tree)
			tree2 := parser.Parse(rendered)
			if tree2 == nil {
				t.Fatalf("Parse returned nil for rendered latex: %s", rendered)
			}

			if !tree.DeepEqWith(tree2, parser.DeepEqCfg{SkipPos: true}) {
				t.Errorf("AST mismatch after parse->render->parse\noriginal:  %s\nrendered:  %s\nreparsed:  %s",
					tree.VisualizeTree(), rendered, tree2.VisualizeTree())
			}
		})
	}
}
