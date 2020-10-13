package generator

import (
	"fmt"
	str "github.com/eden-framework/strings"
	"github.com/fatih/color"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func IsGoFile(filename string) bool {
	return filepath.Ext(filename) == ".go"
}

func IsGoTextFile(filename string) bool {
	return strings.HasSuffix(filepath.Base(filename), "_test.go")
}

func WriteFile(filename string, content string) {
	pwd, _ := os.Getwd()
	if filepath.IsAbs(filename) {
		filename, _ = filepath.Rel(pwd, filename)
	}
	n3 := mustWriteFile(filename, content)
	log.Printf(color.GreenString("Generated file to %s (%d KiB)", color.BlueString(path.Join(pwd, filename)), n3/1024))
}

func PathExist(p string) bool {
	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func GetServiceName(projectName string) string {
	projectName = strings.TrimPrefix(projectName, "srv-")
	projectName = strings.TrimPrefix(projectName, "service-")
	return projectName
}

func GeneratedSuffix(filename string) string {
	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	ext := filepath.Ext(filename)

	if IsGoFile(filename) && IsGoTextFile(filename) {
		base = strings.Replace(base, "_test.go", "__generated_test.go", -1)
	} else {
		base = strings.Replace(base, ext, fmt.Sprintf("__generated%s", ext), -1)

	}
	return fmt.Sprintf("%s/%s", dir, base)
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

func toDefaultTableName(name string) string {
	return str.ToLowerSnakeCase("t_" + name)
}
