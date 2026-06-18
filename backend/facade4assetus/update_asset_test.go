package facade4assetus

import (
	"errors"
	"testing"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
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
