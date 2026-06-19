package extras4assetus

import (
	"time"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/sneat-core-modules/core/extra"
	"github.com/sneat-co/sneat-go-core/geo"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

func init() {
	RegisterAssetExtraFactory(AssetExtraTypeDocument, func() AssetExtra {
		return new(AssetDocumentExtra)
	})
}

var _ extra.Data = (*AssetDocumentExtra)(nil)

// AssetDocumentExtra is the typed extra for document assets. It carries the full
// legacy shape: docType, number, batchNumber, countryID, issuedBy, issuedOn,
// effectiveFrom and expiresOn. Per-doc-type validation is applied from
// standardDocTypesByID (e.g. a passport requires number + validity).
type AssetDocumentExtra struct {
	WithOptionalRegNumberField // legacy "regNumber" alias for a document number

	DocType       const4assetus.Type `json:"docType,omitempty" firestore:"docType,omitempty"`
	Number        string             `json:"number,omitempty" firestore:"number,omitempty"`
	BatchNumber   string             `json:"batchNumber,omitempty" firestore:"batchNumber,omitempty"`
	CountryID     geo.CountryAlpha2  `json:"countryID,omitempty" firestore:"countryID,omitempty"`
	IssuedBy      string             `json:"issuedBy,omitempty" firestore:"issuedBy,omitempty"`
	IssuedOn      string             `json:"issuedOn,omitempty" firestore:"issuedOn,omitempty"`
	EffectiveFrom string             `json:"effectiveFrom,omitempty" firestore:"effectiveFrom,omitempty"`
	ExpiresOn     string             `json:"expiresOn,omitempty" firestore:"expiresOn,omitempty"`
}

// RequiredFields implements extra.Data.
func (v *AssetDocumentExtra) RequiredFields() []string {
	return nil
}

// IndexedFields implements extra.Data; ported from the legacy declaration.
func (v *AssetDocumentExtra) IndexedFields() []string {
	return []string{"expiresOn", "effectiveFrom"}
}

// GetBrief implements extra.Data.
func (v *AssetDocumentExtra) GetBrief() extra.Data {
	return &AssetDocumentExtra{
		WithOptionalRegNumberField: v.WithOptionalRegNumberField,
		DocType:                    v.DocType,
		Number:                     v.Number,
		IssuedOn:                   v.IssuedOn,
		EffectiveFrom:              v.EffectiveFrom,
		ExpiresOn:                  v.ExpiresOn,
	}
}

// Validate returns an error if the document extra is not valid. It validates the
// date fields, the country code, and applies the per-doc-type schema from
// standardDocTypesByID.
func (v *AssetDocumentExtra) Validate() (err error) {
	if err = v.WithOptionalRegNumberField.Validate(); err != nil {
		return err
	}
	if v.CountryID != "" && !geo.IsValidCountryAlpha2(v.CountryID) {
		return validation.NewErrBadRecordFieldValue("countryID", "invalid country alpha-2 code: "+string(v.CountryID))
	}
	if v.IssuedOn != "" {
		if _, err = validate.DateString(v.IssuedOn); err != nil {
			return validation.NewErrBadRecordFieldValue("issuedOn", err.Error())
		}
	}
	var effectiveFrom, expiresOn time.Time
	if v.EffectiveFrom != "" {
		if effectiveFrom, err = validate.DateString(v.EffectiveFrom); err != nil {
			return validation.NewErrBadRecordFieldValue("effectiveFrom", err.Error())
		}
	}
	if v.ExpiresOn != "" {
		if expiresOn, err = validate.DateString(v.ExpiresOn); err != nil {
			return validation.NewErrBadRecordFieldValue("expiresOn", err.Error())
		}
	}
	if !effectiveFrom.IsZero() && !expiresOn.IsZero() && expiresOn.Before(effectiveFrom) {
		return validation.NewErrBadRecordFieldValue("expiresOn", "is before `effectiveFrom`")
	}
	if err = v.validateDocTypeSchema(); err != nil {
		return err
	}
	return nil
}

// docNumber returns the document number, accepting either the explicit Number
// field or the legacy regNumber alias.
func (v *AssetDocumentExtra) docNumber() string {
	if v.Number != "" {
		return v.Number
	}
	return v.RegNumber
}

// validity returns the document's validity date (expiresOn).
func (v *AssetDocumentExtra) validity() string {
	return v.ExpiresOn
}

// validateDocTypeSchema applies the per-doc-type validation schema from
// standardDocTypesByID. Document types with no standard schema impose no extra
// requirements.
func (v *AssetDocumentExtra) validateDocTypeSchema() error {
	def, ok := standardDocTypesByID[v.DocType]
	if !ok {
		return nil
	}
	f := def.Fields
	if f.Number != nil && f.Number.Required && v.docNumber() == "" {
		return validation.NewErrRecordIsMissingRequiredField("number")
	}
	if f.ValidTill != nil {
		if f.ValidTill.Required && v.validity() == "" {
			return validation.NewErrRecordIsMissingRequiredField("expiresOn")
		}
		if f.ValidTill.Exclude && v.validity() != "" {
			return validation.NewErrBadRecordFieldValue("expiresOn", "is not allowed for docType "+string(v.DocType))
		}
	}
	if f.IssuedOn != nil && f.IssuedOn.Required && v.IssuedOn == "" {
		return validation.NewErrRecordIsMissingRequiredField("issuedOn")
	}
	return nil
}
