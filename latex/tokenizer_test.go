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

func tokenizeAllLit(input string) []string {
	var tok Tokenizer
	tok.Init(input)
	var result []string
	for !tok.IsEOF() {
		result = append(result, tok.Eat())
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
			output: []Token{NEWLINE},
		},
		{
			desc:   "EOF - empty input",
			input:  "",
			output: []Token{},
		},
		{
			desc:   "newline is skipped",
			input:  "\n",
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
		{
			desc:   "UTF-8 - single 2-byte char",
			input:  "é",
			output: []Token{SYM},
		},
		{
			desc:   "UTF-8 - single 3-byte char",
			input:  "中",
			output: []Token{SYM},
		},
		{
			desc:   "UTF-8 - multiple multi-byte chars",
			input:  "中文",
			output: []Token{SYM, SYM},
		},
		{
			desc:   "UTF-8 - mixed with ASCII",
			input:  "x+é+y",
			output: []Token{VARLIT, SYM, SYM, SYM, VARLIT},
		},
		{
			desc:   "UTF-8 - command then multi-byte",
			input:  "\\pi é",
			output: []Token{CMDSTR, SYM},
		},
		{
			desc:   "UTF-8 - braced text with multi-byte",
			input:  "\\text{é}",
			output: []Token{CMDSTR, LBRACE, SYM, RBRACE},
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

func TestTokenizerUTF8Literals(t *testing.T) {
	testCases := []struct {
		desc   string
		input  string
		output []string
	}{
		{
			desc:   "2-byte char is one token",
			input:  "é",
			output: []string{"é"},
		},
		{
			desc:   "3-byte chars are separate tokens",
			input:  "中文",
			output: []string{"中", "文"},
		},
		{
			desc:   "4-byte emoji is one token",
			input:  "😀",
			output: []string{"😀"},
		},
		{
			desc:   "mixed ASCII and multi-byte",
			input:  "x+é",
			output: []string{"x", "+", "é"},
		},
		{
			desc:   "command then multi-byte",
			input:  "\\pi é",
			output: []string{"\\pi", "é"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := tokenizeAllLit(tC.input)
			if len(result) != len(tC.output) {
				t.Fatalf("expected %d tokens, got %d: %v", len(tC.output), len(result), result)
			}
			for i, want := range tC.output {
				if result[i] != want {
					t.Errorf("token %d: want %q, got %q", i, want, result[i])
				}
			}
		})
	}
}
