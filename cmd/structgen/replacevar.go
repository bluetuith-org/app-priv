package main

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
)

type varReplFunc func(n dst.Node) bool

type varReplaceOptions struct {
	varName, pkgName, filePath string

	Start int    `json:"start"`
	End   int    `json:"end"`
	Code  string `json:"code"`

	fn varReplFunc
}

func (v *varReplaceOptions) set(varname, pkgName, filePath string, fn varReplFunc) {
	v.varName = varname
	v.pkgName = pkgName
	v.filePath = filePath
	v.fn = fn
}

func replaceStructVar(vrepl varReplaceOptions) error {
	fset := token.NewFileSet()

	replDst, err := decorator.ParseFile(fset, vrepl.filePath, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	var origCmpLit *dst.CompositeLit

	dstutil.Apply(replDst, func(c *dstutil.Cursor) bool {
		n := c.Node()
		if n == nil {
			return true
		}

		v, ok := n.(*dst.GenDecl)
		if !ok || v.Tok != token.VAR {
			return true
		}

		for _, spec := range v.Specs {
			vspec, ok := spec.(*dst.ValueSpec)
			if !ok {
				continue
			}

			found := false

			for _, name := range vspec.Names {
				if name.Name == vrepl.varName {
					found = true
					break
				}
			}

			if !found {
				return true
			}

			for _, val := range vspec.Values {
				switch v := val.(type) {
				case *dst.CompositeLit:
					origCmpLit = v

				case *dst.UnaryExpr:
					if lit, ok := v.X.(*dst.CompositeLit); ok {
						origCmpLit = lit
					}

				default:
				}
			}
		}

		return true
	}, nil)

	if origCmpLit == nil {
		return errors.New("var: orig: parse error")
	}

	varDst, err := decorator.ParseFile(
		fset, "cmp_var.go",
		fmt.Sprintf("package %s\n\nvar _ = %s\n", vrepl.pkgName, vrepl.Code),
		parser.AllErrors,
	)
	if err != nil {
		return err
	}

	replCmpLit, ok := varDst.Decls[0].(*dst.GenDecl).Specs[0].(*dst.ValueSpec).Values[0].(*dst.CompositeLit)
	if !ok {
		return errors.New("var: repl: parse error")
	}

	dst.Inspect(replCmpLit, vrepl.fn)
	origCmpLit.Elts = replCmpLit.Elts

	return writeCodeToFile(replDst, vrepl.pkgName, vrepl.filePath, false)
}
