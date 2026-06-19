package dal4assetus

import (
	"context"
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// NewVehicleRecordKey builds the dalgo key for a vehicle record:
// /spaces/{spaceID}/ext/assetus/assets/{assetID}/vehicleRecords/{recordID}.
func NewVehicleRecordKey(spaceID coretypes.SpaceID, assetID, recordID string) *dal.Key {
	return dal.NewKeyWithParentAndID(NewAssetKey(spaceID, assetID), dbo4assetus.VehicleRecordsCollection, recordID)
}

// AppendVehicleRecord inserts a new vehicle record (mileage + fuel reading) into
// an asset's vehicleRecords child collection inside an existing transaction. It
// is INSERT-only.
func AppendVehicleRecord(ctx context.Context, tx dal.ReadwriteTransaction, spaceID coretypes.SpaceID, assetID, recordID string, dbo *dbo4assetus.VehicleRecordDbo) error {
	if err := dbo.Validate(); err != nil {
		return fmt.Errorf("invalid vehicle record: %w", err)
	}
	rec := dal.NewRecordWithData(NewVehicleRecordKey(spaceID, assetID, recordID), dbo)
	if err := tx.Insert(ctx, rec); err != nil {
		return fmt.Errorf("failed to insert vehicle record: %w", err)
	}
	return nil
}
