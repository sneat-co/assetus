package assetusext

import (
	"testing"

	"github.com/sneat-co/assetus/backend/const4assetus"
)

// TestExtension_ReportsAssetusID is a smoke test that the assembled extension
// config exposes the assetus extension ID, proving the module wires up.
func TestExtension_ReportsAssetusID(t *testing.T) {
	cfg := Extension()
	if cfg == nil {
		t.Fatal("Extension() returned nil")
	}
	if got := cfg.ID(); got != const4assetus.ExtensionID {
		t.Errorf("Extension().ID() = %q, want %q", got, const4assetus.ExtensionID)
	}
}
