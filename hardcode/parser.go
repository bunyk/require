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
	fmt.Printf("package %s\n\n", arguments["--package"])
	fmt.Println(`import "github.com/bunyk/require"`)
	fmt.Println("\nfunc init() {")
	for _, filename := range arguments["<filename>"].([]string) {
		processFile(filename)
	}
	fmt.Println("}")
}

func processFile(filename string) {
	src := readFile(filename)
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		log.Fatalf("Parse error: %s", err.Error())
	}

	ast.Walk(Visitor{
		FileSet: fset,
	}, f)
}

type Visitor struct {
	FileSet *token.FileSet
}

func (v Visitor) Error(pos token.Pos, msg string) {
	log.Fatalf("Error at %s: %s", v.FileSet.Position(pos), msg)
}

func (v Visitor) Visit(node ast.Node) ast.Visitor {
	call := isRequireFile(node)
	if call == nil {
		return v // Not yet found what we want, need to walk deeper
	}
	if len(call.Args) != 1 {
		v.Error(call.Lparen, "require.File() call requires one argument")
	}
	arg, ok := call.Args[0].(*ast.BasicLit)
	if !ok {
		v.Error(call.Lparen, "require.File() call should take constant argument")
	}
	if arg.Kind != token.STRING {
		v.Error(arg.ValuePos, "require.File() call should take string argument")
	}
	filename, err := strconv.Unquote(arg.Value)
	if err != nil {
		v.Error(arg.ValuePos, err.Error())
	}
	fmt.Printf("\trequire.SetFile(%#v, %#v)\n", filename, readFile(filename))
	return nil // found what we want, do not walk deeper
}

// checks if syntax tree node is require.File() call, and if not returns nil,
// othwerwise returns the call
func isRequireFile(node ast.Node) *ast.CallExpr {
	// We want find all function calls
	call, ok := node.(*ast.CallExpr)
	if !ok { // Not what we are looking for
		return nil
	}
	// We want function to be selected from package
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok { // not what we are looking for
		return nil
	}
	if isIdentifierNamed(selector.X, "require") &&
		isIdentifierNamed(selector.Sel, "File") {
		return call
	}
	return nil
}

func isIdentifierNamed(node ast.Node, name string) bool {
	ident, ok := node.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == name
}

func readFile(filename string) string {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not read file %s: %s", filename, err.Error())
		return ""
	}
	return string(src)
}
