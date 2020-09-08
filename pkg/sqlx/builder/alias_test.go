package builder_test

import (
	"testing"

	"github.com/onsi/gomega"
	. "github.com/profzone/eden-framework/pkg/sqlx/builder"
	. "github.com/profzone/eden-framework/pkg/sqlx/builder/buidertestingutils"
)

func TestAlias(t *testing.T) {
	t.Run("alias", func(t *testing.T) {
		gomega.NewWithT(t).Expect(Alias(Expr("f_id"), "id")).To(BeExpr("(f_id) AS id"))
	})
}
