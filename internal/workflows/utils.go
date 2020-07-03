package workflows

import (
	"github.com/profzone/eden-framework/internal/project"
)

var (
	BaseShipScript   = "eden ci ship --push"
	BaseDeployScript = "eden ci deploy --env ${CI_ENVIRONMENT_NAME}"
)

var DefaultJobForTest = project.Job{
	Stage:   project.STAGE_TEST,
	Builder: "BUILDER_${PROJECT_PROGRAM_LANGUAGE}",
	Run:     project.Script{"eden ci run test"},
}

var DefaultJobForBuild = project.Job{
	Stage:   project.STAGE_BUILD,
	Builder: "BUILDER_${PROJECT_PROGRAM_LANGUAGE}",
	Run:     project.Script{"eden ci run build"},
}

var DefaultJobForShip = project.Job{
	Stage:   project.STAGE_SHIP,
	Builder: "BUILDER_DOCKER",
	Run:     project.Script{BaseShipScript},
}

var DefaultJobForDeploy = project.Job{
	Stage:   project.STAGE_DEPLOY,
	Builder: "BUILDER_RANCHER",
	Run:     project.Script{BaseDeployScript},
}
