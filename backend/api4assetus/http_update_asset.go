package api4assetus

import (
	"net/http"

	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/assetus/backend/facade4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

var updateAsset = facade4assetus.UpdateAsset

// httpPostUpdateAsset handles POST /v0/assetus/update_asset.
func httpPostUpdateAsset(w http.ResponseWriter, r *http.Request) {
	var request dto4assetus.UpdateAssetRequest
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	response, err := updateAsset(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, &response)
}
