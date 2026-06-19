package dbo4assetus

import (
	"testing"
	"time"

	"github.com/crediterra/money"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/strongo/decimal"
	"github.com/strongo/strongoapp/with"
)

// validCreatedFields builds a populated with.CreatedFields so the embedding
// records (vehicle record, space module entry) pass CreatedFields.Validate.
func validCreatedFields() with.CreatedFields {
	now := time.Now()
	return with.CreatedFields{
		CreatedAtField: with.CreatedAtField{CreatedAt: now},
		CreatedByField: with.CreatedByField{CreatedBy: "user1"},
	}
}

// --- AssetDates.Validate -------------------------------------------------

func TestAssetDatesValidate_PresentDatesRoundTrip(t *testing.T) {
	d := AssetDates{
		DateOfBuild:       "2010-05-01",
		DateOfPurchase:    "2015-06-02",
		DateInsuredTill:   "2030-01-01",
		DateCertifiedTill: "2031-12-31",
	}
	if err := d.Validate(); err != nil {
		t.Fatalf("valid AssetDates rejected: %v", err)
	}
}

func TestAssetDatesValidate_RejectsMalformedDate(t *testing.T) {
	d := AssetDates{DateOfPurchase: "not-a-date"}
	if err := d.Validate(); err == nil {
		t.Error("expected rejection of malformed dateOfPurchase, got nil")
	}
}

// --- AssetLiabilityInfo.Validate ----------------------------------------

func TestAssetLiabilityInfoValidate_RequiresID(t *testing.T) {
	if err := (AssetLiabilityInfo{ID: "  "}).Validate(); err == nil {
		t.Error("expected missing-id error for liability, got nil")
	}
}

func TestAssetLiabilityInfoValidate_OK(t *testing.T) {
	l := AssetLiabilityInfo{ID: "liab1", ServiceTypes: []LiabilityServiceType{"electricity"}}
	if err := l.Validate(); err != nil {
		t.Fatalf("valid liability rejected: %v", err)
	}
}

// --- AssetBase.Validate optional-branch coverage -------------------------

func TestAssetBaseValidate_RejectsEmptyTag(t *testing.T) {
	a := validAssetBase()
	a.Tags = []string{"ok", "   "}
	if err := a.Validate(); err == nil {
		t.Error("expected rejection of empty tag, got nil")
	}
}

func TestAssetBaseValidate_RejectsInvalidCountryID(t *testing.T) {
	a := validAssetBase()
	a.CountryID = "z"
	if err := a.Validate(); err == nil {
		t.Error("expected rejection of invalid countryID, got nil")
	}
}

func TestAssetBaseValidate_RejectsInvalidFinancialDirection(t *testing.T) {
	a := validAssetBase()
	a.FinancialDirection = "sideways"
	if err := a.Validate(); err == nil {
		t.Error("expected rejection of invalid financialDirection, got nil")
	}
}

func TestAssetBaseValidate_RejectsInvalidParentCategory(t *testing.T) {
	a := validAssetBase()
	a.ParentCategoryID = const4assetus.Category("no-such-category")
	if err := a.Validate(); err == nil {
		t.Error("expected rejection of invalid parentCategoryID, got nil")
	}
}

func TestAssetBaseValidate_RejectsInvalidTotalsCurrency(t *testing.T) {
	a := validAssetBase()
	a.Totals = []money.Amount{{Currency: "ZZZ", Value: decimal.NewDecimal64p2(1, 0)}}
	if err := a.Validate(); err == nil {
		t.Error("expected rejection of unknown totals currency, got nil")
	}
}

func TestAssetBaseValidate_RejectsInvalidLiability(t *testing.T) {
	a := validAssetBase()
	a.Liabilities = []AssetLiabilityInfo{{ID: ""}}
	if err := a.Validate(); err == nil {
		t.Error("expected rejection of liability with empty id, got nil")
	}
}

