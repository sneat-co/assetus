package const4assetus

// This file documents the legacy → unified mapping for every legacy enum value
// ported into the unified const4assetus model. The maps below are exhaustive:
// every legacy AssetCategory, AssetStatus, AssetPossession, AssetType,
// AssetDocumentType, EngineType, and FuelType value has exactly one unified
// counterpart. No legacy value is dropped.
//
// Legacy values are kept as untyped string keys because the legacy package
// declared its enums as string type aliases.

// LegacyCategoryToUnified maps every legacy AssetCategory string to its unified
// Category. Notable reconciliations:
//   - legacy "vehicle"    -> CategoryVehicles  (singular -> canonical plural)
//   - legacy "sport_gear" -> CategorySportsEquipment
//   - legacy "misc"       -> CategoryOther
//   - legacy "dwelling"/"document"/"debt"/"undefined" carry over directly.
var LegacyCategoryToUnified = map[string]Category{
	"vehicle":    CategoryVehicles,
	"dwelling":   CategoryDwelling,
	"sport_gear": CategorySportsEquipment,
	"document":   CategoryDocument,
	"misc":       CategoryOther,
	"debt":       CategoryDebt,
	"undefined":  CategoryOther,
}

// LegacyStatusToUnified maps every legacy AssetStatus string to its unified
// Status. The legacy set {active, archived, draft} is a subset of the unified
// union {draft, active, transferred, archived, disposed, lost}.
var LegacyStatusToUnified = map[string]Status{
	"active":   StatusActive,
	"archived": StatusArchived,
	"draft":    StatusDraft,
}

// LegacyPossessionToUnified maps every legacy AssetPossession string to its
// unified Possession. The value sets are identical; the map makes the totality
// explicit.
var LegacyPossessionToUnified = map[string]Possession{
	"unknown":     PossessionUnknown,
	"undisclosed": PossessionUndisclosed,
	"owning":      PossessionOwning,
	"leasing":     PossessionLeasing,
	"renting":     PossessionRenting,
}

// LegacyTypeToUnified maps every legacy AssetType and AssetDocumentType string
// (across all categories: vehicle, real-estate/dwelling, sport_gear, document)
// to its unified Type. Subtype string values are preserved verbatim
// (e.g. kite_board stays kite_board).
var LegacyTypeToUnified = map[string]Type{
	// Vehicle
	"aircraft":   TypeVehicleAircraft,
	"boat":       TypeVehicleBoat,
	"bus":        TypeVehicleBus,
	"car":        TypeVehicleCar,
	"helicopter": TypeVehicleHelicopter,
	"motorcycle": TypeVehicleMotorcycle,
	"truck":      TypeVehicleTruck,
	"van":        TypeVehicleVan,
	// Real estate / dwelling
	"apartment": TypeDwellingApartment,
	"house":     TypeDwellingHouse,
	"office":    TypeDwellingOffice,
	"shop":      TypeDwellingShop,
	"land":      TypeDwellingLand,
	"garage":    TypeDwellingGarage,
	"warehouse": TypeDwellingWarehouse,
	// Sport gear -> sports equipment
	"bicycle":         TypeSportsBicycle,
	"kite":            TypeSportsKite,
	"kite_bar":        TypeSportsKiteBar,
	"kite_board":      TypeSportsKiteBoard,
	"kite_hydrofoil":  TypeSportsKiteHydrofoil,
	"prone_hydrofoil": TypeSportsProneHydrofoil,
	"surf_board":      TypeSportsSurfBoard,
	"wetsuit":         TypeSportsWetsuit,
	"wing":            TypeSportsWing,
	"wing_board":      TypeSportsWingBoard,
	"wing_hydrofoil":  TypeSportsWingHydrofoil,
	// Document
	"passport":        TypeDocumentPassport,
	"id_card":         TypeDocumentIDCard,
	"driving_license": TypeDocumentDrivingLicense,
	"marriage_cert":   TypeDocumentMarriageCert,
	"birth_cert":      TypeDocumentBirthCert,
}

// LegacyEngineTypeToUnified maps every legacy EngineType string to its unified
// EngineType (value sets are identical).
var LegacyEngineTypeToUnified = map[string]EngineType{
	"":           EngineTypeUnknown,
	"other":      EngineTypeOther,
	"combustion": EngineTypeCombustion,
	"electric":   EngineTypeElectric,
	"phev":       EngineTypePHEV,
	"hybrid":     EngineTypeHybrid,
	"steam":      EngineTypeSteam,
}

// LegacyFuelTypeToUnified maps every legacy FuelType string to its unified
// FuelType (value sets are identical).
var LegacyFuelTypeToUnified = map[string]FuelType{
	"":         FuelTypeUnknown,
	"other":    FuelTypeOther,
	"bio":      FuelTypeBio,
	"petrol":   FuelTypePetrol,
	"diesel":   FuelTypeDiesel,
	"hydrogen": FuelTypeHydrogen,
}
