package dal4assetus

import (
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// AssetusSpaceEntry is the loaded assetus module entry for a Space.
type AssetusSpaceEntry = record.DataWithID[coretypes.ExtID, *dbo4assetus.AssetusSpaceDbo]
