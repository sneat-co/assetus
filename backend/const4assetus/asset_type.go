package const4assetus

import (
	"fmt"
	"slices"

	"github.com/strongo/validation"
)

// Type is the optional subtype of an asset within its Category (the second
// level of the two-level taxonomy ported from the legacy model). Each Category
// permits its own set of Types; an out-of-category Type is rejected by
// ValidateType. A Category with no defined subtype set (e.g. books) accepts no
// Type — assets in those categories carry no subtype.
type Type string

// Vehicle subtypes (legacy AssetCategoryVehicle subtypes), mapped under the
// unified CategoryVehicles.
const (
	TypeVehicleAircraft   Type = "aircraft"
	TypeVehicleBoat       Type = "boat"
	TypeVehicleBus        Type = "bus"
	TypeVehicleCar        Type = "car"
	TypeVehicleHelicopter Type = "helicopter"
	TypeVehicleMotorcycle Type = "motorcycle"
	TypeVehicleTruck      Type = "truck"
	TypeVehicleVan        Type = "van"
)

// Dwelling / real-estate subtypes (legacy AssetCategoryDwelling subtypes).
const (
	TypeDwellingApartment Type = "apartment"
	TypeDwellingHouse     Type = "house"
	TypeDwellingOffice    Type = "office"
	TypeDwellingShop      Type = "shop"
	TypeDwellingLand      Type = "land"
	TypeDwellingGarage    Type = "garage"
	TypeDwellingWarehouse Type = "warehouse"
)

// Sports-equipment subtypes (legacy AssetCategorySportGear subtypes), mapped
// under the unified CategorySportsEquipment.
const (
	TypeSportsBicycle        Type = "bicycle"
	TypeSportsKite           Type = "kite"
	TypeSportsKiteBar        Type = "kite_bar"
	TypeSportsKiteBoard      Type = "kite_board"
	TypeSportsKiteHydrofoil  Type = "kite_hydrofoil"
	TypeSportsProneHydrofoil Type = "prone_hydrofoil"
	TypeSportsSurfBoard      Type = "surf_board"
	TypeSportsWetsuit        Type = "wetsuit"
	TypeSportsWing           Type = "wing"
	TypeSportsWingBoard      Type = "wing_board"
	TypeSportsWingHydrofoil  Type = "wing_hydrofoil"
)

// Document subtypes (AssetDocumentType, legacy AssetCategoryDocument subtypes).
const (
	TypeDocumentPassport       Type = "passport"
	TypeDocumentIDCard         Type = "id_card"
	TypeDocumentDrivingLicense Type = "driving_license"
	TypeDocumentMarriageCert   Type = "marriage_cert"
	TypeDocumentBirthCert      Type = "birth_cert"
)

// VehicleTypes is the closed subtype set for CategoryVehicles.
var VehicleTypes = []Type{
	TypeVehicleCar,
	TypeVehicleBus,
	TypeVehicleVan,
	TypeVehicleTruck,
	TypeVehicleMotorcycle,
	TypeVehicleBoat,
	TypeVehicleAircraft,
	TypeVehicleHelicopter,
}

// DwellingTypes is the closed subtype set for CategoryDwelling.
var DwellingTypes = []Type{
	TypeDwellingApartment,
	TypeDwellingHouse,
	TypeDwellingOffice,
	TypeDwellingShop,
	TypeDwellingLand,
	TypeDwellingGarage,
	TypeDwellingWarehouse,
}

// SportsEquipmentTypes is the closed subtype set for CategorySportsEquipment.
var SportsEquipmentTypes = []Type{
	TypeSportsBicycle,
	TypeSportsKite,
	TypeSportsKiteBar,
	TypeSportsKiteBoard,
	TypeSportsKiteHydrofoil,
	TypeSportsProneHydrofoil,
	TypeSportsSurfBoard,
	TypeSportsWetsuit,
	TypeSportsWing,
	TypeSportsWingBoard,
	TypeSportsWingHydrofoil,
}

// DocumentTypes is the closed subtype set for CategoryDocument.
var DocumentTypes = []Type{
	TypeDocumentPassport,
	TypeDocumentIDCard,
	TypeDocumentDrivingLicense,
	TypeDocumentMarriageCert,
	TypeDocumentBirthCert,
}

// typesByCategory maps each Category that has a two-level taxonomy to its
// permitted subtype set. Categories absent from this map have no subtype: a
// non-empty Type is rejected for them (preserving the legacy per-category Type
// validation).
var typesByCategory = map[Category][]Type{
	CategoryVehicles:        VehicleTypes,
	CategoryDwelling:        DwellingTypes,
	CategorySportsEquipment: SportsEquipmentTypes,
	CategoryDocument:        DocumentTypes,
}

// TypesForCategory returns the permitted subtype set for a category, or nil if
// the category has no subtypes.
func TypesForCategory(category Category) []Type {
	return typesByCategory[category]
}

// ValidateType validates an asset's Type against its Category, preserving the
// legacy per-category Type validation: a Type valid for one category (e.g.
// kite_board for sports_equipment) is rejected for another (e.g. document).
//
// Rules:
//   - The category must itself be valid.
//   - If the category has a subtype set, the Type must be a member of it. An
//     empty Type is allowed (subtype is optional).
//   - If the category has no subtype set, a non-empty Type is rejected.
func ValidateType(category Category, t Type) error {
	if err := ValidateCategory(category); err != nil {
		return err
	}
	permitted, hasSubtypes := typesByCategory[category]
	if !hasSubtypes {
		if t != "" {
			return validation.NewErrBadRecordFieldValue("type",
				fmt.Sprintf("category %q does not permit a subtype, got %q", category, t))
		}
		return nil
	}
	if t == "" {
		return nil
	}
	if !slices.Contains(permitted, t) {
		return validation.NewErrBadRecordFieldValue("type",
			fmt.Sprintf("%q is not a valid type for category %q, expected one of %v", t, category, permitted))
	}
	return nil
}
