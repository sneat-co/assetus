package dto4assetus

import (
	"strings"
	"time"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
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
