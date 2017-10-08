package gotestlint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckNameNoError(t *testing.T) {
	testcase := TestCase{
		Testname: "TestFooBarWithSomething",
		FuncCalls: []FuncCall{
			{Name: "fooBar"},
		},
	}
	assert.NoError(t, CheckName(testcase))
}

func TestCheckNameNoMatch(t *testing.T) {
	testcase := TestCase{
		Testname: "TestFooBarWithSomething",
		FuncCalls: []FuncCall{
			{Name: "something"},
			{Name: "blah", Receiver: "Doer"},
		},
	}
	err := CheckName(testcase)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected prefix: TestSomething, TestDoerBlah")
}