// AC financial-fields-round-trip: optional financial fields are preserved with
// their exact monetary amounts and capability flags.
func TestAssetBaseValidate_FinancialFieldsRoundTrip(t *testing.T) {
	a := validAssetBase()
	a.CountryID = "IE"
	a.Possession = const4assetus.PossessionOwning
	a.Type = ""
	a.CanHaveIncome = true
	a.CanHaveExpense = true
	a.FinancialDirection = "income"
	a.Totals = []money.Amount{money.NewAmount(money.CurrencyEUR, decimal.NewDecimal64p2(123, 45))}
	a.Liabilities = []AssetLiabilityInfo{{ID: "liab1", ServiceTypes: []LiabilityServiceType{"gas"}}}
	a.AssetDates = AssetDates{DateOfPurchase: "2020-01-15"}

	if err := a.Validate(); err != nil {
		t.Fatalf("asset with financial fields rejected: %v", err)
	}
	if len(a.Totals) != 1 || a.Totals[0].Currency != money.CurrencyEUR {
		t.Fatalf("totals currency not preserved: %+v", a.Totals)
	}
	if got := a.Totals[0].Value; got != decimal.NewDecimal64p2(123, 45) {
		t.Errorf("totals value = %v, want 123.45", got)
	}
}

// AC possession-default: an asset created without a possession value resolves to
// owning while a supplied possession is preserved.
func TestAssetBase_WithPossessionDefault(t *testing.T) {
	a := validAssetBase()
	if got := a.WithPossessionDefault(); got != const4assetus.PossessionOwning {
		t.Errorf("default possession = %q, want owning", got)
	}
	a.Possession = const4assetus.PossessionRenting
	if got := a.WithPossessionDefault(); got != const4assetus.PossessionRenting {
		t.Errorf("possession = %q, want renting", got)
	}
}

// --- AssetDbo.Validate error branches -----------------------------------

func TestAssetDboValidate_RejectsMissingSpaceID(t *testing.T) {
	dbo := &AssetDbo{AssetBase: validAssetBase()}
	if err := dbo.Validate(); err == nil {
		t.Error("expected rejection of AssetDbo without space IDs, got nil")
	}
}

// --- AssetBrief.Validate -------------------------------------------------

func validAssetBrief() AssetBrief {
	return AssetBrief{
		ID:         "asset1",
		Name:       "Bicycle",
		Category:   const4assetus.CategoryVehicles,
		Condition:  const4assetus.ConditionGood,
		Status:     const4assetus.StatusActive,
		Visibility: const4assetus.VisibilityFamily,
	}
}

func TestAssetBriefValidate_OK(t *testing.T) {
	if err := validAssetBrief().Validate(); err != nil {
		t.Fatalf("valid brief rejected: %v", err)
	}
}

func TestAssetBriefValidate_RequiresID(t *testing.T) {
	b := validAssetBrief()
	b.ID = "  "
	if err := b.Validate(); err == nil {
		t.Error("expected missing-id error, got nil")
	}
}

func TestAssetBriefValidate_RequiresName(t *testing.T) {
	b := validAssetBrief()
	b.Name = ""
	if err := b.Validate(); err == nil {
		t.Error("expected missing-name error, got nil")
	}
}

func TestAssetBriefValidate_RejectsInvalidStatus(t *testing.T) {
	b := validAssetBrief()
	b.Status = "borrowed"
	if err := b.Validate(); err == nil {
		t.Error("expected rejection of invalid brief status, got nil")
	}
}

func TestBriefFromAsset_CopiesFields(t *testing.T) {
	dbo := &AssetDbo{AssetBase: validAssetBase()}
	brief := BriefFromAsset("asset9", dbo)
	if brief.ID != "asset9" || brief.Name != dbo.Name || brief.Category != dbo.Category {
		t.Errorf("brief not built from asset: %+v", brief)
	}
}

// --- AssetHistoryEventBase.Validate -------------------------------------

func TestAssetHistoryEventValidate_RequiresOccurredAt(t *testing.T) {
	base := AssetHistoryEventBase{
		Type:     const4assetus.HistoryEventPurchased,
		ActorRef: "user1",
	}
	if err := base.Validate(); err == nil {
		t.Error("expected missing-occurredAt error, got nil")
	}
}

func TestAssetHistoryEventValidate_RequiresActorRef(t *testing.T) {
	base := AssetHistoryEventBase{
		Type:       const4assetus.HistoryEventPurchased,
		OccurredAt: time.Now(),
	}
	if err := base.Validate(); err == nil {
		t.Error("expected missing-actorRef error, got nil")
	}
}

