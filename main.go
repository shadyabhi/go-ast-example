package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strconv"
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
	case *ast.FuncDecl:
		funcName := stmt.Name.Name
		for _, s := range stmt.Body.List {
			exprStmt, ok := s.(*ast.ExprStmt)
			if !ok {
				continue
			}
			callExpr, ok := exprStmt.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}

			ident := selectorExpr.X.(*ast.Ident)
			if ident.Name == "logrus" {
				locationPos := ident.NamePos
				filename, line, pos := getContext(fset.Position(locationPos).String())
				logrus.Infof("File=%s, Function=%s, LineNo=%d, Pos=%d", filename, funcName, line, pos)
			}
		}
	}
	return true
}

func getContext(s string) (filename string, line, pos int) {
	re := regexp.MustCompile(".*\\/(.*\\.go):(\\d+):(\\d+)")
	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) == 1 {
		line, _ := strconv.Atoi(matches[0][2])
		pos, _ := strconv.Atoi(matches[0][3])
		return matches[0][1], line, pos
	}
	return "", 0, 0
}
