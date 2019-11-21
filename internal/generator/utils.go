package generator

import (
	"fmt"
	"github.com/fatih/color"
	"go/build"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

// https://github.com/golang/lint/blob/master/lint.go#L720
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

func SplitToWords(s string) (entries []string) {
	if !utf8.ValidString(s) {
		return []string{s}
	}
	entries = []string{}
	var runes [][]rune
	lastClass := 0
	class := 0

	// split into fields based on class of unicode character
	for _, r := range s {
		switch true {
		case unicode.IsSpace(r):
			class = 1
		case unicode.IsLower(r):
			class = 2
		case unicode.IsUpper(r):
			class = 3
		case unicode.IsDigit(r):
			class = 4
		default:
			class = 5
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}

	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}

	// construct []string from results
	//
	for _, s := range runes {
		if len(s) > 0 && (unicode.IsDigit(s[0]) || unicode.IsLetter(s[0])) {
			entries = append(entries, string(s))
		}
	}

	return
}

type RewordsReducer func(result string, word string, index int) string

func Rewords(s string, reducer RewordsReducer) string {
	words := SplitToWords(s)

	var result = ""

	for idx, word := range words {
		result = reducer(result, word, idx)
	}

	return result
}

func ToUpperFirst(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func ToCamelCase(s string) string {
	upperString := strings.ToUpper(s)
	if commonInitialisms[upperString] {
		return upperString
	}
	return ToUpperFirst(strings.ToLower(upperString))
}

func ToUpperCamelCase(s string) string {
	return Rewords(s, func(result string, word string, idx int) string {
		return result + ToCamelCase(word)
	})
}

func ToLowerCamelCase(s string) string {
	return Rewords(s, func(result string, word string, idx int) string {
		if idx == 0 {
			return result + strings.ToLower(word)
		}
		return result + ToCamelCase(word)
	})
}

func ToUpperSnakeCase(s string) string {
	return Rewords(s, func(result string, word string, idx int) string {
		newWord := strings.ToUpper(word)
		if idx == 0 || (len(newWord) == 1 && unicode.IsDigit(rune(newWord[0]))) {
			return result + newWord
		}
		return result + "_" + newWord
	})
}

func ToLowerSnakeCase(s string) string {
	return strings.ToLower(ToUpperSnakeCase(s))
}

func ToLowerLinkCase(s string) string {
	return strings.Replace(ToLowerSnakeCase(s), "_", "-", -1)
}

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
