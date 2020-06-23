package drone

type PipelineWorkspace struct {
	Path string `yaml:"path"`
}

func (w *PipelineWorkspace) WithPath(path string) *PipelineWorkspace {
	w.Path = path
	return w
}
