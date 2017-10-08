package gotestlint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindFunctionForTestCaseNoError(t *testing.T) {
	testcase := TestCase{
		Testname: "TestFooBarWithSomething",
		FuncCalls: []Function{
			{Name: "fooBar"},
		},
	}
	function, err := FindFunctionForTestCase(testcase)
	require.NoError(t, err)
	assert.Equal(t, Function{Name: "fooBar"}, function)
}

func TestFindFunctionForTestCaseNoMatch(t *testing.T) {
	testcase := TestCase{
		Testname: "TestFooBarWithSomething",
		FuncCalls: []Function{
			{Name: "something"},
			{Name: "blah", Receiver: "Doer"},
		},
	}
	_, err := FindFunctionForTestCase(testcase)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected prefix: TestSomething, TestDoerBlah")
}
