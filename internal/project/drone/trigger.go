package drone

import "github.com/profzone/eden-framework/internal/project/drone/enums"

type PipelineTriggerIncludeAndExcludeString struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}
type PipelineTriggerIncludeAndExcludeEvent struct {
	Include []enums.DroneCiTriggerEvent `yaml:"include"`
	Exclude []enums.DroneCiTriggerEvent `yaml:"exclude"`
}

type PipelineTrigger struct {
	Branch     *PipelineTriggerIncludeAndExcludeString `yaml:"branch,omitempty"`
	Event      *PipelineTriggerIncludeAndExcludeEvent  `yaml:"event,omitempty"`
	Reference  *PipelineTriggerIncludeAndExcludeString `yaml:"ref,omitempty"`
	Repository *PipelineTriggerIncludeAndExcludeString `yaml:"repo,omitempty"`
	Status     enums.DroneCiTriggerStatus              `yaml:"status,omitempty"`
	Target     *PipelineTriggerIncludeAndExcludeString `yaml:"target,omitempty"`
}

func (t *PipelineTrigger) WithBranchInclude(name string) *PipelineTrigger {
	if t.Branch == nil {
		t.Branch = &PipelineTriggerIncludeAndExcludeString{}
	}
	t.Branch.Include = append(t.Branch.Include, name)
	return t
}

func (t *PipelineTrigger) WithBranchExclude(name string) *PipelineTrigger {
	if t.Branch == nil {
		t.Branch = &PipelineTriggerIncludeAndExcludeString{}
	}
	t.Branch.Exclude = append(t.Branch.Exclude, name)
	return t
}

func (t *PipelineTrigger) WithEventInclude(evt enums.DroneCiTriggerEvent) *PipelineTrigger {
	if t.Event == nil {
		t.Event = &PipelineTriggerIncludeAndExcludeEvent{}
	}
	t.Event.Include = append(t.Event.Include, evt)
	return t
}

func (t *PipelineTrigger) WithEventExclude(evt enums.DroneCiTriggerEvent) *PipelineTrigger {
	if t.Event == nil {
		t.Event = &PipelineTriggerIncludeAndExcludeEvent{}
	}
	t.Event.Exclude = append(t.Event.Exclude, evt)
	return t
}

func (t *PipelineTrigger) WithRefInclude(name string) *PipelineTrigger {
	if t.Reference == nil {
		t.Reference = &PipelineTriggerIncludeAndExcludeString{}
	}
	t.Reference.Include = append(t.Reference.Include, name)
	return t
}

func (t *PipelineTrigger) WithRefExclude(name string) *PipelineTrigger {
	if t.Reference == nil {
		t.Reference = &PipelineTriggerIncludeAndExcludeString{}
	}
	t.Reference.Exclude = append(t.Reference.Exclude, name)
	return t
}

func (t *PipelineTrigger) WithRepoInclude(name string) *PipelineTrigger {
	if t.Repository == nil {
		t.Repository = &PipelineTriggerIncludeAndExcludeString{}
	}
	t.Repository.Include = append(t.Repository.Include, name)
	return t
}

func (t *PipelineTrigger) WithRepoExclude(name string) *PipelineTrigger {
	if t.Repository == nil {
		t.Repository = &PipelineTriggerIncludeAndExcludeString{}
	}
	t.Repository.Exclude = append(t.Repository.Exclude, name)
	return t
}

func (t *PipelineTrigger) WithStatus(status enums.DroneCiTriggerStatus) *PipelineTrigger {
	t.Status = status
	return t
}

func (t *PipelineTrigger) WithTargetInclude(name string) *PipelineTrigger {
	if t.Target == nil {
		t.Target = &PipelineTriggerIncludeAndExcludeString{}
	}
	t.Target.Include = append(t.Target.Include, name)
	return t
}

func (t *PipelineTrigger) WithTargetExclude(name string) *PipelineTrigger {
	if t.Target == nil {
		t.Target = &PipelineTriggerIncludeAndExcludeString{}
	}
	t.Target.Exclude = append(t.Target.Exclude, name)
	return t
}
