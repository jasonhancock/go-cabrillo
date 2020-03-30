package cabrillo

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestRST(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		tests := []struct {
			description string
			input       RST
			expected    string
		}{
			{"with tone", RST{Readability: 5, Strength: 9, Tone: 9}, "599"},
			{"without tone", RST{Readability: 5, Strength: 9, Tone: 0}, "59"},
		}

		for _, tt := range tests {
			t.Run(tt.description, func(t *testing.T) {
				require.Equal(t, tt.input.String(), tt.expected)
			})
		}
	})

	t.Run("NewRST", func(t *testing.T) {
		tests := []struct {
			description string
			input       string
			expected    string
			err         error
		}{
			{"readability-strength", "59", "59", nil},
			{"readability-strength-tone", "599", "599", nil},
			{"invalid-length-short", "5", "", errors.New("invalid RST report length")},
			{"invalid-length-long", "5999", "", errors.New("invalid RST report length")},
			{"invalid-readability", "a9", "", errors.New("parsing readability digit")},
			{"invalid-strength", "5a", "", errors.New("parsing strength digit")},
			{"invalid-tone", "59a", "", errors.New("parsing tone digit")},
		}

		for _, tt := range tests {
			t.Run(tt.description, func(t *testing.T) {
				rst, err := NewRST(tt.input)

				if tt.err == nil {
					require.NoError(t, err)
					require.Equal(t, tt.expected, rst.String())
					return
				}

				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err.Error())
			})
		}
	})
}

func TestQSO(t *testing.T) {
	t.Run("NewQSO", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			qso, err := NewQSO("QSO:  7030 CW 2017-11-25 2134 K1IR          599 5      IQ3R          599 15\n")
			require.NoError(t, err)
			require.Equal(t, "7030", qso.Frequency)
			require.Equal(t, "CW", qso.Mode)
			require.Equal(t, "201711252134", qso.Timestamp.Format("200601021504"))

			require.Equal(t, "K1IR", qso.TxInfo.Callsign)
			require.Equal(t, "599", qso.TxInfo.SignalReport.String())
			require.Equal(t, "5", qso.TxInfo.Exchange)

			require.Equal(t, "IQ3R", qso.RxInfo.Callsign)
			require.Equal(t, "599", qso.RxInfo.SignalReport.String())
			require.Equal(t, "15", qso.RxInfo.Exchange)

			require.Equal(t, 0, qso.Transmitter)
		})

		t.Run("normal-with-operator", func(t *testing.T) {
			qso, err := NewQSO("QSO:  7250 PH 2000-10-26 0711 AA1ZZZ          59  05     WA6MIC        59  03     0")
			require.NoError(t, err)
			require.Equal(t, "7250", qso.Frequency)
			require.Equal(t, "PH", qso.Mode)
			require.Equal(t, "200010260711", qso.Timestamp.Format("200601021504"))

			require.Equal(t, "AA1ZZZ", qso.TxInfo.Callsign)
			require.Equal(t, "59", qso.TxInfo.SignalReport.String())
			require.Equal(t, "05", qso.TxInfo.Exchange)

			require.Equal(t, "WA6MIC", qso.RxInfo.Callsign)
			require.Equal(t, "59", qso.RxInfo.SignalReport.String())
			require.Equal(t, "03", qso.RxInfo.Exchange)

			require.Equal(t, 0, qso.Transmitter)
		})
	})
}
