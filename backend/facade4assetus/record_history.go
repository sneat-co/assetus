package facade4assetus

import (
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/assetus/backend/dal4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// RecordHistoryEvent appends a lifecycle event to an asset's append-only
// history. Membership of the owning Space is enforced. Existing events are never
// mutated or removed — this is an INSERT-only operation.
func RecordHistoryEvent(ctx facade.ContextWithUser, request dto4assetus.RecordHistoryEventRequest) (err error) {
	if err = request.Validate(); err != nil {
		return err
	}
	return dal4assetus.RunAssetWorker(ctx, request.SpaceID, request.AssetID,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4assetus.AssetWorkerParams) (err error) {
			if !params.Asset.Record.Exists() {
				return fmt.Errorf("%w: asset %q not found in space %q", dal.ErrRecordNotFound, request.AssetID, request.SpaceID)
			}
			occurredAt := params.Started
			if request.OccurredAt != nil {
				occurredAt = *request.OccurredAt
			}
			dbo := &dbo4assetus.AssetHistoryEventDbo{
				AssetHistoryEventBase: dbo4assetus.AssetHistoryEventBase{
					Type:       request.Type,
					OccurredAt: occurredAt,
					ActorRef:   params.UserID(),
					Note:       request.Note,
				},
			}
			return dal4assetus.AppendHistoryEvent(ctx, tx, request.SpaceID, request.AssetID, newHistoryEventID(), dbo)
		})
}

// GetHistory reads an asset's append-only history, ordered oldest-first.
// Membership of the owning Space is enforced.
func GetHistory(ctx facade.ContextWithUser, request dto4assetus.GetHistoryRequest) (response dto4assetus.GetHistoryResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	// Enforce membership (and asset existence) via the worker before reading.
	err = dal4assetus.RunAssetWorker(ctx, request.SpaceID, request.AssetID,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4assetus.AssetWorkerParams) (err error) {
			if !params.Asset.Record.Exists() {
				return fmt.Errorf("%w: asset %q not found in space %q", dal.ErrRecordNotFound, request.AssetID, request.SpaceID)
			}
			return nil
		})
	if err != nil {
		return
	}
	db, err := facade.GetSneatDB(ctx)
	if err != nil {
		return response, fmt.Errorf("failed to get db: %w", err)
	}
	events, err := dal4assetus.ListAssetHistory(ctx, db, request.SpaceID, request.AssetID)
	if err != nil {
		return response, err
	}
	response.AssetID = request.AssetID
	response.Events = make([]dto4assetus.HistoryEventItem, 0, len(events))
	for _, e := range events {
		response.Events = append(response.Events, toHistoryEventItem(e.ID, e.Dbo))
	}
	return response, nil
}

func toHistoryEventItem(id string, dbo *dbo4assetus.AssetHistoryEventDbo) dto4assetus.HistoryEventItem {
	item := dto4assetus.HistoryEventItem{
		ID:         id,
		Type:       dbo.Type,
		OccurredAt: dbo.OccurredAt,
		ActorRef:   dbo.ActorRef,
		Note:       dbo.Note,
	}
	item.FromOwner = toOwnerRefDTO(dbo.FromOwner)
	item.ToOwner = toOwnerRefDTO(dbo.ToOwner)
	return item
}

func toOwnerRefDTO(ref *dbo4assetus.OwnerRef) *dto4assetus.OwnerRefDTO {
	if ref == nil {
		return nil
	}
	return &dto4assetus.OwnerRefDTO{
		SpaceID:   string(ref.SpaceID),
		SpaceType: string(ref.SpaceType),
		OwnerType: string(ref.OwnerType),
	}
}
