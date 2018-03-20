package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/docopt/docopt-go"
)

var usage string = `Program to hardcode file contents into your go code.

Usage:
  hardcode [--package=<package>] <filename>...
  hardcode -h | --help

Options:
  -h --help     Show this screen.
  --package=<package>    Package name for file [default: resources]
`

func main() {
	arguments, err := docopt.ParseDoc(usage)
	if err != nil {
		fmt.Println(err)
		fmt.Println(usage)
		return
	}
	visitor := Visitor{
		FileSet: token.NewFileSet(),
		Files: make(map[string]string),
		FileSequences: make(map[string][]string),
	}
	for _, filename := range arguments["<filename>"].([]string) {
		processFile(filename, readFile(filename), &visitor)
	}

	fmt.Printf("package %s\n\n", arguments["--package"])
	fmt.Println(`import "github.com/bunyk/require"`)
	fmt.Println("\nfunc init() {")
	fmt.Println("}")
}

func processFile(name, contents string, visitor *Visitor) {
	f, err := parser.ParseFile(visitor.FileSet, name, contents, 0)
	if err != nil {
		log.Fatalf("Parse error: %s", err.Error())
	}
	ast.Walk(*visitor, f)
}

type Visitor struct {
	FileSet *token.FileSet // Code to visit
	Files map[string]string // require.File occurences
	FileSequences map[string][]string // require.FileSequence occurences
}

// Generate error in some part of file
func (v Visitor) Error(pos token.Pos, msg string, args ...interface{}) {
	log.Fatalf(
		"Error at %s: %s",
		v.FileSet.Position(pos),
		fmt.Sprintf(msg, args...),
	)
}

func (v Visitor) Visit(node ast.Node) ast.Visitor {
	call, funcName := isRequireFileOrSequence(node)
	if call == nil {
		return v // Not yet found what we want, need to walk deeper
	}
	if len(call.Args) != 1 {
		v.Error(call.Lparen, "require.%s() call requires one argument", funcName)
	}
	arg, ok := call.Args[0].(*ast.BasicLit)
	if !ok {
		v.Error(call.Lparen, "require.%s() call should take constant argument", funcName)
	}
	if arg.Kind != token.STRING {
		v.Error(arg.ValuePos, "require.%s() call should take string argument", funcName)
	}
	filename, err := strconv.Unquote(arg.Value)
	if err != nil {
		v.Error(arg.ValuePos, err.Error())
	}
	if funcName == "File" {
		v.Files[filename] = readFile(filename)
		return nil
	}
	if funcName == "FileSequence" {
		fmt.Println("TODO: Should load following files:", filename)
	}
	return nil // found what we want, do not walk deeper
}

// checks if syntax tree node is require.File() or require.FileSequence() call,
// and if not returns nil, othwerwise returns the call and the function name
func isRequireFileOrSequence(node ast.Node) (*ast.CallExpr, string) {
	// We want find all function calls
	call, ok := node.(*ast.CallExpr)
	if !ok { // Not what we are looking for
		return nil, ""
	}
	// We want function to be selected from package
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok { // not what we are looking for
		return nil, ""
	}
	if getIdentifierName(selector.X) != "require" {
		return nil, ""
	}
	funcName := getIdentifierName(selector.Sel)
	if (funcName == "File" || funcName == "FileSequence") {
		return call, funcName
	}
	return nil, ""
}

// Return identifier name, or empty string if that is not identifier
func getIdentifierName(node ast.Node) string {
	ident, ok := node.(*ast.Ident)
	if !ok {
		return ""
	}
	return ident.Name
}

func readFile(filename string) string {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not read file %s: %s", filename, err.Error())
		return ""
	}
	return string(src)
}
