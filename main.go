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
	logrusImported := checkImported(file.Imports)

	var filename, funcName string
	var line, pos int

	if logrusImported {
		logrus.Debugf("logrus imports found")
		ast.Inspect(file, func(n ast.Node) bool {
			filename, funcName, line, pos = astWalker(file, fset, n)
			// Not all nodes have logrus
			if filename != "" {
				logrus.Infof("File=%s, Function=%s, LineNo=%d, Pos=%d", filename, funcName, line, pos)
			}
			return true
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
func astWalker(file *ast.File, fset *token.FileSet, n ast.Node) (filename, funcName string, line, pos int) {
	switch stmt := n.(type) {
	case *ast.FuncDecl:
		funcName = stmt.Name.Name
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
				logMsg := callExpr.Args[0].(*ast.BasicLit)
				_, _, pos = getContext(fset.Position(logMsg.ValuePos).String())
				filename, line, _ = getContext(fset.Position(locationPos).String())
			}
		}
	}
	return filename, funcName, line, pos
}

func getContext(s string) (filename string, line, pos int) {
	re := regexp.MustCompile(`.*\/(.*\.go):(\d+):(\d+)`)
	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) == 1 {
		line, _ := strconv.Atoi(matches[0][2])
		pos, _ := strconv.Atoi(matches[0][3])
		return matches[0][1], line, pos
	}
	return "", 0, 0
}
