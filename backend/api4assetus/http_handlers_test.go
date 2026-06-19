package api4assetus

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
)

// --- Auth bypass scaffolding -------------------------------------------------
//
// The handlers call apicore.VerifyAuthenticatedRequestAndDecodeBody with
// verify.DefaultJsonWithAuthRequired. That function is composed of
// apicore.VerifyRequestAndCreateUserContext (the auth step) followed by
// apicore.DecodeRequestBody (the JSON decode step). We override the auth step
// with a stub that returns an authenticated user context, leaving the real
// decode in place so body-decoding behaviour is still exercised. This mirrors
// the pattern used by sneat-go-backend's own api handler tests.

type mockUserContext struct {
	facade.UserContext
	userID string
}

func (m mockUserContext) GetUserID() string { return m.userID }

type mockContextWithUser struct {
	facade.ContextWithUser
	ctx  context.Context
	user facade.UserContext
}

func (m mockContextWithUser) User() facade.UserContext { return m.user }
func (m mockContextWithUser) Value(key any) any        { return m.ctx.Value(key) }

// authAsUser overrides the auth step so requests are treated as authenticated.
func authAsUser(t *testing.T) {
	t.Helper()
	old := apicore.VerifyRequestAndCreateUserContext
	apicore.VerifyRequestAndCreateUserContext = func(w http.ResponseWriter, r *http.Request, options verify.RequestOptions) (facade.ContextWithUser, error) {
		return mockContextWithUser{
			ctx:  t.Context(),
			user: mockUserContext{userID: "u1"},
		}, nil
	}
	t.Cleanup(func() { apicore.VerifyRequestAndCreateUserContext = old })
}

// authRejected overrides the auth step so requests are treated as
// unauthenticated, exactly as the real verifier would for a missing token.
func authRejected(t *testing.T) {
	t.Helper()
	old := apicore.VerifyRequestAndCreateUserContext
	apicore.VerifyRequestAndCreateUserContext = func(w http.ResponseWriter, r *http.Request, options verify.RequestOptions) (facade.ContextWithUser, error) {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, facade.ErrUnauthenticated
	}
	t.Cleanup(func() { apicore.VerifyRequestAndCreateUserContext = old })
}

// Valid request bodies that pass each DTO's Validate() so decoding reaches the
// facade. Each must satisfy the corresponding *Request.Validate() rules.
const (
	validCreateAssetJSON   = `{"spaceID":"space1","name":"My Asset","category":"electronics"}`
	validUpdateAssetJSON   = `{"spaceID":"space1","assetID":"asset1","name":"My Asset","category":"electronics","condition":"good","visibility":"private"}`
	validRemoveAssetJSON   = `{"spaceID":"space1","assetID":"asset1"}`
	validTransferAssetJSON = `{"spaceID":"space1","assetID":"asset1","toSpaceID":"space2"}`
	validRecordEventJSON   = `{"spaceID":"space1","assetID":"asset1","type":"repaired"}`
	validVehicleRecJSON    = `{"spaceID":"space1","assetID":"asset1","mileage":1000,"mileageUnit":"km"}`
)

func newPostRequest(path, body string) *http.Request {
	return httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
}

// errBoom is a generic facade error; HandleError maps it to 500.
var errBoom = errors.New("boom")

// =============================================================================
// Routes
// =============================================================================

func TestRegisterHttpRoutes_registersAllEightEndpointsUnderAssetusPrefix(t *testing.T) {
	type reg struct {
		method, path string
	}
	var got []reg
	handle := func(method, path string, _ http.HandlerFunc) {
		if !strings.HasPrefix(path, "/v0/assetus/") {
			t.Errorf("path %q does not start with /v0/assetus/", path)
		}
		got = append(got, reg{method, path})
	}

	RegisterHttpRoutes(handle)

	want := []reg{
		{http.MethodPost, "/v0/assetus/create_asset"},
		{http.MethodGet, "/v0/assetus/asset"},
		{http.MethodPost, "/v0/assetus/update_asset"},
		{http.MethodPost, "/v0/assetus/remove_asset"},
		{http.MethodPost, "/v0/assetus/transfer_asset"},
		{http.MethodPost, "/v0/assetus/record_history_event"},
		{http.MethodGet, "/v0/assetus/asset_history"},
		{http.MethodPost, "/v0/assetus/create_vehicle_record"},
	}
	if len(got) != len(want) {
		t.Fatalf("registered %d routes, want %d: %+v", len(got), len(want), got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("route %d = %+v, want %+v", i, got[i], w)
		}
	}
}

