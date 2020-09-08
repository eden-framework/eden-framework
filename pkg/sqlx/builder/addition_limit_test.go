package builder_test

import (
	"testing"

	"github.com/onsi/gomega"
	. "github.com/profzone/eden-framework/pkg/sqlx/builder"
	. "github.com/profzone/eden-framework/pkg/sqlx/builder/buidertestingutils"
)

func TestLimit(t *testing.T) {
	table := T("T")

	t.Run("select limit", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Select(nil).
				From(
					table,
					Where(
						Col("F_a").Eq(1),
					),
					Limit(1),
				),
		).To(BeExpr(`
SELECT * FROM T
WHERE f_a = ?
LIMIT 1
`, 1))
	})
	t.Run("select without limit", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Select(nil).
				From(
					table,
					Where(
						Col("F_a").Eq(1),
					),
					Limit(-1),
				),
		).To(BeExpr(`
SELECT * FROM T
WHERE f_a = ?
`, 1,
		))
	})
	t.Run("select limit and offset", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Select(nil).
				From(
					table,
					Where(
						Col("F_a").Eq(1),
					),
					Limit(1).Offset(200),
				),
		).To(BeExpr(`
SELECT * FROM T
WHERE f_a = ?
LIMIT 1 OFFSET 200
`,
			1,
		))
	})
}
