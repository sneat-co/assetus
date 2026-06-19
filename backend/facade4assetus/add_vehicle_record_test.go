package facade4assetus

import (
	"context"
	"errors"
	"testing"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dal4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/decimal"
)

func seedVehicleAsset(t *testing.T, spaceID coretypes.SpaceID) string {
	t.Helper()
	created, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		Name:         "Family Car",
		Category:     const4assetus.CategoryVehicles,
		Condition:    const4assetus.ConditionGood,
	})
	if err != nil {
		t.Fatalf("seed CreateAsset failed: %v", err)
	}
	return created.ID
}

// AddVehicleRecord appends a record with mileage + fuel that round-trips through
// the facade into the asset's vehicleRecords child collection.
func TestAddVehicleRecord_RoundTrip(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	db := newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)
	assetID := seedVehicleAsset(t, spaceID)

	resp, err := AddVehicleRecord(userCtx(testUserID), dto4assetus.AddVehicleRecordRequest{
		SpaceRequest:   spaceRequest(spaceID),
		AssetID:        assetID,
		FuelVolume:     42.50,
		FuelVolumeUnit: "litre",
		FuelCost:       73.20,
		Currency:       "EUR",
		Mileage:        123456,
		MileageUnit:    "km",
	})
	if err != nil {
		t.Fatalf("AddVehicleRecord failed: %v", err)
	}
	if resp.ID == "" {
		t.Fatal("expected a non-empty vehicle record ID")
	}

	// Re-read the persisted record from the vehicleRecords child collection.
	rec := new(dbo4assetus.VehicleRecordDbo)
	key := dal4assetus.NewVehicleRecordKey(spaceID, assetID, resp.ID)
	record := dal.NewRecordWithData(key, rec)
	if err = db.Get(context.Background(), record); err != nil {
		t.Fatalf("failed to read persisted vehicle record: %v", err)
	}

	if rec.Mileage == nil {
		t.Fatal("expected mileage portion to be set")
	}
	if rec.Mileage.Value != 123456 {
		t.Errorf("mileage value = %d, want 123456", rec.Mileage.Value)
	}
	if rec.Mileage.Unit != "km" {
		t.Errorf("mileage unit = %q, want km", rec.Mileage.Unit)
	}
	if rec.Fuel == nil {
		t.Fatal("expected fuel portion to be set")
	}
	if rec.Fuel.Volume != decimal.NewDecimal64p2FromFloat64(42.50) {
		t.Errorf("fuel volume = %v, want 42.50", rec.Fuel.Volume)
	}
	if rec.Fuel.Unit != "litre" {
		t.Errorf("fuel unit = %q, want litre", rec.Fuel.Unit)
	}
	if rec.Fuel.FuelCost != decimal.NewDecimal64p2FromFloat64(73.20) {
		t.Errorf("fuel cost = %v, want 73.20", rec.Fuel.FuelCost)
	}
	if rec.Fuel.Currency != "EUR" {
		t.Errorf("fuel currency = %q, want EUR", rec.Fuel.Currency)
	}
	if rec.Fuel.Amount == nil {
		t.Fatal("expected fuel amount to be set when currency supplied")
	}
	if string(rec.Fuel.Amount.Currency) != "EUR" {
		t.Errorf("amount currency = %q, want EUR", rec.Fuel.Amount.Currency)
	}
	if rec.CreatedBy != testUserID {
		t.Errorf("createdBy = %q, want %q", rec.CreatedBy, testUserID)
	}
}

// A non-member cannot add a vehicle record.
func TestAddVehicleRecord_NonMemberRejected(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	_ = newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)
	assetID := seedVehicleAsset(t, spaceID)

	_, err := AddVehicleRecord(userCtx("intruder"), dto4assetus.AddVehicleRecordRequest{
		SpaceRequest: spaceRequest(spaceID),
		AssetID:      assetID,
		Mileage:      100,
		MileageUnit:  "km",
	})
	if err == nil {
		t.Fatal("expected non-member to be rejected")
	}
	if !errors.Is(err, facade.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got: %v", err)
	}
}

// Adding a record to a non-existent asset is rejected.
func TestAddVehicleRecord_AssetNotFound(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	_ = newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)

	_, err := AddVehicleRecord(userCtx(testUserID), dto4assetus.AddVehicleRecordRequest{
		SpaceRequest: spaceRequest(spaceID),
		AssetID:      "missing",
		Mileage:      100,
		MileageUnit:  "km",
	})
	if err == nil {
		t.Fatal("expected not-found rejection")
	}
	if !errors.Is(err, dal.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound, got: %v", err)
	}
}
