package find

import (
	"testing"
)

func TestSampleCaseCallsFunctionInCorrectFile(t *testing.T) {
	sample := Sample()
	t.Fatal(sample)
}
