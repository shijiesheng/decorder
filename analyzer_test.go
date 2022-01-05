package decorder

import (
	"fmt"
	"os"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

//nolint:paralleltest
func TestAll(t *testing.T) {
	analysistest.Run(t, testdata(), Analyzer, "a")
}

//nolint:paralleltest
func TestCustomDecOrder(t *testing.T) {
	optBak := decOrder
	defer func() { decOrder = optBak }()
	decOrder = "func ,const,   var ,type"
	analysistest.Run(t, testdata(), Analyzer, "customDecOrderAll")
}

//nolint:paralleltest
func TestCustomDecOrderAll(t *testing.T) {
	optBak := decOrder
	defer func() { decOrder = optBak }()
	decOrder = "const,var"
	analysistest.Run(t, testdata(), Analyzer, "customDecOrder")
}

//nolint:paralleltest
func TestDisabledInitFuncFirstCheck(t *testing.T) {
	optBak := disableInitFuncFirstCheck
	defer func() { disableInitFuncFirstCheck = optBak }()
	disableInitFuncFirstCheck = true
	analysistest.Run(t, testdata(), Analyzer, "disabledInitFuncFirstCheck")
}

//nolint:paralleltest
func TestDisabledDecNumCheck(t *testing.T) {
	optBak := disableDecNumCheck
	defer func() { disableDecNumCheck = optBak }()
	disableDecNumCheck = true
	analysistest.Run(t, testdata(), Analyzer, "disabledDecNumCheck")
}

//nolint:paralleltest
func TestDisabledDecOrderCheck(t *testing.T) {
	optBak := disableDecOrderCheck
	defer func() { disableDecOrderCheck = optBak }()
	disableDecOrderCheck = true
	analysistest.Run(t, testdata(), Analyzer, "disabledDecOrderCheck")
}

func testdata() string {
	wd, _ := os.Getwd()
	return fmt.Sprintf("%s/testdata", wd)
}
