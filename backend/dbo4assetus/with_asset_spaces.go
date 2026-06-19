package dbo4assetus

import (
	"github.com/strongo/validation"
)

// WithAssetSpaces is the multi-space association ported from the legacy
// briefs4assetus/dbo4assetus. It maps spaceID -> a per-space set of asset
// briefs, so a single asset record can be associated with multiple spaces.
//
// This is ADDITIVE to the MVP single owning space (carried via
// dbmodels.WithSpaceIDs on AssetDbo): it does not replace the owning space.
// Task 4 will fully wire multi-space + relationships; this lands the field/struct
// as a seam.
type WithAssetSpaces struct {
	Spaces map[string]*AssetusSpaceBrief `json:"spaces,omitempty" firestore:"spaces,omitempty"`
}

// AssetusSpaceBrief is the per-space projection of an asset's briefs.
type AssetusSpaceBrief struct {
	Assets AssetBriefs `json:"assets,omitempty" firestore:"assets,omitempty"`
}

// Validate returns an error if not valid.
func (v *WithAssetSpaces) Validate() error {
	for id, spaceBrief := range v.Spaces {
		if id == "" {
			return validation.NewErrBadRecordFieldValue("spaces", "spaceID can not be empty string")
		}
		if spaceBrief == nil {
			return validation.NewErrBadRecordFieldValue("spaces."+id, "can not be nil")
		}
		for assetID, brief := range spaceBrief.Assets {
			if assetID == "" {
				return validation.NewErrBadRecordFieldValue("spaces."+id+".assets", "assetID can not be empty string")
			}
			if brief == nil {
				return validation.NewErrBadRecordFieldValue("spaces."+id+".assets."+assetID, "can not be nil")
			}
			if err := brief.Validate(); err != nil {
				return validation.NewErrBadRecordFieldValue("spaces."+id+".assets."+assetID, err.Error())
			}
		}
	}
	return nil
}
