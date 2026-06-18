package dto4assetus

import (
	"strings"
	"time"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/strongo/validation"
)

// RecordHistoryEventRequest appends a lifecycle event to an asset's append-only
// history (e.g. a Repaired event).
type RecordHistoryEventRequest struct {
	dto4spaceus.SpaceRequest
	AssetID    string                         `json:"assetID"`
	Type       const4assetus.HistoryEventType `json:"type"`
	OccurredAt *time.Time                     `json:"occurredAt,omitempty"`
	Note       string                         `json:"note,omitempty"`
}

// Validate validates the request. Transfer events are recorded by the transfer
// facade, not via this generic endpoint.
func (v RecordHistoryEventRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.AssetID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("assetID")
	}
	if err := const4assetus.ValidateHistoryEventType(v.Type); err != nil {
		return err
	}
	if v.Type == const4assetus.HistoryEventTransferred {
		return validation.NewErrBadRequestFieldValue("type", "transfer events are recorded by the transfer operation, not this endpoint")
	}
	return nil
}

// HistoryEventItem is a single history event in a read response.
type HistoryEventItem struct {
	ID         string                         `json:"id"`
	Type       const4assetus.HistoryEventType `json:"type"`
	OccurredAt time.Time                      `json:"occurredAt"`
	ActorRef   string                         `json:"actorRef"`
	Note       string                         `json:"note,omitempty"`
	FromOwner  *OwnerRefDTO                   `json:"fromOwner,omitempty"`
	ToOwner    *OwnerRefDTO                   `json:"toOwner,omitempty"`
}

// OwnerRefDTO mirrors dbo4assetus.OwnerRef for transport.
type OwnerRefDTO struct {
	SpaceID   string `json:"spaceID"`
	SpaceType string `json:"spaceType,omitempty"`
	OwnerType string `json:"ownerType,omitempty"`
}

// GetHistoryRequest reads an asset's history.
type GetHistoryRequest struct {
	dto4spaceus.SpaceRequest
	AssetID string `json:"assetID"`
}

// Validate validates the request.
func (v GetHistoryRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.AssetID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("assetID")
	}
	return nil
}

// GetHistoryResponse returns an asset's append-only history, ordered oldest-first.
type GetHistoryResponse struct {
	AssetID string             `json:"assetID"`
	Events  []HistoryEventItem `json:"events"`
}
