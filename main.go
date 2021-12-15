package main

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var Analyzer = &analysis.Analyzer{
	Name:     "addlint",
	Doc:      "reports integer additions",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var c = 4

const (
	a = 1
	b = 2
)

const hello = 123

func init() {
}

func main() {
	singlechecker.Main(Analyzer)
}

func init() {

}

func run(pass *analysis.Pass) (interface{}, error) {
	runDeclNumCheck(pass)
	runInitFuncFirstCheck(pass)

	return nil, nil
}

func runInitFuncFirstCheck(pass *analysis.Pass) {
	detective := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nonInitFound := false

	detective.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		dec, ok := n.(*ast.FuncDecl)
		if !ok {
			return
		}

		if dec.Name.Name == "init" {
			if nonInitFound {
				pass.Reportf(dec.Pos(), "init func must be the first function in file")
			}
		} else {
			nonInitFound = true
		}
	})
}

func runDeclNumCheck(pass *analysis.Pass) {
	detective := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	ts := []token.Token{token.TYPE, token.CONST, token.VAR}

	counts := map[token.Token]int{}
	for _, t := range ts {
		counts[t] = 0
	}

	var funcPoss []struct {
		start token.Pos
		end   token.Pos
	}

	detective.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		dec, ok := n.(*ast.FuncDecl)
		if !ok {
			return
		}

		funcPoss = append(funcPoss, struct {
			start token.Pos
			end   token.Pos
		}{start: dec.Pos(), end: dec.End()})
	})

	detective.Preorder([]ast.Node{(*ast.GenDecl)(nil)}, func(n ast.Node) {
		dec, ok := n.(*ast.GenDecl)
		if !ok {
			return
		}

		for _, poss := range funcPoss {
			if poss.start < dec.Pos() && poss.end > dec.Pos() {
				return
			}
		}

		for _, t := range ts {
			if dec.Tok == t {
				counts[t]++

				if counts[t] > 1 {
					pass.Reportf(dec.Pos(), "multiple \"%s\" declarations are not allowed; use parentheses instead", t.String())
				}
			}
		}
	})
}
