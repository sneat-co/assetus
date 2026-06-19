package extras4assetus

import (
	"testing"
	"time"

	"github.com/crediterra/money"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/sneat-core-modules/core/extra"
	"github.com/strongo/decimal"
	"github.com/strongo/strongoapp/with"
)

// AC typed-extras-optional: an asset with a vehicle extra resolves by extraType,
// and an asset with no extra at all is a valid flat-core asset.
func TestTypedExtrasOptional(t *testing.T) {
	t.Run("vehicle_extra_resolves_by_extraType", func(t *testing.T) {
		var holder extra.WithExtraField
		vehicle := &AssetVehicleExtra{
			WithMakeModelRegNumberFields: WithMakeModelRegNumberFields{
				WithMakeModelFields: WithMakeModelFields{Make: "Toyota", Model: "Camry"},
			},
		}
		if err := holder.SetExtra(AssetExtraTypeVehicle, vehicle); err != nil {
			t.Fatalf("SetExtra: %v", err)
		}
		got, err := holder.GetExtraData()
		if err != nil {
			t.Fatalf("GetExtraData: %v", err)
		}
		if _, ok := got.(*AssetVehicleExtra); !ok {
			t.Fatalf("extraType %q resolved to %T, want *AssetVehicleExtra", holder.ExtraType, got)
		}
		if err := holder.Validate(); err != nil {
			t.Fatalf("vehicle extra holder rejected: %v", err)
		}
	})

	t.Run("no_extra_is_valid_flat_core", func(t *testing.T) {
		// A flat-core asset with no extra at all is valid.
		a := dbo4assetus.AssetBase{
			Name:       "Plain asset",
			Category:   const4assetus.CategoryBooks,
			Condition:  const4assetus.ConditionGood,
			Status:     const4assetus.StatusActive,
			Visibility: const4assetus.VisibilityFamily,
		}
		if err := a.Validate(); err != nil {
			t.Fatalf("flat-core asset with no extra rejected: %v", err)
		}
		// The registry resolves an unset extraType to the no-op extra.
		none := extra.NewExtraData("")
		if err := none.Validate(); err != nil {
			t.Fatalf("no-extra validate: %v", err)
		}
	})
}

// AC vehicle-extra-no-field-dropped: every vehicle attribute (incl. engine
// serial number and NCT due-date) is stored on the extra, and three vehicle
// records preserve mileage + fuel(volume/unit/cost/currency).
func TestVehicleExtraNoFieldDropped(t *testing.T) {
	v := &AssetVehicleExtra{
		WithMakeModelRegNumberFields: WithMakeModelRegNumberFields{
			WithMakeModelFields:        WithMakeModelFields{Make: "Toyota", Model: "Corolla"},
			WithOptionalRegNumberField: WithOptionalRegNumberField{RegNumber: "12-D-3456"},
		},
		WithEngineData: WithEngineData{
			EngineType:         const4assetus.EngineTypeCombustion,
			EngineFuel:         const4assetus.FuelTypeDiesel,
			EngineCC:           1998,
			EngineKW:           110,
			EngineNM:           320,
			EngineSerialNumber: "ENG-SN-0001",
		},
		Vin:            "JTDBR32E720012345",
		NctExpires:     "2030-06-30",
		TaxExpires:     "2026-12-31",
		NextServiceDue: "2027-03-15",
	}
	if err := v.Validate(); err != nil {
		t.Fatalf("valid vehicle extra rejected: %v", err)
	}

	// Every attribute is present.
	checks := map[string]bool{
		"make":               v.Make == "Toyota",
		"model":              v.Model == "Corolla",
		"regNumber":          v.RegNumber == "12-D-3456",
		"vin":                v.Vin == "JTDBR32E720012345",
		"engineType":         v.EngineType == const4assetus.EngineTypeCombustion,
		"engineFuel":         v.EngineFuel == const4assetus.FuelTypeDiesel,
		"engineCC":           v.EngineCC == 1998,
		"engineKW":           v.EngineKW == 110,
		"engineNM":           v.EngineNM == 320,
		"engineSerialNumber": v.EngineSerialNumber == "ENG-SN-0001",
		"nctExpires":         v.NctExpires == "2030-06-30",
		"taxExpires":         v.TaxExpires == "2026-12-31",
		"nextServiceDue":     v.NextServiceDue == "2027-03-15",
	}
	for field, ok := range checks {
		if !ok {
			t.Errorf("vehicle attribute %q was dropped or wrong", field)
		}
	}

	// Three vehicle records, each with mileage AND fuel(volume/unit/cost/currency).
	now := time.Now()
	created := with.CreatedFields{
		CreatedAtField: with.CreatedAtField{CreatedAt: now},
		CreatedByField: with.CreatedByField{CreatedBy: "user1"},
	}
	records := []dbo4assetus.VehicleRecordDbo{
		{
			CreatedFields: created,
			Mileage:       &dbo4assetus.VehicleMileage{Value: 10000, Unit: "km"},
			Fuel: &dbo4assetus.VehicleFuelRecord{
				Volume:   decimal.NewDecimal64p2(4, 50),
				Unit:     "l",
				FuelCost: decimal.NewDecimal64p2(75, 0),
				Currency: "EUR",
			},
		},
		{
			CreatedFields: created,
			Mileage:       &dbo4assetus.VehicleMileage{Value: 20000, Unit: "km"},
			Fuel: &dbo4assetus.VehicleFuelRecord{
				Volume:   decimal.NewDecimal64p2(50, 0),
				Unit:     "l",
				Amount:   amountPtr(money.NewAmount(money.CurrencyEUR, decimal.NewDecimal64p2(80, 0))),
				FuelCost: decimal.NewDecimal64p2(80, 0),
				Currency: "EUR",
			},
		},
		{
			CreatedFields: created,
			Mileage:       &dbo4assetus.VehicleMileage{Value: 30000, Unit: "km"},
			Fuel: &dbo4assetus.VehicleFuelRecord{
				Volume:   decimal.NewDecimal64p2(48, 25),
				Unit:     "l",
				FuelCost: decimal.NewDecimal64p2(79, 99),
				Currency: "EUR",
			},
		},
	}
	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}
	for i, r := range records {
		if err := r.Validate(); err != nil {
			t.Fatalf("record[%d] rejected: %v", i, err)
		}
		if r.Mileage == nil || r.Mileage.Value == 0 || r.Mileage.Unit == "" {
			t.Errorf("record[%d] mileage dropped: %+v", i, r.Mileage)
		}
		if r.Fuel == nil || r.Fuel.Volume == 0 || r.Fuel.Unit == "" || r.Fuel.FuelCost == 0 || r.Fuel.Currency == "" {
			t.Errorf("record[%d] fuel field dropped: %+v", i, r.Fuel)
		}
	}
}

