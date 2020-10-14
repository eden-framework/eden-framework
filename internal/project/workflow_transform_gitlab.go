package project

import (
	"fmt"
	"github.com/eden-framework/context"
	str "github.com/eden-framework/strings"
	"strings"
)

var DefaultCICache = CICache{
	Key:       "${CI_BUILD_REF}_${CI_BUILD_REF_NAME}",
	UnTracked: true,
}

func (w *Workflow) ToCIConfig(project *Project) *CIConfig {
	ciConfig := NewCIConfig().
		WithCache(DefaultCICache).
		WithStages(STAGE_TEST, STAGE_BUILD, STAGE_SHIP, STAGE_DEPLOY)

	for branch, branchFlow := range w.BranchFlows {
		if !branchFlow.Skip {
			for _, job := range branchFlow.Jobs {
				if !job.Skip {
					envVars := context.EnvVars{}
					envVars.LoadFromEnviron()

					ciJob := NewCIJob(job.Stage).
						WithTags(project.Group).
						WithEnv(branchFlow.Env["GOENV"]).
						WithImage(fmt.Sprintf("${%s}/${%s}", EnvKeyDockerRegistryKey, strings.ToUpper(envVars.Parse(job.Builder)))).
						WithArtifacts(job.Artifacts).
						WithScript(job.Run...)

					if branch != "*" {
						ciJob = ciJob.WithOnly(branch)
					}

					ciConfig = ciConfig.AddJob(
						fmt.Sprintf("%s_%s", str.ToLowerCamelCase(branch), job.Stage),
						ciJob,
					)
				}
			}
		}
	}

	return ciConfig
}
