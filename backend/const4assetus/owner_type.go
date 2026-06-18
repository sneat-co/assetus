package const4assetus

import (
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// OwnerType is the type of an asset's owner, derived from the owning Space's
// type. There is no separate owner entity in the MVP — owner == the owning
// Space — so OwnerType is always derived, never stored as a source of truth.
type OwnerType string

const (
	OwnerTypeIndividual   OwnerType = "individual"
	OwnerTypeFamily       OwnerType = "family"
	OwnerTypeSportsClub   OwnerType = "sports_club"
	OwnerTypeCommunity    OwnerType = "community"
	OwnerTypeSchool       OwnerType = "school"
	OwnerTypeOrganisation OwnerType = "organisation"
	// OwnerTypeUnknown is returned for space types that do not map to a known
	// owner type.
	OwnerTypeUnknown OwnerType = ""
)

// New spaceus space types required for Community and School owner types. These
// are NOT yet present in coretypes (a spaceus precondition owned externally);
// declared here as string literals so DeriveOwnerType maps them the moment
// spaceus ships them. See backstage assetus-mvp Open Questions.
const (
	spaceTypeCommunity coretypes.SpaceType = "community"
	spaceTypeSchool    coretypes.SpaceType = "school"
)

// DeriveOwnerType maps a Space type to the asset owner type:
//
//	private   -> Individual
//	family    -> Family
//	club      -> SportsClub
//	company   -> Organisation
//	community -> Community   (once the spaceus precondition ships)
//	school    -> School      (once the spaceus precondition ships)
//
// group/space and any unknown type return OwnerTypeUnknown.
func DeriveOwnerType(spaceType coretypes.SpaceType) OwnerType {
	switch spaceType {
	case coretypes.SpaceTypePrivate:
		return OwnerTypeIndividual
	case coretypes.SpaceTypeFamily:
		return OwnerTypeFamily
	case coretypes.SpaceTypeClub:
		return OwnerTypeSportsClub
	case coretypes.SpaceTypeCompany:
		return OwnerTypeOrganisation
	case spaceTypeCommunity:
		return OwnerTypeCommunity
	case spaceTypeSchool:
		return OwnerTypeSchool
	default:
		return OwnerTypeUnknown
	}
}
