package generator

import (
	"github.com/eden-framework/eden-framework/pkg/generator/files"
)

type EntryPointPlugins interface {
	GenerateEntryPoint(opt ServiceOption, cwd string) string
}

type FilePlugins interface {
	GenerateFilePoint(opt ServiceOption, cwd string) []*files.GoFile
}
