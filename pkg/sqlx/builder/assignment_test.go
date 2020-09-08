package builder_test

import (
	"context"
	"testing"

	"github.com/onsi/gomega"
	. "github.com/profzone/eden-framework/pkg/sqlx/builder"
	. "github.com/profzone/eden-framework/pkg/sqlx/builder/buidertestingutils"
)

func TestAssignment(t *testing.T) {
	t.Run("ColumnsAndValues", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			ColumnsAndValues(Cols("a", "b"), 1, 2, 3, 4).Ex(ContextWithToggles(context.Background(), Toggles{
				ToggleUseValues: true,
			})),
		).To(BeExpr("(a,b) VALUES (?,?),(?,?)", 1, 2, 3, 4))
	})
}
