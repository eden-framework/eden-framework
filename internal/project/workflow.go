package project

import (
	"sort"

	"gopkg.in/yaml.v2"
)

var PresetWorkflows = Workflows{}

func RegisterWorkFlow(name string, flow *Workflow) {
	PresetWorkflows[name] = flow
}

const (
	STAGE_TEST   = "test"
	STAGE_BUILD  = "build"
	STAGE_SHIP   = "ship"
	STAGE_DEPLOY = "deploy"
)

type Workflows map[string]*Workflow

func (bs Workflows) List() (list []string) {
	for name := range bs {
		list = append(list, name)
	}
	sort.Strings(list)
	return
}

type Workflow struct {
	Extends     string      `yaml:"extends,omitempty"`
	BranchFlows BranchFlows `yaml:"branch_flows,inline"`
}

func (w Workflow) GetAvailableEnv() (envs []string) {
	for _, flow := range w.BranchFlows {
		if flow.Env != "" {
			envs = append(envs, flow.Env)
		}
	}
	sort.Strings(envs)
	return
}

func (w Workflow) TryExtendsOrSetDefaults() *Workflow {
	if w.Extends != "" {
		if presetFlow, ok := PresetWorkflows[w.Extends]; ok {
			w.BranchFlows = presetFlow.BranchFlows.Merge(w.BranchFlows)
			w.Extends = ""
		}
	}

	for _, branchFlow := range w.BranchFlows {
		for stage, job := range branchFlow.Jobs {
			switch stage {
			case STAGE_TEST, STAGE_BUILD:
				if job.Builder == "" {
					job.Builder = "BUILDER_${PROJECT_PROGRAM_LANGUAGE}"
				}
			case STAGE_SHIP:
				if job.Builder == "" {
					job.Builder = "BUILDER_DOCKER"
				}
			case STAGE_DEPLOY:
				if job.Builder == "" {
					job.Builder = "BUILDER_RANCHER"
				}
			}
		}
	}

	return &w
}

type BranchFlows map[string]BranchFlow

func (branchFlows BranchFlows) Merge(nextBranches BranchFlows) BranchFlows {
	finalBranchFlows := BranchFlows{}
	for name, branch := range branchFlows {
		if nextBranch, ok := nextBranches[name]; ok {
			finalBranchFlows[name] = branch.Merge(&nextBranch)
			delete(nextBranches, name)
		} else {
			finalBranchFlows[name] = branch
		}
	}
	for name, branch := range nextBranches {
		finalBranchFlows[name] = branch
	}

	for name, branch := range finalBranchFlows {
		if branch.Extends != "" {
			if branchFlowNow, ok := finalBranchFlows[branch.Extends]; ok {
				branch.Extends = ""
				finalBranchFlows[name] = branchFlowNow.Merge(&branch)
			}
		}
	}

	return finalBranchFlows
}

type BranchFlow struct {
	Skip    bool   `yaml:"skip,omitempty"`
	Extends string `yaml:"extends,omitempty"`
	Env     string `yaml:"env,omitempty"`
	Jobs    Jobs   `yaml:"jobs,inline,omitempty"`
}

func (branchFlow BranchFlow) MarshalYAML() (interface{}, error) {
	if branchFlow.Skip {
		return "skip", nil
	}

	// sortable config
	slice := yaml.MapSlice{}

	if branchFlow.Env != "" {
		slice = append(slice, yaml.MapItem{
			Key:   "env",
			Value: branchFlow.Env,
		})
	}

	slice = append(slice, branchFlow.Jobs.ToYAMLMapSlice()...)

	return slice, nil
}

func (branchFlow *BranchFlow) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	err := unmarshal(&v)
	if err == nil && v == "skip" {
		branchFlow.Skip = true
		return nil
	}

	type BranchFlowAlias BranchFlow
	b := BranchFlowAlias(*branchFlow)
	if err := unmarshal(&b); err != nil {
		return err
	}
	*branchFlow = BranchFlow(b)
	return nil
}

func (branchFlow BranchFlow) Merge(nextBranchFlow *BranchFlow) BranchFlow {
	if nextBranchFlow.Env != "" {
		branchFlow.Env = nextBranchFlow.Env
	}
	if nextBranchFlow.Extends != "" {
		branchFlow.Extends = nextBranchFlow.Extends
	}
	if len(nextBranchFlow.Jobs) > 0 {
		branchFlow.Jobs = branchFlow.Jobs.Merge(nextBranchFlow.Jobs)
	}
	return branchFlow
}

type Jobs map[string]Job

func (jobs Jobs) ToYAMLMapSlice() yaml.MapSlice {
	slice := yaml.MapSlice{}
	keys := []string{
		STAGE_TEST,
		STAGE_BUILD,
		STAGE_SHIP,
		STAGE_DEPLOY,
	}
	for _, key := range keys {
		if j, ok := jobs[key]; ok {
			slice = append(slice, yaml.MapItem{
				Key:   key,
				Value: j,
			})
		}
	}
	return slice
}

func (jobs Jobs) MarshalYAML() (interface{}, error) {
	return jobs.ToYAMLMapSlice(), nil
}

func (jobs Jobs) Merge(nextJobs Jobs) Jobs {
	finalJobs := Jobs{}
	for name, job := range jobs {
		if nextJob, ok := nextJobs[name]; ok {
			finalJobs[name] = job.Merge(&nextJob)
			delete(nextJobs, name)
		} else {
			finalJobs[name] = job
		}
	}
	for name, job := range nextJobs {
		finalJobs[name] = job
	}
	return finalJobs
}

type Job struct {
	Skip      bool        `yaml:"skip,omitempty"`
	Builder   string      `yaml:"builder,omitempty"`
	Run       Script      `yaml:"run,omitempty"`
	Artifacts *CIArtifact `yaml:"artifacts,omitempty"`
}

func (job Job) MarshalYAML() (interface{}, error) {
	if job.Skip {
		return "skip", nil
	}
	return job, nil
}

func (job *Job) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	err := unmarshal(&v)
	if err == nil && v == "skip" {
		job.Skip = true
		return nil
	}
	type JobAlias Job
	b := JobAlias(*job)
	if err := unmarshal(&b); err != nil {
		return err
	}
	*job = Job(b)
	return nil
}

func (job Job) Merge(nextJob *Job) Job {
	if nextJob.Skip {
		return Job{
			Skip: nextJob.Skip,
		}
	}
	if nextJob.Builder != "" {
		job.Builder = nextJob.Builder
	}
	if !nextJob.Run.IsZero() {
		job.Run = nextJob.Run
	}
	if nextJob.Artifacts != nil {
		job.Artifacts = nextJob.Artifacts
	}
	return job
}
