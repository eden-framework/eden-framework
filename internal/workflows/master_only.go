package workflows

import (
	"gitee.com/eden-framework/eden-framework/internal/project"
)

func init() {
	project.RegisterWorkFlow("master-only", MasterOnly)
}

var MasterOnly = &project.Workflow{
	BranchFlows: project.BranchFlows{
		"master": {
			Env: map[string]string{
				"GOENV": "STAGING",
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
