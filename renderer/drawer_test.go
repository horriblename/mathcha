package renderer

import (
	"testing"

	parser "github.com/horriblename/mathcha/latex"
)

func TestPrerender(t *testing.T) {
	testCases := []struct {
		desc   string
		input  string
		expect string
	}{
		{
			desc:   "NumberLit - single digit",
			input:  "5",
			expect: "5",
		},
		{
			desc:   "NumberLit - multiple digits",
			input:  "123",
			expect: "123",
		},
		{
			desc:   "VarLit - single letter",
			input:  "x",
			expect: "x",
		},
		{
			desc:   "VarLit - multiple letters",
			input:  "xyz",
			expect: "xyz",
		},
		{
			desc:   "SimpleOpLit - plus sign",
			input:  "+",
			expect: " + ",
		},
		{
			desc:   "SimpleOpLit - multiple symbols",
			input:  "+-=",
			expect: " +  -  = ",
		},
		{
			desc:   "SimpleCmdLit - math symbol",
			input:  "\\times",
			expect: "×",
		},
		{
			desc:   "SimpleCmdLit - greek letter",
			input:  "\\pi",
			expect: "π",
		},
		{
			desc:   "CompositeExpr - simple braces",
			input:  "{x}",
			expect: "x",
		},
		{
			desc:   "CompositeExpr - nested",
			input:  "{a + b}",
			expect: "a + b",
		},
		{
			desc:   "SuperExpr - superscript",
			input:  "x^2",
			expect: " 2\nx ",
		},
		{
			desc:   "SubExpr - subscript",
			input:  "x_1",
			expect: "x \n 1",
		},
		{
			desc:   "Cmd1ArgExpr - sqrt",
			input:  "\\sqrt{x}",
			expect: "⎷x",
		},
		{
			desc:   "Cmd1ArgExpr - underline",
			input:  "\\underline{x}",
			expect: "x",
		},
		{
			desc:   "Cmd2ArgExpr - frac",
			input:  "\\frac{1}{2}",
			expect: "1\n─\n2",
		},
		{
			desc:   "Cmd2ArgExpr - binom",
			input:  "\\binom{a}{b}",
			expect: "[unimplemented command container]",
		},
		{
			desc:   "ParenCompExpr - left right parentheses",
			input:  "\\left( x \\right)",
			expect: "(x)",
		},
		{
			desc:   "ParenCompExpr - brackets",
			input:  "\\left[ x \\right]",
			expect: "[x]",
		},
		{
			desc:   "TextContainer - text command",
			input:  "\\text{hello}",
			expect: "hello",
		},
		{
			desc:   "EnvExpr - matrix environment",
			input:  `\begin{matrix} a & b \\ c & d \end{matrix}`,
			expect: "⎡a b⎤\n⎣c d⎦",
		},
		{
			desc:   "EnvExpr - single cell",
			input:  `\begin{matrix} x \end{matrix}`,
			expect: "[x]",
		},
		{
			desc: "EnvExpr - empty cells align correctly",
			input: `\begin{matrix}
	1 & 2 & 3 & 4 & 5\\
	a &   &   &   & d
\end{matrix}`,
			expect: "⎡1 2 3 4 5⎤\n⎣a       d⎦",
		},
		{
			desc: "EnvExpr - empty cells with uneven columns align correctly",
			input: `\begin{matrix}
	1 & 2 & 3 & 4 & 5\\
	a &  & d
\end{matrix}`,
			expect: "⎡1 2 3 4 5⎤\n⎣a d      ⎦",
		},
		{
			desc:   "Combined - simple expression",
			input:  "x + 1",
			expect: "x + 1",
		},
	}

	r := New(false)
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tree := parser.Parse(tC.input)
			if tree == nil {
				t.Fatal("Parse returned nil")
			}
			out, _ := r.Prerender(tree)
			if out != tC.expect {
				t.Errorf("got:  %q\nwant: %q", out, tC.expect)
			}
		})
	}
}
