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
			Env: map[string]string{
				"GOENV": "PROD",
			},
			Jobs: project.Jobs{
				project.STAGE_TEST:  DefaultJobForTest,
				project.STAGE_BUILD: DefaultJobForBuild,
				project.STAGE_SHIP:  DefaultJobForShip,
				project.STAGE_DEPLOY: DefaultJobForDeploy.Merge(&project.Job{
					Run: project.Script{
						"eden ci deploy --env=STAGING",
						"eden ci deploy --env=TEST",
						"eden ci deploy --env=DEMO",
					},
				}),
			},
		},
		"feature/*": {
			Env: map[string]string{
				"GOENV": "STAGING",
			},
			Jobs: project.Jobs{
				project.STAGE_TEST:  DefaultJobForTest,
				project.STAGE_BUILD: DefaultJobForBuild,
				project.STAGE_SHIP:  DefaultJobForShip,
				project.STAGE_DEPLOY: DefaultJobForDeploy.Merge(&project.Job{
					Run: project.Script{
						"eden ci deploy",
					},
				}),
			},
		},
		"test/feature/*": {
			Extends: `feature/*`,
			Env: map[string]string{
				"GOENV": "TEST",
			},
		},
		`demo/feature/*`: {
			Extends: `feature/*`,
			Env: map[string]string{
				"GOENV": "DEMO",
			},
		},
	},
}
