package workflows

import (
	"github.com/profzone/eden-framework/internal/project"
)

func init() {
	project.RegisterWorkFlow("master-only", MasterOnly)
}

var MasterOnly = &project.Workflow{
	BranchFlows: project.BranchFlows{
		"master": {
			Env: "STAGING",
			Jobs: project.Jobs{
				project.STAGE_TEST:   DefaultJobForTest,
				project.STAGE_BUILD:  DefaultJobForBuild,
				project.STAGE_SHIP:   DefaultJobForShip,
				project.STAGE_DEPLOY: DefaultJobForDeploy,
			},
		},
	},
}
