package extras4assetus

import (
	"strings"

	"github.com/strongo/validation"
)

// WithOptionalRegNumberField carries an optional registration/serial number,
// ported from the legacy extras4assetus.WithOptionalRegNumberField.
type WithOptionalRegNumberField struct {
	RegNumber string `json:"regNumber,omitempty" firestore:"regNumber,omitempty" dalgo:"regNumber,omitempty"`
}

// Validate returns an error if the reg number has leading/trailing spaces.
func (v *WithOptionalRegNumberField) Validate() error {
	if regNumber := strings.TrimSpace(v.RegNumber); regNumber != v.RegNumber {
		return validation.NewErrBadRecordFieldValue("regNumber", "should not have leading or trailing spaces")
	}
	return nil
}

// WithMakeModelFields carries a vehicle make and model, ported from the legacy
// extras4assetus.WithMakeModelFields.
type WithMakeModelFields struct {
	Make  string `json:"make,omitempty" firestore:"make,omitempty"`
	Model string `json:"model,omitempty" firestore:"model,omitempty"`
}

// GenerateTitleFromMakeModelAndRegNumber builds an asset title from make/model
// plus an optional reg number.
func (v *WithMakeModelFields) GenerateTitleFromMakeModelAndRegNumber(regNumber string) string {
	title := make([]string, 0, 4)
	if v.Make != "" {
		title = append(title, v.Make)
	}
	if v.Model != "" {
		title = append(title, v.Model)
	}
	if regNumber != "" {
		title = append(title, "#", regNumber)
	}
	if len(title) == 0 {
		return ""
	}
	return strings.Join(title, " ")
}

// Validate requires make and model and rejects leading/trailing spaces.
func (v *WithMakeModelFields) Validate() error {
	if makeValue := strings.TrimSpace(v.Make); makeValue == "" {
		return validation.NewErrRecordIsMissingRequiredField("make")
	} else if makeValue != v.Make {
		return validation.NewErrBadRecordFieldValue("make", "should not have leading or trailing spaces")
	}
	if model := strings.TrimSpace(v.Model); model == "" {
		return validation.NewErrRecordIsMissingRequiredField("model")
	} else if model != v.Model {
		return validation.NewErrBadRecordFieldValue("model", "should not have leading or trailing spaces")
	}
	return nil
}

// WithMakeModelRegNumberFields combines make/model with an optional reg number.
type WithMakeModelRegNumberFields struct {
	WithMakeModelFields
	WithOptionalRegNumberField
}

// Validate validates both the make/model and the reg number.
func (v *WithMakeModelRegNumberFields) Validate() error {
	if err := v.WithMakeModelFields.Validate(); err != nil {
		return err
	}
	if err := v.WithOptionalRegNumberField.Validate(); err != nil {
		return err
	}
	return nil
}