func TestAssetHistoryEventDboValidate_DelegatesToBase(t *testing.T) {
	dbo := &AssetHistoryEventDbo{
		AssetHistoryEventBase: AssetHistoryEventBase{
			Type:       const4assetus.HistoryEventPurchased,
			OccurredAt: time.Now(),
			ActorRef:   "user1",
		},
	}
	if err := dbo.Validate(); err != nil {
		t.Fatalf("valid history event dbo rejected: %v", err)
	}
	dbo.ActorRef = ""
	if err := dbo.Validate(); err == nil {
		t.Error("expected history event dbo to reject missing actorRef, got nil")
	}
}

// --- AssetusSpaceDbo.Validate -------------------------------------------

func TestAssetusSpaceDboValidate_OK(t *testing.T) {
	dbo := &AssetusSpaceDbo{
		CreatedFields: validCreatedFields(),
		Assets:        AssetBriefs{"asset1": ptr(validAssetBrief())},
	}
	if err := dbo.Validate(); err != nil {
		t.Fatalf("valid space module entry rejected: %v", err)
	}
}

func TestAssetusSpaceDboValidate_RejectsNilBrief(t *testing.T) {
	dbo := &AssetusSpaceDbo{
		CreatedFields: validCreatedFields(),
		Assets:        AssetBriefs{"asset1": nil},
	}
	if err := dbo.Validate(); err == nil {
		t.Error("expected rejection of nil brief, got nil")
	}
}

func TestAssetusSpaceDboValidate_RejectsInvalidBrief(t *testing.T) {
	bad := validAssetBrief()
	bad.Name = ""
	dbo := &AssetusSpaceDbo{
		CreatedFields: validCreatedFields(),
		Assets:        AssetBriefs{"asset1": &bad},
	}
	if err := dbo.Validate(); err == nil {
		t.Error("expected rejection of invalid brief, got nil")
	}
}

// --- WithAssetSpaces.Validate -------------------------------------------

func TestWithAssetSpacesValidate_OK(t *testing.T) {
	w := &WithAssetSpaces{
		Spaces: map[string]*AssetusSpaceBrief{
			"spaceA": {Assets: AssetBriefs{"asset1": ptr(validAssetBrief())}},
		},
	}
	if err := w.Validate(); err != nil {
		t.Fatalf("valid asset spaces rejected: %v", err)
	}
}

func TestWithAssetSpacesValidate_RejectsEmptySpaceID(t *testing.T) {
	w := &WithAssetSpaces{Spaces: map[string]*AssetusSpaceBrief{"": {}}}
	if err := w.Validate(); err == nil {
		t.Error("expected rejection of empty spaceID, got nil")
	}
}

func TestWithAssetSpacesValidate_RejectsNilSpaceBrief(t *testing.T) {
	w := &WithAssetSpaces{Spaces: map[string]*AssetusSpaceBrief{"spaceA": nil}}
	if err := w.Validate(); err == nil {
		t.Error("expected rejection of nil space brief, got nil")
	}
}

func TestWithAssetSpacesValidate_RejectsEmptyAssetID(t *testing.T) {
	w := &WithAssetSpaces{
		Spaces: map[string]*AssetusSpaceBrief{
			"spaceA": {Assets: AssetBriefs{"": ptr(validAssetBrief())}},
		},
	}
	if err := w.Validate(); err == nil {
		t.Error("expected rejection of empty assetID, got nil")
	}
}

func TestWithAssetSpacesValidate_RejectsNilAssetBrief(t *testing.T) {
	w := &WithAssetSpaces{
		Spaces: map[string]*AssetusSpaceBrief{
			"spaceA": {Assets: AssetBriefs{"asset1": nil}},
		},
	}
	if err := w.Validate(); err == nil {
		t.Error("expected rejection of nil asset brief, got nil")
	}
}

func TestWithAssetSpacesValidate_RejectsInvalidAssetBrief(t *testing.T) {
	bad := validAssetBrief()
	bad.ID = ""
	w := &WithAssetSpaces{
		Spaces: map[string]*AssetusSpaceBrief{
			"spaceA": {Assets: AssetBriefs{"asset1": &bad}},
		},
	}
	if err := w.Validate(); err == nil {
		t.Error("expected rejection of invalid asset brief, got nil")
	}
}

