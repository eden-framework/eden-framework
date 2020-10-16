package operator

import "github.com/eden-framework/eden-framework/internal/generator/importer"

var annotationMap = map[string]Annotation{}

func registerAnnotation(extensionKey string, annotation Annotation) {
	annotationMap[extensionKey] = annotation
}

func getAnnotation(extensionKey string) Annotation {
	return annotationMap[extensionKey]
}

const (
	CmdGenerateGetRevertID = "generate_get_revert_id"
	CmdGenerateInterface   = "generate_interface"
	CmdGenerateImplement   = "generate_implement"
)

type Annotation interface {
	ID() string
	Extension() string
	SetArgs(args ...string)
	Importer() *importer.PackageImporter
	Run(cmd string, op Op) string
}
