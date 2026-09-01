package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
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

type genOptions struct {
	pkgName string

	currentFilePath, currentStructName           string
	newFilePath, newStructName, newStructComment string
}

type genTempl interface {
	AppendAccessor(s string) (string, bool)
	GetPartialCode() string
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

	modifyStruct(
		newModStructArg(fset, cfgDst, tpl, g).
			setNodes(parentNode, tspec, stype),
	)

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

func modifyStruct(v *modStructArg) {
	v.ts.Name.Name = v.newStructName

	decs := v.parentNode.Decorations()
	decs.Before = dst.EmptyLine
	decs.After = dst.EmptyLine
	decs.Start.Clear()
	decs.Start.Append(v.newStructComment)

	genStructAndAccessors(v.st, v.tpl, v.genOptions, "")

	f, err := decorator.ParseFile(v.fset, "part.go", v.tpl.GetPartialCode(), parser.AllErrors)
	if err != nil {
		exit(err)
	}

	v.cfgFile.Decls = append(v.cfgFile.Decls, f.Decls...)
}

func genStructAndAccessors(st *dst.StructType, tpl genTempl, opt *genOptions, prefix string) {
	for _, field := range st.Fields.List {
		for _, name := range field.Names {
			name := name.Name

			if prefix != "" {
				name = prefix + "." + name
			}

			switch ft := field.Type.(type) {
			case *dst.StructType:
				genStructAndAccessors(ft, tpl, opt, name)
				continue

			case *dst.Ident:
				typeName, ok := tpl.AppendAccessor(name)
				if !ok || typeName == "" {
					continue
				}

				field.Type = &dst.Ident{
					Name: typeName,
				}
			}
		}
	}
}

func deleteNodeFromCursor(c *dstutil.Cursor) {
	if c.Index() >= 0 {
		c.Delete()
	}
}

func prepareAllImports(pkgName string) (impDir string, closeFunc func(), err error) {
	impDir, err = os.MkdirTemp("", "bluetuith-themegen-*")
	if err != nil {
		return "", nil, err
	}

	projPath, err := projDirPath()
	if err != nil {
		return "", nil, err
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedImports | packages.NeedFiles,
		Dir:  projPath,
	}, "./...")
	if err != nil {
		return "", nil, err
	}

	var sb strings.Builder

	fmt.Fprintf(&sb, "package %s\n\nimport (\n", pkgName)

	importMap := make(map[string]string)
	for _, pkg := range pkgs {
		for _, imp := range pkg.Imports {
			importMap[imp.PkgPath] = imp.Name
		}
	}

	for importPath := range importMap {
		fmt.Fprintf(&sb, `"%s"`, importPath)
		sb.WriteString("\n")
	}

	sb.WriteString(")")

	fp := path.Join(impDir, "imports.go")
	file, err := openFileToWrite(fp)
	if err != nil {
		return "", nil, err
	}

	file.WriteString(sb.String())

	if err := file.Close(); err != nil {
		return "", nil, err
	}

	return impDir, func() {
		if err := os.Remove(fp); err != nil {
			panic(err)
		}

		if err := os.Remove(impDir); err != nil {
			panic(err)
		}
	}, nil
}

func writeCodeToFile(generatedDst *dst.File, pkgName, filePath string, writeGenComment bool) error {
	fileName := filepath.Base(filePath)
	impDir, closeFunc, err := prepareAllImports(pkgName)
	if err != nil {
		return err
	}
	defer closeFunc()

	var b bytes.Buffer
	defer b.Reset()

	decorator.Fprint(&b, generatedDst)

	proc, err := imports.Process(filepath.Join(impDir, fileName), b.Bytes(), nil)
	if err != nil {
		return err
	}

	file, err := openFileToWrite(filePath)
	if err != nil {
		return err
	}

	if writeGenComment {
		_, err = file.WriteString("// Code generated by structgen. DO NOT EDIT.\n\n")
		if err != nil {
			return err
		}
	}

	_, err = file.Write(proc)
	if err != nil {
		return err
	}

	return file.Close()
}

func projDirPath() (string, error) {
	cmd := exec.Command("go", "env", "GOMOD")

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	str := string(out)
	str = strings.TrimSpace(str)

	if str == "" {
		return "", errors.New("GOMOD not found")
	}

	return filepath.Dir(filepath.Clean(str)), nil
}

func openFileToWrite(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	return file, err
}

func execute(cmd string, params ...string) ([]byte, error) {
	c := exec.Command(cmd, params...)

	return c.Output()
}

//revive:disable
func exit(msg any) {
	fmt.Println(msg)

	os.Exit(1)
}
