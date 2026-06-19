package extras4assetus

import (
	"github.com/sneat-co/assetus/backend/const4assetus"
)

// WithEngineData carries a vehicle's engine data, ported from the legacy
// extras4assetus.WithEngineData. The EngineSerialNumber field is sourced from
// the legacy frontend DTO (IEngine.engineSerialNumber). Engine type / fuel
// validation reuses const4assetus.ValidateEngineFuel.
type WithEngineData struct {
	EngineType         const4assetus.EngineType `json:"engineType,omitempty" firestore:"engineType,omitempty"`
	EngineFuel         const4assetus.FuelType   `json:"engineFuel,omitempty" firestore:"engineFuel,omitempty"`
	EngineCC           int                      `json:"engineCC,omitempty" firestore:"engineCC,omitempty"`                     // Engine volume in cubic centimetres
	EngineKW           int                      `json:"engineKW,omitempty" firestore:"engineKW,omitempty"`                     // Engine power in kilowatts
	EngineNM           int                      `json:"engineNM,omitempty" firestore:"engineNM,omitempty"`                     // Engine torque in Newton metres
	EngineSerialNumber string                   `json:"engineSerialNumber,omitempty" firestore:"engineSerialNumber,omitempty"` // Engine serial number
}

// Validate returns an error if the engine type / fuel combination is invalid.
func (v WithEngineData) Validate() error {
	return const4assetus.ValidateEngineFuel(v.EngineType, v.EngineFuel)
}
