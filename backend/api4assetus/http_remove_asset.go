package api4assetus

import (
	"net/http"

	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/assetus/backend/facade4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

var removeAsset = facade4assetus.RemoveAsset

// httpPostRemoveAsset handles POST /v0/assetus/remove_asset. Soft-archive by
// default; set hardDelete=true for an explicit permanent delete.
func httpPostRemoveAsset(w http.ResponseWriter, r *http.Request) {
	var request dto4assetus.RemoveAssetRequest
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = removeAsset(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, nil)
}
