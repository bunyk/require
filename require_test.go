package require

import "testing"

func TestHelloWorld(t *testing.T) {
	if require.File("fixtures/file.txt") != "Hello\nworld!\n" {
		t.Fatal("Not able to open file")
	}
}
