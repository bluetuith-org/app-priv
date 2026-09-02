package main

import (
	"errors"
	"go/parser"
	"go/token"
	"slices"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

type fnStructOpts struct {
	fnName   string
	pkgName  string
	filePath string

	accMap           map[string]struct{}
	receiverToAppend string
}

type assignStmt struct {
	pos  int
	root string
	n    []dst.Stmt
}

type assignNode struct {
	pos     int
	assign  dst.Stmt
	baseAcc string
}

func genStructAccessorFn(f fnStructOpts) error {
	fset := token.NewFileSet()

	fnDst, err := decorator.ParseFile(fset, f.filePath, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	fnNode, err := findFuncByName(fnDst, f.fnName)
	if err != nil {
		return err
	}

	var (
		assignCount int
		maxIdx      int
		stmtList    []assignStmt
	)

	for r, n := range assignMapFn(fnNode, &f, true, func(idx int, stmt dst.Stmt, _ string) {
		decs := stmt.Decorations()
		decs.Before = dst.None
		decs.After = dst.None

		assignCount++
		maxIdx = idx
	}) {
		_, id, ok := strings.Cut(r, ".")
		if !ok {
			continue
		}

		st, ok := getAccNode(&f, id, n.baseAcc, n.pos)
		if ok {
			stmtList = append(stmtList, st)
		}
	}

	for v := range f.accMap {
		id, _, ok := strings.Cut(v, ".")
		if !ok {
			continue
		}

		st, ok := getAccNode(&f, id, f.receiverToAppend, maxIdx)
		if ok {
			stmtList = append(stmtList, st)
		}
	}

	offset := 0
	for _, stmt := range stmtList {
		fnNode.Body.List = slices.Insert(fnNode.Body.List, stmt.pos+offset, stmt.n...)
		offset += len(stmt.n)
		assignCount += len(stmt.n)
	}

	slices.SortStableFunc(fnNode.Body.List, func(a, b dst.Stmt) int {
		a1, ok := a.(*dst.AssignStmt)
		if !ok || a1.Tok == token.DEFINE {
			return 0
		}

		b1, ok := b.(*dst.AssignStmt)
		if !ok || b1.Tok == token.DEFINE {
			return 0
		}

		accSep1 := buildAccFromAssign(a1)
		accSep2 := buildAccFromAssign(b1)

		rootAcc1 := getRootAcc(accSep1)
		rootAcc2 := getRootAcc(accSep2)

		if len(rootAcc1) != len(rootAcc2) {
			return strings.Compare(rootAcc1, rootAcc2)
		}

		fullAcc1 := strings.Join(accSep1, ".")
		fullAcc2 := strings.Join(accSep2, ".")

		dotCount1 := strings.Count(fullAcc1, ".")
		dotCount2 := strings.Count(fullAcc2, ".")

		if dotCount1 != dotCount2 {
			return dotCount1 - dotCount2
		}

		return strings.Compare(fullAcc1, fullAcc2)
	})

	var (
		curRoot string
		curStmt dst.Stmt
	)

	assignMapFn(fnNode, &f, false, func(_ int, stmt dst.Stmt, root string) {
		endNode := curRoot == root && assignCount-1 == 0

		if curRoot != "" && curRoot != root || endNode {
			s := curStmt
			if endNode {
				s = stmt
			}

			decs := s.Decorations()
			decs.Before = dst.None
			decs.After = dst.EmptyLine
		}

		curRoot = root
		curStmt = stmt

		assignCount--
	})

	return writeCodeToFile(fnDst, f.pkgName, f.filePath, false)
}

func findFuncByName(fnDst *dst.File, name string) (*dst.FuncDecl, error) {
	var fnNode *dst.FuncDecl

	dst.Inspect(fnDst, func(n dst.Node) bool {
		if n == nil {
			return false
		}

		fnv, ok := n.(*dst.FuncDecl)
		if !ok {
			return true
		}

		if fnv.Name.Name == name {
			fnNode = fnv

			return false
		}

		return true
	})

	if fnNode == nil {
		return nil, errors.New("fn not found")
	}

	return fnNode, nil
}

func assignMapFn(fnNode *dst.FuncDecl, f *fnStructOpts, getMap bool, fn func(idx int, stmt dst.Stmt, root string)) map[string]assignNode {
	var assignMap map[string]assignNode
	if getMap {
		assignMap = make(map[string]assignNode)
	}

	for idx, stmt := range fnNode.Body.List {
		assign, ok := stmt.(*dst.AssignStmt)
		if !ok || assign.Tok == token.DEFINE {
			continue
		}

		var accSep []string

		dst.Inspect(assign, func(n dst.Node) bool {
			switch v := n.(type) {
			case *dst.Ident:
				accSep = append(accSep, v.Name)

			default:
			}

			return true
		})

		baseAcc := accSep[0]
		root := getRootAcc(accSep)

		acc := strings.Join(accSep[1:], ".")
		delete(f.accMap, acc)

		if fn != nil {
			fn(idx, assign, root)
		}

		if assignMap == nil {
			continue
		}

		assignMap[root] = assignNode{
			pos:     idx + 1,
			assign:  assign,
			baseAcc: baseAcc,
		}
	}

	return assignMap
}

func getRootAcc(accSep []string) string {
	return strings.Join(accSep[0:min(2, len(accSep)-1)], ".")
}

func getAccNode(f *fnStructOpts, prefix, baseAcc string, curPos int) (n assignStmt, ok bool) {
	sortedAcc := make([]string, 0, 2+len(baseAcc)+len(prefix))

	for k := range f.accMap {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		sortedAcc = append(sortedAcc, baseAcc+"."+k)
		delete(f.accMap, k)
	}
	if len(sortedAcc) == 0 {
		return n, false
	}

	slices.Sort(sortedAcc)

	n.pos = curPos

	for _, s := range sortedAcc {
		expr := buildSelectorExpr(s)

		var b dst.Stmt = &dst.AssignStmt{
			Tok: token.ASSIGN,
			Lhs: []dst.Expr{expr},
			Rhs: []dst.Expr{&dst.BasicLit{Kind: token.VAR, Value: `__UNDEFINED__`}},
		}

		n.n = append(n.n, b)
	}

	return n, true
}

func buildAccFromAssign(a *dst.AssignStmt) []string {
	var accSep []string

	dst.Inspect(a, func(n dst.Node) bool {
		switch v := n.(type) {
		case *dst.Ident:
			accSep = append(accSep, v.Name)

		default:
		}

		return true
	})

	return accSep
}
