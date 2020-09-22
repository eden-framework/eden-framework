package project

import (
	"fmt"
	"sort"
)

var RegisteredBuilders = Builders{}

func RegisterBuilder(name string, builder *Builder) {
	builder.name = name
	RegisteredBuilders[name] = builder
	SetEnv(name, builder.Image)
}

type FullImage string

func (i FullImage) String() string {
	return fmt.Sprintf("${%s}/%s", EnvKeyDockerRegistryKey, string(i))
}

type Builder struct {
	name            string
	ProgramLanguage string
	Image           string
	WorkingDir      string
}

type Builders map[string]*Builder

func (bs Builders) GetBuilderBy(programLanguage string) *Builder {
	for _, b := range bs {
		if b.ProgramLanguage == programLanguage {
			return b
		}
	}
	return nil
}

func (bs Builders) SupportProgramLanguages() (list []string) {
	for _, b := range bs {
		if b.ProgramLanguage != "" {
			list = append(list, b.ProgramLanguage)
		}
	}
	sort.Strings(list)
	return
}
