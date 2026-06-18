package dal4assetus

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// HistoryEventEntry is a loaded history event with its key/ID.
type HistoryEventEntry = struct {
	ID  string
	Dbo *dbo4assetus.AssetHistoryEventDbo
}

// NewHistoryEventKey builds the dalgo key for a history event:
// /spaces/{spaceID}/ext/assetus/assets/{assetID}/history/{eventID}.
func NewHistoryEventKey(spaceID coretypes.SpaceID, assetID, eventID string) *dal.Key {
	return dal.NewKeyWithParentAndID(NewAssetKey(spaceID, assetID), dbo4assetus.HistoryCollection, eventID)
}

// historyCollectionRef is the dalgo collection ref for an asset's history.
func historyCollectionRef(spaceID coretypes.SpaceID, assetID string) dal.CollectionRef {
	return dal.NewCollectionRef(dbo4assetus.HistoryCollection, "", NewAssetKey(spaceID, assetID))
}

// AppendHistoryEvent inserts a new history event for an asset inside an existing
// transaction. It is INSERT-only: existing events are never mutated or removed.
func AppendHistoryEvent(ctx context.Context, tx dal.ReadwriteTransaction, spaceID coretypes.SpaceID, assetID, eventID string, dbo *dbo4assetus.AssetHistoryEventDbo) error {
	if err := dbo.Validate(); err != nil {
		return fmt.Errorf("invalid history event: %w", err)
	}
	rec := dal.NewRecordWithData(NewHistoryEventKey(spaceID, assetID, eventID), dbo)
	if err := tx.Insert(ctx, rec); err != nil {
		return fmt.Errorf("failed to insert history event: %w", err)
	}
	return nil
}

// ListAssetHistory reads all history events for an asset, ordered by OccurredAt
// ascending (stable).
func ListAssetHistory(ctx context.Context, db dal.DB, spaceID coretypes.SpaceID, assetID string) ([]HistoryEventEntry, error) {
	q := dal.From(historyCollectionRef(spaceID, assetID)).NewQuery().
		SelectIntoRecord(func() dal.Record {
			return dal.NewRecordWithIncompleteKey(dbo4assetus.HistoryCollection, reflect.String, new(dbo4assetus.AssetHistoryEventDbo))
		})
	records, err := dal.ExecuteQueryAndReadAllToRecords(ctx, q, db)
	if err != nil {
		return nil, fmt.Errorf("failed to list asset history: %w", err)
	}
	events := make([]HistoryEventEntry, 0, len(records))
	for _, rec := range records {
		id, _ := rec.Key().ID.(string)
		events = append(events, HistoryEventEntry{ID: id, Dbo: rec.Data().(*dbo4assetus.AssetHistoryEventDbo)})
	}
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].Dbo.OccurredAt.Before(events[j].Dbo.OccurredAt)
	})
	return events, nil
}
