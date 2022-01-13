package latex

import (
	"fmt"
)

type Parser struct {
	eh        ErrorHandler
	tokenizer Tokenizer
	// Next token
	pos       Pos      // token position
	tok       Token    // one token look-ahead
	lit       string   // token literal
	expecting []string // FIXME new type?

	// Non-syntactic parser control
	exprLev  int // depth in tree of current position
	treeRoot *UnboundCompExpr
}

func (p *Parser) Init(src string) {
	//eh := func(pos Pos, msg string) { p.errors = append(p.errors, msg) }
	p.tokenizer.Init(src /*, eh*/)
	p.next()
	p.treeRoot = p.parseTopLevel()
}

func (p *Parser) GetTree() *UnboundCompExpr { return p.treeRoot }

func (p *Parser) next() {
	p.tok = p.tokenizer.Peek()
	if !p.tokenizer.IsEOF() {
		p.lit = p.tokenizer.Eat()
	} else {
		p.lit = p.tokenizer.curr // TODO add method to tokenizer instead of directly accessing curr?
	}
	// println("next():p.tok:", p.tok.String(), " p.lit:", p.lit,
	// 	" t.IsEOF:", p.tokenizer.IsEOF(), " p.IsEOF:", p.IsEOF(),
	// 	" depth:", p.exprLev)
}

// the tokenizer is always one token ahead, we can use its tok value
// to look ahead
func (p *Parser) lookahead() Token { return p.tokenizer.Peek() }

// Note that the parser's EOF is separate from the tokenizer's.
// the Parser's EOF should arrive one iteration of Parser.next()
// later than the tokenizer
func (p *Parser) IsEOF() bool { return p.tok == EOF }

// Expect a closing expression, when such an expression is encountered,
// The parser will attempt to close off the matching expression
func (p *Parser) expect(lit string) {
	p.expecting = append(p.expecting, lit)
}

func (p *Parser) dropExpect(lit string) {
	if p.expecting[len(p.expecting)-1] != lit {
		p.eh.AddErr(ERR_MISSING_CLOSE, "Expected '"+p.expecting[len(p.expecting)-1]+
			"', got '"+lit+"' instead")
	} else {
		p.expecting = p.expecting[0 : len(p.expecting)-1]
	}
}

func (p *Parser) matchExpectation(lit string) bool {
	if p.exprLev <= 0 {
		return false
	}
	for i := len(p.expecting); i > 0; i-- {
		if p.expecting[i] == lit {
			return true
		}
	}
	return false
}

func (p *Parser) parseTopLevel() *UnboundCompExpr {
	tree := new(UnboundCompExpr)
	for !p.IsEOF() {
		tree.AppendChildren(p.parseGenericOnce())
	}
	// println(tree.VisualizeTree())
	return tree
}

// parse one token
func (p *Parser) parseGenericOnce() Expr {
	// println("--\nParser.parseGeneric(): p.lit is \"", p.lit,
	// 	"\", token ", p.tok.String(), " depth: ", p.exprLev)
	switch p.tok {
	case CMDSTR:
		return p.parseStringCmd()
	case CMDSYM:
		return p.parseSymbolCmd()
	case NUM:
		return p.parseNumLit()
	case VARLIT:
		return p.parseVarLit()
	case SYM:
		return p.parseSimpleOpLit()
	case LBRACE:
		return p.parseCompositeExpr()
	case CARET:
		return p.parseSuperExpr()
	case UNDERSCORE:
		return p.parseSubExpr()
	case RBRACE:
		if p.matchExpectation(p.lit) {
			return &EmptyExpr{}
		}
		// FIXME case RBRACE is to catch unclosed paired objects (e.g. \left \right)
		// this error handling is not the best
		p.eh.AddErr(ERR_UNMATCHED_CLOSE, fmt.Sprintf("before cursor %d, at token %s of type %s",
			p.tokenizer.Cursor, p.lit, p.tok.String()))
	}
	// println("BadExpr!")
	p.next()
	return &BadExpr{}
}

func (p *Parser) parseStringCmd() Expr {
	kind := MatchLatexCmd(p.lit)

	var leaf Expr

	switch {
	case kind.TakesRawStrArg():
		leaf = p.parseTextCommand(kind)
	case kind.IsVanillaSym():
		leaf = &(SimpleCmdLit{Source: p.lit, Type: kind})
		p.next()
	case kind.TakesOneArg():
		leaf = p.parseCmd1Arg(kind)
	case kind.TakesTwoArg():
		leaf = p.parseCmd2Arg(kind)
	case kind.IsEnclosing():
		leaf = p.parseCmdEnclosing(kind)
	case kind == CMD_UNKNOWN:
		leaf = &(UnknownCmdLit{Source: p.lit})
		p.next()
	default:
		// this shouldn't be triggered
		leaf = &(BadExpr{})
		p.next()
	}

	// p.next() // FIXME next() should not be called here, but beware to call it
	// appropriately from within the above parse functions
	return leaf
}

