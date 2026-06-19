package facade4assetus

import (
	"fmt"

	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/assetus/backend/dal4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/decimal"
)

// AddVehicleRecord appends a vehicle record (a mileage and/or fuel reading) to a
// vehicle asset's vehicleRecords child collection. Membership of the owning
// Space is enforced by the underlying asset worker. The record is INSERT-only.
// Ported from the legacy facade4assetus.AddVehicleRecord (which was a TODO stub)
// and mapped onto the unified dbo4assetus.VehicleRecordDbo carrying fuel volume/
// unit/amount + fuelCost/currency and mileage value/unit.
func AddVehicleRecord(ctx facade.ContextWithUser, request dto4assetus.AddVehicleRecordRequest) (response dto4assetus.AddVehicleRecordResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4assetus.RunAssetWorker(ctx, request.SpaceID, request.AssetID,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4assetus.AssetWorkerParams) (err error) {
			if !params.Asset.Record.Exists() {
				return fmt.Errorf("%w: asset %q not found in space %q", dal.ErrRecordNotFound, request.AssetID, request.SpaceID)
			}
			now := params.Started
			userID := params.UserID()

			dbo := vehicleRecordFromRequest(request)
			dbo.CreatedAt = now
			dbo.CreatedBy = userID

			recordID := newVehicleRecordID()
			if err = dal4assetus.AppendVehicleRecord(ctx, tx, request.SpaceID, request.AssetID, recordID, dbo); err != nil {
				return fmt.Errorf("failed to append vehicle record: %w", err)
			}
			response.ID = recordID
			return nil
		})
	if err != nil {
		return dto4assetus.AddVehicleRecordResponse{}, err
	}
	return
}

// vehicleRecordFromRequest maps the fuel-bearing request payload onto the
// unified VehicleRecordDbo. The fuel portion is set only when any fuel field is
// supplied; the mileage portion only when a mileage value is supplied.
func vehicleRecordFromRequest(request dto4assetus.AddVehicleRecordRequest) *dbo4assetus.VehicleRecordDbo {
	dbo := new(dbo4assetus.VehicleRecordDbo)

	if request.FuelVolume != 0 || request.FuelCost != 0 || request.FuelVolumeUnit != "" || request.Currency != "" {
		fuel := &dbo4assetus.VehicleFuelRecord{
			Volume:   decimal.NewDecimal64p2FromFloat64(float64(request.FuelVolume)),
			Unit:     request.FuelVolumeUnit,
			FuelCost: decimal.NewDecimal64p2FromFloat64(float64(request.FuelCost)),
			Currency: request.Currency,
		}
		if request.Currency != "" {
			amount := money.NewAmount(money.CurrencyCode(request.Currency), fuel.FuelCost)
			fuel.Amount = &amount
		}
		dbo.Fuel = fuel
	}

	if request.Mileage != 0 || request.MileageUnit != "" {
		dbo.Mileage = &dbo4assetus.VehicleMileage{
			Value: int(request.Mileage),
			Unit:  request.MileageUnit,
		}
	}

	return dbo
}
