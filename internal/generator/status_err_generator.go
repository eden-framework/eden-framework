package generator

import (
	"github.com/eden-framework/courier/status_error"
	"github.com/eden-framework/eden-framework/internal/generator/files"
	"github.com/eden-framework/eden-framework/internal/generator/scanner"
	"github.com/eden-framework/packagex"
	"github.com/sirupsen/logrus"
	"go/types"
	"os"
	"path"
	"strconv"
	"strings"
)

type StatusErrGenerator struct {
	pkg              *packagex.Package
	packageName      string
	statusErrorCodes status_error.StatusErrorCodeMap
}

func (s *StatusErrGenerator) Finally() {
	panic("implement me")
}

func NewStatusErrGenerator() *StatusErrGenerator {
	return &StatusErrGenerator{
		statusErrorCodes: status_error.StatusErrorCodeMap{},
	}
}

func (s *StatusErrGenerator) Load(path string) {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsExist(err) {
			logrus.Panicf("entry path does not exist: %s", path)
		}
	}
	pkg, err := packagex.LoadFrom("", path)
	if err != nil {
		logrus.Panic(err)
	}

	s.pkg = pkg
}

func (s *StatusErrGenerator) Pick() {
	for ident, obj := range s.pkg.TypesInfo.Defs {
		if constObj, ok := obj.(*types.Const); ok {
			if strings.HasPrefix(constObj.Type().String(), scanner.PkgImportPathStatusErr) {
				key := constObj.Name()
				if key == "_" {
					continue
				}
				for _, f := range s.pkg.Syntax {
					scan := scanner.NewCommentScanner(s.pkg.Fset, f)
					doc := scan.CommentsOf(ident)
					if doc == "" {
						continue
					}

					code, _ := strconv.ParseInt(constObj.Val().String(), 10, 64)
					msg, desc, canBeErrTalk := ParseStatusErrorDesc(doc)

					s.statusErrorCodes.Register(key, code, msg, desc, canBeErrTalk)
				}
			}
		}
	}
	s.packageName = s.pkg.Name
}

func (s *StatusErrGenerator) Output(outputPath string) Outputs {
	outputs := Outputs{}
	errCodeFile := files.NewErrCodeFile(s.packageName, s.statusErrorCodes)
	outputs.Add(GeneratedSuffix(path.Join(outputPath, "status_err_codes.go")), errCodeFile.String())
	return outputs
}

func ParseStatusErrorDesc(str string) (msg string, desc string, canBeErrTalk bool) {
	lines := strings.Split(str, "\n")
	firstLine := strings.Split(lines[0], "@errTalk")

	if len(firstLine) > 1 {
		canBeErrTalk = true
		msg = strings.TrimSpace(firstLine[1])
	} else {
		canBeErrTalk = false
		msg = strings.TrimSpace(firstLine[0])
	}

	if len(lines) > 1 {
		desc = strings.TrimSpace(strings.Join(lines[1:], "\n"))
	}
	return
}
