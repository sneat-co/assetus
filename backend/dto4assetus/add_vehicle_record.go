package dto4assetus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

// AddVehicleRecordRequest is the request to append a vehicle record (a mileage
// and/or fuel reading) to a vehicle asset's vehicleRecords child collection.
// Ported from the legacy dto4assetus.AddVehicleRecordRequest, preserving the
// fuel-bearing payload (fuelVolume/fuelVolumeUnit/fuelCost/currency/mileage/
// mileageUnit).
type AddVehicleRecordRequest struct {
	dto4spaceus.SpaceRequest
	AssetID        string  `json:"assetID"`
	FuelVolume     float32 `json:"fuelVolume,omitempty"`
	FuelVolumeUnit string  `json:"fuelVolumeUnit,omitempty"`
	FuelCost       float32 `json:"fuelCost,omitempty"`
	Currency       string  `json:"currency,omitempty"`
	Mileage        float32 `json:"mileage,omitempty"`
	MileageUnit    string  `json:"mileageUnit,omitempty"`
}

// Validate validates the request.
func (v AddVehicleRecordRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if err := validate.RecordID(v.AssetID); err != nil {
		return validation.NewErrBadRequestFieldValue("assetID", err.Error())
	}
	return nil
}

// AddVehicleRecordResponse is returned on successful vehicle-record creation.
type AddVehicleRecordResponse struct {
	ID string `json:"id"`
}
