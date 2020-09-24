package plugin

import (
	"github.com/eden-framework/eden-framework/pkg/generator"
	"github.com/eden-framework/eden-framework/pkg/generator/files"
)

type ServicePlugins interface {
	NewApplicationGenerationPoint(opt generator.ServiceOption) string
	FileGenerationPoint(opt generator.ServiceOption) *files.GoFile
}
