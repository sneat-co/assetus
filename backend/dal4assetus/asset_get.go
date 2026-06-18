package dal4assetus

import (
	"context"

	"github.com/dal-go/dalgo/dal"
)

// GetAssetByID loads an asset record.
func GetAssetByID(ctx context.Context, getter dal.ReadSession, asset AssetEntry) error {
	return getter.Get(ctx, asset.Record)
}

// GetAssetForUpdate loads an asset record inside a read-write transaction.
func GetAssetForUpdate(ctx context.Context, tx dal.ReadwriteTransaction, asset AssetEntry) error {
	return GetAssetByID(ctx, tx, asset)
}
