package ident

import (
	"errors"

	"github.com/rogpeppe/godef/go/ast"
	"github.com/rogpeppe/godef/go/parser"
	"github.com/rogpeppe/godef/go/types"
)

func lookup(filepath string, offset int) (Definition, error) {
	def := Definition{}

	f, err := parser.ParseFile(fileset, filepath, nil, 0, getScope(filepath), nil)
	if err != nil {
		return def, err
	}

	containsOffset := func(node ast.Node) bool {
		from := fileset.Position(node.Pos()).Offset
		to := fileset.Position(node.End()).Offset
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

	pos := getDefPosition(ident)
	if pos == nil {
		return def, errors.New("could not find definition of identifier")
	}

	obj, _ := types.ExprType(ident, types.DefaultImporter, fileset)
	def.Name = obj.Name
	def.Position = *pos
	return def, nil
}
