package const4assetus

import (
	"fmt"
	"slices"

	"github.com/strongo/validation"
)

// HistoryEventType is the closed, write-validated type of an append-only asset
// history event.
type HistoryEventType string

const (
	HistoryEventPurchased   HistoryEventType = "purchased"
	HistoryEventRepaired    HistoryEventType = "repaired"
	HistoryEventTransferred HistoryEventType = "transferred"
	HistoryEventSold        HistoryEventType = "sold"
	HistoryEventDonated     HistoryEventType = "donated"
	HistoryEventLost        HistoryEventType = "lost"
)

// HistoryEventTypes is the closed set of valid history event types.
var HistoryEventTypes = []HistoryEventType{
	HistoryEventPurchased,
	HistoryEventRepaired,
	HistoryEventTransferred,
	HistoryEventSold,
	HistoryEventDonated,
	HistoryEventLost,
}

// IsValidHistoryEventType reports whether v is a member of the closed set.
func IsValidHistoryEventType(v HistoryEventType) bool {
	return slices.Contains(HistoryEventTypes, v)
}

// ValidateHistoryEventType returns an error if v is not a valid event type.
func ValidateHistoryEventType(v HistoryEventType) error {
	if v == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if !IsValidHistoryEventType(v) {
		return validation.NewErrBadRecordFieldValue("type",
			fmt.Sprintf("unknown history event type %q, expected one of %v", v, HistoryEventTypes))
	}
	return nil
}
