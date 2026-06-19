package const4assetus

import (
	"slices"
	"testing"
)

// AC category-superset: legacy sport_gear+kite_board maps to sports_equipment
// with kite_board preserved; MVP books maps to books with no subtype; and the
// legacy↔unified mapping table has an entry for every legacy AssetCategory,
// AssetType, EngineType, and FuelType value.
func TestCategorySuperset_LegacyMappingTotality(t *testing.T) {
	// sport_gear -> sports_equipment, kite_board subtype preserved.
	if got := LegacyCategoryToUnified["sport_gear"]; got != CategorySportsEquipment {
		t.Errorf("sport_gear maps to %q, want %q", got, CategorySportsEquipment)
	}
	if got := LegacyTypeToUnified["kite_board"]; got != TypeSportsKiteBoard {
		t.Errorf("kite_board maps to %q, want %q", got, TypeSportsKiteBoard)
	}
	if string(LegacyTypeToUnified["kite_board"]) != "kite_board" {
		t.Errorf("kite_board subtype value not preserved: %q", LegacyTypeToUnified["kite_board"])
	}
	if err := ValidateType(CategorySportsEquipment, TypeSportsKiteBoard); err != nil {
		t.Errorf("kite_board should be valid for sports_equipment: %v", err)
	}

	// books -> books, no subtype.
	if !slices.Contains(Categories, CategoryBooks) {
		t.Fatal("books missing from unified Categories")
	}
	if TypesForCategory(CategoryBooks) != nil {
		t.Errorf("books should have no subtypes, got %v", TypesForCategory(CategoryBooks))
	}
	if err := ValidateType(CategoryBooks, ""); err != nil {
		t.Errorf("books with no subtype should be valid: %v", err)
	}

	// vehicle -> vehicles, misc -> other reconciliations.
	if got := LegacyCategoryToUnified["vehicle"]; got != CategoryVehicles {
		t.Errorf("vehicle maps to %q, want %q", got, CategoryVehicles)
	}
	if got := LegacyCategoryToUnified["misc"]; got != CategoryOther {
		t.Errorf("misc maps to %q, want %q", got, CategoryOther)
	}

	// Totality: every legacy category/status/possession value maps to a valid
	// unified value, and every type/engine/fuel value is present.
	for legacy, unified := range LegacyCategoryToUnified {
		if !IsValidCategory(unified) {
			t.Errorf("legacy category %q maps to invalid unified category %q", legacy, unified)
		}
	}
	// debt retained as a category value.
	if !IsValidCategory(CategoryDebt) {
		t.Error("debt must be retained as a Category value")
	}
	// Every legacy AssetType / AssetDocumentType maps to a unified Type that is
	// permitted by some category (no engine/fuel/subtype value dropped).
	allUnifiedTypes := slices.Concat(VehicleTypes, DwellingTypes, SportsEquipmentTypes, DocumentTypes)
	for legacy, unified := range LegacyTypeToUnified {
		if !slices.Contains(allUnifiedTypes, unified) {
			t.Errorf("legacy type %q maps to unknown unified type %q", legacy, unified)
		}
	}
	if len(LegacyTypeToUnified) != len(allUnifiedTypes) {
		t.Errorf("legacy type mapping has %d entries, unified has %d types", len(LegacyTypeToUnified), len(allUnifiedTypes))
	}
	// Engine/fuel totality: no engine/fuel enum value dropped.
	for legacy, unified := range LegacyEngineTypeToUnified {
		if !slices.Contains(EngineTypes, unified) {
			t.Errorf("legacy engine type %q maps to unknown unified engine type %q", legacy, unified)
		}
	}
	if len(LegacyEngineTypeToUnified) != len(EngineTypes) {
		t.Errorf("legacy engine mapping has %d entries, unified has %d", len(LegacyEngineTypeToUnified), len(EngineTypes))
	}
	for legacy, unified := range LegacyFuelTypeToUnified {
		if !IsKnownFuelType(unified) {
			t.Errorf("legacy fuel type %q maps to unknown unified fuel type %q", legacy, unified)
		}
	}
	if len(LegacyFuelTypeToUnified) != len(FuelTypes) {
		t.Errorf("legacy fuel mapping has %d entries, unified has %d", len(LegacyFuelTypeToUnified), len(FuelTypes))
	}
}

