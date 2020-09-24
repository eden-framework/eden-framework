package project

import (
	"fmt"
	"github.com/eden-framework/eden-framework/internal/project/drone"
	"github.com/eden-framework/eden-framework/pkg/executil"
	str "github.com/eden-framework/eden-framework/pkg/strings"
	"strings"
)

func (w *Workflow) ToDroneConfig(p *Project) *drone.CIDronePipelineDocker {
	config := drone.NewCIDronePipelineDocker()
	vol := drone.NewPipelineVolume().
		WithName("temp").
		WithTemp()
	config.Volumes = append(config.Volumes, *vol)

	for branch, branchFlow := range w.BranchFlows {
		if branchFlow.Skip {
			continue
		}

		for _, job := range branchFlow.Jobs {
			if job.Skip {
				continue
			}

			envVars := executil.EnvVars{}
			envVars.LoadFromEnviron()

			vol := drone.NewPipelineStepVolume().
				WithName("temp").
				WithPath("/go")

			image := fmt.Sprintf("${%s}/${%s}", EnvKeyDockerRegistryKey, strings.ToUpper(envVars.Parse(job.Builder)))
			image = envVars.Parse(image)
			step := drone.NewPipelineStep().
				WithName(fmt.Sprintf("%s_%s", str.ToLowerCamelCase(branch), job.Stage)).
				WithEnvs(branchFlow.Env).
				WithImage(image).
				WithVolume(*vol).
				WithCommands(job.Run...)

			if branch != "*" {
				step.WithBranchInclude(branch)
			}

			config.Steps = append(config.Steps, *step)
		}
	}

	return config
}
