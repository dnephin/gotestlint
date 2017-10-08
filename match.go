package gotestlint

import (
	"github.com/pkg/errors"
	"strings"
)

// FindFunctionForTestCase compares the test case name to the functions called by the test
// case. If the name doesn't match any of the function calls then return an error.
func FindFunctionForTestCase(testcase TestCase) (Function, error) {
	prefixes := []string{}
	for _, funcCall := range testcase.FuncCalls {
		if strings.HasPrefix(testcase.Testname, funcCall.ExpectedTestPrefix()) {
			return funcCall, nil
		}
		prefixes = append(prefixes, funcCall.ExpectedTestPrefix())
	}
	return Function{}, errors.Errorf(
		"test name %q does not contain any expected prefix: %s",
		testcase.Testname, strings.Join(prefixes, ", "))
}

// CheckFileMatch checks that the function tested by
func CheckFileMatch(testcase TestCase, functions []Function) error {
	funcUnderTest, err := FindFunctionForTestCase(testcase)
	if err != nil {
		return err
	}
	for _, function := range functions {
		if funcUnderTest.Equal(function) {
			return nil
		}
	}
	return errors.Errorf("test for %s is in the wrong file", funcUnderTest)
}
