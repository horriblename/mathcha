package latex

import "testing"

func tokenizeAll(input string) []Token {
	var tok Tokenizer
	tok.Init(input)
	var result []Token
	for !tok.IsEOF() {
		result = append(result, tok.Peek())
		tok.Eat()
	}
	return result
}

func TestTokenizer(t *testing.T) {
	testCases := []struct {
		desc   string
		input  string
		output []Token
	}{
		{
			desc:   "NUM - single digit",
			input:  "5",
			output: []Token{NUM},
		},
		{
			desc:   "NUM - multiple digits",
			input:  "123",
			output: []Token{NUM, NUM, NUM},
		},
		{
			desc:   "VARLIT - single letter",
			input:  "x",
			output: []Token{VARLIT},
		},
		{
			desc:   "VARLIT - multiple letters",
			input:  "xyz",
			output: []Token{VARLIT, VARLIT, VARLIT},
		},
		{
			desc:   "SYM - plus sign",
			input:  "+",
			output: []Token{SYM},
		},
		{
			desc:   "SYM - asterisk",
			input:  "*",
			output: []Token{SYM},
		},
		{
			desc:   "LBRACE - left brace",
			input:  "{",
			output: []Token{LBRACE},
		},
		{
			desc:   "RBRACE - right brace",
			input:  "}",
			output: []Token{RBRACE},
		},
		{
			desc:   "CARET - caret",
			input:  "^",
			output: []Token{CARET},
		},
		{
			desc:   "UNDERSCORE - underscore",
			input:  "_",
			output: []Token{UNDERSCORE},
		},
		{
			desc:   "AMPERSAND - ampersand",
			input:  "&",
			output: []Token{AMPERSAND},
		},
		{
			desc:   "COMMENT - comment",
			input:  "% this is a comment",
			output: []Token{COMMENT, VARLIT, VARLIT, VARLIT, VARLIT, VARLIT, VARLIT, VARLIT, VARLIT, VARLIT, VARLIT, VARLIT, VARLIT, VARLIT, VARLIT},
		},
		{
			desc:   "CMDSTR - command with letters",
			input:  "\\frac",
			output: []Token{CMDSTR},
		},
		{
			desc:   "CMDSTR - another command",
			input:  "\\text",
			output: []Token{CMDSTR},
		},
		{
			desc:   "CMDSYM - command with symbol",
			input:  "\\^",
			output: []Token{CMDSYM},
		},
		{
			desc:   "CMDSYM - backslash escape (or newline, I still have not checked what this does)",
			input:  "\\\\",
			output: []Token{CMDSYM},
		},
		{
			desc:   "EOF - empty input",
			input:  "",
			output: []Token{},
		},
		{
			desc:  "comprehensive - mixed tokens",
			input: "\\frac{1}{2} + x^2 * y_1 & z",
			output: []Token{
				CMDSTR, LBRACE, NUM, RBRACE,
				LBRACE, NUM, RBRACE,
				SYM, VARLIT,
				CARET, NUM,
				SYM, VARLIT,
				UNDERSCORE, NUM,
				AMPERSAND, VARLIT,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := tokenizeAll(tC.input)
			if len(result) != len(tC.output) {
				t.Fatalf("expected %d tokens, got %d: %v", len(tC.output), len(result), result)
			}
			for i, want := range tC.output {
				if result[i] != want {
					t.Errorf("token %d: want %v, got %v", i, want, result[i])
				}
			}
		})
	}
}
