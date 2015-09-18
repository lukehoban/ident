package ident

import "github.com/rogpeppe/godef/go/token"

// Definition is the code position that defines either a
// - func
// - var
// - interface or
// - struct field
type Definition struct {
	Name string
	token.Position
}

// Reference is a code position referencing a Definition
type Reference struct {
	token.Position
}

// Lookup the Definition of the identifier at the given byte offset
// in filepath
func Lookup(filepath string, offset int) (Definition, error) {
	return lookup(filepath, offset)
}

// FindReferences starts scanning valid go source files for references
// to Definition. If path is a folder, all contained files will be
// scanned. If recursive is true, all subfolders will be scanned recursively.
// Alternatively path can be a single go file.
// Found references and errors will be reported asynchronous.
func (def Definition) FindReferences(path string, recursive bool) (chan Reference, chan error) {
	return def.findReferences(path, recursive)
}
