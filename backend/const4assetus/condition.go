package const4assetus

import (
	"fmt"
	"slices"

	"github.com/strongo/validation"
)

// Condition is a closed, write-validated set describing the physical condition
// of an asset.
type Condition string

const (
	ConditionNew         Condition = "new"
	ConditionExcellent   Condition = "excellent"
	ConditionGood        Condition = "good"
	ConditionFair        Condition = "fair"
	ConditionNeedsRepair Condition = "needs_repair"
	ConditionBroken      Condition = "broken"
)

// Conditions is the closed set of valid conditions.
var Conditions = []Condition{
	ConditionNew,
	ConditionExcellent,
	ConditionGood,
	ConditionFair,
	ConditionNeedsRepair,
	ConditionBroken,
}

// IsValidCondition reports whether v is a member of the closed condition set.
func IsValidCondition(v Condition) bool {
	return slices.Contains(Conditions, v)
}

// ValidateCondition returns an error if v is not a valid condition. A condition
// is required on every asset (set by the creator).
func ValidateCondition(v Condition) error {
	if v == "" {
		return validation.NewErrRecordIsMissingRequiredField("condition")
	}
	if !IsValidCondition(v) {
		return validation.NewErrBadRecordFieldValue("condition",
			fmt.Sprintf("unknown condition %q, expected one of %v", v, Conditions))
	}
	return nil
}
