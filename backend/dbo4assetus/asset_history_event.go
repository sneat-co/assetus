package dbo4assetus

import (
	"strings"
	"time"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/strongo/validation"
)

// HistoryCollection is the Firestore child collection name for an asset's
// append-only history: /spaces/{spaceID}/ext/assetus/assets/{assetID}/history/{eventID}.
const HistoryCollection = "history"

// AssetHistoryEventBase carries the fields of a history event. Embedded by
// AssetHistoryEventDbo (the persisted, append-only record).
type AssetHistoryEventBase struct {
	Type       const4assetus.HistoryEventType `json:"type" firestore:"type"`
	OccurredAt time.Time                      `json:"occurredAt" firestore:"occurredAt"`
	// ActorRef identifies the acting member (the user who caused the event).
	ActorRef string `json:"actorRef" firestore:"actorRef"`
	Note     string `json:"note,omitempty" firestore:"note,omitempty"`

	// FromOwner / ToOwner are populated only for Transferred events, recording
	// the prior and new owner so the prior owner is preserved, never overwritten.
	FromOwner *OwnerRef `json:"fromOwner,omitempty" firestore:"fromOwner,omitempty"`
	ToOwner   *OwnerRef `json:"toOwner,omitempty" firestore:"toOwner,omitempty"`
}

// Validate returns an error if the event is not valid.
func (v AssetHistoryEventBase) Validate() error {
	if err := const4assetus.ValidateHistoryEventType(v.Type); err != nil {
		return err
	}
	if v.OccurredAt.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("occurredAt")
	}
	if strings.TrimSpace(v.ActorRef) == "" {
		return validation.NewErrRecordIsMissingRequiredField("actorRef")
	}
	if v.Type == const4assetus.HistoryEventTransferred {
		if v.FromOwner == nil {
			return validation.NewErrRecordIsMissingRequiredField("fromOwner")
		}
		if v.ToOwner == nil {
			return validation.NewErrRecordIsMissingRequiredField("toOwner")
		}
	}
	return nil
}

// AssetHistoryEventDbo is a persisted, append-only history event.
type AssetHistoryEventDbo struct {
	AssetHistoryEventBase
}

// Validate returns an error if the record is not valid.
func (v *AssetHistoryEventDbo) Validate() error {
	return v.AssetHistoryEventBase.Validate()
}
