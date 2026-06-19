package dbo4assetus

import (
	"fmt"
	"strings"
	"time"

	"github.com/crediterra/money"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/sneat-go-core/geo"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

// AssetsCollection is the Firestore collection name for assets under the
// assetus space module: /spaces/{spaceID}/ext/assetus/assets/{assetID}.
const AssetsCollection = "assets"

// MonetaryAmount is an optional monetary metadata value (purchase price /
// estimated value). Reserved-simple: a currency code plus a decimal value.
type MonetaryAmount struct {
	Currency string  `json:"currency,omitempty" firestore:"currency,omitempty"`
	Value    float64 `json:"value,omitempty" firestore:"value,omitempty"`
}

// AssetDates carries the optional ISO-date (YYYY-MM-DD) fields ported from the
// legacy AssetDates. Each field is optional.
type AssetDates struct {
	DateOfBuild       string `json:"dateOfBuild,omitempty" firestore:"dateOfBuild,omitempty"`
	DateOfPurchase    string `json:"dateOfPurchase,omitempty" firestore:"dateOfPurchase,omitempty"`
	DateInsuredTill   string `json:"dateInsuredTill,omitempty" firestore:"dateInsuredTill,omitempty"`
	DateCertifiedTill string `json:"dateCertifiedTill,omitempty" firestore:"dateCertifiedTill,omitempty"`
}

// Validate returns an error if any present date is not a valid date string.
// All fields are optional: an empty AssetDates is valid.
func (v AssetDates) Validate() error {
	for name, value := range map[string]string{
		"dateOfBuild":       v.DateOfBuild,
		"dateOfPurchase":    v.DateOfPurchase,
		"dateInsuredTill":   v.DateInsuredTill,
		"dateCertifiedTill": v.DateCertifiedTill,
	} {
		if value == "" {
			continue
		}
		if _, err := validate.DateString(value); err != nil {
			return validation.NewErrBadRecordFieldValue(name, err.Error())
		}
	}
	return nil
}

// GeoPoint is an optional geo coordinate ported as the unified "geo value" for
// an asset (e.g. a dwelling location). Both fields default to zero and the whole
// field is optional.
type GeoPoint struct {
	Lat float64 `json:"lat" firestore:"lat"`
	Lng float64 `json:"lng" firestore:"lng"`
}

// LiabilityServiceType is a service type linked to an asset's liabilities,
// ported from the legacy LiabilityServiceType union (electricity, gas, nct…).
// Kept as an open string to preserve the legacy value set without re-declaring
// the full union here.
type LiabilityServiceType string

// AssetLiabilityInfo is the asset-side linkage to a liability record, ported
// from the legacy frontend AssetLiabilityInfo DTO.
type AssetLiabilityInfo struct {
	ID           string                 `json:"id" firestore:"id"`
	ServiceTypes []LiabilityServiceType `json:"serviceTypes,omitempty" firestore:"serviceTypes,omitempty"`
}

// Validate returns an error if the liability linkage is missing its required ID.
func (v AssetLiabilityInfo) Validate() error {
	if strings.TrimSpace(v.ID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	return nil
}

// AssetBase carries the user-editable fields of an asset. It is embedded by
// AssetDbo (the persisted record) and reused by the create/edit facades.
type AssetBase struct {
	Name        string                   `json:"name" firestore:"name"`
	Description string                   `json:"description,omitempty" firestore:"description,omitempty"`
	Category    const4assetus.Category   `json:"category" firestore:"category"`
	Condition   const4assetus.Condition  `json:"condition" firestore:"condition"`
	Status      const4assetus.Status     `json:"status" firestore:"status"`
	Visibility  const4assetus.Visibility `json:"visibility" firestore:"visibility"`

	// Optional metadata.
	AcquisitionDate *time.Time      `json:"acquisitionDate,omitempty" firestore:"acquisitionDate,omitempty"`
	PurchasePrice   *MonetaryAmount `json:"purchasePrice,omitempty" firestore:"purchasePrice,omitempty"`
	EstimatedValue  *MonetaryAmount `json:"estimatedValue,omitempty" firestore:"estimatedValue,omitempty"`
	Location        string          `json:"location,omitempty" firestore:"location,omitempty"`
	Notes           string          `json:"notes,omitempty" firestore:"notes,omitempty"`
	Tags            []string        `json:"tags,omitempty" firestore:"tags,omitempty"`

	// Photos is reserved for a fast-follow and is NOT implemented in the MVP.
	// It is declared so the schema does not need reshaping when photo
	// upload/storage lands.
	Photos []string `json:"photos,omitempty" firestore:"photos,omitempty"`

	// --- Ported legacy fields (all OPTIONAL) -----------------------------
	// These fields are merged in from the legacy AssetBrief/AssetBaseDbo and
	// frontend DTO. Every one is optional (omitempty); an MVP-shaped asset that
	// supplies none of them is still valid.

	// CountryID is the optional ISO-3166 alpha-2 country code of the asset.
	// In the legacy model this was a required, query-only field; in the unified
	// model it is optional.
	CountryID geo.CountryAlpha2 `json:"countryID,omitempty" firestore:"countryID,omitempty"`

	// Type is the optional subtype within the asset's Category (e.g. car,
	// apartment, passport). Validated against Category via const4assetus.
	Type const4assetus.Type `json:"type,omitempty" firestore:"type,omitempty"`

	// Possession describes how the asset is possessed (owning/leasing/renting…).
	// Optional; when unset it defaults to owning — see WithPossessionDefault.
	Possession const4assetus.Possession `json:"possession,omitempty" firestore:"possession,omitempty"`

	// ParentCategoryID is the optional parent category of the asset.
	ParentCategoryID const4assetus.Category `json:"parentCategoryID,omitempty" firestore:"parentCategoryID,omitempty"`

	// YearOfBuild is the optional year an asset was built/manufactured.
	YearOfBuild *int `json:"yearOfBuild,omitempty" firestore:"yearOfBuild,omitempty"`

	// IsRequest marks the asset as a request rather than an owned record.
	IsRequest bool `json:"isRequest,omitempty" firestore:"isRequest,omitempty"`

	// Geo is the optional geo coordinate of the asset.
	Geo *GeoPoint `json:"geo,omitempty" firestore:"geo,omitempty"`

	// AssetDates carries the optional build/purchase/insured/certified dates.
	AssetDates

	// WithCustomFields carries the optional per-asset custom fields
	// (fieldsStr/fieldsInt/fieldsDate/fieldsAmount) ported from the legacy model.
	dbmodels.WithCustomFields

	// --- Ported financial fields (all OPTIONAL) --------------------------
	// OWNER DECISION (Task 2): the legacy financial dimension lives as OPTIONAL
	// fields on the core asset, NOT a separate module. Disposition for the
	// capability-coverage table:
	//   - per-asset/group totals          -> Totals (optional)
	//   - canHaveIncome / canHaveExpense  -> CanHaveIncome / CanHaveExpense
	//   - income/expense direction        -> FinancialDirection
	//   - debt category                   -> const4assetus.CategoryDebt (Category)
	//   - asset-side liability linkage    -> Liabilities / NotUsedServiceTypes
	//                                        ([]AssetLiabilityInfo)

	// Totals holds the optional per-asset financial totals (ITotalsHolder).
	Totals []money.Amount `json:"totals,omitempty" firestore:"totals,omitempty"`

	// CanHaveIncome / CanHaveExpense are the optional income/expense capability
	// flags ported from the legacy IAssetCategory.
	CanHaveIncome  bool `json:"canHaveIncome,omitempty" firestore:"canHaveIncome,omitempty"`
	CanHaveExpense bool `json:"canHaveExpense,omitempty" firestore:"canHaveExpense,omitempty"`

	// FinancialDirection is the optional income/expense direction of the asset's
	// financial flow ("income" or "expense").
	FinancialDirection string `json:"financialDirection,omitempty" firestore:"financialDirection,omitempty"`

	// Liabilities is the optional asset-side linkage to liability records.
	Liabilities []AssetLiabilityInfo `json:"liabilities,omitempty" firestore:"liabilities,omitempty"`

	// NotUsedServiceTypes lists liability service types explicitly not used by
	// this asset, ported from the legacy notUsedServiceTypes.
	NotUsedServiceTypes []LiabilityServiceType `json:"notUsedServiceTypes,omitempty" firestore:"notUsedServiceTypes,omitempty"`

	// --- Ported relationship fields (all OPTIONAL) -----------------------
	// WithAssetRelationships carries the optional group membership, parent/
	// sub-asset nesting, asset linking (sameAssetID/relatedAs) and per-asset
	// member info ported from the legacy AssetBrief/AssetDbo and frontend DTO.
	WithAssetRelationships
}

// WithPossessionDefault returns the asset's Possession, defaulting to
// const4assetus.PossessionOwning when it is unset. The MVP create/normalize
// path should call this so an asset created without a possession value resolves
// to owning. (Task 5 wires this into the facade; this helper is the seam.)
func (v AssetBase) WithPossessionDefault() const4assetus.Possession {
	if v.Possession == "" {
		return const4assetus.PossessionOwning
	}
	return v.Possession
}

// Validate enforces the closed-enum write boundary plus required fields. This is
// the single enforcement point reused by the create and update facades — out-of-set
// values (e.g. status "borrowed"/"reserved") are rejected before persistence.
func (v AssetBase) Validate() error {
	if strings.TrimSpace(v.Name) == "" {
		return validation.NewErrRecordIsMissingRequiredField("name")
	}
	if err := const4assetus.ValidateCategory(v.Category); err != nil {
		return err
	}
	if err := const4assetus.ValidateCondition(v.Condition); err != nil {
		return err
	}
	if err := const4assetus.ValidateStatus(v.Status); err != nil {
		return err
	}
	if err := const4assetus.ValidateVisibility(v.Visibility); err != nil {
		return err
	}
	for i, tag := range v.Tags {
		if strings.TrimSpace(tag) == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("tags[%d]", i), "tag must not be empty")
		}
	}

	// Optional ported fields: validated only when present.
	if v.CountryID != "" && !geo.IsValidCountryAlpha2(v.CountryID) {
		return validation.NewErrBadRecordFieldValue("countryID",
			fmt.Sprintf("invalid country alpha-2 code: %q", v.CountryID))
	}
	if err := const4assetus.ValidateType(v.Category, v.Type); err != nil {
		return err
	}
	if v.Possession != "" {
		if err := const4assetus.ValidatePossession(v.Possession, false); err != nil {
			return err
		}
	}
	if v.ParentCategoryID != "" {
		if err := const4assetus.ValidateCategory(v.ParentCategoryID); err != nil {
			return err
		}
	}
	switch v.FinancialDirection {
	case "", "income", "expense":
		// OK
	default:
		return validation.NewErrBadRecordFieldValue("financialDirection",
			fmt.Sprintf("expected income or expense, got %q", v.FinancialDirection))
	}
	if err := v.AssetDates.Validate(); err != nil {
		return err
	}
	if err := v.WithCustomFields.Validate(); err != nil {
		return err
	}
	for i, liability := range v.Liabilities {
		if err := liability.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("liabilities[%d]", i), err.Error())
		}
	}
	for i, amount := range v.Totals {
		if err := amount.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("totals[%d]", i), err.Error())
		}
	}
	if err := v.WithAssetRelationships.Validate(); err != nil {
		return err
	}
	return nil
}

