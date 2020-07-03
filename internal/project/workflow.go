package project

import (
	"fmt"
	"github.com/imdario/mergo"
	"sort"
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

func (w Workflow) GetAvailableEnv() (envs map[string]string) {
	for _, flow := range w.BranchFlows {
		if flow.Env != nil && len(flow.Env) > 0 {
			_ = mergo.Merge(&envs, flow.Env, mergo.WithOverride)
		}
	}
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
		for _, job := range branchFlow.Jobs {
			switch job.Stage {
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
	Skip    bool              `yaml:"skip,omitempty"`
	Extends string            `yaml:"extends,omitempty"`
	Env     map[string]string `yaml:"env,omitempty"`
	Jobs    Jobs              `yaml:"jobs,inline,omitempty"`
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
	if nextBranchFlow.Env != nil && len(nextBranchFlow.Env) > 0 {
		branchFlow.Env = nextBranchFlow.Env
		_ = mergo.Merge(&branchFlow.Env, nextBranchFlow.Env, mergo.WithOverride)
	}
	if nextBranchFlow.Extends != "" {
		branchFlow.Extends = nextBranchFlow.Extends
	}
	if len(nextBranchFlow.Jobs) > 0 {
		branchFlow.Jobs = branchFlow.Jobs.Merge(nextBranchFlow.Jobs)
	}
	return branchFlow
}

type Jobs []Job

func (jobs Jobs) Merge(nextJobs Jobs) Jobs {
	finalJobs := Jobs{}
	for _, job := range jobs {
		index, nextJob := nextJobs.Find(job.Stage)
		if nextJob != nil {
			finalJobs = append(finalJobs, nextJob.Merge(&job))
			nextJobs, _ = nextJobs.Remove(index)
		} else {
			finalJobs = append(finalJobs, job)
		}
	}
	for name, job := range nextJobs {
		finalJobs[name] = job
	}
	return finalJobs
}

func (jobs Jobs) Find(stage string) (int, *Job) {
	for i, j := range jobs {
		if j.Stage == stage {
			return i, &j
		}
	}
	return -1, nil
}

func (jobs Jobs) Remove(index int) (Jobs, error) {
	if index == 0 {
		return jobs[1:], nil
	} else if index == len(jobs)-1 {
		return jobs[:len(jobs)-2], nil
	} else if index > len(jobs)-1 {
		return jobs, fmt.Errorf("index out of range: %d, length: %d", index, len(jobs)-1)
	} else {
		finalJobs := Jobs{}
		finalJobs = append(finalJobs, jobs[:index-1]...)
		finalJobs = append(finalJobs, jobs[index+1:]...)
		return finalJobs, nil
	}
}

type Job struct {
	Stage     string      `yaml:"stage,omitempty"`
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
