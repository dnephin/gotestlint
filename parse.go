package gotestlint

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Directory of go source files
type Directory struct {
	TestCases map[string][]TestCase
	Functions map[string][]Function
}

func newDirectory() *Directory {
	return &Directory{
		TestCases: make(map[string][]TestCase),
		Functions: make(map[string][]Function),
	}
}

// TestCase stores the name of the testcase function, the filename of the file
// which contains the testcase, and a list of all the functions that are called
// by the testcase.
type TestCase struct {
	Testname  string
	FuncCalls []Function
}

// Function stores data about a function call
type Function struct {
	Receiver string
	Name     string
}

func (fc Function) String() string {
	if fc.Receiver == "" {
		return fc.Name
	}
	return fmt.Sprintf("%s.%s", fc.Receiver, fc.Name)
}

// ExpectedTestPrefix returns the expected prefix for a test case that tests the
// function
// TODO: move to a function in match
func (fc Function) ExpectedTestPrefix() string {
	return fmt.Sprintf("Test%s%s", strings.Title(fc.Receiver), strings.Title(fc.Name))
}

// ParseDirectory returns a list of TestCases found in the test files in
// a directory
func ParseDirectory(path string) (*Directory, error) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, path, nil, parser.DeclarationErrors)
	if err != nil {
		return nil, err
	}

	directory := newDirectory()
	for _, pkg := range packages {
		// TODO: tmp hack around multiple packages
		if strings.HasSuffix(pkg.Name, "_test") {
			continue
		}
		getTestCases(directory, pkg)
		getFunctions(directory, pkg)
	}
	return directory, nil
}

func getTestCases(directory *Directory, pkg *ast.Package) {
	for filename, file := range pkg.Files {
		if !isTestFile(filename) {
			continue
		}

		testcases := []TestCase{}
		for _, obj := range file.Scope.Objects {
			if !isTestFunc(obj) {
				continue
			}

			testcase := TestCase{Testname: obj.Name}
			for _, stmt := range obj.Decl.(*ast.FuncDecl).Body.List {
				calls := funcCallsFromTestCase(stmt)
				testcase.FuncCalls = append(testcase.FuncCalls, calls...)
			}
			testcases = append(testcases, testcase)
		}
		directory.TestCases[filename] = testcases
	}
}

func isTestFile(filename string) bool {
	return strings.HasSuffix(filename, "_test.go")
}

func isTestFunc(obj *ast.Object) bool {
	return strings.HasPrefix(obj.Name, "Test") && obj.Kind == ast.Fun
}

func funcCallsFromTestCase(stmt ast.Node) []Function {
	visitor := &testCaseFuncCallVisitor{}
	ast.Walk(visitor, stmt)
	return visitor.calls
}

type testCaseFuncCallVisitor struct {
	calls []Function
}

func (v *testCaseFuncCallVisitor) Visit(node ast.Node) ast.Visitor {
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

func (v *testCaseFuncCallVisitor) appendCall(name, receiver string) {
	v.calls = append(v.calls, Function{Name: name, Receiver: receiver})
}

func getFunctions(directory *Directory, pkg *ast.Package) {
	for filename, file := range pkg.Files {
		if isTestFile(filename) {
			continue
		}

		visitor := &fileFunctionDeclVisitor{}
		ast.Walk(visitor, file)
		directory.Functions[filename] = visitor.calls
	}
}

type fileFunctionDeclVisitor struct {
	calls []Function
}

func (v *fileFunctionDeclVisitor) Visit(node ast.Node) ast.Visitor {
	switch typed := node.(type) {
	case *ast.FuncDecl:
		v.calls = append(v.calls, functionFromDecl(typed))
		return nil
	}
	return v
}

func functionFromDecl(decl *ast.FuncDecl) Function {
	function := Function{
		Name: decl.Name.Name,
	}
	if decl.Recv != nil {
		function.Receiver = fmt.Sprintf("%s", decl.Recv.List[0].Type)
	}
	return function
}
