package gotestlint

import (
	"github.com/pkg/errors"
	"strings"
)

// CheckName compares the test case name to the functions called by the test
// case. If the name doesn't match any of the function calls then return an error.
func CheckName(testcase TestCase) error {
	prefixes := []string{}
	for _, funcCall := range testcase.FuncCalls {
		if strings.HasPrefix(testcase.Testname, funcCall.ExpectedTestPrefix()) {
			return nil
		}
		prefixes = append(prefixes, funcCall.ExpectedTestPrefix())
	}
	return errors.Errorf("test name %q does not contain any expected prefix: %s",
		testcase.Testname, strings.Join(prefixes, ", "))
}
