package facade4assetus

import (
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/assetus/backend/dal4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// GetAsset reads an asset and exposes its owner: the owning Space plus the owner
// type derived from the Space type (private→Individual, family→Family,
// club→SportsClub, company→Organisation; community→Community, school→School once
// those spaceus space types ship). Membership of the owning Space is enforced.
func GetAsset(ctx facade.ContextWithUser, request dto4assetus.GetAssetRequest) (response dto4assetus.GetAssetResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4assetus.RunAssetWorker(ctx, request.SpaceID, request.AssetID,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4assetus.AssetWorkerParams) (err error) {
			if !params.Asset.Record.Exists() {
				return fmt.Errorf("%w: asset %q not found in space %q", dal.ErrRecordNotFound, request.AssetID, request.SpaceID)
			}
			response.ID = request.AssetID
			response.Asset = params.Asset.Data
			response.Owner = dbo4assetus.NewOwnerRef(request.SpaceID, params.Space.Data.Type)
			return nil
		})
	if err != nil {
		return dto4assetus.GetAssetResponse{}, err
	}
	return
}
