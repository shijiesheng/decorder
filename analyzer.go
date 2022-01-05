package decorder

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	Name = "decorder"

	FlagDo    = "dec-order"
	FlagDdnc  = "disable-dec-num-check"
	FlagDdoc  = "disable-dec-order-check"
	FlagDiffc = "disable-init-func-first-check"
)

var (
	Analyzer = &analysis.Analyzer{
		Name: Name,
		Doc:  "check declaration order and count of types, constants, variables and functions",
		Run:  run,
	}

	decOrder                  string
	disableDecNumCheck        bool
	disableDecOrderCheck      bool
	disableInitFuncFirstCheck bool
)

//nolint:lll
func init() {
	Analyzer.Flags.StringVar(&decOrder, FlagDo, "type,const,var,func", "define the required order of types, constants, variables and functions declarations inside a file")
	Analyzer.Flags.BoolVar(&disableDecNumCheck, FlagDdnc, false, "option to disable check for number of e.g. var declarations inside file")
	Analyzer.Flags.BoolVar(&disableDecOrderCheck, FlagDdoc, false, "option to disable check for order of declarations inside file")
	Analyzer.Flags.BoolVar(&disableInitFuncFirstCheck, FlagDiffc, false, "option to disable check that init function is always first function in file")
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, runDeclNumCheck(pass))

		if !disableInitFuncFirstCheck {
			ast.Inspect(f, runInitFuncFirstCheck(pass))
		}
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

//nolint:funlen,gocognit
func runDeclNumCheck(pass *analysis.Pass) func(ast.Node) bool {
	ts := []token.Token{token.TYPE, token.CONST, token.VAR, token.FUNC}

	tsMap := map[string]token.Token{}
	counts := map[token.Token]int{}
	for _, t := range ts {
		counts[t] = 0
		tsMap[t.String()] = t
	}

	var funcPoss []struct {
		start token.Pos
		end   token.Pos
	}

	dos := strings.Split(decOrder, ",")

	for i, do := range dos {
		dos[i] = strings.TrimSpace(do)
	}

	check := func(t token.Token) (string, bool) {
		for i, do := range dos {
			if do == t.String() {
				for j := i + 1; j < len(dos); j++ {
					if counts[tsMap[dos[j]]] > 0 {
						return dos[j], false
					}
				}
				return "", true
			}
		}

		return "", true
	}

	return func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			funcPoss = append(funcPoss, struct {
				start token.Pos
				end   token.Pos
			}{start: fn.Pos(), end: fn.End()})

			counts[token.FUNC]++

			if !disableDecOrderCheck {
				o, c := check(token.FUNC)
				if !c {
					pass.Reportf(fn.Pos(), "%s must not be placed after %s", token.FUNC.String(), o)
				}
			}

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

		if !disableDecNumCheck {
			for _, t := range ts {
				if dn.Tok == t {
					counts[t]++

					if counts[t] > 1 {
						pass.Reportf(dn.Pos(), "multiple \"%s\" declarations are not allowed; use parentheses instead", t.String())
					}
				}
			}
		}

		if !disableDecOrderCheck {
			o, c := check(dn.Tok)
			if !c {
				pass.Reportf(dn.Pos(), "%s must not be placed after %s", dn.Tok.String(), o)
			}
		}

		return true
	}
}
