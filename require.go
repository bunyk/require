package require

var Files map[string]string
var FileSequences map[string][]string

// File receives filename as parameter and returns contents of that file as a string
func File(filename string) string {
	return Files[filename]
}

// File sequence receives file name pattern as parameter and returns slice of files
// that match that pattern. Pattern syntax is described here: https://golang.org/pkg/path/filepath/#Match
func FileSequence(namePattern string) []string {
	return FileSequences[namePattern]
}
