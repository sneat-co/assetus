package api4assetus

import (
	"github.com/sneat-co/sneat-go-core/extension"
)

// RegisterHttpRoutes registers assetus HTTP routes. Capability handlers
// (create/manage/remove/transfer/get) are added here as they land.
func RegisterHttpRoutes(handle extension.HTTPHandleFunc) {
	_ = handle
}
