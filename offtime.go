package cabrillo

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

// OffTime is used to indicate off-time.
// OFFTIME: 2002-03-22 0300 2002-03-22 0743
// yyyy-mm-dd nnnn yyyy-mm-dd nnnn
// -----begin----- ------end------
type OffTime struct {
	Begin time.Time
	End   time.Time
}

// parseOffTime expects to be sent a string like:
// 2002-03-22 0300 2002-03-22 0743
func parseOffTime(str string) (OffTime, error) {
	pieces := strings.Fields(str)
	if len(pieces) != 4 {
		return OffTime{}, errors.New("invalid number of fields in offtime")
	}

	const format = "2006-01-02 1504"

	var ot OffTime
	var err error
	ot.Begin, err = time.Parse(format, pieces[0]+" "+pieces[1])
	if err != nil {
		return OffTime{}, errors.Wrap(err, "parsing begin time")
	}

	ot.End, err = time.Parse(format, pieces[2]+" "+pieces[3])
	if err != nil {
		return OffTime{}, errors.Wrap(err, "parsing end time")
	}

	return ot, nil
}
