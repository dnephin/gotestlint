package gotestlint

import (
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
	FuncCalls []*ast.CallExpr
}

func findAllCalls(pkg *ast.Package) []TestCase {
	all := []TestCase{}
	for filename, file := range pkg.Files {
		if !strings.HasSuffix(filename, "_test.go") {
			continue
		}

		for _, obj := range file.Scope.Objects {
			if !isTestFunc(obj) {
				continue
			}

			testcase := TestCase{Filename: filename, Testname: obj.Name}
			for _, stmt := range obj.Decl.(*ast.FuncDecl).Body.List {
				visitor := &astVisitor{}
				ast.Walk(visitor, stmt)
				testcase.FuncCalls = append(testcase.FuncCalls, visitor.calls...)
			}

			all = append(all, testcase)
		}
	}
	return all
}

func isTestFunc(obj *ast.Object) bool {
	return strings.HasPrefix(obj.Name, "Test") && obj.Kind == ast.Fun
}

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

type astVisitor struct {
	calls []*ast.CallExpr
}

func (v *astVisitor) Visit(node ast.Node) ast.Visitor {
	switch typed := node.(type) {
	case *ast.CallExpr:
		if isCallWithoutSelector(typed.Fun) {
			v.appendCall(typed)
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

func (v *astVisitor) appendCall(expr *ast.CallExpr) {
	v.calls = append(v.calls, expr)
}
