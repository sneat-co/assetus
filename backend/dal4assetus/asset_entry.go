package dal4assetus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// AssetEntry is a loaded asset record with its key/ID.
type AssetEntry = record.DataWithID[string, *dbo4assetus.AssetDbo]

// NewAssetKey builds the dalgo key for an asset:
// /spaces/{spaceID}/ext/assetus/assets/{assetID}.
func NewAssetKey(spaceID coretypes.SpaceID, assetID string) *dal.Key {
	spaceModuleKey := dbo4spaceus.NewSpaceModuleKey(spaceID, const4assetus.ExtensionID)
	return dal.NewKeyWithParentAndID(spaceModuleKey, dbo4assetus.AssetsCollection, assetID)
}

// NewAssetEntry builds an empty AssetEntry addressed by (spaceID, assetID).
func NewAssetEntry(spaceID coretypes.SpaceID, assetID string) (asset AssetEntry) {
	key := NewAssetKey(spaceID, assetID)
	asset.ID = assetID
	asset.Key = key
	asset.Data = new(dbo4assetus.AssetDbo)
	asset.Record = dal.NewRecordWithData(key, asset.Data)
	return
}