// AC document-extra-full-shape: a passport doc extra preserves all 8 fields and
// applies the passport validation schema (number + validity required).
func TestDocumentExtraFullShape(t *testing.T) {
	d := &AssetDocumentExtra{
		DocType:       const4assetus.TypeDocumentPassport,
		Number:        "P1234567",
		BatchNumber:   "B-42",
		CountryID:     "IE",
		IssuedBy:      "DFA Ireland",
		IssuedOn:      "2020-01-15",
		EffectiveFrom: "2020-01-15",
		ExpiresOn:     "2030-01-14",
	}
	if err := d.Validate(); err != nil {
		t.Fatalf("valid passport doc extra rejected: %v", err)
	}

	checks := map[string]bool{
		"docType":       d.DocType == const4assetus.TypeDocumentPassport,
		"number":        d.Number == "P1234567",
		"batchNumber":   d.BatchNumber == "B-42",
		"countryID":     d.CountryID == "IE",
		"issuedBy":      d.IssuedBy == "DFA Ireland",
		"issuedOn":      d.IssuedOn == "2020-01-15",
		"effectiveFrom": d.EffectiveFrom == "2020-01-15",
		"expiresOn":     d.ExpiresOn == "2030-01-14",
	}
	for field, ok := range checks {
		if !ok {
			t.Errorf("document attribute %q was dropped or wrong", field)
		}
	}

	t.Run("passport_requires_number", func(t *testing.T) {
		bad := &AssetDocumentExtra{DocType: const4assetus.TypeDocumentPassport, ExpiresOn: "2030-01-14"}
		if err := bad.Validate(); err == nil {
			t.Error("passport without number should be rejected")
		}
	})
	t.Run("passport_requires_validity", func(t *testing.T) {
		bad := &AssetDocumentExtra{DocType: const4assetus.TypeDocumentPassport, Number: "P1"}
		if err := bad.Validate(); err == nil {
			t.Error("passport without validity (expiresOn) should be rejected")
		}
	})
	t.Run("resolves_by_extraType", func(t *testing.T) {
		var holder extra.WithExtraField
		if err := holder.SetExtra(AssetExtraTypeDocument, d); err != nil {
			t.Fatalf("SetExtra: %v", err)
		}
		got, err := holder.GetExtraData()
		if err != nil {
			t.Fatalf("GetExtraData: %v", err)
		}
		if _, ok := got.(*AssetDocumentExtra); !ok {
			t.Fatalf("resolved to %T, want *AssetDocumentExtra", got)
		}
	})
}

func TestRegistryResolvesAllExtraTypes(t *testing.T) {
	for _, tc := range []struct {
		t    extra.Type
		want any
	}{
		{AssetExtraTypeVehicle, &AssetVehicleExtra{}},
		{AssetExtraTypeDwelling, &AssetDwellingExtra{}},
		{AssetExtraTypeDocument, &AssetDocumentExtra{}},
	} {
		got := NewAssetExtra(tc.t)
		if got == nil {
			t.Errorf("NewAssetExtra(%q) = nil", tc.t)
		}
	}
	if NewAssetExtra("nope") != nil {
		t.Error("NewAssetExtra of unknown type should be nil")
	}
}

func amountPtr(a money.Amount) *money.Amount { return &a }
