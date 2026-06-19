package facade4assetus

import (
	"context"
	"errors"
	"testing"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dal4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/assetus/backend/extras4assetus"
	"github.com/sneat-co/sneat-core-modules/core/extra"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

func seedAsset(t *testing.T, spaceID coretypes.SpaceID) string {
	t.Helper()
	created, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		Name:         "Original Name",
		Category:     const4assetus.CategoryBooks,
		Condition:    const4assetus.ConditionGood,
	})
	if err != nil {
		t.Fatalf("seed CreateAsset failed: %v", err)
	}
	return created.ID
}

// AC: member-edits-asset — name & category update; history is unchanged.
func TestUpdateAsset_MemberEdits_HistoryUnchanged(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	_ = newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)
	id := seedAsset(t, spaceID)

	before, err := GetHistory(userCtx(testUserID), dto4assetus.GetHistoryRequest{SpaceRequest: spaceRequest(spaceID), AssetID: id})
	if err != nil {
		t.Fatalf("GetHistory before failed: %v", err)
	}

	resp, err := UpdateAsset(userCtx(testUserID), dto4assetus.UpdateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		AssetID:      id,
		Name:         "Updated Name",
		Category:     const4assetus.CategoryElectronics,
		Condition:    const4assetus.ConditionGood,
		Visibility:   const4assetus.VisibilityFamily,
	})
	if err != nil {
		t.Fatalf("UpdateAsset failed: %v", err)
	}
	if resp.Asset.Name != "Updated Name" || resp.Asset.Category != const4assetus.CategoryElectronics {
		t.Errorf("update not applied: name=%q category=%q", resp.Asset.Name, resp.Asset.Category)
	}

	// AC: history-is-append-only (after a metadata edit) — history unchanged.
	after, err := GetHistory(userCtx(testUserID), dto4assetus.GetHistoryRequest{SpaceRequest: spaceRequest(spaceID), AssetID: id})
	if err != nil {
		t.Fatalf("GetHistory after failed: %v", err)
	}
	if len(after.Events) != len(before.Events) {
		t.Errorf("history length changed by edit: before=%d after=%d", len(before.Events), len(after.Events))
	}
	if len(after.Events) > 0 && after.Events[0].ID != before.Events[0].ID {
		t.Errorf("first history event mutated by edit")
	}
}

// AC: condition-can-be-updated — Good -> Needs Repair persists.
func TestUpdateAsset_ConditionUpdated(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	_ = newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)
	id := seedAsset(t, spaceID)

	resp, err := UpdateAsset(userCtx(testUserID), dto4assetus.UpdateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		AssetID:      id,
		Name:         "Original Name",
		Category:     const4assetus.CategoryBooks,
		Condition:    const4assetus.ConditionNeedsRepair,
		Visibility:   const4assetus.VisibilityFamily,
	})
	if err != nil {
		t.Fatalf("UpdateAsset failed: %v", err)
	}
	if resp.Asset.Condition != const4assetus.ConditionNeedsRepair {
		t.Errorf("condition = %q, want needs_repair", resp.Asset.Condition)
	}
}

