package workflows

import (
	"fmt"
	"github.com/profzone/eden-framework/internal/project"
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
				project.STAGE_TEST:  DefaultJobForTest,
				project.STAGE_BUILD: DefaultJobForBuild,
				project.STAGE_SHIP: DefaultJobForShip.Merge(&project.Job{
					Run: project.Script{fmt.Sprintf("%s --latest", BaseShipScript)},
				}),
				project.STAGE_DEPLOY: DefaultJobForDeploy.Merge(&project.Job{
					Run: project.Script{fmt.Sprintf("%s --latest", BaseDeployScript)},
				}),
			},
		},
		`release/*`: {
			Env: map[string]string{
				"GOENV": "TEST",
			},
			Jobs: project.Jobs{
				project.STAGE_TEST:  DefaultJobForTest,
				project.STAGE_BUILD: DefaultJobForBuild,
				project.STAGE_SHIP: DefaultJobForShip.Merge(&project.Job{
					Run: project.Script{fmt.Sprintf("%s --suffix ${CI_ENVIRONMENT_NAME}", BaseShipScript)},
				}),
				project.STAGE_DEPLOY: DefaultJobForDeploy.Merge(&project.Job{
					Run: project.Script{fmt.Sprintf("%s --suffix ${CI_ENVIRONMENT_NAME}", BaseDeployScript)},
				}),
			},
		},
		"master": {
			Env: map[string]string{
				"GOENV": "DEMO",
			},
			Jobs: project.Jobs{
				project.STAGE_TEST:   DefaultJobForTest,
				project.STAGE_BUILD:  DefaultJobForBuild,
				project.STAGE_SHIP:   DefaultJobForShip,
				project.STAGE_DEPLOY: DefaultJobForDeploy,
			},
		},
	},
}
