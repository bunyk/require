package main

import (
	"testing" 
	"go/token"
)

func TestProcessFile(t *testing.T) {
	visitor := Visitor{
		FileSet: token.NewFileSet(),
		Files: make(map[string]string),
		FileSequences: make(map[string][]string),
	}
	processFile("test.go", `package test
	func init() {
		a := require.File("fixtures/test.json")
		b := require.FileSequence("fixtures/plan*.txt")
	}`, &visitor)
	t.Fatal("No tests")
}
