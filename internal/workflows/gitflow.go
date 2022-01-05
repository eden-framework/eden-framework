package workflows

import (
	"fmt"
	"gitee.com/eden-framework/eden-framework/internal/project"
)

func init() {
	project.RegisterWorkFlow("gitflow", Gitflow)
}

var Gitflow = &project.Workflow{
	BranchFlows: project.BranchFlows{
		"develop": {
			Env: map[string]string{
				"GOENV": "STAGING",
			},
			Jobs: project.Jobs{
				DefaultJobForTest,
				DefaultJobForBuild,
				DefaultJobForShip.Merge(&project.Job{
					Run: project.Script{fmt.Sprintf("%s --latest", BaseShipScript)},
				}),
				DefaultJobForDeploy.Merge(&project.Job{
					Run: project.Script{fmt.Sprintf("%s --latest", BaseDeployScript)},
				}),
			},
		},
		`release/*`: {
			Env: map[string]string{
				"GOENV": "TEST",
			},
			Jobs: project.Jobs{
				DefaultJobForTest,
				DefaultJobForBuild,
				DefaultJobForShip.Merge(&project.Job{
					Run: project.Script{fmt.Sprintf("%s --suffix ${CI_ENVIRONMENT_NAME}", BaseShipScript)},
				}),
				DefaultJobForDeploy.Merge(&project.Job{
					Run: project.Script{fmt.Sprintf("%s --suffix ${CI_ENVIRONMENT_NAME}", BaseDeployScript)},
				}),
			},
		},
		"master": {
			Env: map[string]string{
				"GOENV": "DEMO",
			},
			Jobs: project.Jobs{
				DefaultJobForTest,
				DefaultJobForBuild,
				DefaultJobForShip,
				DefaultJobForDeploy,
			},
		},
	},
}
