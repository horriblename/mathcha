package latex

type DeepEqCfg struct {
	SkipPos bool
}

func (x *BadExpr) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *BadExpr) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*BadExpr)
	if !ok {
		return false
	}
	if !cfg.SkipPos {
		if x.From != o.From || x.To != o.To {
			return false
		}
	}
	return x.source == o.source
}

func (x *EmptyExpr) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *EmptyExpr) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*EmptyExpr)
	if !ok {
		return false
	}
	if !cfg.SkipPos {
		if x.From != o.From || x.To != o.To {
			return false
		}
	}
	return x.Type == o.Type
}

func (x *NumberLit) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *NumberLit) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*NumberLit)
	if !ok {
		return false
	}
	if !cfg.SkipPos {
		if x.From != o.From || x.To != o.To {
			return false
		}
	}
	return x.Source == o.Source
}

func (x *VarLit) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *VarLit) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*VarLit)
	if !ok {
		return false
	}
	if !cfg.SkipPos {
		if x.From != o.From || x.To != o.To {
			return false
		}
	}
	return x.Source == o.Source
}

func (x *CompositeExpr) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *CompositeExpr) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*CompositeExpr)
	if !ok {
		return false
	}
	if x.Lbrace != o.Lbrace || x.Rbrace != o.Rbrace {
		return false
	}
	if len(x.Elts) != len(o.Elts) {
		return false
	}
	for i, el := range x.Elts {
		if !el.DeepEqWith(o.Elts[i], cfg) {
			return false
		}
	}
	return true
}

func (x *UnboundCompExpr) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *UnboundCompExpr) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*UnboundCompExpr)
	if !ok {
		return false
	}
	if !cfg.SkipPos {
		if x.From != o.From || x.To != o.To {
			return false
		}
	}
	if len(x.Elts) != len(o.Elts) {
		return false
	}
	for i, el := range x.Elts {
		if !el.DeepEqWith(o.Elts[i], cfg) {
			return false
		}
	}
	return true
}

func (x *ParenCompExpr) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *ParenCompExpr) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*ParenCompExpr)
	if !ok {
		return false
	}
	if !cfg.SkipPos {
		if x.From != o.From || x.To != o.To {
			return false
		}
	}
	if x.Left != o.Left || x.Right != o.Right {
		return false
	}
	if len(x.Elts) != len(o.Elts) {
		return false
	}
	for i, el := range x.Elts {
		if !el.DeepEqWith(o.Elts[i], cfg) {
			return false
		}
	}
	return true
}

func (x *EnvExpr) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *EnvExpr) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*EnvExpr)
	if !ok {
		return false
	}
	if !cfg.SkipPos {
		if x.From != o.From || x.To != o.To {
			return false
		}
	}
	if x.Name != o.Name {
		return false
	}
	if len(x.Elts) != len(o.Elts) {
		return false
	}
	for i, row := range x.Elts {
		if len(row) != len(o.Elts[i]) {
			return false
		}
		for j, cell := range row {
			if !cell.DeepEqWith(o.Elts[i][j], cfg) {
				return false
			}
		}
	}
	return true
}

func (x *SimpleOpLit) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *SimpleOpLit) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*SimpleOpLit)
	if !ok {
		return false
	}
	if !cfg.SkipPos {
		if x.From != o.From || x.To != o.To {
			return false
		}
	}
	return x.Source == o.Source
}

func (x *IncompleteCmdLit) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *IncompleteCmdLit) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*IncompleteCmdLit)
	if !ok {
		return false
	}
	if cfg.SkipPos {
		if x.To != o.To {
			return false
		}
	}
	return x.Backslash == o.Backslash && x.Source == o.Source
}

func (x *UnknownCmdLit) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *UnknownCmdLit) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*UnknownCmdLit)
	if !ok {
		return false
	}
	if cfg.SkipPos {
		if x.To != o.To {
			return false
		}
	}
	return x.Backslash == o.Backslash && x.Source == o.Source
}

func (x RawRuneLit) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x RawRuneLit) DeepEqWith(other Expr, _ DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(RawRuneLit)
	if !ok {
		return false
	}
	return x == o
}

func (x *SimpleCmdLit) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *SimpleCmdLit) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*SimpleCmdLit)
	if !ok {
		return false
	}
	if cfg.SkipPos {
		if x.Backslash != o.Backslash || x.To != o.To {
			return false
		}
	}
	return x.Source == o.Source && x.Type == o.Type
}

func (x *SuperExpr) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *SuperExpr) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*SuperExpr)
	if !ok {
		return false
	}
	if x.Symbol != o.Symbol || x.Close != o.Close {
		return false
	}
	if x.X == nil && o.X == nil {
		return true
	}
	if x.X == nil || o.X == nil {
		return false
	}
	return x.X.DeepEqWith(o.X, cfg)
}

func (x *SubExpr) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *SubExpr) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*SubExpr)
	if !ok {
		return false
	}
	if x.Symbol != o.Symbol || x.Close != o.Close {
		return false
	}
	if x.X == nil && o.X == nil {
		return true
	}
	if x.X == nil || o.X == nil {
		return false
	}
	return x.X.DeepEqWith(o.X, cfg)
}

func (x *Cmd1ArgExpr) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *Cmd1ArgExpr) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*Cmd1ArgExpr)
	if !ok {
		return false
	}
	if cfg.SkipPos {
		if x.Backslash != o.Backslash || x.To != o.To {
			return false
		}
	}
	if x.Type != o.Type {
		return false
	}
	if x.Arg1 == nil && o.Arg1 == nil {
		return true
	}
	if x.Arg1 == nil || o.Arg1 == nil {
		return false
	}
	return x.Arg1.DeepEqWith(o.Arg1, cfg)
}

func (x *Cmd2ArgExpr) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *Cmd2ArgExpr) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*Cmd2ArgExpr)
	if !ok {
		return false
	}
	if cfg.SkipPos {
		if x.Backslash != o.Backslash || x.To != o.To {
			return false
		}
	}
	if x.Type != o.Type {
		return false
	}
	if (x.Arg1 == nil) != (o.Arg1 == nil) {
		return false
	}
	if (x.Arg2 == nil) != (o.Arg2 == nil) {
		return false
	}
	if x.Arg1 != nil && !x.Arg1.DeepEqWith(o.Arg1, cfg) {
		return false
	}
	if x.Arg2 != nil && !x.Arg2.DeepEqWith(o.Arg2, cfg) {
		return false
	}
	return true
}

func (x *TextContainer) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *TextContainer) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*TextContainer)
	if !ok {
		return false
	}
	if cfg.SkipPos {
		if x.CmdText != o.CmdText || x.From != o.From || x.To != o.To {
			return false
		}
	}
	if x.Type != o.Type {
		return false
	}
	if x.Text == nil && o.Text == nil {
		return true
	}
	if x.Text == nil || o.Text == nil {
		return false
	}
	return x.Text.DeepEqWith(o.Text, cfg)
}

func (x *TextStringWrapper) DeepEq(other Expr) bool {
	return x.DeepEqWith(other, DeepEqCfg{})
}

func (x *TextStringWrapper) DeepEqWith(other Expr, cfg DeepEqCfg) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*TextStringWrapper)
	if !ok {
		return false
	}
	if len(x.Runes) != len(o.Runes) {
		return false
	}
	for i, r := range x.Runes {
		if !r.DeepEqWith(o.Runes[i], cfg) {
			return false
		}
	}
	return true
}
