package executil

import (
	"os/exec"
	"testing"
)

func TestRunWithLog(t *testing.T) {
	StdRun(exec.Command("go", "version"))
	//StdRun(exec.Command("go", "test"))
}
