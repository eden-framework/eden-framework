package builder_test

import (
	"testing"

	"github.com/onsi/gomega"
	. "github.com/profzone/eden-framework/pkg/sqlx/builder"
	. "github.com/profzone/eden-framework/pkg/sqlx/builder/buidertestingutils"
)

func TestFunc(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		gomega.NewWithT(t).Expect(Func("")).To(BeExpr(""))
	})
	t.Run("count", func(t *testing.T) {
		gomega.NewWithT(t).Expect(Count()).To(BeExpr("COUNT(1)"))
	})
	t.Run("AVG", func(t *testing.T) {
		gomega.NewWithT(t).Expect(Avg()).To(BeExpr("AVG(*)"))
	})
}
