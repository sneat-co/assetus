package dto4assetus

import (
	"strings"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/strongo/validation"
)

// RemoveAssetRequest removes an asset. The default action is a soft-archive
// (status -> Archived, record and history preserved). HardDelete requests the
// explicit, permanent deletion of the asset record and its history.
type RemoveAssetRequest struct {
	dto4spaceus.SpaceRequest
	AssetID    string `json:"assetID"`
	HardDelete bool   `json:"hardDelete,omitempty"`
}

// Validate validates the request.
func (v RemoveAssetRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.AssetID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("assetID")
	}
	return nil
}
