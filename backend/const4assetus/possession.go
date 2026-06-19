package const4assetus

import (
	"fmt"
	"slices"

	"github.com/strongo/validation"
)

// Possession describes how an asset is possessed (owned, leased, rented, etc.),
// ported from the legacy AssetPossession.
type Possession string

const (
	PossessionUnknown     Possession = "unknown"
	PossessionUndisclosed Possession = "undisclosed"
	PossessionOwning      Possession = "owning"
	PossessionLeasing     Possession = "leasing"
	PossessionRenting     Possession = "renting"
)

// Possessions is the list of all possible possession values.
var Possessions = []Possession{
	PossessionUnknown,
	PossessionUndisclosed,
	PossessionOwning,
	PossessionLeasing,
	PossessionRenting,
}

// ValidatePossession validates a possession value. When required is true an
// empty value is rejected as a missing required field.
func ValidatePossession(v Possession, required bool) error {
	if required && v == "" {
		return validation.NewErrRecordIsMissingRequiredField("possession")
	}
	if !slices.Contains(Possessions, v) {
		return validation.NewErrBadRecordFieldValue("possession",
			fmt.Sprintf("unknown possession %q, expected one of %v", v, Possessions))
	}
	return nil
}
