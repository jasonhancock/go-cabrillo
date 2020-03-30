package cabrillo

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	// maxAddressLines is the maximum number of "ADDRESS:" lines allowed in a log.
	maxAddressLines = 6
	// maxLengthAddress is the maximum length of an "ADDRESS:" field.
	maxLengthAddress = 45
	// maxLengthName is the maximum allowed length of the "NAME:" field.
	maxLengthName = 75
)

var (
	errAddressTooLong      = errors.Errorf("address too long (maximum length %d characters)", maxLengthAddress)
	errNameTooLong         = errors.Errorf("name too long (maximum length %d characters)", maxLengthName)
	errTooManyAddressLines = errors.Errorf("only allowed up to %d ADDRESS lines", maxAddressLines)
)

// Log is a data structure representing an entire Cabrillo formatted Log file.
// Data can be parsed into the structure. The long-term plan is to be able to
// generate a Cabrillo formatted log file from the data structure.
type Log struct {
	Address          Address
	CallSign         string
	Categories       []Category
	Certificate      bool
	ClaimedScore     int
	Club             string
	Contest          string
	CreatedBy        string
	Email            string
	ExtensibleFields []ExtensibleField
	GridLocator      string
	Location         string
	Name             string
	OffTimes         []OffTime
	//TODO: when generating the line, each line is a max of 75 chars long. Use multiple lines. Host station prefixed with `@`
	Operators []string
	QSOs      []QSO
	SoapBox   []string
	Version   string
	XQSOs     []QSO
}

// Category returns the value for the specified category or an empty string if a
// value has not been specified for that category.
func (l *Log) Category(name string) string {
	for i := range l.Categories {
		if l.Categories[i].Name == name {
			return l.Categories[i].Value
		}
	}

	return ""
}

// ExtendedField returns the values for the specified extended field.
func (l *Log) ExtendedField(name string) []string {
	for i := range l.ExtensibleFields {
		if l.ExtensibleFields[i].Name == name {
			return l.ExtensibleFields[i].Values
		}
	}

	return nil
}

// AddCategory adds a category to the existing Log.
func (l *Log) AddCategory(name, value string) error {
	// Validate the category name.
	if _, ok := validCategories[name]; !ok {
		return errors.Errorf("unknown category %q", name)
	}

	for i := range l.Categories {
		if l.Categories[i].Name == name {
			// If a category is specified twice, last entry wins.
			l.Categories[i].Value = value
			return nil
		}
	}

	l.Categories = append(l.Categories, Category{Name: name, Value: value})
	return nil
}

// AddExtensibleField adds an extenisble field to the existing Log. If the field
// has been added previously, appends to the existing field.
func (l *Log) AddExtensibleField(name, value string) {
	for i := range l.ExtensibleFields {
		if l.ExtensibleFields[i].Name == name {
			l.ExtensibleFields[i].Values = append(l.ExtensibleFields[i].Values, value)
			return
		}
	}

	l.ExtensibleFields = append(
		l.ExtensibleFields,
		ExtensibleField{
			Name:   name,
			Values: []string{value},
		},
	)
}

// ExtensibleField represents a line prefixed by "X-" in the log other than
// "X-QSO". Each entry in Values represents a separate line on which the field
// was found. Thus if you had:
// X-COMMENT: some comment
// X-COMMENT: some other comment
// You would end up with one ExtensibleField with Name="COMMENT" with two entries in Values.
type ExtensibleField struct {
	Name   string
	Values []string
}

// Address represents someone's physical address.
type Address struct {
	// Maximum of 6 lines
	Address       []string
	City          string
	StateProvince string
	PostalCode    string
	Country       string
}

