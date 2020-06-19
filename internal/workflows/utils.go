package workflows

import (
	"github.com/profzone/eden-framework/internal/project"
)

var (
	BaseShipScript   = "eden ci ship --push"
	BaseDeployScript = "eden ci deploy --env ${CI_ENVIRONMENT_NAME}"
)

var DefaultJobForTest = project.Job{
	Builder: "BUILDER_${PROJECT_PROGRAM_LANGUAGE}",
	Run:     project.Script{"eden ci run test"},
}

var DefaultJobForBuild = project.Job{
	Builder: "BUILDER_${PROJECT_PROGRAM_LANGUAGE}",
	Run:     project.Script{"eden ci run build"},
}

var DefaultJobForShip = project.Job{
	Builder: "BUILDER_DOCKER",
	Run:     project.Script{BaseShipScript},
}

var DefaultJobForDeploy = project.Job{
	Builder: "BUILDER_RANCHER",
	Run:     project.Script{BaseDeployScript},
}
