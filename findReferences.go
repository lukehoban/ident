package ident

import (
	"os"
	"path"
	"strings"

	"github.com/rogpeppe/godef/go/ast"
	"github.com/rogpeppe/godef/go/parser"
)

func (def Definition) findReferences(searchpath string, recursive bool) (chan Reference, chan error) {
	refs := make(chan Reference)
	errs := make(chan error, 1000)

	// returns true on error and reports it
	failed := func(err error) bool {
		if err != nil {
			select {
			case errs <- err:
			default:
			}
			return true
		}
		return false
	}

	scanAST := func(f ast.Node) {
		check := func(expr ast.Expr) {
			pos := getDefPosition(expr)
			if pos != nil && *pos == def.Position {
				refs <- Reference{fileset.Position(expr.Pos())}
			}
		}

		ast.Inspect(f, func(node ast.Node) bool {
			switch e := node.(type) {
			case *ast.SelectorExpr:
				if e.Sel.Name == def.Name {
					check(e)
				}
			case *ast.Ident:
				if e.Name == def.Name {
					check(e)
				}
			}
			return true
		})
	}

	scanFile := func(filepath string) {
		defer func() {
			if e := recover(); e != nil {
				return
			}
		}()
		f, err := parser.ParseFile(fileset, filepath, nil, 0, getScope(filepath), nil)
		if failed(err) {
			return
		}
		scanAST(f)
	}

	var scanFolder func(dirpath string)
	scanFolder = func(dirpath string) {
		filter := func(fi os.FileInfo) bool {
			return path.Ext(fi.Name()) == ".go"
		}
		defer func() {
			if e := recover(); e != nil {
				return
			}
		}()
		result, err := parser.ParseDir(fileset, dirpath, filter, 0, nil)
		if failed(err) {
			return
		}

		for _, pkg := range result {
			scanAST(pkg)
		}

		if !recursive {
			return
		}

		dir, err := os.Open(dirpath)
		if failed(err) {
			return
		}

		infos, err := dir.Readdir(0)
		if failed(err) {
			return
		}

		for _, fi := range infos {
			if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
				scanFolder(path.Join(dirpath, fi.Name()))
			}
		}
	}

	go func() {
		defer close(refs)
		defer close(errs)

		fi, err := os.Lstat(searchpath)
		if err != nil {
			return
		}
		if fi.IsDir() {
			scanFolder(searchpath)
		} else {
			scanFile(searchpath)
		}
	}()

	return refs, errs
}
