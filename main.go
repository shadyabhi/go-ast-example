package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/Sirupsen/logrus"
)

// ParsedFile is a single file out of a possible set contained in a run
// of deprehend
type ParsedFile struct {
	logrusLine []int
	af         *ast.File
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
}
func main() {
	pf := &ParsedFile{}
	fset := token.NewFileSet() // positions are relative to fset

	var err error
	pf.af, err = parser.ParseFile(fset, "./testcode/code.go", nil, 0)
	if err != nil {
		logrus.Fatalf("Couldn't parse file: %s", err)
	}

	// Print AST for debugging
	// ast.Print(fset, pf.af)
	ifImport := checkImported(pf.af.Imports)
	if ifImport {
		logrus.Infoln("logrus imports found")
		// Let's find line numbers
		ast.Inspect(pf.af, func(n ast.Node) bool {
			return astWalker(fset, n)
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
func astWalker(fset *token.FileSet, n ast.Node) bool {
	switch stmt := n.(type) {
	case *ast.CallExpr:
		expr, ok := stmt.Fun.(*ast.SelectorExpr)
		if !ok {
			break
		}
		ident := expr.X.(*ast.Ident)
		if ident.Name == "logrus" {
			ln := fset.Position(stmt.Pos()).Line
			logrus.Infof("Line number: %d", ln)
		}
	}
	return true
}
