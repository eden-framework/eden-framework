package drone

type PipelineWorkspace struct {
	Path string `yaml:"path" json:"path"`
}

func NewPipelineWorkspace() *PipelineWorkspace {
	return new(PipelineWorkspace)
}

func (w *PipelineWorkspace) WithPath(path string) *PipelineWorkspace {
	w.Path = path
	return w
}
