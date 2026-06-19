package api4assetus

import (
	"net/http"

	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/assetus/backend/facade4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

var addVehicleRecord = facade4assetus.AddVehicleRecord

// httpPostCreateVehicleRecord handles POST /v0/assetus/create_vehicle_record.
func httpPostCreateVehicleRecord(w http.ResponseWriter, r *http.Request) {
	var request dto4assetus.AddVehicleRecordRequest
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	response, err := addVehicleRecord(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, &response)
}