// AC: optional-legacy-fields-roundtrip (update) — an update sets/replaces the
// rich editable unified fields (incl. a vehicle extra); they persist, the
// ownership-lifecycle status is preserved, and the history is NOT altered.
func TestUpdateAsset_RichFieldsReplacedHistoryUnchanged(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	db := newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)
	id := seedAsset(t, spaceID)

	before, err := GetHistory(userCtx(testUserID), dto4assetus.GetHistoryRequest{SpaceRequest: spaceRequest(spaceID), AssetID: id})
	if err != nil {
		t.Fatalf("GetHistory before failed: %v", err)
	}

	var ef extra.WithExtraField
	if err = ef.SetExtra(extras4assetus.AssetExtraTypeVehicle, &extras4assetus.AssetVehicleExtra{
		WithMakeModelRegNumberFields: extras4assetus.WithMakeModelRegNumberFields{
			WithMakeModelFields: extras4assetus.WithMakeModelFields{Make: "Tesla", Model: "Model 3"},
		},
		Vin: "1HGCM82633A004352",
		WithEngineData: extras4assetus.WithEngineData{
			EngineType: const4assetus.EngineTypeElectric,
		},
	}); err != nil {
		t.Fatalf("SetExtra failed: %v", err)
	}

	year := 2021
	resp, err := UpdateAsset(userCtx(testUserID), dto4assetus.UpdateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		AssetID:      id,
		Name:         "Updated Car",
		Category:     const4assetus.CategoryVehicles,
		Condition:    const4assetus.ConditionGood,
		Visibility:   const4assetus.VisibilityFamily,
		Type:         const4assetus.TypeVehicleCar,
		Possession:   const4assetus.PossessionRenting,
		CountryID:    "IE",
		YearOfBuild:  &year,
		Geo:          &dbo4assetus.GeoPoint{Lat: 1, Lng: 2},
		WithAssetRelationships: dbo4assetus.WithAssetRelationships{
			SameAssetID: "same-1",
		},
		WithExtraField: ef,
	})
	if err != nil {
		t.Fatalf("UpdateAsset failed: %v", err)
	}
	if resp.Asset.Type != const4assetus.TypeVehicleCar || resp.Asset.Possession != const4assetus.PossessionRenting {
		t.Errorf("rich fields not applied: type=%q possession=%q", resp.Asset.Type, resp.Asset.Possession)
	}
	// Status preserved from the seeded asset (active).
	if resp.Asset.Status != const4assetus.StatusActive {
		t.Errorf("status = %q, want active (preserved)", resp.Asset.Status)
	}

	// Re-read to confirm persistence and round-trip of the extra.
	asset := dal4assetus.NewAssetEntry(spaceID, id)
	if err = db.Get(context.Background(), asset.Record); err != nil {
		t.Fatalf("failed to read updated asset: %v", err)
	}
	if asset.Data.SameAssetID != "same-1" {
		t.Errorf("sameAssetID = %q, want same-1", asset.Data.SameAssetID)
	}
	if asset.Data.ExtraType != extras4assetus.AssetExtraTypeVehicle {
		t.Errorf("extraType = %q, want vehicle", asset.Data.ExtraType)
	}

	// History unchanged by the edit.
	after, err := GetHistory(userCtx(testUserID), dto4assetus.GetHistoryRequest{SpaceRequest: spaceRequest(spaceID), AssetID: id})
	if err != nil {
		t.Fatalf("GetHistory after failed: %v", err)
	}
	if len(after.Events) != len(before.Events) {
		t.Errorf("history length changed by edit: before=%d after=%d", len(before.Events), len(after.Events))
	}
}

// AC: non-member-cannot-edit — a non-member is rejected; asset unchanged.
func TestUpdateAsset_NonMemberRejected(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	_ = newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)
	id := seedAsset(t, spaceID)

	_, err := UpdateAsset(userCtx("intruder"), dto4assetus.UpdateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		AssetID:      id,
		Name:         "Hacked",
		Category:     const4assetus.CategoryOther,
		Condition:    const4assetus.ConditionBroken,
		Visibility:   const4assetus.VisibilityPublic,
	})
	if err == nil {
		t.Fatal("expected non-member edit to be rejected")
	}
	if !errors.Is(err, facade.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
	// Asset must be unchanged.
	got, err := GetAsset(userCtx(testUserID), dto4assetus.GetAssetRequest{SpaceRequest: spaceRequest(spaceID), AssetID: id})
	if err != nil {
		t.Fatalf("GetAsset failed: %v", err)
	}
	if got.Asset.Name != "Original Name" {
		t.Errorf("asset was modified by rejected edit: name=%q", got.Asset.Name)
	}
}