// AssetDbo is the persisted asset aggregate at
// /spaces/{spaceID}/ext/assetus/assets/{assetID}.
//
// The owner is the owning Space (carried via WithSpaceIDs); there is no separate
// owner entity and no sharing/availability fields. ext.yardius is intentionally
// absent in the MVP.
// The canonical owning Space — the anchor for lifecycle/history/transfer and
// owner-derivation — is carried via WithSpaceIDs. WithAssetSpaces is an ADDITIVE
// multi-space association: the asset can surface under several spaces while the
// owning Space stays the single authoritative owner.
type AssetDbo struct {
	AssetBase
	dbmodels.WithModified
	dbmodels.WithSpaceIDs
	WithAssetSpaces
}

// Validate returns an error if the record is not valid.
func (v *AssetDbo) Validate() error {
	if err := v.WithSpaceIDs.Validate(); err != nil {
		return err
	}
	if err := v.WithModified.Validate(); err != nil {
		return err
	}
	if err := v.AssetBase.Validate(); err != nil {
		return err
	}
	if err := v.WithAssetSpaces.Validate(); err != nil {
		return err
	}
	return nil
}

// AssetBrief is a denormalized summary of an asset for listing on the module
// space entry. It deliberately carries no sharing/availability fields.
type AssetBrief struct {
	ID         string                   `json:"id" firestore:"id"`
	Name       string                   `json:"name" firestore:"name"`
	Category   const4assetus.Category   `json:"category" firestore:"category"`
	Condition  const4assetus.Condition  `json:"condition" firestore:"condition"`
	Status     const4assetus.Status     `json:"status" firestore:"status"`
	Visibility const4assetus.Visibility `json:"visibility" firestore:"visibility"`
}

// Validate returns an error if the brief is not valid.
func (v AssetBrief) Validate() error {
	if strings.TrimSpace(v.ID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if strings.TrimSpace(v.Name) == "" {
		return validation.NewErrRecordIsMissingRequiredField("name")
	}
	if err := const4assetus.ValidateCategory(v.Category); err != nil {
		return err
	}
	if err := const4assetus.ValidateCondition(v.Condition); err != nil {
		return err
	}
	if err := const4assetus.ValidateStatus(v.Status); err != nil {
		return err
	}
	if err := const4assetus.ValidateVisibility(v.Visibility); err != nil {
		return err
	}
	return nil
}

// BriefFromAsset builds an AssetBrief from a persisted asset.
func BriefFromAsset(id string, asset *AssetDbo) *AssetBrief {
	return &AssetBrief{
		ID:         id,
		Name:       asset.Name,
		Category:   asset.Category,
		Condition:  asset.Condition,
		Status:     asset.Status,
		Visibility: asset.Visibility,
	}
}
