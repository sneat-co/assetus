package facade4assetus

import (
	"fmt"
	"strings"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/assetus/backend/dal4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// UpdateAsset edits an asset's name, description, category, condition,
// visibility, and optional metadata. Membership of the owning Space is enforced
// (non-members are rejected). The ownership-lifecycle status is preserved, and
// the asset's append-only history is NOT touched by an edit.
func UpdateAsset(ctx facade.ContextWithUser, request dto4assetus.UpdateAssetRequest) (response dto4assetus.UpdateAssetResponse, err error) {
	request.Name = strings.TrimSpace(request.Name)
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4assetus.RunAssetWorker(ctx, request.SpaceID, request.AssetID,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4assetus.AssetWorkerParams) (err error) {
			if !params.Asset.Record.Exists() {
				return fmt.Errorf("%w: asset %q not found in space %q", dal.ErrRecordNotFound, request.AssetID, request.SpaceID)
			}
			asset := params.Asset.Data
			// Preserve ownership-lifecycle status, ownership and audit-create fields.
			asset.Name = request.Name
			asset.Description = request.Description
			asset.Category = request.Category
			asset.Condition = request.Condition
			asset.Visibility = request.Visibility
			asset.AcquisitionDate = request.AcquisitionDate
			asset.PurchasePrice = request.PurchasePrice
			asset.EstimatedValue = request.EstimatedValue
			asset.Location = request.Location
			asset.Notes = request.Notes
			asset.Tags = request.Tags
			asset.UpdatedAt = params.Started
			asset.UpdatedBy = params.UserID()

			if err = asset.Validate(); err != nil {
				return fmt.Errorf("updated asset is not valid: %w", err)
			}
			// Full-record replacement of the asset document only; the history
			// child collection is a separate set of records and is untouched.
			if err = tx.Set(ctx, params.Asset.Record); err != nil {
				return fmt.Errorf("failed to update asset record: %w", err)
			}

			// Keep the denormalized brief on the module entry in sync.
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

			response.ID = request.AssetID
			response.Asset = asset
			return nil
		})
	if err != nil {
		return dto4assetus.UpdateAssetResponse{}, err
	}
	return
}
