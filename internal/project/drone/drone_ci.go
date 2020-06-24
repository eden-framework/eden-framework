package drone

import (
	"github.com/profzone/eden-framework/internal/project/drone/enums"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type CIDroneBase struct {
	Kind enums.DroneCiKind `yaml:"kind" json:"kind"`
	Name string            `yaml:"name,omitempty" json:"name,omitempty"`
}

func (b *CIDroneBase) WithKind(k enums.DroneCiKind) *CIDroneBase {
	b.Kind = k
	return b
}

func (b *CIDroneBase) WithName(n string) *CIDroneBase {
	b.Name = n
	return b
}

type CIDroneSecret struct {
	CIDroneBase `yaml:",inline"`
	Data        string `yaml:"data" json:"data"`
}

func NewCIDroneSecret(name string) *CIDroneSecret {
	s := new(CIDroneSecret)
	s.WithKind(enums.DRONE_CI_KIND__secret).WithName(name)
	return s
}

func (s *CIDroneSecret) WithData(d string) *CIDroneSecret {
	s.Data = d
	return s
}

type CIDroneSignature struct {
	CIDroneBase `yaml:",inline"`
	Hmac        string `yaml:"hmac" json:"hmac"`
}

func NewCIDroneSignature() *CIDroneSignature {
	s := new(CIDroneSignature)
	s.WithKind(enums.DRONE_CI_KIND__signature)
	return s
}

func (s *CIDroneSignature) WithHmac(h string) *CIDroneSignature {
	s.Hmac = h
	return s
}

type CIDronePipeline struct {
	CIDroneBase `yaml:",inline"`
	Type        enums.DroneCiType `yaml:"type" json:"type"`
}

func NewCIDronePipeline() *CIDronePipeline {
	p := new(CIDronePipeline)
	p.WithKind(enums.DRONE_CI_KIND__pipeline)
	return p
}

func (p *CIDronePipeline) WithType(t enums.DroneCiType) *CIDronePipeline {
	p.Type = t
	return p
}

type CIDronePipelineDocker struct {
	CIDronePipeline `yaml:",inline"`
	Trigger         *PipelineTrigger   `yaml:"trigger,omitempty" json:"trigger,omitempty"`
	Platform        *PipelinePlatform  `yaml:"platform,omitempty" json:"platform,omitempty"`
	Workspace       *PipelineWorkspace `yaml:"workspace,omitempty" json:"workspace,omitempty"`
	Clone           *PipelineClone     `yaml:"clone,omitempty" json:"clone,omitempty"`
	Steps           []PipelineStep     `yaml:"steps,omitempty" json:"steps,omitempty"`
	Services        []PipelineService  `yaml:"services,omitempty" json:"services,omitempty"`
	Node            map[string]string  `yaml:"node,omitempty" json:"node,omitempty"`
	Volumes         []PipelineVolume   `yaml:"volumes,omitempty" json:"volumes,omitempty"`
}

func NewCIDronePipelineDocker() *CIDronePipelineDocker {
	p := new(CIDronePipelineDocker)
	p.CIDronePipeline = *NewCIDronePipeline()
	p.WithType(enums.DRONE_CI_TYPE__docker)
	return p
}

func (c *CIDronePipelineDocker) WriteToFile() {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(".drone.yml", bytes, os.ModePerm)
}
