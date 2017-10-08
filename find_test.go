package gotestlint

import (
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

	expected := []TestCase{
		{
			Filename:  dir.Join("find_test.go"),
			Testname:  "TestSampleCaseCallsFunctionInCorrectFile",
			FuncCalls: []FuncCall{{Name: "Sample"}},
		},
	}
	assert.Equal(t, expected, testCases)
}
