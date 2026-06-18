package facade4assetus

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

// seedExtraSpace adds a second space to an existing in-memory DB.
func seedExtraSpace(t *testing.T, db dal.DB, spaceID coretypes.SpaceID, spaceType coretypes.SpaceType, userIDs ...string) {
	t.Helper()
	now := time.Now()
	space := dbo4spaceus.NewSpaceEntry(spaceID)
	space.Data.Type = spaceType
	space.Data.Title = "Test " + string(spaceType)
	space.Data.Status = dbmodels.StatusActive
	space.Data.CreatedAt = now
	space.Data.CreatedBy = "seed"
	space.Data.IncreaseVersion(now, "seed")
	space.Data.UserIDs = userIDs
	if err := db.RunReadwriteTransaction(context.Background(), func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Insert(ctx, space.Record)
	}); err != nil {
		t.Fatalf("failed to seed extra space: %v", err)
	}
}

// AC: transfer-reassigns-and-records — transferring from A to B re-owns the
// asset (owner type re-derived from B) and appends a Transferred event recording
// A as prior and B as new owner.
func TestTransferAsset_ReassignsAndRecords(t *testing.T) {
	const spaceA coretypes.SpaceID = "family1"
	const spaceB coretypes.SpaceID = "club1"
	db := newTestDBWithSpace(t, spaceA, coretypes.SpaceTypeFamily, testUserID)
	seedExtraSpace(t, db, spaceB, coretypes.SpaceTypeClub, testUserID)

	id := seedAsset(t, spaceA)

	resp, err := TransferAsset(userCtx(testUserID), dto4assetus.TransferAssetRequest{
		SpaceRequest: spaceRequest(spaceA),
		AssetID:      id,
		ToSpaceID:    spaceB,
	})
	if err != nil {
		t.Fatalf("TransferAsset failed: %v", err)
	}
	if resp.Owner.OwnerType != string(const4assetus.OwnerTypeSportsClub) {
		t.Errorf("new owner type = %q, want sports_club", resp.Owner.OwnerType)
	}

	// Asset now owned by B.
	got, err := GetAsset(userCtx(testUserID), dto4assetus.GetAssetRequest{SpaceRequest: spaceRequest(spaceB), AssetID: id})
	if err != nil {
		t.Fatalf("GetAsset under B failed: %v", err)
	}
	if got.Owner.OwnerType != const4assetus.OwnerTypeSportsClub {
		t.Errorf("owner type under B = %q, want sports_club", got.Owner.OwnerType)
	}

	// The relocated asset's owning Space is now B, not A. (The "no longer
	// present under A's path" guarantee relies on Firestore's full-path
	// scoping; the dalgo2memory adapter keys records by collection-name+id
	// only, so it cannot represent two same-id records under different parents —
	// see backstage NEEDS-REVIEW; covered by the Firestore-emulator path.)
	if slices.Contains(got.Asset.SpaceIDs, spaceA) {
		t.Errorf("relocated asset still lists source space A as owner: %v", got.Asset.SpaceIDs)
	}
	if !slices.Contains(got.Asset.SpaceIDs, spaceB) {
		t.Errorf("relocated asset does not list destination space B as owner: %v", got.Asset.SpaceIDs)
	}

	// A Transferred history event records prior (A) and new (B) owner.
	hist, err := GetHistory(userCtx(testUserID), dto4assetus.GetHistoryRequest{SpaceRequest: spaceRequest(spaceB), AssetID: id})
	if err != nil {
		t.Fatalf("GetHistory under B failed: %v", err)
	}
	var transferred *dto4assetus.HistoryEventItem
	for i := range hist.Events {
		if hist.Events[i].Type == const4assetus.HistoryEventTransferred {
			transferred = &hist.Events[i]
		}
	}
	if transferred == nil {
		t.Fatalf("no Transferred event found; history=%+v", hist.Events)
	}
	if transferred.FromOwner == nil || transferred.FromOwner.SpaceID != string(spaceA) {
		t.Errorf("Transferred.fromOwner = %+v, want spaceID %q", transferred.FromOwner, spaceA)
	}
	if transferred.ToOwner == nil || transferred.ToOwner.SpaceID != string(spaceB) {
		t.Errorf("Transferred.toOwner = %+v, want spaceID %q", transferred.ToOwner, spaceB)
	}
	// Prior history (Purchased on create) is preserved through the relocation.
	var hasPurchased bool
	for _, e := range hist.Events {
		if e.Type == const4assetus.HistoryEventPurchased {
			hasPurchased = true
		}
	}
	if !hasPurchased {
		t.Error("prior Purchased event lost during transfer")
	}
}
