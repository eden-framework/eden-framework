package generator

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"os"
	"path"
	"profzone/eden-framework/internal/generator/api"
	"profzone/eden-framework/internal/generator/scanner"
	"profzone/eden-framework/internal/project"
)

type ApiGenerator struct {
	Api             api.Api
	OperatorScanner *scanner.OperatorScanner
	ModelScanner    *scanner.ModelScanner
	pkgs            []*packages.Package
}

func NewApiGenerator(op *scanner.OperatorScanner, model *scanner.ModelScanner) *ApiGenerator {
	return &ApiGenerator{
		Api:             api.NewApi(),
		OperatorScanner: op,
		ModelScanner:    model,
	}
}

func (a *ApiGenerator) Load(cwd string) {
	entryPath := path.Join(cwd, "cmd")
	_, err := os.Stat(entryPath)
	if err != nil {
		if !os.IsExist(err) {
			logrus.Panicf("entry path does not exist: %s", entryPath)
		}
	}
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedDeps | packages.NeedFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedTypesSizes,
		Dir:  entryPath,
	}

	pkgs, err := packages.Load(cfg)
	if err != nil {
		logrus.Panic(err)
	}

	errs := packages.PrintErrors(pkgs)
	if errs > 0 {
		logrus.Panicf("packages.PrintErrors(a.pkgs) = %d", errs)
	}

	a.pkgs = pkgs

	proj := project.Project{}
	err = proj.UnmarshalFromFile(cwd, "")
	if err != nil {
		logrus.Panic(err)
	}

	a.Api.ServiceName = proj.Name
}

func (a *ApiGenerator) Pick() {
	packages.Visit(a.pkgs, nil, func(i *packages.Package) {
		for _, f := range i.Syntax {
			ast.Walk(a.OperatorScanner, f)
		}
	})
	packages.Visit(a.pkgs, nil, func(i *packages.Package) {
		for _, f := range i.Syntax {
			ast.Walk(a.ModelScanner, f)
		}
	})
}

func (a *ApiGenerator) Output(outputPath string) Outputs {
	data, err := json.MarshalIndent(a.Api, "", "    ")
	if err != nil {
		logrus.Panic(err)
	}
	return Outputs{
		"api.json": string(data),
	}
}
