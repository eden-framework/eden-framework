package operator

import (
	"gitee.com/eden-framework/eden-framework/internal/generator/importer"
	"gitee.com/eden-framework/eden-framework/internal/generator/scanner"
)

func init() {
	registerAnnotation(scanner.XAnnotationRevert, NewAnnotationRevert("Revert", scanner.XAnnotationRevert))
}

type AnnotationRevert struct {
	id        string
	extension string
	ipt       *importer.PackageImporter

	clientName string
	targetFunc string
}

func NewAnnotationRevert(id, extension string) Annotation {
	return &AnnotationRevert{
		id:        id,
		extension: extension,
		ipt:       importer.NewPackageImporter(""),
	}
}

func (a *AnnotationRevert) ID() string {
	return a.id
}

func (a *AnnotationRevert) Extension() string {
	return a.extension
}

func (a *AnnotationRevert) SetArgs(args ...string) {
	if len(args) > 1 && len(args[0]) > 0 {
		a.targetFunc = args[0]
		a.clientName = args[1]
	} else {
		panic("[AnnotationRevert] Revert annotation must have 1 parameter(not empty) at least")
	}
}

func (a *AnnotationRevert) Importer() *importer.PackageImporter {
	return a.ipt
}

func (a *AnnotationRevert) Run(cmd string, op Op) string {
	switch cmd {
	case CmdGenerateGetRevertID:
		return a.generateGetRevertID(op)
	case CmdGenerateInterface:
		return a.generateInterface(op)
	case CmdGenerateImplement:
		return a.generateImplement(op)
	}
	return ""
}

func (a *AnnotationRevert) generateGetRevertID(op Op) string {
	return `
func (r ` + ResponseOf(op.ID()) + `) GetRevertID() uint64 {
	return r.Body.ID
}
`
}

func (a *AnnotationRevert) generateInterface(op Op) string {
	return op.ID() + `(id uint64, metas... ` + a.ipt.Use("gitee.com/eden-framework/courier.Metadata") + `) (err error)`
}

func (a *AnnotationRevert) generateImplement(op Op) string {
	interfaceMethod := a.generateInterface(op)
	return `
func (c ` + a.clientName + `) ` + interfaceMethod + ` {
	req := ` + RequestOf(op.ID()) + `{
		` + op.RevertIDField() + `: id,
	}
	err = c.Request(c.Name + ".` + op.ID() + `", "` + op.Method() + `", "` + op.Path() + `", req, metas...).Do().Err

	return
}
`
}
