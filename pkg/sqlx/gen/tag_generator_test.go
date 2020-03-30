package gen

import (
	"os"
	"testing"

	"github.com/profzone/eden-framework/pkg/codegen"
)

func init() {
	os.Chdir("./test")
}

func TestTagGen(t *testing.T) {
	clientGenerator := TagGenerator{
		WithDefaults: true,
	}
	clientGenerator.StructNames = []string{"User", "User2"}
	codegen.Generate(&clientGenerator)
}
