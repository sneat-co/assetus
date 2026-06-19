package const4assetus

import (
	"fmt"
	"slices"

	"github.com/strongo/validation"
)

// Status is the closed, write-validated ownership-lifecycle status of an asset.
//
// This set is ownership-only by design. Sharing/availability states such as
// "borrowed" or "reserved" are NOT Assetus statuses — they are owned by the
// future Yardius sharing layer. ValidateStatus rejects any value outside this
// set at the write boundary.
type Status string

const (
	// StatusDraft is the pre-active lifecycle state ported from the legacy
	// AssetStatus set ({active, archived, draft}). An asset starts as a draft
	// before it is activated.
	StatusDraft       Status = "draft"
	StatusActive      Status = "active"
	StatusTransferred Status = "transferred"
	StatusArchived    Status = "archived"
	StatusDisposed    Status = "disposed"
	StatusLost        Status = "lost"
)

// Statuses is the closed set of valid ownership-lifecycle statuses. It is the
// exact union of the legacy AssetStatus set {active, archived, draft} and the
// MVP Status set {active, transferred, archived, disposed, lost}, i.e.
// {draft, active, transferred, archived, disposed, lost}.
var Statuses = []Status{
	StatusDraft,
	StatusActive,
	StatusTransferred,
	StatusArchived,
	StatusDisposed,
	StatusLost,
}

// IsValidStatus reports whether v is a member of the closed ownership-status
// set. Sharing/availability values (e.g. "borrowed", "reserved") return false.
func IsValidStatus(v Status) bool {
	return slices.Contains(Statuses, v)
}

// ValidateStatus returns an error if v is not a valid Assetus ownership status.
// A status is required on every asset.
func ValidateStatus(v Status) error {
	if v == "" {
		return validation.NewErrRecordIsMissingRequiredField("status")
	}
	if !IsValidStatus(v) {
		return validation.NewErrBadRecordFieldValue("status",
			fmt.Sprintf("%q is not a valid Assetus status; valid set: %v "+
				"(sharing/availability states like borrowed/reserved are owned by Yardius, not Assetus)", v, Statuses))
	}
	return nil
}