// ParseLog attempts to parse the data from the reader into a Log structure.
func ParseLog(r io.Reader) (Log, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return Log{}, err
	}

	l := Log{
		Certificate: true, // Defaults to yes per the specification.
	}
	lines := strings.Split(string(b), "\n")
	for lineNum, line := range lines {
		lineParts := strings.Fields(line)
		if len(lineParts) < 2 {
			continue
		}

		//originalTag := lineParts[0]
		lineParts[0] = strings.ToUpper(lineParts[0])

		switch lineParts[0] {
		case "ADDRESS:":
			addr := strings.Join(lineParts[1:], " ")
			if len(addr) > maxLengthAddress {
				return Log{}, newLineError(errAddressTooLong, lineNum)
			}
			l.Address.Address = append(l.Address.Address, addr)
			if len(l.Address.Address) > maxAddressLines {
				return Log{}, newLineError(errTooManyAddressLines, lineNum)
			}
		case "ADDRESS-CITY:":
			l.Address.City = strings.Join(lineParts[1:], " ")
		case "ADDRESS-COUNTRY:":
			l.Address.Country = strings.Join(lineParts[1:], " ")
		case "ADDRESS-POSTALCODE:":
			l.Address.PostalCode = strings.Join(lineParts[1:], " ")
		case "ADDRESS-STATE-PROVINCE:":
			l.Address.StateProvince = strings.Join(lineParts[1:], " ")
		case "CALLSIGN:":
			l.CallSign = lineParts[1]
		case "CERTIFICATE:":
			var err error
			l.Certificate, err = parseYN(lineParts[1])
			if err != nil {
				return Log{}, errors.Wrap(err, "parsing CERTIFICATE field")
			}
		case "CLAIMED-SCORE:":
			var err error
			l.ClaimedScore, err = strconv.Atoi(lineParts[1])
			if err != nil {
				return Log{}, errors.Wrap(err, "parsing CLAIMED-SCORE field")
			}
		case "CLUB:":
			l.Club = strings.Join(lineParts[1:], " ")
		case "CONTEST:":
			l.Contest = strings.Join(lineParts[1:], " ")
		case "CREATED-BY:":
			l.CreatedBy = strings.Join(lineParts[1:], " ")
		case "EMAIL:":
			addr, err := mail.ParseAddress(strings.Join(lineParts[1:], " "))
			if err != nil {
				return Log{}, newLineError(errors.Wrap(err, "parsing email address"), lineNum)
			}
			l.Email = addr.Address
		case "END-OF-LOG:":
		case "GRID-LOCATOR:":
			// TODO: should we validate the grid locator?
			l.GridLocator = strings.Join(lineParts[1:], " ")
		case "LOCATION:":
			l.Location = strings.Join(lineParts[1:], " ")
		case "NAME:":
			l.Name = strings.Join(lineParts[1:], " ")
			if len(l.Name) > maxLengthName {
				return Log{}, newLineError(errNameTooLong, lineNum)
			}
		case "OFFTIME:":
			ot, err := parseOffTime(strings.Join(lineParts[1:], " "))
			if err != nil {
				return Log{}, newLineError(err, lineNum)
			}
			l.OffTimes = append(l.OffTimes, ot)
		case "OPERATORS:":
			l.Operators = append(l.Operators, operatorsField(strings.Join(lineParts[1:], " "))...)
		case "QSO:":
			qso, err := NewQSO(line)
			if err != nil {
				return Log{}, newLineError(err, lineNum)
			}
			l.QSOs = append(l.QSOs, qso)
		case "SOAPBOX:":
			// max line length 75
			l.SoapBox = append(l.SoapBox, strings.Join(lineParts[1:], " "))
		case "START-OF-LOG:":
			l.Version = lineParts[1]
		case "X-QSO:":
			qso, err := NewQSO(line)
			if err != nil {
				return Log{}, newLineError(err, lineNum)
			}

			l.XQSOs = append(l.XQSOs, qso)

		default:
			if strings.HasPrefix(lineParts[0], "CATEGORY-") && strings.HasSuffix(lineParts[0], ":") {
				name := strings.TrimSuffix(strings.TrimPrefix(lineParts[0], "CATEGORY-"), ":")
				if err := l.AddCategory(name, strings.Join(lineParts[1:], " ")); err != nil {
					return Log{}, newLineError(err, lineNum)
				}
				continue
			}
			if strings.HasPrefix(lineParts[0], "X-") {
				name := strings.TrimSuffix(strings.TrimPrefix(lineParts[0], "X-"), ":")
				l.AddExtensibleField(name, strings.Join(lineParts[1:], " "))
				continue
			}
			// we shouldn't be calling out to log here. Should this be an error?
			log.Printf("unknown tag %q", line)
		}
	}

	return l, nil
}

// The operators field is a space or comma delimited field
func operatorsField(str string) []string {
	str = strings.ReplaceAll(str, ",", " ")
	pieces := strings.Fields(str)
	var data []string
	for _, v := range pieces {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		data = append(data, v)
	}

	return data
}

func parseYN(str string) (bool, error) {
	str = strings.TrimSpace(strings.ToUpper(str))
	if str == "YES" {
		return true, nil
	}
	if str == "NO" {
		return false, nil
	}

	return false, errors.Errorf("cannot parse %q as either YES or NO", str)
}

type lineError struct {
	error
	lineNumber int
}

func (e lineError) Error() string {
	return fmt.Sprintf("%d: %s", e.lineNumber, e.error.Error())
}

func newLineError(err error, lineNumber int) error {
	return &lineError{
		error:      err,
		lineNumber: lineNumber,
	}
}
