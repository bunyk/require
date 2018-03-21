package main

import (
	"testing" 
	"go/token"
	"reflect"
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

	expectedFiles := map[string]string{
		"fixtures/test.json": "{}\n",
	}
	if !reflect.DeepEqual(visitor.Files, expectedFiles) {
		t.Log("Got files:", visitor.Files)
		t.Fatal("Expected", expectedFiles)
	}

	expectedFileSeq := map[string][]string{
		"fixtures/plan*.txt": []string{
			"plan A\n",
			"plan B\n",
		},
	}
	if !reflect.DeepEqual(visitor.FileSequences, expectedFileSeq) {
		t.Log("Got file sequences:", visitor.FileSequences)
		t.Fatal("Expected", expectedFileSeq)
	}
}
