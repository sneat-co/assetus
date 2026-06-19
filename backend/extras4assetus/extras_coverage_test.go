package extras4assetus

import (
	"testing"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

func TestAssetDocumentExtra_RequiredFields_isNil(t *testing.T) {
	if got := (&AssetDocumentExtra{}).RequiredFields(); got != nil {
		t.Errorf("RequiredFields() = %v, want nil", got)
	}
}

func TestAssetDocumentExtra_IndexedFields_listsValidityDates(t *testing.T) {
	got := (&AssetDocumentExtra{}).IndexedFields()
	want := []string{"expiresOn", "effectiveFrom"}
	if len(got) != len(want) {
		t.Fatalf("IndexedFields() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("IndexedFields()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestAssetDocumentExtra_GetBrief_keepsIdentityAndDatesDropsBatchAndCountry(t *testing.T) {
	d := &AssetDocumentExtra{
		DocType:       const4assetus.TypeDocumentPassport,
		Number:        "P1",
		BatchNumber:   "B-1",
		CountryID:     "IE",
		IssuedBy:      "DFA",
		IssuedOn:      "2020-01-15",
		EffectiveFrom: "2020-01-15",
		ExpiresOn:     "2030-01-14",
	}
	brief, ok := d.GetBrief().(*AssetDocumentExtra)
	if !ok {
		t.Fatalf("GetBrief() returned %T, want *AssetDocumentExtra", d.GetBrief())
	}
	if brief.Number != "P1" || brief.DocType != const4assetus.TypeDocumentPassport {
		t.Errorf("brief dropped identity: %+v", brief)
	}
	if brief.IssuedOn != "2020-01-15" || brief.EffectiveFrom != "2020-01-15" || brief.ExpiresOn != "2030-01-14" {
		t.Errorf("brief dropped dates: %+v", brief)
	}
	if brief.BatchNumber != "" || brief.CountryID != "" || brief.IssuedBy != "" {
		t.Errorf("brief should omit batchNumber/countryID/issuedBy, got %+v", brief)
	}
}

func TestAssetDocumentExtra_Validate_rejectsRegNumberWithSpaces(t *testing.T) {
	d := &AssetDocumentExtra{WithOptionalRegNumberField: WithOptionalRegNumberField{RegNumber: " X "}}
	if err := d.Validate(); err == nil {
		t.Error("regNumber with surrounding spaces should be rejected")
	}
}

func TestAssetDocumentExtra_Validate_rejectsInvalidCountryCode(t *testing.T) {
	d := &AssetDocumentExtra{CountryID: "IRL", DocType: "noschema"}
	if err := d.Validate(); err == nil {
		t.Error("invalid country alpha-2 code should be rejected")
	}
}

func TestAssetDocumentExtra_Validate_rejectsInvalidIssuedOnDate(t *testing.T) {
	d := &AssetDocumentExtra{IssuedOn: "not-a-date", DocType: "noschema"}
	if err := d.Validate(); err == nil {
		t.Error("invalid issuedOn date should be rejected")
	}
}

func TestAssetDocumentExtra_Validate_rejectsInvalidEffectiveFromDate(t *testing.T) {
	d := &AssetDocumentExtra{EffectiveFrom: "13/2020", DocType: "noschema"}
	if err := d.Validate(); err == nil {
		t.Error("invalid effectiveFrom date should be rejected")
	}
}

func TestAssetDocumentExtra_Validate_rejectsInvalidExpiresOnDate(t *testing.T) {
	d := &AssetDocumentExtra{ExpiresOn: "2030-13-40", DocType: "noschema"}
	if err := d.Validate(); err == nil {
		t.Error("invalid expiresOn date should be rejected")
	}
}

func TestAssetDocumentExtra_Validate_rejectsExpiresBeforeEffectiveFrom(t *testing.T) {
	d := &AssetDocumentExtra{
		DocType:       "noschema",
		EffectiveFrom: "2030-01-15",
		ExpiresOn:     "2020-01-14",
	}
	if err := d.Validate(); err == nil {
		t.Error("expiresOn before effectiveFrom should be rejected")
	}
}

func TestAssetDocumentExtra_Validate_drivingLicenceRequiresNumber(t *testing.T) {
	d := &AssetDocumentExtra{DocType: const4assetus.TypeDocumentDrivingLicense, ExpiresOn: "2030-01-14"}
	if err := d.Validate(); err == nil {
		t.Error("driving licence without number should be rejected")
	}
}

func TestAssetDocumentExtra_Validate_drivingLicenceRequiresValidity(t *testing.T) {
	d := &AssetDocumentExtra{DocType: const4assetus.TypeDocumentDrivingLicense, Number: "DL-1"}
	if err := d.Validate(); err == nil {
		t.Error("driving licence without validity (expiresOn) should be rejected")
	}
}

func TestAssetDocumentExtra_Validate_drivingLicenceAcceptsNumberViaRegNumberAlias(t *testing.T) {
	d := &AssetDocumentExtra{
		DocType:                    const4assetus.TypeDocumentDrivingLicense,
		WithOptionalRegNumberField: WithOptionalRegNumberField{RegNumber: "DL-1"},
		ExpiresOn:                  "2030-01-14",
	}
	if err := d.Validate(); err != nil {
		t.Errorf("driving licence with number via regNumber alias should be accepted: %v", err)
	}
}

func TestAssetDocumentExtra_Validate_birthCertRequiresIssuedOn(t *testing.T) {
	d := &AssetDocumentExtra{DocType: const4assetus.TypeDocumentBirthCert, Number: "BC-1"}
	if err := d.Validate(); err == nil {
		t.Error("birth certificate without issuedOn should be rejected")
	}
}

func TestAssetDocumentExtra_Validate_birthCertExcludesValidity(t *testing.T) {
	d := &AssetDocumentExtra{
		DocType:   const4assetus.TypeDocumentBirthCert,
		Number:    "BC-1",
		IssuedOn:  "2000-05-05",
		ExpiresOn: "2030-01-14",
	}
	if err := d.Validate(); err == nil {
		t.Error("birth certificate with expiresOn should be rejected (validity excluded)")
	}
}

func TestAssetDocumentExtra_Validate_birthCertHappyPath(t *testing.T) {
	d := &AssetDocumentExtra{
		DocType:  const4assetus.TypeDocumentBirthCert,
		Number:   "BC-1",
		IssuedOn: "2000-05-05",
	}
	if err := d.Validate(); err != nil {
		t.Errorf("valid birth certificate rejected: %v", err)
	}
}

func TestAssetDocumentExtra_Validate_otherDocTypeRequiresTitleOnlyNoNumber(t *testing.T) {
	// "other" schema constrains only Title (not Number/validity), so a bare
	// "other" document passes the per-doc-type schema.
	d := &AssetDocumentExtra{DocType: "other"}
	if err := d.Validate(); err != nil {
		t.Errorf("`other` doc type without number/validity should pass schema: %v", err)
	}
}

func TestAssetDocumentExtra_Validate_unknownDocTypeImposesNoSchema(t *testing.T) {
	d := &AssetDocumentExtra{DocType: "totally-unknown"}
	if err := d.Validate(); err != nil {
		t.Errorf("unknown doc type should impose no extra requirements: %v", err)
	}
}

func TestAssetDwellingExtra_RequiredFields_isNil(t *testing.T) {
	if got := (&AssetDwellingExtra{}).RequiredFields(); got != nil {
		t.Errorf("RequiredFields() = %v, want nil", got)
	}
}

func TestAssetDwellingExtra_IndexedFields_isNil(t *testing.T) {
	if got := (&AssetDwellingExtra{}).IndexedFields(); got != nil {
		t.Errorf("IndexedFields() = %v, want nil", got)
	}
}

func TestAssetDwellingExtra_GetBrief_keepsBedroomsAndAreaDropsAddress(t *testing.T) {
	d := &AssetDwellingExtra{
		Address:          &dbmodels.Address{},
		NumberOfBedrooms: 3,
		AreaSqM:          85,
	}
	brief, ok := d.GetBrief().(*AssetDwellingExtra)
	if !ok {
		t.Fatalf("GetBrief() returned %T, want *AssetDwellingExtra", d.GetBrief())
	}
	if brief.NumberOfBedrooms != 3 || brief.AreaSqM != 85 {
		t.Errorf("brief dropped bedrooms/area: %+v", brief)
	}
	if brief.Address != nil {
		t.Errorf("brief should omit address, got %+v", brief.Address)
	}
}

func TestAssetDwellingExtra_Validate_happyPath(t *testing.T) {
	d := &AssetDwellingExtra{NumberOfBedrooms: 2, AreaSqM: 70}
	d.RentPrice.Value = 1500
	d.RentPrice.Currency = "EUR"
	if err := d.Validate(); err != nil {
		t.Errorf("valid dwelling rejected: %v", err)
	}
}

func TestAssetDwellingExtra_Validate_rejectsNegativeBedrooms(t *testing.T) {
	d := &AssetDwellingExtra{NumberOfBedrooms: -1}
	if err := d.Validate(); err == nil {
		t.Error("negative numberOfBedrooms should be rejected")
	}
}

func TestAssetDwellingExtra_Validate_rejectsNegativeArea(t *testing.T) {
	d := &AssetDwellingExtra{AreaSqM: -1}
	if err := d.Validate(); err == nil {
		t.Error("negative areaSqM should be rejected")
	}
}

func TestAssetDwellingExtra_Validate_rejectsNegativeRentPrice(t *testing.T) {
	d := &AssetDwellingExtra{}
	d.RentPrice.Value = -100
	if err := d.Validate(); err == nil {
		t.Error("negative rent_price.value should be rejected")
	}
}

func TestAssetDwellingExtra_Validate_rejectsInvalidAddress(t *testing.T) {
	// An address with no countryID is rejected by Address.Validate.
	d := &AssetDwellingExtra{Address: &dbmodels.Address{City: "Dublin"}}
	if err := d.Validate(); err == nil {
		t.Error("dwelling with invalid address (no countryID) should be rejected")
	}
}

func TestAssetVehicleExtra_Validate_rejectsIncompatibleEngineFuel(t *testing.T) {
	v := &AssetVehicleExtra{
		WithMakeModelRegNumberFields: WithMakeModelRegNumberFields{
			WithMakeModelFields: WithMakeModelFields{Make: "Tesla", Model: "Model 3"},
		},
		WithEngineData: WithEngineData{
			EngineType: const4assetus.EngineTypeElectric,
			EngineFuel: const4assetus.FuelTypeDiesel,
		},
	}
	if err := v.Validate(); err == nil {
		t.Error("electric engine with diesel fuel should be rejected")
	}
}

func TestAssetVehicleExtra_RequiredFields_isNil(t *testing.T) {
	if got := (&AssetVehicleExtra{}).RequiredFields(); got != nil {
		t.Errorf("RequiredFields() = %v, want nil", got)
	}
}

func TestAssetVehicleExtra_IndexedFields_listsMakeModelRegVin(t *testing.T) {
	got := (&AssetVehicleExtra{}).IndexedFields()
	want := []string{"make", "model", "make+model", "regNumber", "vin"}
	if len(got) != len(want) {
		t.Fatalf("IndexedFields() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("IndexedFields()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestAssetVehicleExtra_GetBrief_keepsMakeModelVin(t *testing.T) {
	v := &AssetVehicleExtra{
		WithMakeModelRegNumberFields: WithMakeModelRegNumberFields{
			WithMakeModelFields: WithMakeModelFields{Make: "Toyota", Model: "Camry"},
		},
		Vin:        "VIN123",
		NctExpires: "2030-06-30",
	}
	brief, ok := v.GetBrief().(*AssetVehicleExtra)
	if !ok {
		t.Fatalf("GetBrief() returned %T, want *AssetVehicleExtra", v.GetBrief())
	}
	if brief.Make != "Toyota" || brief.Model != "Camry" || brief.Vin != "VIN123" {
		t.Errorf("brief dropped make/model/vin: %+v", brief)
	}
	if brief.NctExpires != "" {
		t.Errorf("brief should omit nctExpires, got %q", brief.NctExpires)
	}
}

func TestAssetVehicleExtra_Validate_rejectsMissingMake(t *testing.T) {
	v := &AssetVehicleExtra{
		WithMakeModelRegNumberFields: WithMakeModelRegNumberFields{
			WithMakeModelFields: WithMakeModelFields{Model: "Camry"},
		},
	}
	if err := v.Validate(); err == nil {
		t.Error("vehicle without make should be rejected")
	}
}

func TestAssetVehicleExtra_Validate_rejectsInvalidNctExpiresDate(t *testing.T) {
	v := &AssetVehicleExtra{
		WithMakeModelRegNumberFields: WithMakeModelRegNumberFields{
			WithMakeModelFields: WithMakeModelFields{Make: "Toyota", Model: "Camry"},
		},
		NctExpires: "not-a-date",
	}
	if err := v.Validate(); err == nil {
		t.Error("invalid nctExpires date should be rejected")
	}
}

func TestAssetVehicleExtra_Validate_rejectsInvalidTaxExpiresDate(t *testing.T) {
	v := &AssetVehicleExtra{
		WithMakeModelRegNumberFields: WithMakeModelRegNumberFields{
			WithMakeModelFields: WithMakeModelFields{Make: "Toyota", Model: "Camry"},
		},
		TaxExpires: "2030-13-40",
	}
	if err := v.Validate(); err == nil {
		t.Error("invalid taxExpires date should be rejected")
	}
}

func TestWithMakeModelFields_Validate_rejectsMissingModel(t *testing.T) {
	v := &WithMakeModelFields{Make: "Toyota"}
	if err := v.Validate(); err == nil {
		t.Error("missing model should be rejected")
	}
}

func TestWithMakeModelFields_Validate_rejectsMakeWithSpaces(t *testing.T) {
	v := &WithMakeModelFields{Make: " Toyota ", Model: "Camry"}
	if err := v.Validate(); err == nil {
		t.Error("make with surrounding spaces should be rejected")
	}
}

func TestWithMakeModelFields_Validate_rejectsModelWithSpaces(t *testing.T) {
	v := &WithMakeModelFields{Make: "Toyota", Model: " Camry "}
	if err := v.Validate(); err == nil {
		t.Error("model with surrounding spaces should be rejected")
	}
}

func TestWithMakeModelFields_GenerateTitle_makeModelAndRegNumber(t *testing.T) {
	v := &WithMakeModelFields{Make: "Toyota", Model: "Camry"}
	if got := v.GenerateTitleFromMakeModelAndRegNumber("12-D-3456"); got != "Toyota Camry # 12-D-3456" {
		t.Errorf("title = %q, want %q", got, "Toyota Camry # 12-D-3456")
	}
}

func TestWithMakeModelFields_GenerateTitle_makeModelOnly(t *testing.T) {
	v := &WithMakeModelFields{Make: "Toyota", Model: "Camry"}
	if got := v.GenerateTitleFromMakeModelAndRegNumber(""); got != "Toyota Camry" {
		t.Errorf("title = %q, want %q", got, "Toyota Camry")
	}
}

func TestWithMakeModelFields_GenerateTitle_emptyWhenAllBlank(t *testing.T) {
	v := &WithMakeModelFields{}
	if got := v.GenerateTitleFromMakeModelAndRegNumber(""); got != "" {
		t.Errorf("title = %q, want empty string", got)
	}
}

func TestWithMakeModelRegNumberFields_Validate_rejectsRegNumberWithSpaces(t *testing.T) {
	v := &WithMakeModelRegNumberFields{
		WithMakeModelFields:        WithMakeModelFields{Make: "Toyota", Model: "Camry"},
		WithOptionalRegNumberField: WithOptionalRegNumberField{RegNumber: " 12-D "},
	}
	if err := v.Validate(); err == nil {
		t.Error("regNumber with surrounding spaces should be rejected")
	}
}

func TestWithMakeModelRegNumberFields_Validate_propagatesMakeModelError(t *testing.T) {
	v := &WithMakeModelRegNumberFields{
		WithOptionalRegNumberField: WithOptionalRegNumberField{RegNumber: "12-D"},
	}
	if err := v.Validate(); err == nil {
		t.Error("missing make/model should be rejected")
	}
}

func TestWithOptionalRegNumberField_Validate_happyPath(t *testing.T) {
	v := &WithOptionalRegNumberField{RegNumber: "12-D-3456"}
	if err := v.Validate(); err != nil {
		t.Errorf("valid regNumber rejected: %v", err)
	}
}

func TestDocTypeSchema_returnsSchemaForPassport(t *testing.T) {
	def, ok := DocTypeSchema(const4assetus.TypeDocumentPassport)
	if !ok {
		t.Fatal("DocTypeSchema(passport) = _, false; want a schema")
	}
	if def.ID != const4assetus.TypeDocumentPassport {
		t.Errorf("schema ID = %q, want %q", def.ID, const4assetus.TypeDocumentPassport)
	}
	if def.Fields.Number == nil || !def.Fields.Number.Required {
		t.Error("passport schema should require number")
	}
}

func TestDocTypeSchema_falseForUnknownDocType(t *testing.T) {
	if _, ok := DocTypeSchema("unknown-doc-type"); ok {
		t.Error("DocTypeSchema(unknown) should return ok=false")
	}
}

func TestRegisterAssetExtraFactory_registersResolvableFactory(t *testing.T) {
	const t1 = "coverage-test-extra"
	RegisterAssetExtraFactory(t1, func() AssetExtra { return &AssetVehicleExtra{} })
	if got := NewAssetExtra(t1); got == nil {
		t.Fatalf("NewAssetExtra(%q) = nil after registration", t1)
	}
	if _, ok := NewAssetExtra(t1).(*AssetVehicleExtra); !ok {
		t.Errorf("registered factory resolved to %T, want *AssetVehicleExtra", NewAssetExtra(t1))
	}
}
