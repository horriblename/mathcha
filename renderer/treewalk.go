// Copied straight from https://cs.opensource.google/go/go/+/refs/tags/go1.17.5:src/go/ast/walk.go
// and modified slightly to fit my use case
package renderer

import (
	parser "github.com/horriblename/mathcha/latex"
)

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node parser.Expr, dim *Dimensions) Visitor
}

// Helper functions for common node lists. They may be empty.

// func walkExprList(v Visitor, list []parser.Expr) {
// 	for _, x := range list {
// 		Walk(v, x, )
// 	}
// }

// Walk traverses an AST and a Dimensions tree in parallel
// in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
//
func Walk(v Visitor, node parser.Expr, dim *Dimensions) {
	if v = v.Visit(node, dim); v == nil {
		return
	}

	// walk children
	switch n := node.(type) {
	// Comments and fields
	case parser.Container:
		for i := range n.Children() {
			Walk(v, n.Children()[i], dim.Children[i])
		}
	}
	v.Visit(nil, dim)
}
