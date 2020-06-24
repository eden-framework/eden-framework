package drone

type PipelineClone struct {
	Depth   int  `yaml:"depth" json:"depth"`
	Disable bool `yaml:"disable" json:"disable"`
}

func NewPipelineClone() *PipelineClone {
	return new(PipelineClone)
}

func (c *PipelineClone) WithDepth(d int) *PipelineClone {
	c.Depth = d
	return c
}

func (c *PipelineClone) SetDisable() *PipelineClone {
	c.Disable = true
	return c
}

func (c *PipelineClone) SetEnable() *PipelineClone {
	c.Disable = false
	return c
}
