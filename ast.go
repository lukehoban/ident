package ident

import (
	"path"

	"code.google.com/p/rog-go/exp/go/ast"
	"code.google.com/p/rog-go/exp/go/parser"
	"code.google.com/p/rog-go/exp/go/types"
)

var g_fileset = types.FileSet
var g_scopes = map[string]*ast.Scope{}

func getScope(filepath string) *ast.Scope {
	dirpath := path.Base(filepath)
	scope, ok := g_scopes[dirpath]
	if !ok {
		scope = ast.NewScope(parser.Universe)
		g_scopes[dirpath] = scope
	}
	return scope
}
