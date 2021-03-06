package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/dnephin/gotestlint"
	"github.com/pkg/errors"
)

type options struct {
	path string
}

func main() {
	opts, err := parseOptions()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error())
	}
	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error())
	}
}

func parseOptions() (options, error) {
	args := os.Args[1:]
	if len(args) != 1 {
		return options{}, errors.New("requires 1 argument")
	}
	return options{path: args[0]}, nil
}

func run(opts options) error {
	calls, err := gotestlint.TestCasesFromDir(opts.path)
	if err != nil {
		return err
	}

	grouped := groupByFilename(calls)
	for filename, calls := range grouped {
		fmt.Printf("File: %s\n", filename)
		for _, call := range calls {
			fmt.Print(formatCalls(call))
		}
	}
	return nil
}

func groupByFilename(calls []gotestlint.TestCase) map[string][]gotestlint.TestCase {
	out := make(map[string][]gotestlint.TestCase)
	for _, call := range calls {
		out[call.Filename] = append(out[call.Filename], call)
	}
	return out
}

func formatCalls(testcalls gotestlint.TestCase) string {
	buf := new(bytes.Buffer)
	buf.WriteString("  " + testcalls.Testname + "\n")
	for _, call := range testcalls.FuncCalls {
		buf.WriteString(fmt.Sprintf("    %s()\n", call.Fun))
	}
	return buf.String()
}