// =============================================================================
// httpPostCreateAsset
// =============================================================================

func TestHttpPostCreateAsset_returns201WithAssetIDOnSuccess(t *testing.T) {
	authAsUser(t)
	old := createAsset
	t.Cleanup(func() { createAsset = old })
	createAsset = func(ctx facade.ContextWithUser, request dto4assetus.CreateAssetRequest) (dto4assetus.CreateAssetResponse, error) {
		return dto4assetus.CreateAssetResponse{ID: "new-asset-1"}, nil
	}

	w := httptest.NewRecorder()
	httpPostCreateAsset(w, newPostRequest("/v0/assetus/create_asset", validCreateAssetJSON))

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "new-asset-1") {
		t.Errorf("body %q does not contain created asset id", w.Body.String())
	}
}

func TestHttpPostCreateAsset_returns500WhenFacadeFails(t *testing.T) {
	authAsUser(t)
	old := createAsset
	t.Cleanup(func() { createAsset = old })
	createAsset = func(ctx facade.ContextWithUser, request dto4assetus.CreateAssetRequest) (dto4assetus.CreateAssetResponse, error) {
		return dto4assetus.CreateAssetResponse{}, errBoom
	}

	w := httptest.NewRecorder()
	httpPostCreateAsset(w, newPostRequest("/v0/assetus/create_asset", validCreateAssetJSON))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "boom") {
		t.Errorf("body %q does not contain facade error message", w.Body.String())
	}
}

func TestHttpPostCreateAsset_returns400WhenBodyIsInvalidJSON(t *testing.T) {
	authAsUser(t)
	w := httptest.NewRecorder()
	httpPostCreateAsset(w, newPostRequest("/v0/assetus/create_asset", "not json"))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}

