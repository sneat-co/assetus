package api4assetus

import (
	"net/http"

	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/assetus/backend/facade4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

var getAsset = facade4assetus.GetAsset

// httpGetAsset handles GET /v0/assetus/asset?spaceID=&assetID=.
func httpGetAsset(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	request := dto4assetus.GetAssetRequest{AssetID: query.Get("assetID")}
	request.SpaceID = coretypes.SpaceID(query.Get("spaceID"))
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	response, err := getAsset(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, &response)
}
