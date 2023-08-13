package cabrillo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
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
		return RST{}, fmt.Errorf("parsing readability digit %q: %w", string(report[0]), err)
	}

	rst.Strength, err = strconv.Atoi(string(report[1]))
	if err != nil {
		return RST{}, fmt.Errorf("parsing strength digit %q: %w", string(report[1]), err)
	}

	if len(report) == 3 {
		rst.Tone, err = strconv.Atoi(string(report[2]))
		if err != nil {
			return RST{}, fmt.Errorf("parsing tone digit %q: %w", string(report[2]), err)
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

// NewQSO parses a line from a cabrillo log into a QSQ struct. exchangeFields
// specifies the number of fields in the exchange. If the exchange was only a
// serial number, this should be set to 1. If the exchange is a name, serial
// number, and QTH all delimited by spaces, set this to 3.
func NewQSO(line string, exchangeFields int) (QSO, error) {
	fields := strings.Fields(line)

	fieldsMin := 9 + exchangeFields*2
	fieldsMax := 10 + exchangeFields*2

	if len(fields) != fieldsMin && len(fields) != fieldsMax {
		return QSO{}, fmt.Errorf(
			"invalid number of fields in QSO. Got %d, expected %d or %d - %q",
			len(fields),
			fieldsMin,
			fieldsMax,
			line,
		)
	}

	qso := QSO{
		Frequency: fields[1],
		Mode:      fields[2],
		TxInfo: Info{
			Callsign: fields[5],
			Exchange: strings.Join(fields[7:7+exchangeFields], " "),
		},
		RxInfo: Info{
			Callsign: fields[7+exchangeFields],
			Exchange: strings.Join(fields[7+exchangeFields+2:7+exchangeFields+2+exchangeFields], " "),
		},
	}

	var err error
	qso.Timestamp, err = time.Parse("2006-01-021504", fields[3]+fields[4])
	if err != nil {
		return QSO{}, err
	}

	qso.TxInfo.SignalReport, err = NewRST(fields[6])
	if err != nil {
		return QSO{}, fmt.Errorf("parsing tx RST: %w", err)
	}

	qso.RxInfo.SignalReport, err = NewRST(fields[7+exchangeFields+1])
	if err != nil {
		return QSO{}, fmt.Errorf("parsing rx RST: %w", err)
	}

	if len(fields) == fieldsMax {
		qso.Transmitter, err = strconv.Atoi(fields[fieldsMax-1])
		if err != nil {
			return QSO{}, fmt.Errorf("parsing transmitter: %w", err)
		}
	}

	return qso, nil
}
