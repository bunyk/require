package require

var Files map[string]string
var FileSequences map[string][]string

// File receives filename as parameter and returns contents of that file as a string
func File(filename string) string {
	return Files[filename]
}

// File sequence receives file name pattern (glob) as parameter and returns slice of files for that pattern
func FileSequence(namePattern string) []string {
	return FileSequences[filename]
}
