// This program verifies that any .go files given on the
// command line don't import packages on a pre-specified
// blocklist.
package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "need at least one file argument")
		os.Exit(1)
	}
	fset := token.NewFileSet()
	for _, filename := range os.Args[1:] {
		f, err := parser.ParseFile(fset, filename, nil, parser.ImportsOnly)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse %v: %v\n", filename, err)
			os.Exit(1)
		}
		for _, imp := range f.Imports {
			// remove leading and trailing double quote characters
			path := imp.Path.Value[1 : len(imp.Path.Value)-1]
			if blocklistPath(path) {
				fmt.Fprintf(os.Stderr, "illegal package in %v: %v\n", filename, path)
				os.Exit(1)
			}
		}
	}
}
