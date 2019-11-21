package executil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseByEnv(t *testing.T) {
	tt := assert.New(t)

	os.Setenv("TEST", "1")

	envVar := EnvVars{}
	envVar.LoadFromEnviron()

	tt.Equal("1", envVar.Parse("${TEST}"))
	tt.Equal("${MISSING}", envVar.Parse("${MISSING}"))
	tt.Equal("${IGNORE}", envVar.Parse("$${IGNORE}"))
}
