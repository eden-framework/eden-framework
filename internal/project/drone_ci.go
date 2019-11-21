package project

type CIWorkspace struct {
	Base string `yaml:"base,omitempty"`
	Path string `yaml:"path,omitempty"`
}

type CIPipelineMap map[string]CIPipeline

type CIPipeline struct {
	Image    string              `yaml:"image"`
	IsPull   bool                `yaml:"pull,omitempty"`
	Commands []string            `yaml:"commands"`
	When     CIPipelineCondition `yaml:"when,omitempty"`
}

func (p *CIPipeline) WithImage(image string) *CIPipeline {
	p.Image = image
	return p
}

func (p *CIPipeline) WithPull(pull bool) *CIPipeline {
	p.IsPull = pull
	return p
}

func (p *CIPipeline) WithCommands(commands ...string) *CIPipeline {
	p.Commands = append(p.Commands, commands...)
	return p
}

func (p *CIPipeline) WithCondition(condition CIPipelineCondition) *CIPipeline {
	p.When = condition
	return p
}

type CIPipelineCondition struct {
	Branch     []string `yaml:"branch,omitempty"`
	Event      []string `yaml:"event,omitempty"`
	Reference  []string `yaml:"reference,omitempty"`
	Repository []string `yaml:"repository,omitempty"`
}

func NewCIPipelineCondition() CIPipelineCondition {
	return CIPipelineCondition{}
}

func (c *CIPipelineCondition) WithBranches(branches ...string) *CIPipelineCondition {
	c.Branch = append(c.Branch, branches...)
	return c
}

func (c *CIPipelineCondition) WithEvents(events ...string) *CIPipelineCondition {
	c.Event = append(c.Event, events...)
	return c
}

func (c *CIPipelineCondition) WithReference(references ...string) *CIPipelineCondition {
	c.Reference = append(c.Reference, references...)
	return c
}

func (c *CIPipelineCondition) WithRepository(repositories ...string) *CIPipelineCondition {
	c.Repository = append(c.Repository, repositories...)
	return c
}

type CIDroneConfig struct {
	Workspace CIWorkspace   `yaml:"workspace"`
	Pipeline  CIPipelineMap `yaml:"pipeline"`
}
