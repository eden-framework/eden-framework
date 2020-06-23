package drone

import "github.com/profzone/eden-framework/internal/project/drone/enums"

type CIDroneBase struct {
	Kind enums.DroneCiKind `yaml:"kind"`
	Name string            `yaml:"name"`
}

type CIDroneSecret struct {
	CIDroneBase
	Data string `yaml:"data"`
}

type CIDroneSignature struct {
	CIDroneBase
	Hmac string `yaml:"hmac"`
}

type CIDronePipeline struct {
	CIDroneBase
	Type enums.DroneCiType `yaml:"type"`
}

type CIDronePipelineDocker struct {
	CIDronePipeline
	Trigger PipelineTrigger `yaml:"trigger"`
}