func TestHttpPostCreateAsset_returns401WhenUnauthenticated(t *testing.T) {
	authRejected(t)
	w := httptest.NewRecorder()
	httpPostCreateAsset(w, newPostRequest("/v0/assetus/create_asset", validCreateAssetJSON))

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

// =============================================================================
// httpGetAsset
// =============================================================================

// httpGetAsset reads its input from the query string and authenticates via
// apicore.VerifyRequestAndCreateUserContext (no body decode), so a GET with
// valid query params + auth reaches the facade and returns 200 with the JSON
// asset.
func TestHttpGetAsset_returns200WithAssetOnSuccess(t *testing.T) {
	authAsUser(t)
	old := getAsset
	t.Cleanup(func() { getAsset = old })
	var gotRequest dto4assetus.GetAssetRequest
	getAsset = func(ctx facade.ContextWithUser, request dto4assetus.GetAssetRequest) (dto4assetus.GetAssetResponse, error) {
		gotRequest = request
		return dto4assetus.GetAssetResponse{ID: "a-42"}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/v0/assetus/asset?spaceID=s-7&assetID=a-42", nil)
	w := httptest.NewRecorder()
	httpGetAsset(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "a-42") {
		t.Errorf("body %q does not contain asset id", w.Body.String())
	}
	if gotRequest.AssetID != "a-42" || string(gotRequest.SpaceID) != "s-7" {
		t.Errorf("facade received request %+v, want assetID=a-42 spaceID=s-7", gotRequest)
	}
}

func TestHttpGetAsset_returns500WhenFacadeFails(t *testing.T) {
	authAsUser(t)
	old := getAsset
	t.Cleanup(func() { getAsset = old })
	getAsset = func(ctx facade.ContextWithUser, request dto4assetus.GetAssetRequest) (dto4assetus.GetAssetResponse, error) {
		return dto4assetus.GetAssetResponse{}, errBoom
	}

	req := httptest.NewRequest(http.MethodGet, "/v0/assetus/asset?spaceID=s-7&assetID=a-42", nil)
	w := httptest.NewRecorder()
	httpGetAsset(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body=%s", w.Code, w.Body.String())
	}
}

func TestHttpGetAsset_returns400WhenQueryParamsMissing(t *testing.T) {
	authAsUser(t)
	old := getAsset
	t.Cleanup(func() { getAsset = old })
	getAsset = func(ctx facade.ContextWithUser, request dto4assetus.GetAssetRequest) (dto4assetus.GetAssetResponse, error) {
		t.Fatalf("facade must not be reached when request validation fails")
		return dto4assetus.GetAssetResponse{}, nil
	}

	// Missing assetID fails GetAssetRequest.Validate().
	req := httptest.NewRequest(http.MethodGet, "/v0/assetus/asset?spaceID=s-7", nil)
	w := httptest.NewRecorder()
	httpGetAsset(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}

func TestHttpGetAsset_returns401WhenUnauthenticated(t *testing.T) {
	authRejected(t)
	req := httptest.NewRequest(http.MethodGet, "/v0/assetus/asset?spaceID=s&assetID=a", nil)
	w := httptest.NewRecorder()
	httpGetAsset(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

// =============================================================================
// httpPostUpdateAsset
// =============================================================================

func TestHttpPostUpdateAsset_returns200WithAssetIDOnSuccess(t *testing.T) {
	authAsUser(t)
	old := updateAsset
	t.Cleanup(func() { updateAsset = old })
	updateAsset = func(ctx facade.ContextWithUser, request dto4assetus.UpdateAssetRequest) (dto4assetus.UpdateAssetResponse, error) {
		return dto4assetus.UpdateAssetResponse{ID: "asset1"}, nil
	}

	w := httptest.NewRecorder()
	httpPostUpdateAsset(w, newPostRequest("/v0/assetus/update_asset", validUpdateAssetJSON))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "asset1") {
		t.Errorf("body %q does not contain asset id", w.Body.String())
	}
}

func TestHttpPostUpdateAsset_returns500WhenFacadeFails(t *testing.T) {
	authAsUser(t)
	old := updateAsset
	t.Cleanup(func() { updateAsset = old })
	updateAsset = func(ctx facade.ContextWithUser, request dto4assetus.UpdateAssetRequest) (dto4assetus.UpdateAssetResponse, error) {
		return dto4assetus.UpdateAssetResponse{}, errBoom
	}

	w := httptest.NewRecorder()
	httpPostUpdateAsset(w, newPostRequest("/v0/assetus/update_asset", validUpdateAssetJSON))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body=%s", w.Code, w.Body.String())
	}
}

func TestHttpPostUpdateAsset_returns400WhenBodyIsInvalidJSON(t *testing.T) {
	authAsUser(t)
	w := httptest.NewRecorder()
	httpPostUpdateAsset(w, newPostRequest("/v0/assetus/update_asset", "not json"))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}

// =============================================================================
// httpPostRemoveAsset
//
// The success path returns http.StatusNoContent (204) with a nil body, which is
// the ReturnJSON contract for a response with no payload (create/update/transfer
// return a body; remove has none).
// =============================================================================

func TestHttpPostRemoveAsset_returns204OnSuccess(t *testing.T) {
	authAsUser(t)
	old := removeAsset
	t.Cleanup(func() { removeAsset = old })
	removeAsset = func(ctx facade.ContextWithUser, request dto4assetus.RemoveAssetRequest) error {
		return nil
	}

	w := httptest.NewRecorder()
	httpPostRemoveAsset(w, newPostRequest("/v0/assetus/remove_asset", validRemoveAssetJSON))

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", w.Code, w.Body.String())
	}
	if w.Body.Len() != 0 {
		t.Errorf("expected empty body for 204, got %q", w.Body.String())
	}
}

func TestHttpPostRemoveAsset_returns500WhenFacadeFails(t *testing.T) {
	authAsUser(t)
	old := removeAsset
	t.Cleanup(func() { removeAsset = old })
	removeAsset = func(ctx facade.ContextWithUser, request dto4assetus.RemoveAssetRequest) error {
		return errBoom
	}

	w := httptest.NewRecorder()
	httpPostRemoveAsset(w, newPostRequest("/v0/assetus/remove_asset", validRemoveAssetJSON))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body=%s", w.Code, w.Body.String())
	}
}

func TestHttpPostRemoveAsset_returns400WhenBodyIsInvalidJSON(t *testing.T) {
	authAsUser(t)
	w := httptest.NewRecorder()
	httpPostRemoveAsset(w, newPostRequest("/v0/assetus/remove_asset", "not json"))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}

// =============================================================================
// httpPostTransferAsset
// =============================================================================

func TestHttpPostTransferAsset_returns200WithAssetIDOnSuccess(t *testing.T) {
	authAsUser(t)
	old := transferAsset
	t.Cleanup(func() { transferAsset = old })
	transferAsset = func(ctx facade.ContextWithUser, request dto4assetus.TransferAssetRequest) (dto4assetus.TransferAssetResponse, error) {
		return dto4assetus.TransferAssetResponse{ID: "asset1"}, nil
	}

	body := `{"spaceID":"space1","assetID":"asset1","toSpaceID":"space2"}`
	w := httptest.NewRecorder()
	httpPostTransferAsset(w, newPostRequest("/v0/assetus/transfer_asset", body))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "asset1") {
		t.Errorf("body %q does not contain asset id", w.Body.String())
	}
}

