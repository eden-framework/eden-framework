package generator

import (
	"github.com/eden-framework/eden-framework/pkg/generator/files"
)

type EntryPointPlugins interface {
	NewApplicationGenerationPoint(opt ServiceOption, cwd string) string
}

type FilePlugins interface {
	FileGenerationPoint(opt ServiceOption, cwd string) *files.GoFile
}
