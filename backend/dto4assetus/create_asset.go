package dto4assetus

import (
	"strings"
	"time"

	"github.com/crediterra/money"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/sneat-core-modules/core/extra"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/geo"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

// CreateAssetRequest is the request to create an asset in a Space. The owner is
// the Space; status defaults to Active (or an optionally supplied draft);
// visibility defaults to the Space's default visibility unless an override is
// supplied.
type CreateAssetRequest struct {
	dto4spaceus.SpaceRequest
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	Category    const4assetus.Category   `json:"category"`
	Condition   const4assetus.Condition  `json:"condition"`
	Visibility  const4assetus.Visibility `json:"visibility,omitempty"` // optional per-asset override

	// Status is OPTIONAL on create: when empty it defaults to active; when
	// supplied it must be a valid ownership status (e.g. draft).
	Status const4assetus.Status `json:"status,omitempty"`

	// Optional metadata.
	AcquisitionDate *time.Time                  `json:"acquisitionDate,omitempty"`
	PurchasePrice   *dbo4assetus.MonetaryAmount `json:"purchasePrice,omitempty"`
	EstimatedValue  *dbo4assetus.MonetaryAmount `json:"estimatedValue,omitempty"`
	Location        string                      `json:"location,omitempty"`
	Notes           string                      `json:"notes,omitempty"`
	Tags            []string                    `json:"tags,omitempty"`

	// --- Optional unified (superset) fields, all settable on create -------
	Type             const4assetus.Type       `json:"type,omitempty"`
	Possession       const4assetus.Possession `json:"possession,omitempty"`
	CountryID        geo.CountryAlpha2        `json:"countryID,omitempty"`
	ParentCategoryID const4assetus.Category   `json:"parentCategoryID,omitempty"`
	YearOfBuild      *int                     `json:"yearOfBuild,omitempty"`
	IsRequest        bool                     `json:"isRequest,omitempty"`
	Geo              *dbo4assetus.GeoPoint    `json:"geo,omitempty"`

	// AssetDates (dateOfBuild/dateOfPurchase/dateInsuredTill/dateCertifiedTill).
	dbo4assetus.AssetDates

	// Custom fields (fieldsStr/fieldsInt/fieldsDate/fieldsAmount).
	dbmodels.WithCustomFields

	// Financial.
	Totals              []money.Amount                     `json:"totals,omitempty"`
	CanHaveIncome       bool                               `json:"canHaveIncome,omitempty"`
	CanHaveExpense      bool                               `json:"canHaveExpense,omitempty"`
	FinancialDirection  string                             `json:"financialDirection,omitempty"`
	Liabilities         []dbo4assetus.AssetLiabilityInfo   `json:"liabilities,omitempty"`
	NotUsedServiceTypes []dbo4assetus.LiabilityServiceType `json:"notUsedServiceTypes,omitempty"`

	// Relationships (groupID/group/parentAssetID/subAssets/sameAssetID/
	// relatedAs/memberIDs/membersInfo).
	dbo4assetus.WithAssetRelationships

	// The OPTIONAL polymorphic typed extra (vehicle/dwelling/document),
	// accepted via the same extraType + extra shape carried by the record.
	extra.WithExtraField
}

// Validate validates the request. Status is OPTIONAL (defaults to active on
// create); when supplied it must be a valid ownership status. Visibility is
// optional (inherited when empty) but, when supplied, must be a valid value.
// All added unified fields are validated only when present, reusing the same
// validators applied by the persisted record.
func (v CreateAssetRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.Name) == "" {
		return validation.NewErrRequestIsMissingRequiredField("name")
	}
	if err := const4assetus.ValidateCategory(v.Category); err != nil {
		return err
	}
	if err := const4assetus.ValidateConditionOptional(v.Condition); err != nil {
		return err
	}
	if v.Visibility != "" {
		if err := const4assetus.ValidateVisibility(v.Visibility); err != nil {
			return err
		}
	}
	if v.Status != "" {
		if err := const4assetus.ValidateStatus(v.Status); err != nil {
			return err
		}
	}
	return nil
}

// CreateAssetResponse is returned on successful asset creation.
type CreateAssetResponse struct {
	ID    string                `json:"id"`
	Asset *dbo4assetus.AssetDbo `json:"asset"`
}
