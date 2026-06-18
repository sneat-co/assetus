package api4assetus

import (
	"net/http"

	"github.com/sneat-co/sneat-go-core/extension"
)

// RegisterHttpRoutes registers assetus HTTP routes. Capability handlers
// (create/manage/remove/transfer/get) are added here as they land.
func RegisterHttpRoutes(handle extension.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/assetus/create_asset", httpPostCreateAsset)
	handle(http.MethodGet, "/v0/assetus/asset", httpGetAsset)
	handle(http.MethodPost, "/v0/assetus/update_asset", httpPostUpdateAsset)
	handle(http.MethodPost, "/v0/assetus/record_history_event", httpPostRecordHistoryEvent)
	handle(http.MethodGet, "/v0/assetus/asset_history", httpGetHistory)
}
