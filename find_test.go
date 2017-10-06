package gotestlint

import (
	"fmt"
	"testing"

	"github.com/gotestyourself/gotestyourself/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindAllCalls(t *testing.T) {
	dir := fs.NewDir(t, "test-find-all-calls", fs.FromDir("./testdata/find"))
	defer dir.Remove()

	testCases, err := TestCasesFromDir(dir.Path())
	require.NoError(t, err)
	assert.Len(t, testCases, 1)
	testCase := testCases[0]
	assert.Equal(t, dir.Join("find_test.go"), testCase.Filename)
	assert.Equal(t, "TestSampleCaseCallsFunctionInCorrectFile", testCase.Testname)
	assert.Len(t, testCase.FuncCalls, 1)
	assert.Equal(t, "Sample", fmt.Sprintf("%s", testCase.FuncCalls[0].Fun))
}
