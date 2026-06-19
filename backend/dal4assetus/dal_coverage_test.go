package dal4assetus

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/dal-go/dalgo/adapters/dalgo2memory"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/with"
)

// newTestDBWithSpace builds an in-memory dalgo DB seeded with a single Space
// record (members = userIDs), wires facade.GetSneatDB to it, and returns the
// DB. Copied from facade4assetus/helpers_test.go so membership-gated worker
// tests in this package have the same setup.
func newTestDBWithSpace(t *testing.T, spaceID coretypes.SpaceID, spaceType coretypes.SpaceType, userIDs ...string) dal.DB {
	t.Helper()
	db := dalgo2memory.NewDB()
	now := time.Now()
	space := dbo4spaceus.NewSpaceEntry(spaceID)
	space.Data.Type = spaceType
	space.Data.Title = "Test " + string(spaceType) + " space"
	space.Data.Status = dbmodels.StatusActive
	space.Data.CreatedAt = now
	space.Data.CreatedBy = "seed"
	space.Data.IncreaseVersion(now, "seed")
	space.Data.UserIDs = userIDs
	if err := space.Data.Validate(); err != nil {
		t.Fatalf("seed space invalid: %v", err)
	}
	ctx := context.Background()
	if err := db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Insert(ctx, space.Record)
	}); err != nil {
		t.Fatalf("failed to seed space: %v", err)
	}
	facade.GetSneatDB = func(context.Context) (dal.DB, error) { return db, nil }
	return db
}

func validHistoryEvent(occurredAt time.Time) *dbo4assetus.AssetHistoryEventDbo {
	return &dbo4assetus.AssetHistoryEventDbo{
		AssetHistoryEventBase: dbo4assetus.AssetHistoryEventBase{
			Type:       const4assetus.HistoryEventRepaired,
			OccurredAt: occurredAt,
			ActorRef:   "user1",
		},
	}
}

func validVehicleRecord() *dbo4assetus.VehicleRecordDbo {
	return &dbo4assetus.VehicleRecordDbo{
		CreatedFields: with.CreatedFields{
			CreatedAtField: with.CreatedAtField{CreatedAt: time.Now()},
			CreatedByField: with.CreatedByField{CreatedBy: "user1"},
		},
		Mileage: &dbo4assetus.VehicleMileage{Value: 1000, Unit: "km"},
	}
}

func TestNewHistoryEventKey_Path(t *testing.T) {
	key := NewHistoryEventKey("space1", "asset1", "evt1")
	if key.ID != "evt1" {
		t.Errorf("key.ID = %v, want evt1", key.ID)
	}
	if key.Collection() != dbo4assetus.HistoryCollection {
		t.Errorf("key.Collection() = %v, want %v", key.Collection(), dbo4assetus.HistoryCollection)
	}
	if got := key.String(); got != "spaces/space1/ext/assetus/assets/asset1/history/evt1" {
		t.Errorf("key path = %q, want spaces/space1/ext/assetus/assets/asset1/history/evt1", got)
	}
}

func TestNewVehicleRecordKey_Path(t *testing.T) {
	key := NewVehicleRecordKey("space1", "asset1", "rec1")
	if key.ID != "rec1" {
		t.Errorf("key.ID = %v, want rec1", key.ID)
	}
	if key.Collection() != dbo4assetus.VehicleRecordsCollection {
		t.Errorf("key.Collection() = %v, want %v", key.Collection(), dbo4assetus.VehicleRecordsCollection)
	}
	if got := key.String(); got != "spaces/space1/ext/assetus/assets/asset1/vehicleRecords/rec1" {
		t.Errorf("key path = %q, want spaces/space1/ext/assetus/assets/asset1/vehicleRecords/rec1", got)
	}
}

func TestAppendHistoryEvent_RoundTrip(t *testing.T) {
	const spaceID coretypes.SpaceID = "space1"
	db := dalgo2memory.NewDB()
	ctx := context.Background()
	occurredAt := time.Now()

	err := db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return AppendHistoryEvent(ctx, tx, spaceID, "asset1", "evt1", validHistoryEvent(occurredAt))
	})
	if err != nil {
		t.Fatalf("AppendHistoryEvent failed: %v", err)
	}

	got := new(dbo4assetus.AssetHistoryEventDbo)
	rec := dal.NewRecordWithData(NewHistoryEventKey(spaceID, "asset1", "evt1"), got)
	if err = db.Get(ctx, rec); err != nil {
		t.Fatalf("failed to read persisted history event: %v", err)
	}
	if got.Type != const4assetus.HistoryEventRepaired {
		t.Errorf("persisted type = %v, want Repaired", got.Type)
	}
	if got.ActorRef != "user1" {
		t.Errorf("persisted actorRef = %q, want user1", got.ActorRef)
	}
}

