package drone

import "gitee.com/eden-framework/eden-framework/internal/project/drone/enums"

type PipelinePlatform struct {
	OS           enums.DroneCiPlatformOs   `yaml:"os" json:"os"`
	Architecture enums.DroneCiPlatformArch `yaml:"arch" json:"arch"`
	Version      int                       `yaml:"version,omitempty" json:"version,omitempty"`
}

func NewPipelinePlatform() *PipelinePlatform {
	return new(PipelinePlatform)
}

func (p *PipelinePlatform) WithOS(os enums.DroneCiPlatformOs) *PipelinePlatform {
	p.OS = os
	return p
}

func (p *PipelinePlatform) WithArchitecture(arch enums.DroneCiPlatformArch) *PipelinePlatform {
	p.Architecture = arch
	return p
}

func (p *PipelinePlatform) WithVersion(v int) *PipelinePlatform {
	p.Version = v
	return p
}
