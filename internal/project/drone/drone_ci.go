package drone

import "github.com/profzone/eden-framework/internal/project/drone/enums"

type CIDroneBase struct {
	Kind enums.DroneCiKind `yaml:"kind" json:"kind"`
	Name string            `yaml:"name" json:"name"`
}

type CIDroneSecret struct {
	CIDroneBase
	Data string `yaml:"data" json:"data"`
}

type CIDroneSignature struct {
	CIDroneBase
	Hmac string `yaml:"hmac" json:"hmac"`
}

type CIDronePipeline struct {
	CIDroneBase
	Type enums.DroneCiType `yaml:"type" json:"type"`
}

type CIDronePipelineDocker struct {
	CIDronePipeline
	Trigger   *PipelineTrigger   `yaml:"trigger,omitempty" json:"trigger,omitempty"`
	Platform  *PipelinePlatform  `yaml:"platform,omitempty" json:"platform,omitempty"`
	Workspace *PipelineWorkspace `yaml:"workspace,omitempty" json:"workspace,omitempty"`
	Clone     *PipelineClone     `yaml:"clone,omitempty" json:"clone,omitempty"`
	Steps     []PipelineStep     `yaml:"steps,omitempty" json:"steps,omitempty"`
	Services  []PipelineService  `yaml:"services,omitempty" json:"services,omitempty"`
	Node      map[string]string  `yaml:"node,omitempty" json:"node,omitempty"`
	Volumes   []PipelineVolume   `yaml:"volumes,omitempty" json:"volumes,omitempty"`
}
