package gotestlint

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// TestCase stores the name of the testcase function, the filename of the file
// which contains the testcase, and a list of all the functions that are called
// by the testcase.
type TestCase struct {
	Filename  string
	Testname  string
	FuncCalls []FuncCall
}

// FuncCall stores data about a function call
type FuncCall struct {
	Receiver string
	Name     string
}

func (fc FuncCall) String() string {
	if fc.Receiver == "" {
		return fc.Name
	}
	return fmt.Sprintf("%s.%s", fc.Receiver, fc.Name)
}

// ExpectedTestPrefix returns the expected prefix for a test case that tests the
// function
func (fc FuncCall) ExpectedTestPrefix() string {
	return fmt.Sprintf("Test%s%s", strings.Title(fc.Receiver), strings.Title(fc.Name))
}

// TestCasesFromDir returns a list of TestCases found in the test files in
// a directory
func TestCasesFromDir(path string) ([]TestCase, error) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, path, nil, parser.DeclarationErrors)
	if err != nil {
		return nil, err
	}
	all := []TestCase{}
	for _, pkg := range packages {
		// TODO: tmp hack around multiple packages
		if strings.HasSuffix(pkg.Name, "_test") {
			continue
		}
		all = append(all, findAllCalls(pkg)...)
	}
	return all, nil
}

func findAllCalls(pkg *ast.Package) []TestCase {
	testcases := []TestCase{}
	for filename, file := range pkg.Files {
		if !isTestFile(filename) {
			continue
		}

		for _, obj := range file.Scope.Objects {
			if !isTestFunc(obj) {
				continue
			}

			testcase := TestCase{Filename: filename, Testname: obj.Name}
			for _, stmt := range obj.Decl.(*ast.FuncDecl).Body.List {
				calls := funcCallsFromTestCase(stmt)
				testcase.FuncCalls = append(testcase.FuncCalls, calls...)
			}
			testcases = append(testcases, testcase)
		}
	}
	return testcases
}

func isTestFile(filename string) bool {
	return strings.HasSuffix(filename, "_test.go")
}

func isTestFunc(obj *ast.Object) bool {
	return strings.HasPrefix(obj.Name, "Test") && obj.Kind == ast.Fun
}

func funcCallsFromTestCase(stmt ast.Node) []FuncCall {
	visitor := &astVisitor{}
	ast.Walk(visitor, stmt)
	return visitor.calls
}

type astVisitor struct {
	calls []FuncCall
}

func (v *astVisitor) Visit(node ast.Node) ast.Visitor {
	switch typed := node.(type) {
	case *ast.CallExpr:
		if isCallWithoutSelector(typed.Fun) {
			v.appendCall(typed.Fun.(*ast.Ident).Name, "")
		}
		// TODO: also find calls with selectors where the selector is the
		// struct that is being tested
	}
	return v
}

func isCallWithoutSelector(expr ast.Expr) bool {
	_, ok := expr.(*ast.Ident)
	return ok
}

func (v *astVisitor) appendCall(name, receiver string) {
	v.calls = append(v.calls, FuncCall{Name: name, Receiver: receiver})
}
