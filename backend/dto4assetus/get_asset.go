package dto4assetus

import (
	"strings"

	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/strongo/validation"
)

// GetAssetRequest is the request to read a single asset.
type GetAssetRequest struct {
	dto4spaceus.SpaceRequest
	AssetID string `json:"assetID"`
}

// Validate validates the request.
func (v GetAssetRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.AssetID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("assetID")
	}
	return nil
}

// GetAssetResponse returns an asset together with its derived owner (Space +
// owner type derived from the Space type).
type GetAssetResponse struct {
	ID    string                `json:"id"`
	Asset *dbo4assetus.AssetDbo `json:"asset"`
	Owner dbo4assetus.OwnerRef  `json:"owner"`
}
