package project

import (
	"fmt"
	"github.com/eden-framework/context"
	"github.com/eden-framework/eden-framework/internal/project/drone"
	"github.com/eden-framework/eden-framework/internal/project/drone/enums"
	str "github.com/eden-framework/strings"
	"strings"
)

func (w *Workflow) ToDroneConfig(p *Project) *drone.CIDronePipelineDocker {
	config := drone.NewCIDronePipelineDocker()
	vol := drone.NewPipelineVolume().
		WithName("temp").
		WithTemp()
	config.Volumes = append(config.Volumes, *vol)

	hostVol := drone.NewPipelineVolume().WithName("host").WithHost(&drone.PipelineVolumeHost{
		Path: "/var/run/docker.sock",
	})
	config.Volumes = append(config.Volumes, *hostVol)

	for branch, branchFlow := range w.BranchFlows {
		if branchFlow.Skip {
			continue
		}

		for _, job := range branchFlow.Jobs {
			if job.Skip {
				continue
			}

			envVars := context.EnvVars{}
			envVars.LoadFromEnviron()

			vol := drone.NewPipelineStepVolume().
				WithName("temp").
				WithPath("/go")
			volHost := drone.NewPipelineStepVolume().WithName("host").WithPath("/var/run/docker.sock")

			image := fmt.Sprintf("${%s}/${%s}", EnvKeyDockerRegistryKey, strings.ToUpper(envVars.Parse(job.Builder)))
			image = envVars.Parse(image)
			step := drone.NewPipelineStep().
				WithName(fmt.Sprintf("%s_%s", str.ToLowerCamelCase(branch), job.Stage)).
				WithEnvs(branchFlow.Env).
				WithImage(image).
				WithVolume(*vol, *volHost).
				WithCommands(job.Run...)
			step.Pull = enums.DRONE_CI_STEP_PULL__always

			if branch != "*" {
				step.WithBranchInclude(branch)
			}

			config.Steps = append(config.Steps, *step)
		}
	}

	return config
}
