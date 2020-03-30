package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	tt := assert.New(t)

	keys := Keys{}

	tt.Equal(0, keys.Len())

	{
		keys.Add(PrimaryKey())
		tt.Equal(1, keys.Len())
	}
}