// --- relationships.go: TitledRecord / SubAssetInfo / AssetGroupInfo ------

func TestTitledRecordValidate_RequiresID(t *testing.T) {
	if err := (TitledRecord{ID: " "}).Validate(); err == nil {
		t.Error("expected missing-id error for titled record, got nil")
	}
	if err := (TitledRecord{ID: "t1", Title: "Title"}).Validate(); err != nil {
		t.Fatalf("valid titled record rejected: %v", err)
	}
}

func TestSubAssetInfoValidate_RejectsMissingID(t *testing.T) {
	s := SubAssetInfo{Type: const4assetus.CategoryVehicles}
	if err := s.Validate(); err == nil {
		t.Error("expected missing-id error for sub-asset, got nil")
	}
}

func TestSubAssetInfoValidate_RejectsInvalidCategory(t *testing.T) {
	s := SubAssetInfo{TitledRecord: TitledRecord{ID: "s1"}, Type: const4assetus.Category("bogus")}
	if err := s.Validate(); err == nil {
		t.Error("expected rejection of invalid sub-asset category, got nil")
	}
}

func TestSubAssetInfoValidate_RejectsInvalidCountryID(t *testing.T) {
	s := SubAssetInfo{TitledRecord: TitledRecord{ID: "s1"}, Type: const4assetus.CategoryVehicles, CountryID: "z"}
	if err := s.Validate(); err == nil {
		t.Error("expected rejection of invalid sub-asset countryID, got nil")
	}
}

func TestSubAssetInfoValidate_RejectsInvalidExpires(t *testing.T) {
	s := SubAssetInfo{TitledRecord: TitledRecord{ID: "s1"}, Type: const4assetus.CategoryVehicles, Expires: "31-12-2030"}
	if err := s.Validate(); err == nil {
		t.Error("expected rejection of malformed expires date, got nil")
	}
}

func TestSubAssetInfoValidate_OK(t *testing.T) {
	s := SubAssetInfo{TitledRecord: TitledRecord{ID: "s1"}, Type: const4assetus.CategoryVehicles, CountryID: "IE", Expires: "2030-12-31"}
	if err := s.Validate(); err != nil {
		t.Fatalf("valid sub-asset rejected: %v", err)
	}
}

func TestAssetGroupInfoValidate_RejectsMissingID(t *testing.T) {
	g := AssetGroupInfo{CategoryID: const4assetus.CategoryVehicles}
	if err := g.Validate(); err == nil {
		t.Error("expected missing-id error for group, got nil")
	}
}

func TestAssetGroupInfoValidate_RejectsInvalidCategory(t *testing.T) {
	g := AssetGroupInfo{TitledRecord: TitledRecord{ID: "g1"}, CategoryID: const4assetus.Category("bogus")}
	if err := g.Validate(); err == nil {
		t.Error("expected rejection of invalid group category, got nil")
	}
}

func TestAssetGroupInfoValidate_RejectsInvalidTotals(t *testing.T) {
	g := AssetGroupInfo{
		TitledRecord: TitledRecord{ID: "g1"},
		Totals:       []money.Amount{{Currency: "ZZZ", Value: decimal.NewDecimal64p2(1, 0)}},
	}
	if err := g.Validate(); err == nil {
		t.Error("expected rejection of group with unknown totals currency, got nil")
	}
}

func TestAssetGroupInfoValidate_OK(t *testing.T) {
	g := AssetGroupInfo{
		TitledRecord: TitledRecord{ID: "g1", Title: "Vehicles"},
		CategoryID:   const4assetus.CategoryVehicles,
		Totals:       []money.Amount{money.NewAmount(money.CurrencyEUR, decimal.NewDecimal64p2(50, 0))},
	}
	if err := g.Validate(); err != nil {
		t.Fatalf("valid group rejected: %v", err)
	}
}

// --- WithAssetRelationships.Validate error branches ---------------------

