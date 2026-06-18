package const4assetus

import (
	"fmt"
	"slices"

	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/strongo/validation"
)

// Visibility is the closed, write-validated visibility of an asset. It defaults
// to the owning Space's default visibility (see DefaultVisibilityForSpaceType)
// and may be overridden per asset.
type Visibility string

const (
	VisibilityPrivate          Visibility = "private"
	VisibilityFamily           Visibility = "family"
	VisibilityFriends          Visibility = "friends"
	VisibilityFriendsOfFriends Visibility = "friends_of_friends"
	VisibilitySpecificSpace    Visibility = "specific_space"
	VisibilityPublic           Visibility = "public"
)

// Visibilities is the closed set of valid visibilities.
var Visibilities = []Visibility{
	VisibilityPrivate,
	VisibilityFamily,
	VisibilityFriends,
	VisibilityFriendsOfFriends,
	VisibilitySpecificSpace,
	VisibilityPublic,
}

// IsValidVisibility reports whether v is a member of the closed visibility set.
func IsValidVisibility(v Visibility) bool {
	return slices.Contains(Visibilities, v)
}

// ValidateVisibility returns an error if v is not a valid visibility. A
// visibility is required on every asset.
func ValidateVisibility(v Visibility) error {
	if v == "" {
		return validation.NewErrRecordIsMissingRequiredField("visibility")
	}
	if !IsValidVisibility(v) {
		return validation.NewErrBadRecordFieldValue("visibility",
			fmt.Sprintf("unknown visibility %q, expected one of %v", v, Visibilities))
	}
	return nil
}

// DefaultVisibilityForSpaceType returns the default visibility an asset inherits
// from its owning Space when no per-asset override is provided.
//
// NOTE: spaceus does not (yet) expose a configurable per-space default
// visibility field, so the default is derived from the Space type. A `family`
// Space defaults to Family visibility; a `private` Space to Private; broader
// owner types (club/company/group) default to Specific Space. When spaceus adds
// a stored per-space default-visibility field, this derivation should defer to
// it.
func DefaultVisibilityForSpaceType(spaceType coretypes.SpaceType) Visibility {
	switch spaceType {
	case coretypes.SpaceTypePrivate:
		return VisibilityPrivate
	case coretypes.SpaceTypeFamily:
		return VisibilityFamily
	default:
		// club, company, group, space and any future community/school types.
		return VisibilitySpecificSpace
	}
}
