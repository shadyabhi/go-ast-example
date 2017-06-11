package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/Sirupsen/logrus"
)

// ParsedFile  keeps info regarding parsing the file in question
type ParsedFile struct {
	af *ast.File
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
}
func main() {
	fset := token.NewFileSet()

	var err error
	file, err := parser.ParseFile(fset, "./testcode/code.go", nil, 0)
	if err != nil {
		logrus.Fatalf("Couldn't parse file: %s", err)
	}

	// Print AST for debugging
	// ast.Print(fset, file)
	ifImport := checkImported(file.Imports)
	if ifImport {
		logrus.Infoln("logrus imports found")
		// Let's find line numbers
		ast.Inspect(file, func(n ast.Node) bool {
			return astWalker(file, fset, n)
		})
	}
}

// checkImported checks if logrus was imported in this file
func checkImported(imports []*ast.ImportSpec) bool {
	for _, imp := range imports {
		if strings.Contains(imp.Path.Value, "logrus") {
			return true
		}
	}
	return false
}

// astWalker walks the AST
func astWalker(file *ast.File, fset *token.FileSet, n ast.Node) bool {
	switch stmt := n.(type) {
	case *ast.CallExpr:
		expr, ok := stmt.Fun.(*ast.SelectorExpr)
		if !ok {
			break
		}
		ident := expr.X.(*ast.Ident)
		if ident.Name == "logrus" {
			ln := fset.Position(stmt.Pos()).Line
			msg := stmt.Args[0].(*ast.BasicLit).Value
			logrus.Infof("At line %d, msg: %s", ln, msg)
		}
	}
	return true
}
