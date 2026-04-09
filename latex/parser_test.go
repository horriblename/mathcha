package latex

import (
	"testing"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		desc   string
		input  string
		expect *UnboundCompExpr
	}{
		{
			desc:  "NumberLit - single digit",
			input: "5",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&NumberLit{From: 0, To: 0, Source: "5"},
				},
			},
		},
		{
			desc:  "NumberLit - multiple digits",
			input: "123",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&NumberLit{From: 0, To: 0, Source: "1"},
					&NumberLit{From: 0, To: 0, Source: "2"},
					&NumberLit{From: 0, To: 0, Source: "3"},
				},
			},
		},
		{
			desc:  "VarLit - single letter",
			input: "x",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&VarLit{From: 0, To: 0, Source: "x"},
				},
			},
		},
		{
			desc:  "VarLit - multiple letters",
			input: "xyz",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&VarLit{From: 0, To: 0, Source: "x"},
					&VarLit{From: 0, To: 0, Source: "y"},
					&VarLit{From: 0, To: 0, Source: "z"},
				},
			},
		},
		{
			desc:  "SimpleOpLit - plus sign",
			input: "+",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&SimpleOpLit{From: 0, To: 0, Source: "+"},
				},
			},
		},
		{
			desc:  "SimpleOpLit - multiple symbols",
			input: "+-=",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&SimpleOpLit{From: 0, To: 0, Source: "+"},
					&SimpleOpLit{From: 0, To: 0, Source: "-"},
					&SimpleOpLit{From: 0, To: 0, Source: "="},
				},
			},
		},
		{
			desc:  "SimpleCmdLit - math symbol",
			input: "\\times",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&SimpleCmdLit{Backslash: 0, Source: "\\times", Type: CMD_times, To: 0},
				},
			},
		},
		{
			desc:  "SimpleCmdLit - greek letter",
			input: "\\pi",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&SimpleCmdLit{Backslash: 0, Source: "\\pi", Type: CMD_pi, To: 0},
				},
			},
		},
		{
			desc:  "CompositeExpr - simple braces",
			input: "{x}",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&CompositeExpr{
						Lbrace: 0,
						Elts: []Expr{
							&VarLit{From: 0, To: 0, Source: "x"},
						},
						Rbrace: 0,
					},
				},
			},
		},
		{
			desc:  "CompositeExpr - nested",
			input: "{a + b}",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&CompositeExpr{
						Lbrace: 0,
						Elts: []Expr{
							&VarLit{From: 0, To: 0, Source: "a"},
							&SimpleOpLit{From: 0, To: 0, Source: "+"},
							&VarLit{From: 0, To: 0, Source: "b"},
						},
						Rbrace: 0,
					},
				},
			},
		},
		{
			desc:  "SuperExpr - superscript",
			input: "x^2",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&VarLit{From: 0, To: 0, Source: "x"},
					&Cmd1ArgExpr{
						Type:      CMD_superscript,
						Backslash: 0,
						Arg1: &CompositeExpr{
							Lbrace: 0,
							Elts: []Expr{
								&NumberLit{From: 0, To: 0, Source: "2"},
							},
							Rbrace: 0,
						},
						To: 0,
					},
				},
			},
		},
		{
			desc:  "SubExpr - subscript",
			input: "x_1",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&VarLit{From: 0, To: 0, Source: "x"},
					&Cmd1ArgExpr{
						Type:      CMD_subscript,
						Backslash: 0,
						Arg1: &CompositeExpr{
							Lbrace: 0,
							Elts: []Expr{
								&NumberLit{From: 0, To: 0, Source: "1"},
							},
							Rbrace: 0,
						},
						To: 0,
					},
				},
			},
		},
		{
			desc:  "Cmd1ArgExpr - sqrt",
			input: "\\sqrt{x}",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&Cmd1ArgExpr{
						Type:      CMD_sqrt,
						Backslash: 0,
						Arg1: &CompositeExpr{
							Lbrace: 0,
							Elts: []Expr{
								&VarLit{From: 0, To: 0, Source: "x"},
							},
							Rbrace: 0,
						},
						To: 0,
					},
				},
			},
		},
		{
			desc:  "Cmd1ArgExpr - underline",
			input: "\\underline{x}",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&Cmd1ArgExpr{
						Type:      CMD_underline,
						Backslash: 0,
						Arg1: &CompositeExpr{
							Lbrace: 0,
							Elts: []Expr{
								&VarLit{From: 0, To: 0, Source: "x"},
							},
							Rbrace: 0,
						},
						To: 0,
					},
				},
			},
		},
		{
			desc:  "Cmd2ArgExpr - frac",
			input: "\\frac{1}{2}",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&Cmd2ArgExpr{
						Type:      CMD_frac,
						Backslash: 0,
						Arg1: &CompositeExpr{
							Lbrace: 0,
							Elts: []Expr{
								&NumberLit{From: 0, To: 0, Source: "1"},
							},
							Rbrace: 0,
						},
						Arg2: &CompositeExpr{
							Lbrace: 0,
							Elts: []Expr{
								&NumberLit{From: 0, To: 0, Source: "2"},
							},
							Rbrace: 0,
						},
						To: 0,
					},
				},
			},
		},
		{
			desc:  "Cmd2ArgExpr - binom",
			input: "\\binom{a}{b}",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&Cmd2ArgExpr{
						Type:      CMD_binom,
						Backslash: 0,
						Arg1: &CompositeExpr{
							Lbrace: 0,
							Elts: []Expr{
								&VarLit{From: 0, To: 0, Source: "a"},
							},
							Rbrace: 0,
						},
						Arg2: &CompositeExpr{
							Lbrace: 0,
							Elts: []Expr{
								&VarLit{From: 0, To: 0, Source: "b"},
							},
							Rbrace: 0,
						},
						To: 0,
					},
				},
			},
		},
		{
			desc:  "ParenCompExpr - left right parentheses",
			input: "\\left( x \\right)",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&ParenCompExpr{
						From:  0,
						Left:  "(",
						Right: ")",
						Elts: []Expr{
							&VarLit{From: 0, To: 0, Source: "x"},
						},
						To: 0,
					},
				},
			},
		},
		{
			desc:  "ParenCompExpr - brackets",
			input: "\\left[ x \\right]",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&ParenCompExpr{
						From:  0,
						Left:  "[",
						Right: "]",
						Elts: []Expr{
							&VarLit{From: 0, To: 0, Source: "x"},
						},
						To: 0,
					},
				},
			},
		},
		{
			desc:  "TextContainer - text command",
			input: "\\text{hello}",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&TextContainer{
						CmdText: 0,
						Type:    CMD_text,
						From:    0,
						To:      0,
						Text: &TextStringWrapper{
							Runes: []Expr{
								RawRuneLit('h'), RawRuneLit('e'), RawRuneLit('l'), RawRuneLit('l'), RawRuneLit('o'),
							},
						},
					},
				},
			},
		},
		{
			desc:  "EnvExpr - matrix environment",
			input: `\begin{matrix} a & b \\ c & d \end{matrix}`,
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&EnvExpr{
						From: 7,
						To:   42,
						Name: ENV_matrix,
						Elts: [][]*UnboundCompExpr{
							{
								{Elts: []Expr{&VarLit{From: 0, To: 0, Source: "a"}}},
								{Elts: []Expr{&VarLit{From: 0, To: 0, Source: "b"}}},
							},
							{
								{Elts: []Expr{&VarLit{From: 0, To: 0, Source: "c"}}},
								{Elts: []Expr{&VarLit{From: 0, To: 0, Source: "d"}}},
							},
						},
					},
				},
			},
		},
		{
			desc:  "EnvExpr - single cell",
			input: `\begin{matrix} x \end{matrix}`,
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&EnvExpr{
						Name: ENV_matrix,
						From: 7,
						To:   29,
						Elts: [][]*UnboundCompExpr{
							{
								{Elts: []Expr{&VarLit{From: 0, To: 0, Source: "x"}}},
							},
						},
					},
				},
			},
		},
		{
			desc:  "Combined - simple expression",
			input: "x + 1",
			expect: &UnboundCompExpr{
				Elts: []Expr{
					&VarLit{From: 0, To: 0, Source: "x"},
					&SimpleOpLit{From: 0, To: 0, Source: "+"},
					&NumberLit{From: 0, To: 0, Source: "1"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tree := Parse(tc.input)
			if tree == nil {
				t.Fatal("Parse returned nil")
			}
			if !tree.DeepEq(tc.expect) {
				if env, ok := tree.Elts[0].(*EnvExpr); ok {
					t.Logf("EnvExpr: %#v", env)
				}
				t.Errorf("parsed tree does not match expected\ngot:      %s\nexpected: %s", tree.VisualizeTree(), tc.expect.VisualizeTree())
			}
		})
	}
}
