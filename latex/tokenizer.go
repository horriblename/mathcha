// The Tokenizer has two main functions: Peek() and Eat(). Peek()
// returns the current token(Of the tokenizer). Each token is a
// (atomic) string which is a single digit(NUM) / letter(VARLIT) /
// symbol(SYM) or a latex command, (e.g. \int) and latex commands
// are seperated into commands formed by a word, CMDSTR, or a
// command formed by a single symbol(e.g. \{), CMDSYM.
package latex

import (
	"fmt"
	re "regexp"
	"strconv"
)

type Pos int

// Token is the set of lexical tokens
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	NOTOKEN     // for Parser.expect, representing 'not expecting anything'
	SINGLETOKEN // for Parser.expect, 'expecting single characters'
	COMMENT

	literal_beg
	// basic type literals
	NUM     // numbers
	VARLIT  // variable string literal, contains only alphabets
	TEXTSTR // string in \text, can be anything: symbols, numbers, letters...
	SYM     // symbols: non-alphabet characters that have no special grammar
	literal_end

	operator_beg
	LPAREN // (
	LBRACK // [
	LBRACE // {
	RPAREN // )
	RBRACK // ]
	RBRACE // }

	CARET      // ^
	UNDERSCORE // _
	AMPERSAND  // &
	operator_end

	command_beg
	CMDSTR
	CMDSYM
	command_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	NUM:     "NUM",
	VARLIT:  "VARLIT",
	TEXTSTR: "TEXTSTR",
	SYM:     "SYM",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	RPAREN: ")",
	RBRACK: "]",
	RBRACE: "}",

	CARET:      "^",
	UNDERSCORE: "_",
	AMPERSAND:  "&",

	CMDSTR: "CMDSTR",
	CMDSYM: "CMDSYM",
}

type tokenizer interface {
	Peek() Token
	Eat() Token
	IsEOF() bool
}

type Tokenizer struct {
	Cursor Pos
	Stream string

	curr   string
	tok    Token
	eof    bool
	buffer []Token

	strCmdRegex *re.Regexp // for matching commands that consist of alphabets
	symCmdRegex *re.Regexp // for matching commands that consist of symbols e.g. "\;"
	numRegex    *re.Regexp
	alpRegex    *re.Regexp
	spaceRegex  *re.Regexp // used to remove whitespaces
}

var ()

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}

	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}

	return s
}

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
//
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

func (t *Tokenizer) Init(stream string) {
	t.Stream = stream
	println(stream)
	t.Cursor = Pos(0)

	// FIXME move to package Init() function
	t.strCmdRegex = re.MustCompile(`^\\[a-zA-Z]+`)
	t.symCmdRegex = re.MustCompile(`^\\[^a-zA-Z0-9]`)
	t.numRegex = re.MustCompile("^[0-9]")    // FIXME do I really need regex here
	t.alpRegex = re.MustCompile("^[a-zA-Z]") // FIXME "
	t.spaceRegex = re.MustCompile("^ +")

	t.Eat()
}

func (t *Tokenizer) Peek() Token { return t.tok }

func (t *Tokenizer) Eat() string {
	t.consumeWhitespaces()
	if t.IsEOF() {
		fmt.Println("\x1b[31mthrow error: Eat() when at EOF\x1b[0m")
		return t.curr
	}
	if int(t.Cursor) >= len(t.Stream)-1 {
		t.eof = true
		t.tok = EOF
		return t.curr
	}
	stream := t.Stream[t.Cursor:len(t.Stream)]
	curr := t.curr
	tok := SYM  // ensure a new token is assigned each call
	length := 1 // default value to catch-all

	defer func() {
		t.curr = stream[0:length]
		t.tok = tok
		t.Cursor = Pos(int(t.Cursor) + length)
		//fmt.Printf("stream '\x1b[31m%s\x1b[0m%s'\n", t.curr, stream[length:])
	}()

	if temp := t.strCmdRegex.FindStringIndex(stream); temp != nil {
		length = temp[1]
		tok = CMDSTR
		return curr
	}
	if temp := t.symCmdRegex.FindStringIndex(stream); temp != nil {
		length = temp[1]
		tok = CMDSYM
		return curr
	}
	if temp := t.numRegex.FindStringIndex(stream); temp != nil {
		length = temp[1]
		tok = NUM
		return curr
	}
	if temp := t.alpRegex.FindStringIndex(stream); temp != nil {
		length = temp[1]
		tok = VARLIT
		return curr
	}
	// match single character tokens
	switch stream[0:1] { // stream[0] is a byte, stream[0:1] is a slice
	case "{":
		tok = LBRACE
		return curr
	case "}":
		tok = RBRACE
		return curr
	case "^":
		tok = CARET
		return curr
	case "_":
		tok = UNDERSCORE
		return curr
	case "&":
		tok = AMPERSAND
		return curr
	case "%":
		tok = COMMENT // FIXME better name?
		return curr
	}

	return curr
}

func (t *Tokenizer) IsEOF() bool { return t.eof }
func (t *Tokenizer) consumeWhitespaces() {
	stream := t.Stream[t.Cursor:]
	if loc := t.spaceRegex.FindStringIndex(stream); loc != nil {
		t.Cursor = Pos(int(t.Cursor) + loc[1])
	}
}