// AC out-of-category-type-rejected: Category=document with Type=kite_board (a
// sports subtype) must be rejected because kite_board is not permitted for
// documents.
func TestOutOfCategoryTypeRejected(t *testing.T) {
	if err := ValidateType(CategoryDocument, TypeSportsKiteBoard); err == nil {
		t.Error("kite_board must be rejected for category document")
	}
	// But kite_board is accepted for sports_equipment.
	if err := ValidateType(CategorySportsEquipment, TypeSportsKiteBoard); err != nil {
		t.Errorf("kite_board must be accepted for sports_equipment: %v", err)
	}
	// A category with no subtype rejects any non-empty type.
	if err := ValidateType(CategoryBooks, "novel"); err == nil {
		t.Error("books must reject a non-empty subtype")
	}
	// Each per-category set rejects another category's subtype.
	if err := ValidateType(CategoryVehicles, TypeDwellingHouse); err == nil {
		t.Error("house must be rejected for vehicles")
	}
	// Valid in-category types pass.
	if err := ValidateType(CategoryVehicles, TypeVehicleCar); err != nil {
		t.Errorf("car must be valid for vehicles: %v", err)
	}
	if err := ValidateType(CategoryDwelling, TypeDwellingHouse); err != nil {
		t.Errorf("house must be valid for dwelling: %v", err)
	}
	if err := ValidateType(CategoryDocument, TypeDocumentPassport); err != nil {
		t.Errorf("passport must be valid for document: %v", err)
	}
}

// AC status-no-value-dropped: the unified Status set must be exactly
// {draft, active, transferred, archived, disposed, lost}.
func TestStatusUnionExact(t *testing.T) {
	want := []Status{StatusDraft, StatusActive, StatusTransferred, StatusArchived, StatusDisposed, StatusLost}
	if len(Statuses) != len(want) {
		t.Fatalf("Statuses has %d values, want %d: %v", len(Statuses), len(want), Statuses)
	}
	for _, s := range want {
		if !slices.Contains(Statuses, s) {
			t.Errorf("unified Status set missing %q", s)
		}
	}
	// Legacy set {active, archived, draft} all present.
	for _, s := range []Status{StatusActive, StatusArchived, StatusDraft} {
		if !slices.Contains(Statuses, s) {
			t.Errorf("legacy status %q dropped from union", s)
		}
	}
	// MVP set {active, transferred, archived, disposed, lost} all present.
	for _, s := range []Status{StatusActive, StatusTransferred, StatusArchived, StatusDisposed, StatusLost} {
		if !slices.Contains(Statuses, s) {
			t.Errorf("MVP status %q dropped from union", s)
		}
	}
}

func TestValidatePossession(t *testing.T) {
	for _, p := range Possessions {
		if err := ValidatePossession(p, false); err != nil {
			t.Errorf("ValidatePossession(%q) unexpected error: %v", p, err)
		}
	}
	if err := ValidatePossession("", true); err == nil {
		t.Error("ValidatePossession(empty, required) expected error")
	}
	if err := ValidatePossession("bogus", false); err == nil {
		t.Error("ValidatePossession(bogus) expected error")
	}
}

// Ported engine↔fuel compatibility matrix.
func TestValidateEngineFuel(t *testing.T) {
	tests := []struct {
		name    string
		engine  EngineType
		fuel    FuelType
		wantErr bool
	}{
		{"combustion_petrol", EngineTypeCombustion, FuelTypePetrol, false},
		{"combustion_electric_fuel", EngineTypeCombustion, "electric", true},
		{"electric_unknown", EngineTypeElectric, FuelTypeUnknown, false},
		{"electric_petrol", EngineTypeElectric, FuelTypePetrol, true},
		{"steam_other", EngineTypeSteam, FuelTypeOther, false},
		{"steam_diesel", EngineTypeSteam, FuelTypeDiesel, true},
		{"hybrid_diesel", EngineTypeHybrid, FuelTypeDiesel, false},
		{"other_bio", EngineTypeOther, FuelTypeBio, false},
		{"invalid_engine", "invalid", FuelTypePetrol, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateEngineFuel(tt.engine, tt.fuel); (err != nil) != tt.wantErr {
				t.Errorf("ValidateEngineFuel(%q,%q) error = %v, wantErr %v", tt.engine, tt.fuel, err, tt.wantErr)
			}
		})
	}
}