func TestHttpPostTransferAsset_returns500WhenFacadeFails(t *testing.T) {
	authAsUser(t)
	old := transferAsset
	t.Cleanup(func() { transferAsset = old })
	transferAsset = func(ctx facade.ContextWithUser, request dto4assetus.TransferAssetRequest) (dto4assetus.TransferAssetResponse, error) {
		return dto4assetus.TransferAssetResponse{}, errBoom
	}

	w := httptest.NewRecorder()
	httpPostTransferAsset(w, newPostRequest("/v0/assetus/transfer_asset", validTransferAssetJSON))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body=%s", w.Code, w.Body.String())
	}
}

func TestHttpPostTransferAsset_returns400WhenBodyIsInvalidJSON(t *testing.T) {
	authAsUser(t)
	w := httptest.NewRecorder()
	httpPostTransferAsset(w, newPostRequest("/v0/assetus/transfer_asset", "not json"))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}

// =============================================================================
// httpPostRecordHistoryEvent
//
// NOTE: like remove_asset, the success path passes a nil response with a
// non-204 status (StatusCreated). StatusCreated is neither StatusNoContent nor
// StatusOK, so ReturnJSON does NOT panic; it attempts to JSON-encode a nil
// response, which yields the literal body "null". We assert 201 here.
// =============================================================================

func TestHttpPostRecordHistoryEvent_returns201OnSuccess(t *testing.T) {
	authAsUser(t)
	old := recordHistoryEvent
	t.Cleanup(func() { recordHistoryEvent = old })
	recordHistoryEvent = func(ctx facade.ContextWithUser, request dto4assetus.RecordHistoryEventRequest) error {
		return nil
	}

	body := `{"spaceID":"space1","assetID":"asset1","type":"repaired"}`
	w := httptest.NewRecorder()
	httpPostRecordHistoryEvent(w, newPostRequest("/v0/assetus/record_history_event", body))

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", w.Code, w.Body.String())
	}
}

func TestHttpPostRecordHistoryEvent_returns500WhenFacadeFails(t *testing.T) {
	authAsUser(t)
	old := recordHistoryEvent
	t.Cleanup(func() { recordHistoryEvent = old })
	recordHistoryEvent = func(ctx facade.ContextWithUser, request dto4assetus.RecordHistoryEventRequest) error {
		return errBoom
	}

	w := httptest.NewRecorder()
	httpPostRecordHistoryEvent(w, newPostRequest("/v0/assetus/record_history_event", validRecordEventJSON))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body=%s", w.Code, w.Body.String())
	}
}

func TestHttpPostRecordHistoryEvent_returns400WhenBodyIsInvalidJSON(t *testing.T) {
	authAsUser(t)
	w := httptest.NewRecorder()
	httpPostRecordHistoryEvent(w, newPostRequest("/v0/assetus/record_history_event", "not json"))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}

// =============================================================================
// httpGetHistory
// =============================================================================

