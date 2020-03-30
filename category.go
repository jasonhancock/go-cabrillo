package cabrillo

import (
	"strings"

	"github.com/pkg/errors"
)

// Various constants for category names.
const (
	CategoryAssisted    = "ASSISTED"
	CategoryBand        = "BAND"
	CategoryMode        = "MODE"
	CategoryOperator    = "OPERATOR"
	CategoryPower       = "POWER"
	CategoryStation     = "STATION"
	CategoryTime        = "TIME"
	CategoryTransmitter = "TRANSMITTER"
	CategoryOverlay     = "OVERLAY"
)

// The list of categories officially recognized in the specification as a map
// for quick lookups.
var validCategories = map[string]struct{}{
	CategoryAssisted:    struct{}{},
	CategoryBand:        struct{}{},
	CategoryMode:        struct{}{},
	CategoryOperator:    struct{}{},
	CategoryPower:       struct{}{},
	CategoryStation:     struct{}{},
	CategoryTime:        struct{}{},
	CategoryTransmitter: struct{}{},
	CategoryOverlay:     struct{}{},
}

// Category is a key-value pair of data.
type Category struct {
	Name  string
	Value string
}

// CategoryValidationRule is a rule to validate a category. Some categories have
// fixed lists of permissible values. This likely needs to be changed into an interface.
type CategoryValidationRule struct {
	Field             string
	PermissibleValues []string
}

// Evaluate runs the rule annd determines if the specific instance of Catetgory
// fulfills the rule.
func (r *CategoryValidationRule) Evaluate(cat Category) error {
	// only evaluate rules for this field
	if cat.Name != r.Field {
		return nil
	}

	if len(r.PermissibleValues) > 0 {
		found := false
		for _, v := range r.PermissibleValues {
			if cat.Value == v {
				found = true
				break
			}
		}
		if !found {
			return errors.Errorf("value %q not in possible values %q", cat.Value, strings.Join(r.PermissibleValues, ","))
		}
	}

	return nil
}

// DefaultRules tries to capture a default list of rules from the Cabrillo spec definition.
func DefaultRules() []CategoryValidationRule {
	return []CategoryValidationRule{
		{
			Field: "ASSISTED",
			PermissibleValues: []string{
				"ASSISTED",
				"NON-ASSISTED",
			},
		},
		{
			Field: "BAND",
			PermissibleValues: []string{
				"ALL",
				"160M",
				"80M",
				"40M",
				"20M",
				"15M",
				"10M",
				"6M",
				"4M",
				"2M",
				"222",
				"432",
				"902",
				"1.2G",
				"2.3G",
				"3.4G",
				"5.7G",
				"10G",
				"24G",
				"47G",
				"75G",
				"123G",
				"134G",
				"241G",
				"Light",
				"VHF-3-BAND",
				"VHF-FM-ONLY",
			},
		},
		{
			Field: "MODE",
			PermissibleValues: []string{
				"CW",
				"DIGI",
				"FM",
				"RTTY",
				"SSB",
				"MIXED",
			},
		},
		{
			Field: "OPERATOR",
			PermissibleValues: []string{
				"SINGLE-OP",
				"MULTI-OP",
				"CHECKLOG",
			},
		},
		{
			Field: "POWER",
			PermissibleValues: []string{
				"HIGH",
				"LOW",
				"QRP",
			},
		},
		{
			// TODO: only required for multi-operator entries
			Field: "TRANSMITTER",
			PermissibleValues: []string{
				"ONE",
				"TWO",
				"LIMITED",
				"UNLIMITED",
				"SWL",
			},
		},
	}
}
