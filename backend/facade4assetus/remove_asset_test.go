package facade4assetus

import (
	"context"
	"errors"
	"testing"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dal4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// AC: soft-archive-preserves-history — status becomes Archived; record and all
// history remain retrievable.
func TestRemoveAsset_SoftArchivePreservesHistory(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	_ = newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)
	id := seedAsset(t, spaceID)

	if err := RemoveAsset(userCtx(testUserID), dto4assetus.RemoveAssetRequest{
		SpaceRequest: spaceRequest(spaceID), AssetID: id, // HardDelete defaults false
	}); err != nil {
		t.Fatalf("soft-archive failed: %v", err)
	}

	got, err := GetAsset(userCtx(testUserID), dto4assetus.GetAssetRequest{SpaceRequest: spaceRequest(spaceID), AssetID: id})
	if err != nil {
		t.Fatalf("asset not retrievable after soft-archive: %v", err)
	}
	if got.Asset.Status != const4assetus.StatusArchived {
		t.Errorf("status = %q, want archived", got.Asset.Status)
	}
	hist, err := GetHistory(userCtx(testUserID), dto4assetus.GetHistoryRequest{SpaceRequest: spaceRequest(spaceID), AssetID: id})
	if err != nil {
		t.Fatalf("history not retrievable after soft-archive: %v", err)
	}
	if len(hist.Events) == 0 {
		t.Error("history was lost on soft-archive, expected it preserved")
	}
}

// AC: hard-delete-removes-record — record and history are permanently removed.
func TestRemoveAsset_HardDeleteRemovesRecordAndHistory(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	db := newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)
	id := seedAsset(t, spaceID)

	if err := RemoveAsset(userCtx(testUserID), dto4assetus.RemoveAssetRequest{
		SpaceRequest: spaceRequest(spaceID), AssetID: id, HardDelete: true,
	}); err != nil {
		t.Fatalf("hard-delete failed: %v", err)
	}

	_, err := GetAsset(userCtx(testUserID), dto4assetus.GetAssetRequest{SpaceRequest: spaceRequest(spaceID), AssetID: id})
	if err == nil {
		t.Error("asset still retrievable after hard-delete")
	} else if !errors.Is(err, dal.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound, got %v", err)
	}
	// History must be gone too (read directly via dal, bypassing the existence gate).
	events, err := dal4assetus.ListAssetHistory(context.Background(), db, spaceID, id)
	if err != nil {
		t.Fatalf("history list failed: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("history not removed on hard-delete: %d events remain", len(events))
	}
}
