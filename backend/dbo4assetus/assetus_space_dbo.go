package dbo4assetus

import (
	"fmt"

	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
)

// AssetBriefs maps assetID -> brief, denormalized onto the module space entry
// for cheap listing without reading every asset document.
type AssetBriefs map[string]*AssetBrief

// AssetusSpaceDbo is the assetus module entry for a Space, persisted at
// /spaces/{spaceID}/ext/assetus.
type AssetusSpaceDbo struct {
	with.CreatedFields
	Assets AssetBriefs `json:"assets,omitempty" firestore:"assets,omitempty"`
}

// Validate returns an error if the module entry is not valid.
func (v *AssetusSpaceDbo) Validate() error {
	for id, brief := range v.Assets {
		if brief == nil {
			return validation.NewErrBadRecordFieldValue("assets", fmt.Sprintf("nil brief for assetID=%s", id))
		}
		if err := brief.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("assets", fmt.Sprintf("invalid brief for assetID=%s: %v", id, err))
		}
	}
	return nil
}
