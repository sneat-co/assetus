package extras4assetus

import (
	"github.com/sneat-co/sneat-core-modules/core/extra"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

func init() {
	RegisterAssetExtraFactory(AssetExtraTypeVehicle, func() AssetExtra {
		return new(AssetVehicleExtra)
	})
}

var _ extra.Data = (*AssetVehicleExtra)(nil)

// AssetVehicleExtra is the typed extra for vehicle assets. It carries
// make/model/regNumber/VIN, the engine data (incl. engineSerialNumber), and the
// plain service/tax/inspection (NCT) due-dates ported from the legacy frontend
// IAssetVehicleExtra.
type AssetVehicleExtra struct {
	WithMakeModelRegNumberFields
	WithEngineData
	Vin string `json:"vin,omitempty" firestore:"vin,omitempty"`

	// Plain due-date values (ISO "YYYY-MM-DD"), ported from the legacy frontend
	// IAssetVehicleExtra nctExpires / taxExpires / nextServiceDue.
	NctExpires     string `json:"nctExpires,omitempty" firestore:"nctExpires,omitempty"`
	TaxExpires     string `json:"taxExpires,omitempty" firestore:"taxExpires,omitempty"`
	NextServiceDue string `json:"nextServiceDue,omitempty" firestore:"nextServiceDue,omitempty"`
}

// RequiredFields implements extra.Data.
func (v *AssetVehicleExtra) RequiredFields() []string {
	return nil
}

// IndexedFields implements extra.Data; ported from the legacy declaration.
func (v *AssetVehicleExtra) IndexedFields() []string {
	return []string{"make", "model", "make+model", "regNumber", "vin"}
}

// GetBrief implements extra.Data.
func (v *AssetVehicleExtra) GetBrief() extra.Data {
	return &AssetVehicleExtra{
		WithMakeModelRegNumberFields: v.WithMakeModelRegNumberFields,
		Vin:                          v.Vin,
	}
}

// Validate returns an error if the vehicle extra is not valid.
func (v *AssetVehicleExtra) Validate() error {
	if err := v.WithMakeModelRegNumberFields.Validate(); err != nil {
		return err
	}
	if err := v.WithEngineData.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("engineData", err.Error())
	}
	for name, value := range map[string]string{
		"nctExpires":     v.NctExpires,
		"taxExpires":     v.TaxExpires,
		"nextServiceDue": v.NextServiceDue,
	} {
		if value == "" {
			continue
		}
		if _, err := validate.DateString(value); err != nil {
			return validation.NewErrBadRecordFieldValue(name, err.Error())
		}
	}
	return nil
}
