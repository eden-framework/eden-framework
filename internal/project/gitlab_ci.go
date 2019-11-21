package project

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type CICache struct {
	UnTracked bool     `yaml:"untracked,omitempty"`
	Key       string   `yaml:"key,omitempty"`
	Paths     []string `yaml:"paths,omitempty"`
}

type CIEnv struct {
	Name string `yaml:"name"`
}

type CIArtifact struct {
	UnTracked bool     `yaml:"untracked,omitempty"`
	When      string   `yaml:"when,omitempty"`
	Name      string   `yaml:"name,omitempty"`
	Paths     []string `yaml:"paths,omitempty"`
	ExpireIn  string   `yaml:"expire_in,omitempty"`
}

func NewCIJob(stage string) *CIJob {
	return &CIJob{
		Stage: stage,
	}
}

type CIJob struct {
	Stage string `yaml:"stage,omitempty"`

	Image        string   `yaml:"image,omitempty"`
	Tags         []string `yaml:"tags,omitempty"`
	Services     []string `yaml:"services,omitempty"`
	Dependencies []string `yaml:"dependencies,omitempty"`

	Variables map[string]string `yaml:"variables,omitempty"`

	Environment CIEnv `yaml:"environment,omitempty"`

	BeforeScript []string `yaml:"before_script,omitempty"`
	AfterScript  []string `yaml:"after_script,omitempty"`
	Script       []string `yaml:"script,omitempty"`

	AllowFailure bool `yaml:"allow_failure,omitempty"`

	Artifacts *CIArtifact `yaml:"artifacts,omitempty"`

	Only     []string `yaml:"only,omitempty"`
	Except   []string `yaml:"except,omitempty"`
	When     string   `yaml:"when,omitempty"`
	Coverage string   `yaml:"coverage,omitempty"`
}

func (c CIJob) WithTags(tags ...string) *CIJob {
	c.Tags = append(c.Tags, tags...)
	return &c
}

func (c CIJob) WithDependencies(dependencies ...string) *CIJob {
	c.Dependencies = dependencies
	return &c
}

func (c CIJob) WithScript(script ...string) *CIJob {
	c.Script = append(c.Script, script...)
	return &c
}

func (c CIJob) WithEnv(name string) *CIJob {
	c.Environment = CIEnv{
		Name: name,
	}
	return &c
}

func (c CIJob) AllowFail() *CIJob {
	c.AllowFailure = true
	return &c
}

func (c CIJob) WithArtifacts(artifact *CIArtifact) *CIJob {
	c.Artifacts = artifact
	return &c
}

func (c CIJob) WithOnly(only ...string) *CIJob {
	c.Only = only
	return &c
}

func (c CIJob) WithExcept(except ...string) *CIJob {
	c.Except = except
	return &c
}

func (c CIJob) WithWhen(when string) *CIJob {
	c.When = when
	return &c
}

func (c CIJob) WithImage(image string) *CIJob {
	c.Image = image
	return &c
}

func (c CIJob) WithVariable(key string, value string) CIJob {
	if c.Variables == nil {
		c.Variables = map[string]string{}
	}
	c.Variables[key] = value
	return c
}

type CIJobMap map[string]*CIJob

func NewCIConfig() *CIConfig {
	return &CIConfig{}
}

type CIConfig struct {
	Stages    []string `yaml:"stages"`
	Cache     CICache  `yaml:"cache,omitempty"`
	CommonJob CIJob    `yaml:"common_job,inline,omitempty"`
	CIJobMap  `yaml:"jobs,inline"`
}

func (c *CIConfig) hasStage(stage string) bool {
	for _, s := range c.Stages {
		if s == stage {
			return true
		}
	}
	return false
}

func prefixJobName(name string) string {
	return "job_" + name
}

func (c CIConfig) WithStages(stages ...string) *CIConfig {
	c.Stages = append(c.Stages, stages...)
	return &c
}

func (c CIConfig) WithCache(cache CICache) *CIConfig {
	c.Cache = cache
	return &c
}

func (c CIConfig) WithCommon(job *CIJob) *CIConfig {
	c.CommonJob = CIJob{
		Image:        job.Image,
		Services:     job.Services,
		Variables:    job.Variables,
		BeforeScript: job.BeforeScript,
		AfterScript:  job.AfterScript,
	}
	return &c
}

func (c CIConfig) AddJob(name string, job *CIJob) *CIConfig {
	if !c.hasStage(job.Stage) {
		panic(fmt.Errorf("missing stage %s", job.Stage))
	}

	if c.CIJobMap == nil {
		c.CIJobMap = map[string]*CIJob{}
	}

	c.CIJobMap[prefixJobName(name)] = job

	return &c
}

func (c CIConfig) AddJobWithDependencies(name string, job *CIJob, dependencies ...string) *CIConfig {
	finalDependencies := []string{}
	for _, d := range dependencies {
		prefixedName := prefixJobName(d)
		if _, ok := c.CIJobMap[prefixedName]; !ok {
			panic(fmt.Errorf("job %s should defined first", d))
		}
		finalDependencies = append(finalDependencies, prefixedName)
	}
	return c.AddJob(name, job.WithDependencies(finalDependencies...))
}

func (c *CIConfig) WriteToFile() {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(".gitlab-ci.yml", bytes, os.ModePerm)
}
