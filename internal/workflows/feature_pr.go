package workflows

import (
	"github.com/profzone/eden-framework/internal/project"
)

func init() {
	project.RegisterWorkFlow("feature-pr", FeaturePR)
}

var FeaturePR = &project.Workflow{
	BranchFlows: project.BranchFlows{
		"master": {
			Jobs: project.Jobs{
				project.STAGE_TEST:  DefaultJobForTest,
				project.STAGE_BUILD: DefaultJobForBuild,
				project.STAGE_SHIP:  DefaultJobForShip,
				project.STAGE_DEPLOY: DefaultJobForDeploy.Merge(&project.Job{
					Run: project.Script{
						"RANCHER_ENVIRONMENT=STAGING rancher-env.sh project deploy",
						"RANCHER_ENVIRONMENT=TEST rancher-env.sh project deploy",
						"RANCHER_ENVIRONMENT=DEMO rancher-env.sh project deploy",
					},
				}),
			},
		},
		`/^feature\/.*$/`: {
			Env: "STAGING",
			Jobs: project.Jobs{
				project.STAGE_TEST:  DefaultJobForTest,
				project.STAGE_BUILD: DefaultJobForBuild,
				project.STAGE_SHIP:  DefaultJobForShip,
				project.STAGE_DEPLOY: DefaultJobForDeploy.Merge(&project.Job{
					Run: project.Script{
						"project deploy",
					},
				}),
			},
		},
		`/^test/feature\/.*$/`: {
			Extends: `/^feature\/.*$/`,
			Env:     "TEST",
		},
		`/^demo/feature\/.*$/`: {
			Extends: `/^feature\/.*$/`,
			Env:     "DEMO",
		},
	},
}
