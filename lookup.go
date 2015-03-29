package ident

import (
	"errors"

	"code.google.com/p/rog-go/exp/go/ast"
	"code.google.com/p/rog-go/exp/go/parser"
	"code.google.com/p/rog-go/exp/go/types"
)

func lookup(filepath string, offset int) (Definition, error) {
	def := Definition{}

	f, err := parser.ParseFile(g_fileset, filepath, nil, 0, getScope(filepath))
	if err != nil {
		return def, err
	}

	containsOffset := func(node ast.Node) bool {
		from := g_fileset.Position(node.Pos()).Offset
		to := g_fileset.Position(node.End()).Offset
		return offset >= from && offset < to
	}

	// traverse the ast tree until we find a node at the given offset position
	var ident ast.Expr
	ast.Inspect(f, func(node ast.Node) bool {
		switch expr := node.(type) {
		case *ast.SelectorExpr:
			if containsOffset(expr) && containsOffset(expr.Sel) {
				ident = expr
			}
		case *ast.Ident:
			if containsOffset(expr) {
				ident = expr
			}
		}
		return ident == nil
	})

	if ident == nil {
		return def, errors.New("no identifier found")
	}

	obj, _ := types.ExprType(ident, types.DefaultImporter)
	if obj == nil {
		return def, errors.New("identifier has no definition")
	}

	def.Name = obj.Name
	def.Position = g_fileset.Position(types.DeclPos(obj))
	return def, nil
}
