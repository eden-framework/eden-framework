package generator

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/eden-framework/eden-framework/internal/project"
	"github.com/eden-framework/eden-framework/pkg/generator/files"
	"path"
)

type expressBool bool

func (e *expressBool) WriteAnswer(field string, value interface{}) error {
	result := value.(core.OptionAnswer)
	switch result.Value {
	case "是":
		*e = true
	case "否":
		*e = false
	}
	return nil
}

type ProjectOption struct {
	Name            string
	Group           string
	Owner           string
	Desc            string
	Version         string
	ProgramLanguage string `survey:"project_language"`
	Workflow        string
	ApolloSupport   expressBool `survey:"apollo_support"`
}

func NewProjectGenerator(opt ProjectOption) *ProjectGenerator {
	p := project.Project{}
	p = p.WithName(opt.Name).
		WithGroup(opt.Group).
		WithOwner(opt.Owner).
		WithDesc(opt.Desc).
		WithLanguage(opt.ProgramLanguage)

	if opt.Version != "" {
		p = p.WithVersion(opt.Version)
	}

	if opt.Workflow != "" && opt.Workflow != "custom" {
		p = p.WithWorkflow(opt.Workflow)
	}

	var withApolloFlag string
	if opt.ApolloSupport {
		withApolloFlag = fmt.Sprintf(" -ldflags \"-X github.com/eden-framework/eden-framework/pkg/conf/apollo.Branch=%s.json\"", files.EnvVarInBash(project.EnvKeyCIBranch))
	}
	p.Scripts = map[string]project.Script{
		"build": []string{
			fmt.Sprintf("go build -v -o ./build/$PROJECT_NAME%s ./cmd", withApolloFlag),
			"eden generate openapi",
		},
		"test": []string{
			"go test ./cmd",
		},
	}

	return &ProjectGenerator{
		project: p,
	}
}

type ProjectGenerator struct {
	project project.Project
}

func (p *ProjectGenerator) Load(path string) {
	p.project.SetEnviron()
}

func (p *ProjectGenerator) Pick() {
}

func (p *ProjectGenerator) Output(outputPath string) Outputs {
	outputs := Outputs{}

	outputs.Add(path.Join(outputPath, "project.yml"), p.project.String())
	outputs.Add(path.Join(outputPath, ".drone.yml"), p.project.Workflow.TryExtendsOrSetDefaults().ToDroneConfig(&p.project).String())

	return outputs
}
