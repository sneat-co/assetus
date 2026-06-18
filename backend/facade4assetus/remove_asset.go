package facade4assetus

import (
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dal4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// RemoveAsset removes an asset in one of two ways. The default (soft-archive)
// sets status to Archived and preserves the record and its full history. An
// explicit HardDelete permanently removes the asset record and all of its
// history events. Membership of the owning Space is enforced for both paths.
func RemoveAsset(ctx facade.ContextWithUser, request dto4assetus.RemoveAssetRequest) (err error) {
	if err = request.Validate(); err != nil {
		return err
	}
	if request.HardDelete {
		return hardDeleteAsset(ctx, request)
	}
	return softArchiveAsset(ctx, request)
}

// softArchiveAsset sets status -> Archived, preserving record and history.
func softArchiveAsset(ctx facade.ContextWithUser, request dto4assetus.RemoveAssetRequest) error {
	return dal4assetus.RunAssetWorker(ctx, request.SpaceID, request.AssetID,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4assetus.AssetWorkerParams) (err error) {
			if !params.Asset.Record.Exists() {
				return fmt.Errorf("%w: asset %q not found in space %q", dal.ErrRecordNotFound, request.AssetID, request.SpaceID)
			}
			asset := params.Asset.Data
			asset.Status = const4assetus.StatusArchived
			asset.UpdatedAt = params.Started
			asset.UpdatedBy = params.UserID()
			if err = tx.Set(ctx, params.Asset.Record); err != nil {
				return fmt.Errorf("failed to archive asset record: %w", err)
			}
			if err = params.GetRecords(ctx, tx); err != nil {
				return fmt.Errorf("failed to load module entry: %w", err)
			}
			if params.SpaceModuleEntry.Record.Exists() {
				brief := dbo4assetus.BriefFromAsset(request.AssetID, asset)
				if params.SpaceModuleEntry.Data.Assets == nil {
					params.SpaceModuleEntry.Data.Assets = make(dbo4assetus.AssetBriefs, 1)
				}
				params.SpaceModuleEntry.Data.Assets[request.AssetID] = brief
				params.AddSpaceModuleUpdates(update.ByFieldPath([]string{"assets", request.AssetID}, brief))
			}
			return nil
		})
}

// hardDeleteAsset permanently removes the asset record and all of its history.
func hardDeleteAsset(ctx facade.ContextWithUser, request dto4assetus.RemoveAssetRequest) error {
	// Gather history event keys up front (queries run outside the write tx).
	db, err := facade.GetSneatDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get db: %w", err)
	}
	events, err := dal4assetus.ListAssetHistory(ctx, db, request.SpaceID, request.AssetID)
	if err != nil {
		return fmt.Errorf("failed to list history for hard-delete: %w", err)
	}
	return dal4assetus.RunAssetWorker(ctx, request.SpaceID, request.AssetID,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4assetus.AssetWorkerParams) (err error) {
			if !params.Asset.Record.Exists() {
				return fmt.Errorf("%w: asset %q not found in space %q", dal.ErrRecordNotFound, request.AssetID, request.SpaceID)
			}
			// Delete every history event.
			for _, e := range events {
				key := dal4assetus.NewHistoryEventKey(request.SpaceID, request.AssetID, e.ID)
				if err = tx.Delete(ctx, key); err != nil {
					return fmt.Errorf("failed to delete history event %q: %w", e.ID, err)
				}
			}
			// Delete the asset record itself.
			if err = tx.Delete(ctx, params.Asset.Key); err != nil {
				return fmt.Errorf("failed to delete asset record: %w", err)
			}
			// Drop the denormalized brief from the module entry.
			if err = params.GetRecords(ctx, tx); err != nil {
				return fmt.Errorf("failed to load module entry: %w", err)
			}
			if params.SpaceModuleEntry.Record.Exists() {
				if params.SpaceModuleEntry.Data.Assets != nil {
					delete(params.SpaceModuleEntry.Data.Assets, request.AssetID)
				}
				params.AddSpaceModuleUpdates(update.ByFieldPath([]string{"assets", request.AssetID}, update.DeleteField))
			}
			return nil
		})
}
