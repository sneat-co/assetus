package facade4assetus

import (
	"fmt"
	"strings"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dal4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

// CreateAsset creates an asset owned by the request's Space. Membership of the
// Space is enforced by the underlying module worker (non-members are rejected
// before any write). The owner is the Space, status defaults to Active,
// condition is set by the creator, and visibility defaults to the Space's
// default visibility unless an override is supplied.
func CreateAsset(ctx facade.ContextWithUser, request dto4assetus.CreateAssetRequest) (response dto4assetus.CreateAssetResponse, err error) {
	request.Name = strings.TrimSpace(request.Name)
	if err = request.Validate(); err != nil {
		return
	}
	assetID := newAssetID()
	err = dal4assetus.RunAssetWorker(ctx, request.SpaceID, assetID,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4assetus.AssetWorkerParams) (err error) {
			if err = params.GetRecords(ctx, tx); err != nil {
				return fmt.Errorf("failed to load space module records: %w", err)
			}

			visibility := request.Visibility
			if visibility == "" {
				visibility = const4assetus.DefaultVisibilityForSpaceType(params.Space.Data.Type)
			}

			now := params.Started
			userID := params.UserID()

			asset := params.Asset.Data
			asset.AssetBase = dbo4assetus.AssetBase{
				Name:            request.Name,
				Description:     request.Description,
				Category:        request.Category,
				Condition:       request.Condition,
				Status:          const4assetus.StatusActive,
				Visibility:      visibility,
				AcquisitionDate: request.AcquisitionDate,
				PurchasePrice:   request.PurchasePrice,
				EstimatedValue:  request.EstimatedValue,
				Location:        request.Location,
				Notes:           request.Notes,
				Tags:            request.Tags,
			}
			asset.SpaceIDs = []coretypes.SpaceID{request.SpaceID}
			asset.CreatedAt = now
			asset.CreatedBy = userID
			asset.UpdatedAt = now
			asset.UpdatedBy = userID

			if err = asset.Validate(); err != nil {
				return fmt.Errorf("formed asset is not valid: %w", err)
			}
			if err = tx.Insert(ctx, params.Asset.Record); err != nil {
				return fmt.Errorf("failed to insert asset record: %w", err)
			}

			brief := dbo4assetus.BriefFromAsset(assetID, asset)
			if params.SpaceModuleEntry.Data.Assets == nil {
				params.SpaceModuleEntry.Data.Assets = make(dbo4assetus.AssetBriefs, 1)
			}
			params.SpaceModuleEntry.Data.Assets[assetID] = brief
			if params.SpaceModuleEntry.Record.Exists() {
				params.AddSpaceModuleUpdates(update.ByFieldPath([]string{"assets", assetID}, brief))
			} else {
				params.SpaceModuleEntry.Data.CreatedAt = now
				params.SpaceModuleEntry.Data.CreatedBy = userID
				if err = tx.Insert(ctx, params.SpaceModuleEntry.Record); err != nil {
					return fmt.Errorf("failed to insert assetus module entry: %w", err)
				}
			}

			response.ID = assetID
			response.Asset = asset
			return nil
		})
	if err != nil {
		return dto4assetus.CreateAssetResponse{}, err
	}
	return
}
