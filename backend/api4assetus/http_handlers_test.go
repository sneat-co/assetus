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

// NOTE: httpGetAsset is registered as a GET endpoint but verifies/decodes the
// request with apicore.VerifyAuthenticatedRequestAndDecodeBody. In
// sneat-go-core v0.55.4, apicore.DecodeRequestBody only supports POST/PUT/DELETE
// and returns 405 Method Not Allowed for GET. The facade (getAsset) is
// therefore never reached for a GET request — a genuine production defect: the
// endpoint can never return a successful asset read over its registered method.
// We assert the real 405 behaviour to document the defect; the query-param
// parsing above the verify call is still executed for coverage.
func TestHttpGetAsset_returns405BecauseDecodeRejectsGetMethod(t *testing.T) {
	authAsUser(t)
	old := getAsset
	t.Cleanup(func() { getAsset = old })
	getAsset = func(ctx facade.ContextWithUser, request dto4assetus.GetAssetRequest) (dto4assetus.GetAssetResponse, error) {
		t.Fatalf("facade must not be reached: GET is rejected by DecodeRequestBody")
		return dto4assetus.GetAssetResponse{}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/v0/assetus/asset?spaceID=s-7&assetID=a-42", nil)
	w := httptest.NewRecorder()
	httpGetAsset(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405; body=%s", w.Code, w.Body.String())
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
// NOTE: the success path calls apicore.ReturnJSON(..., http.StatusOK, nil, nil).
// In sneat-go-core v0.55.4 ReturnJSON panics when response is nil and the
// status is StatusOK ("expected to be http.StatusNoContent=204"). The handler
// therefore cannot return a successful response without panicking — a genuine
// production defect surfaced by these tests. We assert the panic to both
// document the defect and exercise the success branch for coverage.
// =============================================================================

func TestHttpPostRemoveAsset_successPathPanicsDueToNilResponseWithStatusOK(t *testing.T) {
	authAsUser(t)
	old := removeAsset
	t.Cleanup(func() { removeAsset = old })
	removeAsset = func(ctx facade.ContextWithUser, request dto4assetus.RemoveAssetRequest) error {
		return nil
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic from ReturnJSON(StatusOK, nil response), got none")
		}
		if !strings.Contains(strings.ToLower(asString(r)), "statusok") &&
			!strings.Contains(asString(r), "204") {
			t.Errorf("unexpected panic value: %v", r)
		}
	}()

	w := httptest.NewRecorder()
	httpPostRemoveAsset(w, newPostRequest("/v0/assetus/remove_asset", validRemoveAssetJSON))
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

// NOTE: same GET-vs-DecodeRequestBody defect as httpGetAsset — DecodeRequestBody
// returns 405 for GET, so getHistory is never reached over the registered GET
// method. We assert the real 405 behaviour.
func TestHttpGetHistory_returns405BecauseDecodeRejectsGetMethod(t *testing.T) {
	authAsUser(t)
	old := getHistory
	t.Cleanup(func() { getHistory = old })
	getHistory = func(ctx facade.ContextWithUser, request dto4assetus.GetHistoryRequest) (dto4assetus.GetHistoryResponse, error) {
		t.Fatalf("facade must not be reached: GET is rejected by DecodeRequestBody")
		return dto4assetus.GetHistoryResponse{}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/v0/assetus/asset_history?spaceID=s-3&assetID=a-9", nil)
	w := httptest.NewRecorder()
	httpGetHistory(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405; body=%s", w.Code, w.Body.String())
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

func asString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	if e, ok := v.(error); ok {
		return e.Error()
	}
	return ""
}
