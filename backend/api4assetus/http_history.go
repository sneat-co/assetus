package api4assetus

import (
	"net/http"

	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/assetus/backend/facade4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

var (
	recordHistoryEvent = facade4assetus.RecordHistoryEvent
	getHistory         = facade4assetus.GetHistory
)

// httpPostRecordHistoryEvent handles POST /v0/assetus/record_history_event.
func httpPostRecordHistoryEvent(w http.ResponseWriter, r *http.Request) {
	var request dto4assetus.RecordHistoryEventRequest
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = recordHistoryEvent(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, nil)
}

// httpGetHistory handles GET /v0/assetus/asset_history?spaceID=&assetID=.
func httpGetHistory(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	request := dto4assetus.GetHistoryRequest{AssetID: query.Get("assetID")}
	request.SpaceID = coretypes.SpaceID(query.Get("spaceID"))
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.Request(verify.AuthenticationRequired(true)))
	if err != nil {
		return
	}
	if err = request.Validate(); err != nil {
		apicore.ReturnError(ctx, w, r, err)
		return
	}
	response, err := getHistory(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, &response)
}
