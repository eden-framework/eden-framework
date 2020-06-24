package drone

import (
	"github.com/imdario/mergo"
	"github.com/profzone/eden-framework/internal/project/drone/enums"
)

type PipelineStep struct {
	Name        string                   `yaml:"name" json:"name"`
	Image       string                   `yaml:"image" json:"image"`
	Pull        enums.DroneCiStepPull    `yaml:"pull,omitempty" json:"pull,omitempty"`
	Commands    []string                 `yaml:"commands,omitempty" json:"commands,omitempty"`
	Environment map[string]string        `yaml:"environment,omitempty" json:"environment,omitempty"`
	Settings    map[string]interface{}   `yaml:"settings,omitempty" json:"settings,omitempty"`
	When        *PipelineTrigger         `yaml:"when,omitempty" json:"when,omitempty"`
	Failure     enums.DroneCiStepFailure `yaml:"failure,omitempty" json:"failure,omitempty"`
	Detach      bool                     `yaml:"detach,omitempty" json:"detach,omitempty"`
	Privileged  bool                     `yaml:"privileged,omitempty" json:"privileged,omitempty"`
	DependsOn   []string                 `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	Volumes     []PipelineStepVolume     `yaml:"volumes,omitempty" json:"volumes,omitempty"`
}

func NewPipelineStep() *PipelineStep {
	return new(PipelineStep)
}

func (s *PipelineStep) WithName(n string) *PipelineStep {
	s.Name = n
	return s
}

func (s *PipelineStep) WithImage(img string) *PipelineStep {
	s.Image = img
	return s
}

func (s *PipelineStep) WithCommands(cmd ...string) *PipelineStep {
	s.Commands = append(s.Commands, cmd...)
	return s
}

func (s *PipelineStep) WithEnvs(envs map[string]string) *PipelineStep {
	_ = mergo.Merge(&s.Environment, envs)
	return s
}

func (s *PipelineStep) WithEnv(key string, value string) *PipelineStep {
	s.Environment[key] = value
	return s
}

func (s *PipelineStep) WithSetting(key string, value interface{}) *PipelineStep {
	s.Settings[key] = value
	return s
}

func (s *PipelineStep) WithFailureIgnore() *PipelineStep {
	s.Failure = enums.DRONE_CI_STEP_FAILURE__ignore
	return s
}

func (s *PipelineStep) WithFailureNotIgnore() *PipelineStep {
	s.Failure = enums.DRONE_CI_STEP_FAILURE_UNKNOWN
	return s
}

func (s *PipelineStep) WithDetach() *PipelineStep {
	s.Detach = true
	return s
}

func (s *PipelineStep) WithNoDetach() *PipelineStep {
	s.Detach = false
	return s
}

func (s *PipelineStep) WithPrivileged() *PipelineStep {
	s.Privileged = true
	return s
}

func (s *PipelineStep) WithNoPrivileged() *PipelineStep {
	s.Privileged = false
	return s
}

func (s *PipelineStep) WithBranchInclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = NewPipelineTrigger()
	}
	s.When.WithBranchInclude(name)
	return s
}

func (s *PipelineStep) WithBranchExclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = NewPipelineTrigger()
	}
	s.When.WithBranchExclude(name)
	return s
}

func (s *PipelineStep) WithEventInclude(evt enums.DroneCiTriggerEvent) *PipelineStep {
	if s.When == nil {
		s.When = NewPipelineTrigger()
	}
	s.When.WithEventInclude(evt)
	return s
}

func (s *PipelineStep) WithEventExclude(evt enums.DroneCiTriggerEvent) *PipelineStep {
	if s.When == nil {
		s.When = NewPipelineTrigger()
	}
	s.When.WithEventExclude(evt)
	return s
}

func (s *PipelineStep) WithRefInclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = NewPipelineTrigger()
	}
	s.When.WithRefInclude(name)
	return s
}

func (s *PipelineStep) WithRefExclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = NewPipelineTrigger()
	}
	s.When.WithRefExclude(name)
	return s
}

func (s *PipelineStep) WithRepoInclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = NewPipelineTrigger()
	}
	s.When.WithRepoInclude(name)
	return s
}

func (s *PipelineStep) WithRepoExclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = NewPipelineTrigger()
	}
	s.When.WithRepoExclude(name)
	return s
}

func (s *PipelineStep) WithStatus(status enums.DroneCiTriggerStatus) *PipelineStep {
	if s.When == nil {
		s.When = NewPipelineTrigger()
	}
	s.When.WithStatus(status)
	return s
}

func (s *PipelineStep) WithTargetInclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = NewPipelineTrigger()
	}
	s.When.WithTargetInclude(name)
	return s
}

func (s *PipelineStep) WithTargetExclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = NewPipelineTrigger()
	}
	s.When.WithTargetExclude(name)
	return s
}

func (s *PipelineStep) WithDependsOn(deps ...string) *PipelineStep {
	s.DependsOn = append(s.DependsOn, deps...)
	return s
}

func (s *PipelineStep) WithVolume(volume ...PipelineStepVolume) *PipelineStep {
	s.Volumes = append(s.Volumes, volume...)
	return s
}
