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

	directory, err := ParseDirectory(dir.Path())
	require.NoError(t, err)

	expected := &Directory{
		TestCases: map[string][]TestCase{
			dir.Join("find_test.go"): {
				{
					Testname:  "TestSampleCaseCallsFunctionInCorrectFile",
					FuncCalls: []Function{{Name: "Sample"}},
				},
			},
		},
		Functions: map[string][]Function{
			dir.Join("find.go"): {
				{Name: "Sample"},
				{Receiver: "Something", Name: "Do"},
			},
		},
	}
	assert.Equal(t, expected, directory)
}
