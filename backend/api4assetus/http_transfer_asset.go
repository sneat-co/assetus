package api4assetus

import (
	"net/http"

	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/assetus/backend/facade4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

var transferAsset = facade4assetus.TransferAsset

// httpPostTransferAsset handles POST /v0/assetus/transfer_asset.
func httpPostTransferAsset(w http.ResponseWriter, r *http.Request) {
	var request dto4assetus.TransferAssetRequest
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	response, err := transferAsset(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, &response)
}
