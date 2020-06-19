package project

import (
	"fmt"
	"github.com/profzone/eden-framework/pkg/executil"
	str "github.com/profzone/eden-framework/pkg/strings"
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
			for stage, job := range branchFlow.Jobs {
				if !job.Skip {
					envVars := executil.EnvVars{}
					envVars.LoadFromEnviron()

					ciJob := NewCIJob(stage).
						WithTags(project.Group).
						WithEnv(branchFlow.Env).
						WithImage(fmt.Sprintf("${%s}/${%s}", DOCKER_REGISTRY_KEY, strings.ToUpper(envVars.Parse(job.Builder)))).
						WithArtifacts(job.Artifacts).
						WithScript(job.Run...)

					if branch != "*" {
						ciJob = ciJob.WithOnly(branch)
					}

					ciConfig = ciConfig.AddJob(
						fmt.Sprintf("%s_%s", str.ToLowerCamelCase(branch), stage),
						ciJob,
					)
				}
			}
		}
	}

	return ciConfig
}