func TestAppendHistoryEvent_RejectsInvalidDbo(t *testing.T) {
	const spaceID coretypes.SpaceID = "space1"
	db := dalgo2memory.NewDB()
	ctx := context.Background()

	// Missing OccurredAt/ActorRef -> Validate must reject before insert.
	invalid := &dbo4assetus.AssetHistoryEventDbo{
		AssetHistoryEventBase: dbo4assetus.AssetHistoryEventBase{
			Type: const4assetus.HistoryEventRepaired,
		},
	}
	err := db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return AppendHistoryEvent(ctx, tx, spaceID, "asset1", "evt1", invalid)
	})
	if err == nil {
		t.Fatal("expected invalid history event to be rejected")
	}
}

func TestListAssetHistory_OrdersByOccurredAtAscending(t *testing.T) {
	const spaceID coretypes.SpaceID = "space1"
	db := dalgo2memory.NewDB()
	ctx := context.Background()

	earlier := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	later := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Insert out of order: the later event first, the earlier event second.
	err := db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		if e := AppendHistoryEvent(ctx, tx, spaceID, "asset1", "evtLater", validHistoryEvent(later)); e != nil {
			return e
		}
		return AppendHistoryEvent(ctx, tx, spaceID, "asset1", "evtEarlier", validHistoryEvent(earlier))
	})
	if err != nil {
		t.Fatalf("seed history failed: %v", err)
	}

	events, err := ListAssetHistory(ctx, db, spaceID, "asset1")
	if err != nil {
		t.Fatalf("ListAssetHistory failed: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("len(events) = %d, want 2", len(events))
	}
	if !events[0].Dbo.OccurredAt.Equal(earlier) {
		t.Errorf("first event OccurredAt = %v, want earliest %v", events[0].Dbo.OccurredAt, earlier)
	}
	if !events[1].Dbo.OccurredAt.Equal(later) {
		t.Errorf("second event OccurredAt = %v, want latest %v", events[1].Dbo.OccurredAt, later)
	}
	if events[0].ID != "evtEarlier" {
		t.Errorf("first event ID = %q, want evtEarlier", events[0].ID)
	}
}

func TestAppendVehicleRecord_RoundTrip(t *testing.T) {
	const spaceID coretypes.SpaceID = "space1"
	db := dalgo2memory.NewDB()
	ctx := context.Background()

	err := db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return AppendVehicleRecord(ctx, tx, spaceID, "asset1", "rec1", validVehicleRecord())
	})
	if err != nil {
		t.Fatalf("AppendVehicleRecord failed: %v", err)
	}

	got := new(dbo4assetus.VehicleRecordDbo)
	rec := dal.NewRecordWithData(NewVehicleRecordKey(spaceID, "asset1", "rec1"), got)
	if err = db.Get(ctx, rec); err != nil {
		t.Fatalf("failed to read persisted vehicle record: %v", err)
	}
	if got.Mileage == nil {
		t.Fatal("expected mileage to be persisted")
	}
	if got.Mileage.Value != 1000 {
		t.Errorf("persisted mileage value = %d, want 1000", got.Mileage.Value)
	}
}

func TestAppendVehicleRecord_RejectsInvalidDbo(t *testing.T) {
	const spaceID coretypes.SpaceID = "space1"
	db := dalgo2memory.NewDB()
	ctx := context.Background()

	// Missing CreatedFields -> Validate must reject before insert.
	invalid := &dbo4assetus.VehicleRecordDbo{
		Mileage: &dbo4assetus.VehicleMileage{Value: 1, Unit: "km"},
	}
	err := db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return AppendVehicleRecord(ctx, tx, spaceID, "asset1", "rec1", invalid)
	})
	if err == nil {
		t.Fatal("expected invalid vehicle record to be rejected")
	}
}

