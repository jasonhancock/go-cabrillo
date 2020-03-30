package cabrillo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {
	t.Run("cq-ww-dx.log", func(t *testing.T) {
		fh, err := os.Open("testdata/cq-ww-dx.log")
		require.NoError(t, err)
		defer fh.Close()

		_, err = ParseLog(fh)
		require.NoError(t, err)
	})

	t.Run("k1ir.log", func(t *testing.T) {
		fh, err := os.Open("testdata/k1ir.log")
		require.NoError(t, err)
		defer fh.Close()

		_, err = ParseLog(fh)
		require.NoError(t, err)
	})

	t.Run("allfields.log", func(t *testing.T) {
		fh, err := os.Open("testdata/allfields.log")
		require.NoError(t, err)
		defer fh.Close()

		l, err := ParseLog(fh)
		require.NoError(t, err)

		require.Equal(t, "AG6K", l.CallSign)
		require.Equal(t, 1234, l.ClaimedScore)
		require.Equal(t, "some club somewhereville USA", l.Club)
		require.Equal(t, "someContest", l.Contest)
		require.Equal(t, "Some Logging Software", l.CreatedBy)
		require.Equal(t, "DM14cc", l.GridLocator)
		require.Equal(t, "WMA", l.Location)
		require.Equal(t, "John Doe", l.Name)
		require.Equal(t, "3.0", l.Version)
		require.False(t, l.Certificate)

		require.Equal(t, "NON-ASSISTED", l.Category(CategoryAssisted))
		require.Equal(t, "ALL", l.Category(CategoryBand))
		require.Equal(t, "PH", l.Category(CategoryMode))
		require.Equal(t, "SINGLE-OP", l.Category(CategoryOperator))
		require.Equal(t, "HIGH", l.Category(CategoryPower))
		require.Equal(t, "FIXED", l.Category(CategoryStation))
		require.Equal(t, "24-HOURS", l.Category(CategoryTime))
		require.Equal(t, "ONE", l.Category(CategoryTransmitter))
		require.Equal(t, "ROOKIE", l.Category(CategoryOverlay))
		require.Equal(t, "", l.Category("UNKNOWN-CATEGORY"))

		require.Len(t, l.Operators, 6)
		require.Contains(t, l.Operators, "AAAA")
		require.Contains(t, l.Operators, "@AAAC")
		require.Contains(t, l.Operators, "ZZZX")

		require.Len(t, l.SoapBox, 2)
		require.Equal(t, "Put your comments here.", l.SoapBox[0])
		require.Equal(t, "Use multiple lines if needed.", l.SoapBox[1])

		require.Len(t, l.Address.Address, 3)
		require.Equal(t, "Address line 1", l.Address.Address[0])
		require.Equal(t, "Address line 2", l.Address.Address[1])
		require.Equal(t, "Address line 3", l.Address.Address[2])
		require.Equal(t, "Some City", l.Address.City)
		require.Equal(t, "CA", l.Address.StateProvince)
		require.Equal(t, "12345", l.Address.PostalCode)
		require.Equal(t, "USA", l.Address.Country)

		vals := l.ExtendedField("FIELD1")
		require.Len(t, vals, 2)
		require.Equal(t, "line 1", vals[0])
		require.Equal(t, "line 2", vals[1])

		vals = l.ExtendedField("FIELD3")
		require.Len(t, vals, 0)

		require.Len(t, l.OffTimes, 4)
		for _, v := range l.OffTimes {
			require.False(t, v.Begin.IsZero())
			require.False(t, v.End.IsZero())
		}

		require.Len(t, l.QSOs, 4)
		require.Equal(t, "3799", l.QSOs[0].Frequency)
		require.Equal(t, "PH", l.QSOs[0].Mode)

		require.Len(t, l.XQSOs, 1)
		require.Equal(t, "7250", l.XQSOs[0].Frequency)
		require.Equal(t, "PH", l.XQSOs[0].Mode)
	})
}
