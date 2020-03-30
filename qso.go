package cabrillo

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// QSO represents a contact.
type QSO struct {
	Frequency   string
	Mode        string
	Timestamp   time.Time
	TxInfo      Info
	RxInfo      Info
	Transmitter int
}

// RST is a signal report.
type RST struct {
	Readability int
	Strength    int
	Tone        int
}

// NewRST parses a signal report ("59" or "599", etc.) and returns it as an RST.
func NewRST(report string) (RST, error) {
	if len(report) < 2 || len(report) > 3 {
		return RST{}, errors.New("invalid RST report length")
	}
	var rst RST

	var err error
	rst.Readability, err = strconv.Atoi(string(report[0]))
	if err != nil {
		return RST{}, errors.Wrapf(err, "parsing readability digit %q", string(report[0]))
	}

	rst.Strength, err = strconv.Atoi(string(report[1]))
	if err != nil {
		return RST{}, errors.Wrapf(err, "parsing strength digit %q", string(report[1]))
	}

	if len(report) == 3 {
		rst.Tone, err = strconv.Atoi(string(report[2]))
		if err != nil {
			return RST{}, errors.Wrapf(err, "parsing tone digit %q", string(report[2]))
		}
	}

	return rst, nil
}

// String fullfills the stringer interface.
func (r RST) String() string {
	s := fmt.Sprintf("%d%d", r.Readability, r.Strength)
	if r.Tone > 0 {
		s += fmt.Sprintf("%d", r.Tone)
	}
	return s
}

// Info stores information about the sender or receiver participating in a QSO.
type Info struct {
	Callsign     string
	SignalReport RST
	Exchange     string
}

// NewQSO parses a line from a cabrillo log into a
func NewQSO(line string) (QSO, error) {
	fields := strings.Fields(line)

	if len(fields) != 11 && len(fields) != 12 {
		return QSO{}, errors.Errorf("invalid number of fields in QSO. Got %d, expected 11 or 12 - %q", len(fields), line)
	}

	qso := QSO{
		Frequency: fields[1],
		Mode:      fields[2],
		TxInfo: Info{
			Callsign: fields[5],
			Exchange: fields[7],
		},
		RxInfo: Info{
			Callsign: fields[8],
			Exchange: fields[10],
		},
	}

	var err error
	qso.Timestamp, err = time.Parse("2006-01-021504", fields[3]+fields[4])
	if err != nil {
		return QSO{}, err
	}

	qso.TxInfo.SignalReport, err = NewRST(fields[6])
	if err != nil {
		return QSO{}, errors.Wrap(err, "parsing tx RST")
	}

	qso.RxInfo.SignalReport, err = NewRST(fields[9])
	if err != nil {
		return QSO{}, errors.Wrap(err, "parsing rx RST")
	}

	if len(fields) == 12 {
		qso.Transmitter, err = strconv.Atoi(fields[11])
		if err != nil {
			return QSO{}, errors.Wrap(err, "parsing transmitter")
		}
	}

	return qso, nil
}
