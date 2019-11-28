package format

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type ImportGroups [][]string

func getAstString(fileSet *token.FileSet, node ast.Node) string {
	buf := &bytes.Buffer{}
	if err := format.Node(buf, fileSet, node); err != nil {
		panic(err)
	}
	return buf.String()
}

func Format(filename string, src []byte) []byte {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filename, src, parser.ParseComments)
	if err != nil {
		panic(fmt.Errorf("errors %s in %s", err.Error(), filename))
	}
	buf := &bytes.Buffer{}
	if err := format.Node(buf, fileSet, file); err != nil {
		panic(fmt.Errorf("errors %s in %s", err.Error(), filename))
	}
	return buf.Bytes()
}

func Process(filename string, src []byte) ([]byte, error) {
	cwd, _ := os.Getwd()
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	ast.SortImports(fileSet, file)

	formattedCode := getAstString(fileSet, file)

	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if genDecl.Tok != token.IMPORT {
				break
			}

			importsCode := getAstString(fileSet, genDecl)

			importGroups := make(ImportGroups, 4)
			for _, spec := range genDecl.Specs {
				importSpec := spec.(*ast.ImportSpec)
				importPath, _ := strconv.Unquote(importSpec.Path.Value)
				pkgs, err := packages.Load(nil, importPath)
				if err != nil {
					panic(fmt.Errorf("errors %s in %s", err.Error(), filename))
				}
				if strings.Contains(pkgs[0].GoFiles[0], runtime.GOROOT()) {
					// libexec
					importGroups[0] = append(importGroups[0], getAstString(fileSet, importSpec))
				} else {
					if strings.HasPrefix(pkgs[0].GoFiles[0], cwd) {
						importGroups[3] = append(importGroups[3], getAstString(fileSet, importSpec))
					} else {
						if strings.HasPrefix(pkgs[0].PkgPath, "profzone") {
							importGroups[2] = append(importGroups[2], getAstString(fileSet, importSpec))
						} else {
							importGroups[1] = append(importGroups[1], getAstString(fileSet, importSpec))
						}
					}
				}
			}

			buf := &bytes.Buffer{}

			buf.WriteString("import (\n")
			for _, importGroup := range importGroups {
				for _, code := range importGroup {
					buf.WriteString(code + "\n")
				}
				buf.WriteString("\n")
			}
			buf.WriteString(")")
			formattedCode = strings.Replace(formattedCode, importsCode, buf.String(), -1)
		}
	}

	return Format(filename, []byte(formattedCode)), nil
}
