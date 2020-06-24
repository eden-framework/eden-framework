package drone

type PipelineVolumeHost struct {
	Path string `yaml:"path" json:"path"`
}

func (vh *PipelineVolumeHost) WithPath(p string) *PipelineVolumeHost {
	vh.Path = p
	return vh
}

type PipelineVolume struct {
	Name string              `yaml:"name" json:"name"`
	Host *PipelineVolumeHost `yaml:"host,omitempty" json:"host,omitempty"`
	Temp *struct{}           `yaml:"temp,omitempty" json:"temp,omitempty"`
}

func NewPipelineVolume() *PipelineVolume {
	return &PipelineVolume{}
}

func (v *PipelineVolume) WithName(n string) *PipelineVolume {
	v.Name = n
	return v
}

func (v *PipelineVolume) WithHost(h *PipelineVolumeHost) *PipelineVolume {
	v.Host = h
	v.Temp = nil
	return v
}

func (v *PipelineVolume) WithTemp() *PipelineVolume {
	v.Host = nil
	v.Temp = &struct{}{}
	return v
}

type PipelineStepVolume struct {
	Name string `yaml:"name" json:"name"`
	Path string `yaml:"path" json:"path"`
}

func NewPipelineStepVolume() *PipelineStepVolume {
	return &PipelineStepVolume{}
}

func (v *PipelineStepVolume) WithName(n string) *PipelineStepVolume {
	v.Name = n
	return v
}

func (v *PipelineStepVolume) WithPath(p string) *PipelineStepVolume {
	v.Path = p
	return v
}
