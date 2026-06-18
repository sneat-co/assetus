package dal4assetus

import (
	"fmt"
	"testing"

	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func TestNewAssetKey_Path(t *testing.T) {
	key := NewAssetKey("space1", "asset1")
	if key.ID != "asset1" {
		t.Errorf("key.ID = %v, want asset1", key.ID)
	}
	if key.Collection() != dbo4assetus.AssetsCollection {
		t.Errorf("key.Collection() = %v, want %v", key.Collection(), dbo4assetus.AssetsCollection)
	}
	// Parent chain must be /spaces/space1/ext/assetus.
	parent := key.Parent()
	if parent == nil {
		t.Fatal("asset key has no parent (expected the assetus module key)")
	}
	if got := fmt.Sprintf("%v", parent.ID); got != "assetus" {
		t.Errorf("module key ID = %v, want assetus", got)
	}
}

func TestNewAssetEntry(t *testing.T) {
	var spaceID coretypes.SpaceID = "space1"
	entry := NewAssetEntry(spaceID, "asset1")
	if entry.ID != "asset1" {
		t.Errorf("entry.ID = %v, want asset1", entry.ID)
	}
	if entry.Data == nil {
		t.Error("entry.Data is nil")
	}
	if entry.Record == nil {
		t.Error("entry.Record is nil")
	}
}
