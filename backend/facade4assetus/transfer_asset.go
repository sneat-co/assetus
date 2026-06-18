package facade4assetus

import (
	"context"
	"fmt"
	"time"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dal4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

// TransferAsset transfers an asset from its owning Space to a destination Space.
// The acting user must be a member of the owning Space. After transfer the
// asset's owning Space — and therefore its derived owner type — is the
// destination, and a Transferred history event is appended recording the prior
// owner (source) and the new owner (destination). The prior owner is preserved
// in history, never silently overwritten. The asset record and its full history
// are relocated to the destination Space's collection.
func TransferAsset(ctx facade.ContextWithUser, request dto4assetus.TransferAssetRequest) (response dto4assetus.TransferAssetResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	db, err := facade.GetSneatDB(ctx)
	if err != nil {
		return response, fmt.Errorf("failed to get db: %w", err)
	}

	// Resolve the destination Space type (for the new owner derivation).
	destSpace := dbo4spaceus.NewSpaceEntry(request.ToSpaceID)
	if err = db.Get(ctx, destSpace.Record); err != nil {
		return response, fmt.Errorf("failed to read destination space %q: %w", request.ToSpaceID, err)
	}
	destSpaceType := destSpace.Data.Type

	// Capture the source history before relocation (queries run outside the tx).
	srcHistory, err := dal4assetus.ListAssetHistory(ctx, db, request.SpaceID, request.AssetID)
	if err != nil {
		return response, fmt.Errorf("failed to read source history: %w", err)
	}

	err = dal4assetus.RunAssetWorker(ctx, request.SpaceID, request.AssetID,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4assetus.AssetWorkerParams) (err error) {
			if !params.Asset.Record.Exists() {
				return fmt.Errorf("%w: asset %q not found in space %q", dal.ErrRecordNotFound, request.AssetID, request.SpaceID)
			}
			fromOwner := dbo4assetus.NewOwnerRef(request.SpaceID, params.Space.Data.Type)
			toOwner := dbo4assetus.NewOwnerRef(request.ToSpaceID, destSpaceType)

			now := params.Started
			userID := params.UserID()

			// Build the relocated asset data under the destination space.
			newAsset := params.Asset.Data
			newAsset.SpaceIDs = []coretypes.SpaceID{request.ToSpaceID}
			newAsset.UpdatedAt = now
			newAsset.UpdatedBy = userID
			if err = newAsset.Validate(); err != nil {
				return fmt.Errorf("relocated asset is not valid: %w", err)
			}

			// Delete the source asset and its history FIRST so the relocation
			// can reuse the same asset/event IDs at the destination path.
			for _, e := range srcHistory {
				if err = tx.Delete(ctx, dal4assetus.NewHistoryEventKey(request.SpaceID, request.AssetID, e.ID)); err != nil {
					return fmt.Errorf("failed to delete source history event %q: %w", e.ID, err)
				}
			}
			if err = tx.Delete(ctx, params.Asset.Key); err != nil {
				return fmt.Errorf("failed to delete source asset: %w", err)
			}

			// Insert the relocated asset under the destination space.
			destAsset := dal4assetus.NewAssetEntry(request.ToSpaceID, request.AssetID)
			destAsset.Data = newAsset
			destAsset.Record = dal.NewRecordWithData(destAsset.Key, newAsset)
			if err = tx.Insert(ctx, destAsset.Record); err != nil {
				return fmt.Errorf("failed to insert relocated asset: %w", err)
			}

			// Relocate existing history events, preserving them.
			for _, e := range srcHistory {
				if err = dal4assetus.AppendHistoryEvent(ctx, tx, request.ToSpaceID, request.AssetID, e.ID, e.Dbo); err != nil {
					return fmt.Errorf("failed to relocate history event %q: %w", e.ID, err)
				}
			}
			// Append the Transferred event recording prior and new owner.
			transferred := &dbo4assetus.AssetHistoryEventDbo{
				AssetHistoryEventBase: dbo4assetus.AssetHistoryEventBase{
					Type:       const4assetus.HistoryEventTransferred,
					OccurredAt: now,
					ActorRef:   userID,
					FromOwner:  &fromOwner,
					ToOwner:    &toOwner,
				},
			}
			if err = dal4assetus.AppendHistoryEvent(ctx, tx, request.ToSpaceID, request.AssetID, newHistoryEventID(), transferred); err != nil {
				return fmt.Errorf("failed to append Transferred event: %w", err)
			}

			// Drop the brief from the source module entry.
			if err = params.GetRecords(ctx, tx); err != nil {
				return fmt.Errorf("failed to load source module entry: %w", err)
			}
			if params.SpaceModuleEntry.Record.Exists() {
				if params.SpaceModuleEntry.Data.Assets != nil {
					delete(params.SpaceModuleEntry.Data.Assets, request.AssetID)
				}
				params.AddSpaceModuleUpdates(update.ByFieldPath([]string{"assets", request.AssetID}, update.DeleteField))
			}

			// Upsert the brief onto the destination module entry.
			if err = upsertDestBrief(ctx, tx, request.ToSpaceID, request.AssetID, newAsset, now, userID); err != nil {
				return fmt.Errorf("failed to update destination module entry: %w", err)
			}

			response.ID = request.AssetID
			response.Owner = dto4assetus.OwnerRefDTO{
				SpaceID:   string(toOwner.SpaceID),
				SpaceType: string(toOwner.SpaceType),
				OwnerType: string(toOwner.OwnerType),
			}
			return nil
		})
	if err != nil {
		return dto4assetus.TransferAssetResponse{}, err
	}
	return
}

// upsertDestBrief adds the relocated asset's brief to the destination space's
// assetus module entry, creating the entry if it does not yet exist.
func upsertDestBrief(ctx context.Context, tx dal.ReadwriteTransaction, spaceID coretypes.SpaceID, assetID string, asset *dbo4assetus.AssetDbo, now time.Time, actorRef string) error {
	moduleKey := dbo4spaceus.NewSpaceModuleKey(spaceID, const4assetus.ExtensionID)
	data := new(dbo4assetus.AssetusSpaceDbo)
	rec := dal.NewRecordWithData(moduleKey, data)
	if err := tx.Get(ctx, rec); err != nil && !dal.IsNotFound(err) {
		return fmt.Errorf("failed to read destination module entry: %w", err)
	}
	brief := dbo4assetus.BriefFromAsset(assetID, asset)
	if data.Assets == nil {
		data.Assets = make(dbo4assetus.AssetBriefs, 1)
	}
	data.Assets[assetID] = brief
	if rec.Exists() {
		return tx.Set(ctx, rec)
	}
	data.CreatedAt = now
	data.CreatedBy = actorRef
	return tx.Insert(ctx, rec)
}
