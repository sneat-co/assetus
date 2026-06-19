package dbo4assetus

import (
	"github.com/crediterra/money"
	"github.com/strongo/decimal"
	"github.com/strongo/strongoapp/with"
)

// VehicleRecordsCollection is the Firestore child collection name for a
// vehicle asset's records (mileage + fuel readings):
// /spaces/{spaceID}/ext/assetus/assets/{assetID}/vehicleRecords/{recordID}.
const VehicleRecordsCollection = "vehicleRecords"

// VehicleFuelRecord is the fuel portion of a vehicle record, ported from the
// legacy dbo4assetus.VehicleFuelRecord plus the fuel cost/currency carried by
// the legacy dal4assetus.Mileage / dto4assetus.AddVehicleRecordRequest payload.
type VehicleFuelRecord struct {
	Volume   decimal.Decimal64p2 `json:"volume,omitempty" firestore:"volume,omitempty"`
	Unit     string              `json:"unit,omitempty" firestore:"unit,omitempty"`
	Amount   *money.Amount       `json:"amount,omitempty" firestore:"amount,omitempty"`
	FuelCost decimal.Decimal64p2 `json:"fuelCost,omitempty" firestore:"fuelCost,omitempty"`
	Currency string              `json:"currency,omitempty" firestore:"currency,omitempty"`
}

// Validate returns an error if the fuel record is not valid.
func (v *VehicleFuelRecord) Validate() error {
	if v == nil {
		return nil
	}
	if v.Amount != nil {
		if err := v.Amount.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// VehicleMileage is the mileage portion of a vehicle record, ported from the
// legacy dbo4assetus.VehicleMileage.
type VehicleMileage struct {
	Value int    `json:"value" firestore:"value"`
	Unit  string `json:"unit" firestore:"unit"`
}

// Validate returns an error if the mileage record is not valid.
func (v *VehicleMileage) Validate() error {
	return nil
}

// VehicleRecordDbo is a persisted vehicle record (a mileage and/or fuel reading)
// in the asset's vehicleRecords child collection. Ported from the legacy
// dbo4assetus.VehicleRecordDbo, preserving both mileage and fuel data.
type VehicleRecordDbo struct {
	with.CreatedFields                    // Mandatory field
	Fuel               *VehicleFuelRecord `json:"fuel,omitempty" firestore:"fuel,omitempty"`
	Mileage            *VehicleMileage    `json:"mileage,omitempty" firestore:"mileage,omitempty"`
}

// Validate returns an error if the record is not valid.
func (v VehicleRecordDbo) Validate() error {
	if err := v.CreatedFields.Validate(); err != nil {
		return err
	}
	if err := v.Fuel.Validate(); err != nil {
		return err
	}
	if err := v.Mileage.Validate(); err != nil {
		return err
	}
	return nil
}
