package dbo4assetus

import (
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// OwnerRef is a read-time projection of an asset's owner. The owner is always
// the owning Space (owner == Space); OwnerType is derived from the Space type
// and is never a stored source of truth.
type OwnerRef struct {
	SpaceID   coretypes.SpaceID       `json:"spaceID" firestore:"spaceID"`
	SpaceType coretypes.SpaceType     `json:"spaceType,omitempty" firestore:"spaceType,omitempty"`
	OwnerType const4assetus.OwnerType `json:"ownerType,omitempty" firestore:"ownerType,omitempty"`
}

// NewOwnerRef builds an OwnerRef for a space, deriving the owner type from the
// space type.
func NewOwnerRef(spaceID coretypes.SpaceID, spaceType coretypes.SpaceType) OwnerRef {
	return OwnerRef{
		SpaceID:   spaceID,
		SpaceType: spaceType,
		OwnerType: const4assetus.DeriveOwnerType(spaceType),
	}
}
