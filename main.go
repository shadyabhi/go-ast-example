package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	fset := token.NewFileSet()

	var err error
	file, err := parser.ParseFile(fset, "./testcode/code.go", nil, 0)
	if err != nil {
		logrus.Fatalf("Couldn't parse file: %s", err)
	}

	astFile, err := os.OpenFile("./ast", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		logrus.Fatalf("Couldn't open file to write ast: %s", err)
	}
	// Print AST for debugging
	err = ast.Fprint(astFile, fset, file, func(name string, value reflect.Value) bool {
		return true
	})
	if err != nil {
		logrus.Fatalf("Error saving AST to file: %s", err)
	}
	astFile.Close()

	logrusName, logrusImported := checkImported(file.Imports)

	var filename, funcName string
	var line, pos int

	if logrusImported {
		logrus.Debugf("logrus imports found")
		ast.Inspect(file, func(n ast.Node) bool {
			filename, funcName, line, pos = astWalker(logrusName, fset, n)
			// Not all nodes have logrus
			if filename != "" {
				logrus.Debugf("File=%s, Function=%s, LineNo=%d, Pos=%d", filename, funcName, line, pos)
			}
			return true
		})
	}
	fset = token.NewFileSet()
	if err = printer.Fprint(os.Stdout, fset, file); err != nil {
		logrus.Fatalf("Error printing new AST")
	}
}

// checkImported checks if logrus was imported in this file
func checkImported(imports []*ast.ImportSpec) (string, bool) {
	var importName string
	for _, imp := range imports {
		if strings.Contains(imp.Path.Value, "logrus") {
			importName = imp.Name.Name
			return importName, true
		}
	}
	return "", false
}

// astWalker walks the AST
func astWalker(logrusName string, fset *token.FileSet, n ast.Node) (filename, funcName string, line, pos int) {
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
			if ident.Name == logrusName {
				locationPos := ident.NamePos
				logMsg := callExpr.Args[0].(*ast.BasicLit)
				_, _, pos = getContext(fset.Position(logMsg.ValuePos).String())
				filename, line, _ = getContext(fset.Position(locationPos).String())
				newLogMsg := &ast.BasicLit{
					ValuePos: logMsg.ValuePos,
					Kind:     token.STRING,
					Value:    fmt.Sprintf("\"%s:%s:%d - ", filename, funcName, line) + logMsg.Value[1:],
				}
				*logMsg = *newLogMsg
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
