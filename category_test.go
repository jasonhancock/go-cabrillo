package cabrillo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCategoryRules(t *testing.T) {
	t.Run("PermissibleValues", func(t *testing.T) {
		rule := CategoryValidationRule{
			Field:             "BLAH",
			PermissibleValues: []string{"VALUE1", "VALUE2"},
		}

		t.Run("matches value", func(t *testing.T) {
			cat := Category{
				Name:  "BLAH",
				Value: "VALUE2",
			}

			require.NoError(t, rule.Evaluate(cat))
		})

		t.Run("doesn't match value", func(t *testing.T) {
			cat := Category{
				Name:  "BLAH",
				Value: "VALUE3",
			}

			require.Error(t, rule.Evaluate(cat))
		})

		t.Run("doesn't match name", func(t *testing.T) {
			cat := Category{
				Name:  "BLAHBLAH",
				Value: "VALUE3",
			}

			require.NoError(t, rule.Evaluate(cat))
		})
	})
}