// seedAsset inserts a minimal valid asset directly into the DB and returns it.
func seedAsset(t *testing.T, db dal.DB, spaceID coretypes.SpaceID, assetID string) {
	t.Helper()
	entry := NewAssetEntry(spaceID, assetID)
	now := time.Now()
	entry.Data.AssetBase = dbo4assetus.AssetBase{
		Name:       "Book",
		Category:   const4assetus.CategoryBooks,
		Condition:  const4assetus.ConditionGood,
		Status:     const4assetus.StatusActive,
		Visibility: const4assetus.VisibilityFamily,
	}
	entry.Data.SpaceIDs = []coretypes.SpaceID{spaceID}
	entry.Data.CreatedAt = now
	entry.Data.CreatedBy = "user1"
	entry.Data.UpdatedAt = now
	entry.Data.UpdatedBy = "user1"
	if err := entry.Data.Validate(); err != nil {
		t.Fatalf("seed asset invalid: %v", err)
	}
	if err := db.RunReadwriteTransaction(context.Background(), func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Insert(ctx, entry.Record)
	}); err != nil {
		t.Fatalf("failed to seed asset: %v", err)
	}
}

func TestGetAssetByID_LoadsPersistedAsset(t *testing.T) {
	const spaceID coretypes.SpaceID = "space1"
	db := dalgo2memory.NewDB()
	seedAsset(t, db, spaceID, "asset1")

	asset := NewAssetEntry(spaceID, "asset1")
	if err := GetAssetByID(context.Background(), db, asset); err != nil {
		t.Fatalf("GetAssetByID failed: %v", err)
	}
	if asset.Data.Name != "Book" {
		t.Errorf("loaded asset Name = %q, want Book", asset.Data.Name)
	}
}

func TestGetAssetForUpdate_LoadsPersistedAsset(t *testing.T) {
	const spaceID coretypes.SpaceID = "space1"
	db := dalgo2memory.NewDB()
	seedAsset(t, db, spaceID, "asset1")

	asset := NewAssetEntry(spaceID, "asset1")
	err := db.RunReadwriteTransaction(context.Background(), func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return GetAssetForUpdate(ctx, tx, asset)
	})
	if err != nil {
		t.Fatalf("GetAssetForUpdate failed: %v", err)
	}
	if asset.Data.Category != const4assetus.CategoryBooks {
		t.Errorf("loaded asset Category = %q, want Books", asset.Data.Category)
	}
}

func TestRunAssetWorker_MemberCreatesAsset(t *testing.T) {
	const spaceID coretypes.SpaceID = "space1"
	db := newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, "user1")
	ctx := facade.NewContextWithUserID(context.Background(), "user1")

	const assetID = "asset1"
	err := RunAssetWorker(ctx, spaceID, assetID, func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *AssetWorkerParams) error {
		if err := params.GetRecords(ctx, tx); err != nil {
			return fmt.Errorf("failed to load module records: %w", err)
		}
		now := params.Started
		asset := params.Asset.Data
		asset.AssetBase = dbo4assetus.AssetBase{
			Name:       "Worker Book",
			Category:   const4assetus.CategoryBooks,
			Condition:  const4assetus.ConditionGood,
			Status:     const4assetus.StatusActive,
			Visibility: const4assetus.VisibilityFamily,
		}
		asset.SpaceIDs = []coretypes.SpaceID{spaceID}
		asset.CreatedAt = now
		asset.CreatedBy = params.UserID()
		asset.UpdatedAt = now
		asset.UpdatedBy = params.UserID()
		if err := asset.Validate(); err != nil {
			return fmt.Errorf("formed asset invalid: %w", err)
		}
		return tx.Insert(ctx, params.Asset.Record)
	})
	if err != nil {
		t.Fatalf("RunAssetWorker (member) failed: %v", err)
	}

	loaded := NewAssetEntry(spaceID, assetID)
	if err = GetAssetByID(context.Background(), db, loaded); err != nil {
		t.Fatalf("created asset was not persisted: %v", err)
	}
	if loaded.Data.Name != "Worker Book" {
		t.Errorf("persisted asset Name = %q, want Worker Book", loaded.Data.Name)
	}
}

func TestRunAssetWorker_NonMemberRejected(t *testing.T) {
	const spaceID coretypes.SpaceID = "space1"
	_ = newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, "user1")
	ctx := facade.NewContextWithUserID(context.Background(), "intruder")

	called := false
	err := RunAssetWorker(ctx, spaceID, "asset1", func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *AssetWorkerParams) error {
		called = true
		return nil
	})
	if err == nil {
		t.Fatal("expected non-member to be rejected")
	}
	if !errors.Is(err, facade.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got: %v", err)
	}
	if called {
		t.Error("worker body must not run for a non-member")
	}
}
