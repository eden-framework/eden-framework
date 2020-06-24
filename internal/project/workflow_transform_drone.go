package project

import (
	"fmt"
	"github.com/profzone/eden-framework/internal/project/drone"
	"github.com/profzone/eden-framework/pkg/executil"
	str "github.com/profzone/eden-framework/pkg/strings"
	"strings"
)

func (w *Workflow) ToDroneConfig(p *Project) *drone.CIDronePipelineDocker {
	config := drone.NewCIDronePipelineDocker()

	for branch, branchFlow := range w.BranchFlows {
		if branchFlow.Skip {
			continue
		}

		for stage, job := range branchFlow.Jobs {
			if job.Skip {
				continue
			}

			envVars := executil.EnvVars{}
			envVars.LoadFromEnviron()

			step := drone.NewPipelineStep().
				WithName(fmt.Sprintf("%s_%s", str.ToLowerCamelCase(branch), stage)).
				WithEnvs(branchFlow.Env).
				WithImage(fmt.Sprintf("${%s}/${%s}", DOCKER_REGISTRY_KEY, strings.ToUpper(envVars.Parse(job.Builder)))).
				WithCommands(job.Run...)

			if branch != "*" {
				step.WithBranchInclude(branch)
			}

			config.Steps = append(config.Steps, *step)
		}
	}

	return config
}
