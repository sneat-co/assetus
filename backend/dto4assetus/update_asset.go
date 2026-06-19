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

// UpdateAssetRequest updates an asset's editable fields (PUT semantics over the
// editable set). The ownership-lifecycle status is NOT editable here — it is
// changed only by the remove (soft-archive) and transfer operations. Editing
// MUST NOT alter the asset's history.
type UpdateAssetRequest struct {
	dto4spaceus.SpaceRequest
	AssetID     string                   `json:"assetID"`
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	Category    const4assetus.Category   `json:"category"`
	Condition   const4assetus.Condition  `json:"condition"`
	Visibility  const4assetus.Visibility `json:"visibility"`

	// Optional metadata (full replacement of the editable metadata set).
	AcquisitionDate *time.Time                  `json:"acquisitionDate,omitempty"`
	Location        string                      `json:"location,omitempty"`
	Notes           string                      `json:"notes,omitempty"`
	Tags            []string                    `json:"tags,omitempty"`
	PurchasePrice   *dbo4assetus.MonetaryAmount `json:"purchasePrice,omitempty"`
	EstimatedValue  *dbo4assetus.MonetaryAmount `json:"estimatedValue,omitempty"`

	// --- Optional unified (superset) editable fields, full-replacement ----
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

// Validate validates the request.
func (v UpdateAssetRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.AssetID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("assetID")
	}
	if strings.TrimSpace(v.Name) == "" {
		return validation.NewErrRequestIsMissingRequiredField("name")
	}
	if err := const4assetus.ValidateCategory(v.Category); err != nil {
		return err
	}
	if err := const4assetus.ValidateCondition(v.Condition); err != nil {
		return err
	}
	if err := const4assetus.ValidateVisibility(v.Visibility); err != nil {
		return err
	}
	return nil
}

// UpdateAssetResponse returns the updated asset.
type UpdateAssetResponse struct {
	ID    string                `json:"id"`
	Asset *dbo4assetus.AssetDbo `json:"asset"`
}
