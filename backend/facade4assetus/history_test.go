package facade4assetus

import (
	"testing"
	"time"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// AC: history-is-append-only — record Purchased then Repaired; both present in
// order; appending a later event does not alter earlier ones.
func TestHistory_AppendOnlyOrdered(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	_ = newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)

	created, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		Name:         "Drill",
		Category:     const4assetus.CategoryTools,
		Condition:    const4assetus.ConditionGood,
	})
	if err != nil {
		t.Fatalf("CreateAsset failed: %v", err)
	}

	// Creation records a Purchased event.
	hist, err := GetHistory(userCtx(testUserID), dto4assetus.GetHistoryRequest{
		SpaceRequest: spaceRequest(spaceID), AssetID: created.ID,
	})
	if err != nil {
		t.Fatalf("GetHistory failed: %v", err)
	}
	if len(hist.Events) != 1 || hist.Events[0].Type != const4assetus.HistoryEventPurchased {
		t.Fatalf("expected [Purchased] after create, got %+v", hist.Events)
	}
	purchasedID := hist.Events[0].ID

	// Append a Repaired event (later timestamp).
	later := time.Now().Add(time.Hour)
	if err = RecordHistoryEvent(userCtx(testUserID), dto4assetus.RecordHistoryEventRequest{
		SpaceRequest: spaceRequest(spaceID), AssetID: created.ID,
		Type: const4assetus.HistoryEventRepaired, OccurredAt: &later, Note: "fixed the chuck",
	}); err != nil {
		t.Fatalf("RecordHistoryEvent failed: %v", err)
	}

	hist, err = GetHistory(userCtx(testUserID), dto4assetus.GetHistoryRequest{
		SpaceRequest: spaceRequest(spaceID), AssetID: created.ID,
	})
	if err != nil {
		t.Fatalf("GetHistory (2) failed: %v", err)
	}
	if len(hist.Events) != 2 {
		t.Fatalf("expected 2 events, got %d: %+v", len(hist.Events), hist.Events)
	}
	if hist.Events[0].Type != const4assetus.HistoryEventPurchased || hist.Events[0].ID != purchasedID {
		t.Errorf("first event mutated: got %+v, want unchanged Purchased %q", hist.Events[0], purchasedID)
	}
	if hist.Events[1].Type != const4assetus.HistoryEventRepaired {
		t.Errorf("second event = %q, want repaired", hist.Events[1].Type)
	}
}

// Transfer events cannot be recorded via the generic record endpoint.
func TestRecordHistoryEvent_RejectsTransferType(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	_ = newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)
	created, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID), Name: "Bike", Category: const4assetus.CategorySportsEquipment, Condition: const4assetus.ConditionGood,
	})
	if err != nil {
		t.Fatalf("CreateAsset failed: %v", err)
	}
	if err = RecordHistoryEvent(userCtx(testUserID), dto4assetus.RecordHistoryEventRequest{
		SpaceRequest: spaceRequest(spaceID), AssetID: created.ID, Type: const4assetus.HistoryEventTransferred,
	}); err == nil {
		t.Error("expected rejection of Transferred via generic record endpoint")
	}
}
