package dto4assetus

import (
	"strings"
	"time"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/strongo/validation"
)

// CreateAssetRequest is the request to create an asset in a Space. The owner is
// the Space; status defaults to Active; visibility defaults to the Space's
// default visibility unless an override is supplied.
type CreateAssetRequest struct {
	dto4spaceus.SpaceRequest
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	Category    const4assetus.Category   `json:"category"`
	Condition   const4assetus.Condition  `json:"condition"`
	Visibility  const4assetus.Visibility `json:"visibility,omitempty"` // optional per-asset override

	// Optional metadata.
	AcquisitionDate *time.Time                  `json:"acquisitionDate,omitempty"`
	PurchasePrice   *dbo4assetus.MonetaryAmount `json:"purchasePrice,omitempty"`
	EstimatedValue  *dbo4assetus.MonetaryAmount `json:"estimatedValue,omitempty"`
	Location        string                      `json:"location,omitempty"`
	Notes           string                      `json:"notes,omitempty"`
	Tags            []string                    `json:"tags,omitempty"`
}

// Validate validates the request. Status is not accepted from the caller (it
// always defaults to Active on create). Visibility is optional (inherited when
// empty) but, when supplied, must be a valid value.
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
	if err := const4assetus.ValidateCondition(v.Condition); err != nil {
		return err
	}
	if v.Visibility != "" {
		if err := const4assetus.ValidateVisibility(v.Visibility); err != nil {
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
