package tests

import (
	"github.com/bunyk/require"
	"testing"
)

func TestFilesPresense(t *testing.T) {
	if require.File("fixtures/file.txt") != "Hello\nworld!\n" {
		t.Fatal("Not able to open file fixtures/file.txt")
	}
	if require.FileSequence("fixtures/file*.txt") != "2\n" {
		t.Fatal("Not able to open file fixtures/file2.txt")
	}
}
