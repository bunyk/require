package tests

import "github.com/bunyk/require"

func init() {
	require.SetFile("fixtures/file2.txt", "2\n")
	require.SetFile("fixtures/file.txt", "Hello\nworld!\n")
}
