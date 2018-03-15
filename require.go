package require

var files map[string]string

// File receives filename as parameter and returns contents of that file as a string
func File(filename string) string {
	return files[filename]
}

// SetFile sets contents for the file by name. Calls to this should be generated
func SetFile(filename, content string) {
	if files == nil {
		files = make(map[string]string)
	}
	files[filename] = content
}
