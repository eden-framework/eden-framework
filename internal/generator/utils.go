package generator

import (
	"github.com/fatih/color"
	"go/build"
	"log"
	"os"
	"path"
	"path/filepath"
)

func GetPackageImportPath(dir string) string {
	pkg, err := build.ImportDir(dir, build.FindOnly)
	if err != nil {
		panic(err)
	}
	return pkg.ImportPath
}

func IsGoFile(filename string) bool {
	return filepath.Ext(filename) == ".go"
}

func WriteFile(filename string, content string) {
	pwd, _ := os.Getwd()
	if filepath.IsAbs(filename) {
		filename, _ = filepath.Rel(pwd, filename)
	}
	n3 := mustWriteFile(filename, content)
	log.Printf(color.GreenString("Generated file to %s (%d KiB)", color.BlueString(path.Join(pwd, filename)), n3/1024))
}

func mustWriteFile(filename string, content string) int {
	dir := filepath.Dir(filename)

	if dir != "" {
		os.MkdirAll(dir, os.ModePerm)
	}

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n3, err := f.WriteString(content)
	if err != nil {
		panic(err)
	}
	f.Sync()

	return n3
}
