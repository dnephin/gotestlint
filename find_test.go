package gotestlint

import "testing"

func TestFindAllCalls(t *testing.T) {
	findInDir("/go/src/github.com/gotestyourself/gotestyourself/skip")
}
