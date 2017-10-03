package gotestlint

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type TestCalls struct {
	Filename string
	Testname string
	Calls    []*ast.CallExpr
}

func findAllCalls(pkg *ast.Package) []TestCalls {
	fmt.Println("PACKAGE: ", pkg.Name)

	all := []TestCalls{}
	for filename, file := range pkg.Files {
		if !strings.HasSuffix(filename, "_test.go") {
			// TODO: handle these
			continue
		}
		fmt.Println("FILENAME: ", filename)

		for _, obj := range file.Scope.Objects {
			if !isTestFunc(obj) {
				continue
			}

			calls := TestCalls{
				Filename: filename,
				Testname: obj.Name,
			}
			fmt.Printf("OBJECT: %s\n", obj.Name)
			for _, stmt := range obj.Decl.(*ast.FuncDecl).Body.List {
				visitor := &astVisitor{}
				ast.Walk(visitor, stmt)
				calls.Calls = visitor.calls
			}

			all = append(all, calls)
		}
	}
	return all
}

func isTestFunc(obj *ast.Object) bool {
	return strings.HasPrefix(obj.Name, "Test") && obj.Kind == ast.Fun
}

func findInDir(path string) error {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, path, nil, parser.DeclarationErrors)
	if err != nil {
		return err
	}
	for _, pkg := range packages {
		// TODO: tmp hack around multiple packages
		if strings.HasSuffix(pkg.Name, "_test") {
			continue
		}
		findAllCalls(pkg)
	}
	return nil
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
