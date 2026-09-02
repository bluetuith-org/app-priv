package main

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
)

type genOptions struct {
	pkgName string

	currentFilePath, currentStructName           string
	newFilePath, newStructName, newStructComment string
}

type genTempl interface {
	AppendAccessor(s string, tag string) (genTemplRet, error)
	GetPartialCode() string
}

type genTemplRet struct {
	typeName   string
	found      bool
	removeTags bool
}

func newGenTemplRet(typeName string, found, removeTags bool) genTemplRet {
	return genTemplRet{typeName: typeName, found: found, removeTags: removeTags}
}

func emptyGenTemplRet() genTemplRet {
	return genTemplRet{}
}

func generateStruct(tpl genTempl, g *genOptions) error {
	fset := token.NewFileSet()

	cfgDst, err := decorator.ParseFile(fset, g.currentFilePath, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	var parentNode dst.Node
	var tspec *dst.TypeSpec
	var stype *dst.StructType

	dstutil.Apply(cfgDst, func(c *dstutil.Cursor) bool {
		n := c.Node()
		if n == nil {
			return true
		}

		switch v := n.(type) {
		case *dst.GenDecl:
			switch v.Tok {
			case token.TYPE:
				for _, spec := range v.Specs {
					ts, st, exists := matchStructType(spec, g.currentStructName)
					switch {
					case !exists:
						deleteNodeFromCursor(c)
						return true

					case exists:
						tspec = ts
						stype = st
						parentNode = v

						return true
					}
				}

			default:
				deleteNodeFromCursor(c)
			}

		case *dst.FuncDecl:
			if !matchFuncRecvType(v, g.newStructName) {
				deleteNodeFromCursor(c)
			}

		default:
		}

		return true
	}, nil)

	if err := modifyStruct(
		newModStructArg(fset, cfgDst, tpl, g).
			setNodes(parentNode, tspec, stype),
	); err != nil {
		return err
	}

	return writeCodeToFile(cfgDst, g.pkgName, g.newFilePath, true)
}

func matchStructType(n dst.Node, name string) (ts *dst.TypeSpec, st *dst.StructType, exists bool) {
	typeSpec, ok := n.(*dst.TypeSpec)
	if !ok {
		return nil, nil, false
	}

	st, ok = typeSpec.Type.(*dst.StructType)
	if !ok {
		return nil, nil, false
	}

	return typeSpec, st, typeSpec.Name.Name == name
}

func matchFuncRecvType(fnType *dst.FuncDecl, recvTypeName string) (exists bool) {
	recv := fnType.Recv
	if recv == nil || len(recv.List) == 0 {
		return false
	}

	for _, rt := range recv.List {
		expr := rt.Type
		if expr == nil {
			continue
		}

		switch v := expr.(type) {
		case *dst.Ident:
			if v.Name == recvTypeName {
				return true
			}

		case *dst.StarExpr:
			if ident, ok := v.X.(*dst.Ident); ok && ident.Name == recvTypeName {
				return true
			}
		}
	}

	return false
}

type modStructArg struct {
	fset *token.FileSet

	cfgFile    *dst.File
	parentNode dst.Node

	ts *dst.TypeSpec
	st *dst.StructType

	tpl genTempl
	*genOptions
}

func newModStructArg(fset *token.FileSet, cfgFile *dst.File, tpl genTempl, g *genOptions) *modStructArg {
	return &modStructArg{fset: fset, cfgFile: cfgFile, tpl: tpl, genOptions: g}
}

func (m *modStructArg) setNodes(parentNode dst.Node, ts *dst.TypeSpec, st *dst.StructType) *modStructArg {
	m.parentNode = parentNode
	m.ts = ts
	m.st = st

	return m
}

func modifyStruct(v *modStructArg) error {
	v.ts.Name.Name = v.newStructName

	decs := v.parentNode.Decorations()
	decs.Before = dst.EmptyLine
	decs.After = dst.EmptyLine
	decs.Start.Clear()
	decs.Start.Append(v.newStructComment)

	if err := genStructAndAccessors(v.st, v.tpl, v.genOptions, ""); err != nil {
		fmt.Println(err)
		return err
	}

	f, err := decorator.ParseFile(v.fset, "part.go", v.tpl.GetPartialCode(), parser.AllErrors)
	if err != nil {
		return err
	}

	v.cfgFile.Decls = append(v.cfgFile.Decls, f.Decls...)

	return nil
}

func genStructAndAccessors(st *dst.StructType, tpl genTempl, opt *genOptions, prefix string) error {
	for _, field := range st.Fields.List {
		for _, name := range field.Names {
			name := name.Name

			if prefix != "" {
				name = prefix + "." + name
			}

			switch ft := field.Type.(type) {
			case *dst.StructType:
				if err := genStructAndAccessors(ft, tpl, opt, name); err != nil {
					return err
				}

				continue

			case *dst.Ident:
				tag := ""
				if field.Tag != nil {
					tag = field.Tag.Value
				}

				ret, err := tpl.AppendAccessor(name, tag)
				if err != nil {
					return err
				}
				if ret.removeTags {
					field.Tag = nil
				}
				if !ret.found || ret.typeName == "" {
					continue
				}

				field.Type = &dst.Ident{
					Name: ret.typeName,
				}
			}
		}
	}

	return nil
}