func TestWithAssetRelationshipsValidate_RejectsInvalidGroup(t *testing.T) {
	w := WithAssetRelationships{Group: &AssetGroupInfo{CategoryID: const4assetus.CategoryVehicles}}
	if err := w.Validate(); err == nil {
		t.Error("expected rejection of group missing id, got nil")
	}
}

func TestWithAssetRelationshipsValidate_RejectsInvalidSubAsset(t *testing.T) {
	w := WithAssetRelationships{SubAssets: []SubAssetInfo{{Type: const4assetus.CategoryVehicles}}}
	if err := w.Validate(); err == nil {
		t.Error("expected rejection of sub-asset missing id, got nil")
	}
}

func TestWithAssetRelationshipsValidate_RejectsEmptyMemberID(t *testing.T) {
	w := WithAssetRelationships{MemberIDs: []string{"m1", "  "}}
	if err := w.Validate(); err == nil {
		t.Error("expected rejection of empty member ID, got nil")
	}
}

func TestWithAssetRelationshipsValidate_RejectsInvalidMemberInfo(t *testing.T) {
	w := WithAssetRelationships{MembersInfo: []TitledRecord{{ID: ""}}}
	if err := w.Validate(); err == nil {
		t.Error("expected rejection of member info missing id, got nil")
	}
}

// --- vehicle_record_dbo.go ----------------------------------------------

func TestVehicleFuelRecordValidate_NilIsOK(t *testing.T) {
	var fuel *VehicleFuelRecord
	if err := fuel.Validate(); err != nil {
		t.Fatalf("nil fuel record rejected: %v", err)
	}
}

func TestVehicleFuelRecordValidate_RejectsInvalidAmount(t *testing.T) {
	fuel := &VehicleFuelRecord{Amount: &money.Amount{Currency: "ZZZ", Value: decimal.NewDecimal64p2(10, 0)}}
	if err := fuel.Validate(); err == nil {
		t.Error("expected rejection of fuel record with unknown currency, got nil")
	}
}

func TestVehicleFuelRecordValidate_OK(t *testing.T) {
	amount := money.NewAmount(money.CurrencyEUR, decimal.NewDecimal64p2(60, 0))
	fuel := &VehicleFuelRecord{
		Volume:   decimal.NewDecimal64p2(45, 0),
		Unit:     "L",
		Amount:   &amount,
		FuelCost: decimal.NewDecimal64p2(1, 50),
		Currency: "EUR",
	}
	if err := fuel.Validate(); err != nil {
		t.Fatalf("valid fuel record rejected: %v", err)
	}
}

func TestVehicleMileageValidate_OK(t *testing.T) {
	m := &VehicleMileage{Value: 120000, Unit: "km"}
	if err := m.Validate(); err != nil {
		t.Fatalf("valid mileage rejected: %v", err)
	}
}

func TestVehicleRecordDboValidate_OK(t *testing.T) {
	amount := money.NewAmount(money.CurrencyEUR, decimal.NewDecimal64p2(60, 0))
	rec := VehicleRecordDbo{
		CreatedFields: validCreatedFields(),
		Fuel: &VehicleFuelRecord{
			Volume: decimal.NewDecimal64p2(45, 0),
			Unit:   "L",
			Amount: &amount,
		},
		Mileage: &VehicleMileage{Value: 120000, Unit: "km"},
	}
	if err := rec.Validate(); err != nil {
		t.Fatalf("valid vehicle record rejected: %v", err)
	}
}

func TestVehicleRecordDboValidate_RejectsMissingCreatedFields(t *testing.T) {
	rec := VehicleRecordDbo{Mileage: &VehicleMileage{Value: 1, Unit: "km"}}
	if err := rec.Validate(); err == nil {
		t.Error("expected rejection of vehicle record without created fields, got nil")
	}
}

func TestVehicleRecordDboValidate_RejectsInvalidFuel(t *testing.T) {
	rec := VehicleRecordDbo{
		CreatedFields: validCreatedFields(),
		Fuel:          &VehicleFuelRecord{Amount: &money.Amount{Currency: "ZZZ", Value: decimal.NewDecimal64p2(1, 0)}},
	}
	if err := rec.Validate(); err == nil {
		t.Error("expected rejection of vehicle record with invalid fuel, got nil")
	}
}

func ptr[T any](v T) *T { return &v }
