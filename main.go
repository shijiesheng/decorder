package main

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var (
	Analyzer = &analysis.Analyzer{
		Name:     "addlint",
		Doc:      "reports integer additions",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	decOrder string
)

func init() {
	Analyzer.Flags.StringVar(&decOrder, "decorder", "type,const,var,func", "define the order of types, constants, variables and functions declarations inside a file")
}

func main() {
	singlechecker.Main(Analyzer)
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, runDeclNumCheck(pass))
		ast.Inspect(f, runInitFuncFirstCheck(pass))
	}

	return nil, nil
}

func runInitFuncFirstCheck(pass *analysis.Pass) func(ast.Node) bool {
	nonInitFound := false

	return func(n ast.Node) bool {
		dec, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		if dec.Name.Name == "init" {
			if nonInitFound {
				pass.Reportf(dec.Pos(), "init func must be the first function in file")
			}
		} else {
			nonInitFound = true
		}

		return true
	}
}

func runDeclNumCheck(pass *analysis.Pass) func(ast.Node) bool {
	ts := []token.Token{token.TYPE, token.CONST, token.VAR}

	counts := map[token.Token]int{}
	for _, t := range ts {
		counts[t] = 0
	}

	var funcPoss []struct {
		start token.Pos
		end   token.Pos
	}

	return func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			funcPoss = append(funcPoss, struct {
				start token.Pos
				end   token.Pos
			}{start: fn.Pos(), end: fn.End()})

			return true
		}

		dn, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		for _, poss := range funcPoss {
			if poss.start < dn.Pos() && poss.end > dn.Pos() {
				return true
			}
		}

		for _, t := range ts {
			if dn.Tok == t {
				counts[t]++

				if counts[t] > 1 {
					pass.Reportf(dn.Pos(), "multiple \"%s\" declarations are not allowed; use parentheses instead", t.String())
				}
			}
		}

		return true
	}
}
