package facade4assetus

import (
	"github.com/strongo/random"
)

// assetIDLength is the length of a generated asset ID.
const assetIDLength = 6

// newAssetID generates a random asset ID.
func newAssetID() string {
	return random.ID(assetIDLength)
}

// newHistoryEventID generates a random history-event ID.
func newHistoryEventID() string {
	return random.ID(assetIDLength)
}
