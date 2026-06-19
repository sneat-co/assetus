package const4assetus

import (
	"fmt"
	"slices"

	"github.com/strongo/validation"
)

// EngineType is the type of a vehicle's engine, ported from the legacy
// extras4assetus.EngineType.
type EngineType string

const (
	EngineTypeUnknown    EngineType = ""
	EngineTypeOther      EngineType = "other"
	EngineTypeCombustion EngineType = "combustion"
	EngineTypeElectric   EngineType = "electric"
	EngineTypePHEV       EngineType = "phev"
	EngineTypeHybrid     EngineType = "hybrid"
	EngineTypeSteam      EngineType = "steam"
)

// EngineTypes is the list of known engine types.
var EngineTypes = []EngineType{
	EngineTypeUnknown,
	EngineTypeOther,
	EngineTypeCombustion,
	EngineTypeElectric,
	EngineTypePHEV,
	EngineTypeHybrid,
	EngineTypeSteam,
}

// FuelType is the type of fuel an engine consumes, ported from the legacy
// extras4assetus.FuelType.
type FuelType string

const (
	FuelTypeUnknown  FuelType = ""
	FuelTypeOther    FuelType = "other"
	FuelTypeBio      FuelType = "bio"
	FuelTypePetrol   FuelType = "petrol"
	FuelTypeDiesel   FuelType = "diesel"
	FuelTypeHydrogen FuelType = "hydrogen"
)

// FuelTypes is the list of known fuel types.
var FuelTypes = []FuelType{
	FuelTypeUnknown,
	FuelTypeOther,
	FuelTypePetrol,
	FuelTypeDiesel,
	FuelTypeHydrogen,
	FuelTypeBio,
}

// IsKnownFuelType reports whether v is a known fuel type.
func IsKnownFuelType(v FuelType) bool {
	return slices.Contains(FuelTypes, v)
}

// ValidateEngineFuel validates the compatibility of an EngineType with a
// FuelType, preserving the legacy engine↔fuel compatibility matrix:
//
//   - unknown/other engine: any known fuel type.
//   - combustion/hybrid/phev: petrol, diesel, hydrogen, unknown or other.
//   - electric: unknown or other only.
//   - steam: unknown or other only.
//
// Any other engine type is rejected.
func ValidateEngineFuel(engineType EngineType, fuel FuelType) error {
	switch engineType {
	case EngineTypeUnknown, EngineTypeOther:
		if !IsKnownFuelType(fuel) {
			return validation.NewErrBadRecordFieldValue("fuelType", fmt.Sprintf("unknown fuel type: %s", fuel))
		}
	case EngineTypeCombustion, EngineTypeHybrid, EngineTypePHEV:
		switch fuel {
		case FuelTypePetrol, FuelTypeDiesel, FuelTypeHydrogen, FuelTypeUnknown, FuelTypeOther:
			// OK
		default:
			return validation.NewErrBadRecordFieldValue("fuelType", fmt.Sprintf("unknown fuel type: %s", fuel))
		}
	case EngineTypeElectric:
		switch fuel {
		case FuelTypeUnknown, FuelTypeOther:
			// OK
		default:
			return validation.NewErrBadRecordFieldValue("fuelType", fmt.Sprintf("unknown fuel type: %s", fuel))
		}
	case EngineTypeSteam:
		switch fuel {
		case FuelTypeUnknown, FuelTypeOther:
			// OK
		default:
			return validation.NewErrBadRecordFieldValue("fuelType", fmt.Sprintf("unknown fuel type: %s", fuel))
		}
	default:
		return validation.NewErrBadRecordFieldValue("engineType", "unknown engine type: "+string(engineType))
	}
	return nil
}
