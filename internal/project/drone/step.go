package drone

import (
	"github.com/profzone/eden-framework/internal/project/drone/enums"
)

type PipelineStep struct {
	Name        string                   `yaml:"name"`
	Image       string                   `yaml:"image"`
	Commands    []string                 `yaml:"commands,omitempty"`
	Environment map[string]interface{}   `yaml:"environment,omitempty"`
	Settings    map[string]interface{}   `yaml:"settings,omitempty"`
	When        *PipelineTrigger         `yaml:"when,omitempty"`
	Failure     enums.DroneCiStepFailure `yaml:"failure,omitempty"`
	Detach      bool                     `yaml:"detach,omitempty"`
	Privileged  bool                     `yaml:"privileged,omitempty"`
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

func (s *PipelineStep) WithEnv(key string, value interface{}) *PipelineStep {
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
		s.When = &PipelineTrigger{}
	}
	s.When.WithBranchInclude(name)
	return s
}

func (s *PipelineStep) WithBranchExclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = &PipelineTrigger{}
	}
	s.When.WithBranchExclude(name)
	return s
}

func (s *PipelineStep) WithEventInclude(evt enums.DroneCiTriggerEvent) *PipelineStep {
	if s.When == nil {
		s.When = &PipelineTrigger{}
	}
	s.When.WithEventInclude(evt)
	return s
}

func (s *PipelineStep) WithEventExclude(evt enums.DroneCiTriggerEvent) *PipelineStep {
	if s.When == nil {
		s.When = &PipelineTrigger{}
	}
	s.When.WithEventExclude(evt)
	return s
}

func (s *PipelineStep) WithRefInclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = &PipelineTrigger{}
	}
	s.When.WithRefInclude(name)
	return s
}

func (s *PipelineStep) WithRefExclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = &PipelineTrigger{}
	}
	s.When.WithRefExclude(name)
	return s
}

func (s *PipelineStep) WithRepoInclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = &PipelineTrigger{}
	}
	s.When.WithRepoInclude(name)
	return s
}

func (s *PipelineStep) WithRepoExclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = &PipelineTrigger{}
	}
	s.When.WithRepoExclude(name)
	return s
}

func (s *PipelineStep) WithStatus(status enums.DroneCiTriggerStatus) *PipelineStep {
	s.When.WithStatus(status)
	return s
}

func (s *PipelineStep) WithTargetInclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = &PipelineTrigger{}
	}
	s.When.WithTargetInclude(name)
	return s
}

func (s *PipelineStep) WithTargetExclude(name string) *PipelineStep {
	if s.When == nil {
		s.When = &PipelineTrigger{}
	}
	s.When.WithTargetExclude(name)
	return s
}
