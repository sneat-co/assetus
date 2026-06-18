package dal4assetus

import (
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

// AssetWorkerParams is passed to an AssetWorker: the membership-gated module
// space params plus the loaded asset and a slot to accumulate asset updates.
type AssetWorkerParams struct {
	*dal4spaceus.ModuleSpaceWorkerParams[*dbo4assetus.AssetusSpaceDbo]
	Asset        AssetEntry
	AssetUpdates []update.Update
}

// AssetWorker is the unit of work run inside RunAssetWorker.
type AssetWorker = func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *AssetWorkerParams) (err error)

// RunAssetWorker runs worker inside a membership-gated read-write transaction
// for a single asset. Membership of the owning Space is enforced by
// dal4spaceus.RunModuleSpaceWorkerWithUserCtx (non-members are rejected with
// facade.ErrUnauthorized before the worker runs). The asset is loaded for
// update; a not-found asset is left as a non-existing record for the worker to
// handle (e.g. create).
func RunAssetWorker(ctx facade.ContextWithUser, spaceID coretypes.SpaceID, assetID string, worker AssetWorker) (err error) {
	return dal4spaceus.RunModuleSpaceWorkerWithUserCtx(ctx, spaceID, const4assetus.ExtensionID, new(dbo4assetus.AssetusSpaceDbo),
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, spaceWorkerParams *dal4spaceus.ModuleSpaceWorkerParams[*dbo4assetus.AssetusSpaceDbo]) (err error) {
			params := AssetWorkerParams{
				ModuleSpaceWorkerParams: spaceWorkerParams,
				Asset:                   NewAssetEntry(spaceID, assetID),
			}
			if assetID != "" {
				if err = GetAssetForUpdate(ctx, tx, params.Asset); err != nil && !dal.IsNotFound(err) {
					return fmt.Errorf("failed to get asset for update: %w", err)
				}
			}
			if err = worker(ctx, tx, &params); err != nil {
				return err
			}
			if updateCount := len(params.AssetUpdates); updateCount > 0 {
				if err = tx.Update(ctx, params.Asset.Key, params.AssetUpdates); err != nil {
					return fmt.Errorf("failed to apply %d asset updates: %w", updateCount, err)
				}
			}
			return nil
		})
}
