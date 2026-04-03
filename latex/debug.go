package latex

func (x *BadExpr) DeepEq(other Expr) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*BadExpr)
	if !ok {
		return false
	}
	return x.source == o.source && x.From == o.From && x.To == o.To
}

func (x *EmptyExpr) DeepEq(other Expr) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*EmptyExpr)
	if !ok {
		return false
	}
	return x.From == o.From && x.To == o.To && x.Type == o.Type
}

func (x *NumberLit) DeepEq(other Expr) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*NumberLit)
	if !ok {
		return false
	}
	return x.Source == o.Source && x.From == o.From && x.To == o.To
}

func (x *VarLit) DeepEq(other Expr) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*VarLit)
	if !ok {
		return false
	}
	return x.Source == o.Source && x.From == o.From && x.To == o.To
}

func (x *CompositeExpr) DeepEq(other Expr) bool {
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
		if !el.DeepEq(o.Elts[i]) {
			return false
		}
	}
	return true
}

func (x *UnboundCompExpr) DeepEq(other Expr) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*UnboundCompExpr)
	if !ok {
		return false
	}
	if x.From != o.From || x.To != o.To {
		return false
	}
	if len(x.Elts) != len(o.Elts) {
		return false
	}
	for i, el := range x.Elts {
		if !el.DeepEq(o.Elts[i]) {
			return false
		}
	}
	return true
}

func (x *ParenCompExpr) DeepEq(other Expr) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*ParenCompExpr)
	if !ok {
		return false
	}
	if x.From != o.From || x.To != o.To || x.Left != o.Left || x.Right != o.Right {
		return false
	}
	if len(x.Elts) != len(o.Elts) {
		return false
	}
	for i, el := range x.Elts {
		if !el.DeepEq(o.Elts[i]) {
			return false
		}
	}
	return true
}

func (x *EnvExpr) DeepEq(other Expr) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*EnvExpr)
	if !ok {
		return false
	}
	if x.Name != o.Name || x.From != o.From || x.To != o.To {
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
			if !cell.DeepEq(o.Elts[i][j]) {
				return false
			}
		}
	}
	return true
}

func (x *SimpleOpLit) DeepEq(other Expr) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*SimpleOpLit)
	if !ok {
		return false
	}
	return x.Source == o.Source && x.From == o.From && x.To == o.To
}

func (x *IncompleteCmdLit) DeepEq(other Expr) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*IncompleteCmdLit)
	if !ok {
		return false
	}
	return x.Backslash == o.Backslash && x.Source == o.Source && x.To == o.To
}

func (x *UnknownCmdLit) DeepEq(other Expr) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*UnknownCmdLit)
	if !ok {
		return false
	}
	return x.Backslash == o.Backslash && x.Source == o.Source && x.To == o.To
}

func (x RawRuneLit) DeepEq(other Expr) bool {
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
	if other == nil {
		return false
	}
	o, ok := other.(*SimpleCmdLit)
	if !ok {
		return false
	}
	return x.Backslash == o.Backslash && x.Source == o.Source && x.Type == o.Type && x.To == o.To
}

func (x *SuperExpr) DeepEq(other Expr) bool {
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
	return x.X.DeepEq(o.X)
}

func (x *SubExpr) DeepEq(other Expr) bool {
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
	return x.X.DeepEq(o.X)
}

func (x *Cmd1ArgExpr) DeepEq(other Expr) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*Cmd1ArgExpr)
	if !ok {
		return false
	}
	if x.Type != o.Type || x.Backslash != o.Backslash || x.To != o.To {
		return false
	}
	if x.Arg1 == nil && o.Arg1 == nil {
		return true
	}
	if x.Arg1 == nil || o.Arg1 == nil {
		return false
	}
	return x.Arg1.DeepEq(o.Arg1)
}

func (x *Cmd2ArgExpr) DeepEq(other Expr) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*Cmd2ArgExpr)
	if !ok {
		return false
	}
	if x.Type != o.Type || x.Backslash != o.Backslash || x.To != o.To {
		return false
	}
	if (x.Arg1 == nil) != (o.Arg1 == nil) {
		return false
	}
	if (x.Arg2 == nil) != (o.Arg2 == nil) {
		return false
	}
	if x.Arg1 != nil && !x.Arg1.DeepEq(o.Arg1) {
		return false
	}
	if x.Arg2 != nil && !x.Arg2.DeepEq(o.Arg2) {
		return false
	}
	return true
}

func (x *TextContainer) DeepEq(other Expr) bool {
	if other == nil {
		return false
	}
	o, ok := other.(*TextContainer)
	if !ok {
		return false
	}
	if x.CmdText != o.CmdText || x.Type != o.Type || x.From != o.From || x.To != o.To {
		return false
	}
	if x.Text == nil && o.Text == nil {
		return true
	}
	if x.Text == nil || o.Text == nil {
		return false
	}
	return x.Text.DeepEq(o.Text)
}

func (x *TextStringWrapper) DeepEq(other Expr) bool {
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
		if !r.DeepEq(o.Runes[i]) {
			return false
		}
	}
	return true
}
