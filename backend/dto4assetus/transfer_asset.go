package dto4assetus

import (
	"strings"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/strongo/validation"
)

// TransferAssetRequest transfers an asset from its owning Space (SpaceID) to a
// destination Space (ToSpaceID). The acting user must be a member of the owning
// Space.
type TransferAssetRequest struct {
	dto4spaceus.SpaceRequest
	AssetID   string            `json:"assetID"`
	ToSpaceID coretypes.SpaceID `json:"toSpaceID"`
}

// Validate validates the request.
func (v TransferAssetRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.AssetID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("assetID")
	}
	if strings.TrimSpace(string(v.ToSpaceID)) == "" {
		return validation.NewErrRequestIsMissingRequiredField("toSpaceID")
	}
	if v.ToSpaceID == v.SpaceID {
		return validation.NewErrBadRequestFieldValue("toSpaceID", "destination space must differ from the owning space")
	}
	return nil
}

// TransferAssetResponse returns the relocated asset and its new owner.
type TransferAssetResponse struct {
	ID    string      `json:"id"`
	Owner OwnerRefDTO `json:"owner"`
}