// FIXME merge into parseStringCmd?
func (p *Parser) parseSymbolCmd() Expr {
	leaf := SimpleCmdLit{
		Source: p.lit,
	}
	p.next()
	return &leaf
}

func (p *Parser) parseNumLit() Expr {
	leaf := NumberLit{
		Source: p.lit,
	}
	p.next()
	return &leaf
}

func (p *Parser) parseVarLit() Expr {
	leaf := VarLit{
		Source: p.lit,
	}
	p.next()
	return &leaf
}

func (p *Parser) parseSimpleOpLit() Expr {
	leaf := SimpleOpLit{
		Source: p.lit,
	}
	p.next()
	return &leaf
}

func (p *Parser) parseCompositeExpr() Expr {
	p.exprLev++
	p.expect("}")
	p.next() // skip "{"
	node := new(CompositeExpr)
	for !p.IsEOF() && p.tok != RBRACE {
		node.AppendChildren(p.parseGenericOnce())
		// println("add child to node; depth: ", p.exprLev)
	}
	if p.IsEOF() {
		// FIXME error handling
		panic("expecting '}' got EOF")
	}
	p.next() // skip "}"
	p.dropExpect("}")
	p.exprLev--
	return node
}

func (p *Parser) parseSuperExpr() Expr {
	p.exprLev++
	p.next() // skip "^"
	node := new(SuperExpr)
	node.X = p.parseGenericOnce()
	p.exprLev--
	return node
}

// this should be merged with parseSuperExpr(), and maybe even parseFuncCmd1
func (p *Parser) parseSubExpr() Expr {
	p.exprLev++
	p.next() // skip "_"
	node := new(SubExpr)
	node.X = p.parseGenericOnce()

	p.exprLev--
	return node
}

func (p *Parser) parseTextCommand(kind LatexCmd) Expr {
	p.exprLev++
	p.next() // skip command
	node := &TextContainer{Text: new(RawStringLit)}
	if p.tok != LBRACE {
		node.Text.Text = p.lit
		p.next()
		return node
	}

	node.Text.Text = p.tokenizer.SkipToDelimiter("}")

	p.next() // skip }
	p.exprLev--
	return node
}

// parse a Command that takes one arguement
func (p *Parser) parseCmd1Arg(kind LatexCmd) Expr {
	p.exprLev++
	p.next() // skip command
	node := &Cmd1ArgExpr{Type: kind}
	node.Arg1 = p.parseGenericOnce()

	p.exprLev--
	return node
}

// parse a Command that takes two arguement
func (p *Parser) parseCmd2Arg(kind LatexCmd) Expr {
	p.exprLev++
	p.next() // skip "\command"
	node := &Cmd2ArgExpr{Type: kind}
	node.Arg1 = p.parseGenericOnce()
	node.Arg2 = p.parseGenericOnce()

	p.exprLev--
	return node
}

func (p *Parser) parseCmdEnclosing(kind LatexCmd) Expr {
	p.exprLev++
	p.expect("\\right")
	p.next() // skip "\left"
	node := new(ParenCompExpr)
	switch p.lit {
	case "(", "[", "\\{":
	default:
		panic("\\left expected '(', '[' or '\\{' but got " + p.lit)
	}
	node.Left = p.lit
	p.next() // skip left parenthesis e.g. "("
	for !p.IsEOF() && p.lit != "\\right" {
		node.AppendChildren(p.parseGenericOnce())
	}
	if p.IsEOF() {
		// FIXME error handling
		panic("expecting `\\right` got EOF")
	}
	p.next() // skip "\right"
	switch {
	case node.Left == "(" && p.lit != ")":
		panic("\\right expected ')' but got " + p.lit)
	case node.Left == "[" && p.lit != "]":
		panic("\\right expected ']' but got " + p.lit)
	case node.Left == "\\{" && p.lit != "\\}":
		panic("\\right expected '\\}' but got " + p.lit)
	}
	node.Right = p.lit
	p.next()
	p.dropExpect("\\right")
	p.exprLev--
	return node
}