// httpGetHistory reads its input from the query string and authenticates via
// apicore.VerifyRequestAndCreateUserContext (no body decode), so a GET with
// valid query params + auth reaches the facade and returns 200 with the JSON
// history.
func TestHttpGetHistory_returns200WithHistoryOnSuccess(t *testing.T) {
	authAsUser(t)
	old := getHistory
	t.Cleanup(func() { getHistory = old })
	var gotRequest dto4assetus.GetHistoryRequest
	getHistory = func(ctx facade.ContextWithUser, request dto4assetus.GetHistoryRequest) (dto4assetus.GetHistoryResponse, error) {
		gotRequest = request
		return dto4assetus.GetHistoryResponse{
			AssetID: "a-9",
			Events:  []dto4assetus.HistoryEventItem{{ID: "ev-1"}},
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/v0/assetus/asset_history?spaceID=s-3&assetID=a-9", nil)
	w := httptest.NewRecorder()
	httpGetHistory(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "ev-1") {
		t.Errorf("body %q does not contain history event id", w.Body.String())
	}
	if gotRequest.AssetID != "a-9" || string(gotRequest.SpaceID) != "s-3" {
		t.Errorf("facade received request %+v, want assetID=a-9 spaceID=s-3", gotRequest)
	}
}

func TestHttpGetHistory_returns500WhenFacadeFails(t *testing.T) {
	authAsUser(t)
	old := getHistory
	t.Cleanup(func() { getHistory = old })
	getHistory = func(ctx facade.ContextWithUser, request dto4assetus.GetHistoryRequest) (dto4assetus.GetHistoryResponse, error) {
		return dto4assetus.GetHistoryResponse{}, errBoom
	}

	req := httptest.NewRequest(http.MethodGet, "/v0/assetus/asset_history?spaceID=s-3&assetID=a-9", nil)
	w := httptest.NewRecorder()
	httpGetHistory(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body=%s", w.Code, w.Body.String())
	}
}

func TestHttpGetHistory_returns400WhenQueryParamsMissing(t *testing.T) {
	authAsUser(t)
	old := getHistory
	t.Cleanup(func() { getHistory = old })
	getHistory = func(ctx facade.ContextWithUser, request dto4assetus.GetHistoryRequest) (dto4assetus.GetHistoryResponse, error) {
		t.Fatalf("facade must not be reached when request validation fails")
		return dto4assetus.GetHistoryResponse{}, nil
	}

	// Missing assetID fails GetHistoryRequest.Validate().
	req := httptest.NewRequest(http.MethodGet, "/v0/assetus/asset_history?spaceID=s-3", nil)
	w := httptest.NewRecorder()
	httpGetHistory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}

func TestHttpGetHistory_returns401WhenUnauthenticated(t *testing.T) {
	authRejected(t)
	req := httptest.NewRequest(http.MethodGet, "/v0/assetus/asset_history?spaceID=s&assetID=a", nil)
	w := httptest.NewRecorder()
	httpGetHistory(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

// =============================================================================
// httpPostCreateVehicleRecord
// =============================================================================

func TestHttpPostCreateVehicleRecord_returns201WithIDOnSuccess(t *testing.T) {
	authAsUser(t)
	old := addVehicleRecord
	t.Cleanup(func() { addVehicleRecord = old })
	addVehicleRecord = func(ctx facade.ContextWithUser, request dto4assetus.AddVehicleRecordRequest) (dto4assetus.AddVehicleRecordResponse, error) {
		return dto4assetus.AddVehicleRecordResponse{ID: "rec-1"}, nil
	}

	body := `{"spaceID":"space1","assetID":"asset1","mileage":1000,"mileageUnit":"km"}`
	w := httptest.NewRecorder()
	httpPostCreateVehicleRecord(w, newPostRequest("/v0/assetus/create_vehicle_record", body))

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "rec-1") {
		t.Errorf("body %q does not contain record id", w.Body.String())
	}
}

func TestHttpPostCreateVehicleRecord_returns500WhenFacadeFails(t *testing.T) {
	authAsUser(t)
	old := addVehicleRecord
	t.Cleanup(func() { addVehicleRecord = old })
	addVehicleRecord = func(ctx facade.ContextWithUser, request dto4assetus.AddVehicleRecordRequest) (dto4assetus.AddVehicleRecordResponse, error) {
		return dto4assetus.AddVehicleRecordResponse{}, errBoom
	}

	w := httptest.NewRecorder()
	httpPostCreateVehicleRecord(w, newPostRequest("/v0/assetus/create_vehicle_record", validVehicleRecJSON))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body=%s", w.Code, w.Body.String())
	}
}

func TestHttpPostCreateVehicleRecord_returns400WhenBodyIsInvalidJSON(t *testing.T) {
	authAsUser(t)
	w := httptest.NewRecorder()
	httpPostCreateVehicleRecord(w, newPostRequest("/v0/assetus/create_vehicle_record", "not json"))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}
