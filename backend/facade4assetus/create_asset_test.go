package facade4assetus

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dal4assetus"
	"github.com/sneat-co/assetus/backend/dbo4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/assetus/backend/extras4assetus"
	"github.com/sneat-co/sneat-core-modules/core/extra"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

// AC: member-creates-asset — a member of a family Space creates an asset; it is
// persisted owned by that Space with status=Active, condition as given, and
// visibility inherited from the Space default (Family).
func TestCreateAsset_MemberCreates(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	db := newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)

	resp, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		Name:         "Harry Potter Collection",
		Category:     const4assetus.CategoryBooks,
		Condition:    const4assetus.ConditionGood,
	})
	if err != nil {
		t.Fatalf("CreateAsset failed: %v", err)
	}
	if resp.ID == "" {
		t.Fatal("expected a non-empty asset ID")
	}
	// Re-read the persisted record.
	asset := dal4assetus.NewAssetEntry(spaceID, resp.ID)
	if err = db.Get(context.Background(), asset.Record); err != nil {
		t.Fatalf("failed to read created asset: %v", err)
	}
	if asset.Data.Status != const4assetus.StatusActive {
		t.Errorf("status = %q, want active", asset.Data.Status)
	}
	if asset.Data.Condition != const4assetus.ConditionGood {
		t.Errorf("condition = %q, want good", asset.Data.Condition)
	}
	if asset.Data.Visibility != const4assetus.VisibilityFamily {
		t.Errorf("visibility = %q, want family (inherited)", asset.Data.Visibility)
	}
	if !slices.Contains(asset.Data.SpaceIDs, spaceID) {
		t.Errorf("asset not owned by space %q (spaceIDs=%v)", spaceID, asset.Data.SpaceIDs)
	}
}

// AC: condition-optional-visibility-default — a legacy asset that carries no
// condition, created in a Space whose default visibility is family, with no
// condition and no visibility supplied: condition remains unset (valid) and
// visibility resolves to family (the owning Space default).
func TestCreateAsset_ConditionOptionalVisibilityDefault(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	db := newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)

	resp, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		Name:         "Legacy Item Without Condition",
		Category:     const4assetus.CategoryOther,
		// No Condition and no Visibility supplied.
	})
	if err != nil {
		t.Fatalf("CreateAsset failed: %v", err)
	}
	if resp.Asset.Condition != "" {
		t.Errorf("condition = %q, want unset", resp.Asset.Condition)
	}
	if resp.Asset.Visibility != const4assetus.VisibilityFamily {
		t.Errorf("visibility = %q, want family (Space default)", resp.Asset.Visibility)
	}
	// Re-read the persisted record to confirm it round-trips as valid.
	asset := dal4assetus.NewAssetEntry(spaceID, resp.ID)
	if err = db.Get(context.Background(), asset.Record); err != nil {
		t.Fatalf("failed to read created asset: %v", err)
	}
	if asset.Data.Condition != "" {
		t.Errorf("persisted condition = %q, want unset", asset.Data.Condition)
	}
	if asset.Data.Visibility != const4assetus.VisibilityFamily {
		t.Errorf("persisted visibility = %q, want family", asset.Data.Visibility)
	}
}

// AC: non-member-cannot-create — a non-member is rejected and no record persists.
func TestCreateAsset_NonMemberRejected(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	db := newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, "someone-else")

	resp, err := CreateAsset(userCtx("intruder"), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		Name:         "Sneaky Asset",
		Category:     const4assetus.CategoryOther,
		Condition:    const4assetus.ConditionGood,
	})
	if err == nil {
		t.Fatal("expected non-member create to be rejected")
	}
	if !errors.Is(err, facade.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got: %v", err)
	}
	if resp.ID != "" {
		// Best-effort: ensure nothing was written under a returned ID.
		asset := dal4assetus.NewAssetEntry(spaceID, resp.ID)
		if getErr := db.Get(context.Background(), asset.Record); getErr == nil {
			t.Error("an asset record was persisted despite rejection")
		}
	}
}

// AC: visibility-inherits-and-overrides — one asset inherits the Space default
// (Family), another overrides to Private.
func TestCreateAsset_VisibilityInheritAndOverride(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	_ = newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)

	inherited, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		Name:         "Inherited Visibility Item",
		Category:     const4assetus.CategoryGames,
		Condition:    const4assetus.ConditionNew,
	})
	if err != nil {
		t.Fatalf("CreateAsset (inherited) failed: %v", err)
	}
	if inherited.Asset.Visibility != const4assetus.VisibilityFamily {
		t.Errorf("inherited visibility = %q, want family", inherited.Asset.Visibility)
	}

	override, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		Name:         "Private Item",
		Category:     const4assetus.CategoryElectronics,
		Condition:    const4assetus.ConditionExcellent,
		Visibility:   const4assetus.VisibilityPrivate,
	})
	if err != nil {
		t.Fatalf("CreateAsset (override) failed: %v", err)
	}
	if override.Asset.Visibility != const4assetus.VisibilityPrivate {
		t.Errorf("override visibility = %q, want private", override.Asset.Visibility)
	}
}

// vehicleExtraRequest builds an extra.WithExtraField carrying a vehicle extra.
func vehicleExtraRequest(t *testing.T) extra.WithExtraField {
	t.Helper()
	var ef extra.WithExtraField
	vehicle := &extras4assetus.AssetVehicleExtra{
		WithMakeModelRegNumberFields: extras4assetus.WithMakeModelRegNumberFields{
			WithMakeModelFields: extras4assetus.WithMakeModelFields{Make: "VW", Model: "Golf"},
		},
		Vin: "WVWZZZ1JZXW000001",
		WithEngineData: extras4assetus.WithEngineData{
			EngineType: const4assetus.EngineTypeCombustion,
			EngineFuel: const4assetus.FuelTypePetrol,
			EngineCC:   1984,
		},
	}
	if err := ef.SetExtra(extras4assetus.AssetExtraTypeVehicle, vehicle); err != nil {
		t.Fatalf("SetExtra failed: %v", err)
	}
	return ef
}

// AC: optional-legacy-fields-roundtrip + leasing-asset-representable — a create
// with the full rich unified field set (incl. a leasing possession and a vehicle
// extra with engine data) persists with every field set and round-trips on read.
func TestCreateAsset_FullRichFieldSet(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	db := newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)

	year := 2018
	acq := time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC)
	req := dto4assetus.CreateAssetRequest{
		SpaceRequest:       spaceRequest(spaceID),
		Name:               "Family Car",
		Description:        "The leased family car",
		Category:           const4assetus.CategoryVehicles,
		Condition:          const4assetus.ConditionGood,
		Visibility:         const4assetus.VisibilityFamily,
		AcquisitionDate:    &acq,
		PurchasePrice:      &dbo4assetus.MonetaryAmount{Currency: "EUR", Value: 25000},
		EstimatedValue:     &dbo4assetus.MonetaryAmount{Currency: "EUR", Value: 18000},
		Location:           "Garage",
		Notes:              "Service due soon",
		Tags:               []string{"car", "family"},
		Type:               const4assetus.TypeVehicleCar,
		Possession:         const4assetus.PossessionLeasing,
		CountryID:          "IE",
		ParentCategoryID:   const4assetus.CategoryVehicles,
		YearOfBuild:        &year,
		IsRequest:          false,
		Geo:                &dbo4assetus.GeoPoint{Lat: 53.3, Lng: -6.2},
		AssetDates:         dbo4assetus.AssetDates{DateOfPurchase: "2020-05-01"},
		CanHaveExpense:     true,
		FinancialDirection: "expense",
		WithAssetRelationships: dbo4assetus.WithAssetRelationships{
			GroupID:   "g1",
			MemberIDs: []string{"m1", "m2"},
		},
		WithExtraField: vehicleExtraRequest(t),
	}
	resp, err := CreateAsset(userCtx(testUserID), req)
	if err != nil {
		t.Fatalf("CreateAsset failed: %v", err)
	}

	// Re-read the persisted record to confirm a full round-trip.
	asset := dal4assetus.NewAssetEntry(spaceID, resp.ID)
	if err = db.Get(context.Background(), asset.Record); err != nil {
		t.Fatalf("failed to read created asset: %v", err)
	}
	got := asset.Data
	if got.Type != const4assetus.TypeVehicleCar {
		t.Errorf("type = %q, want car", got.Type)
	}
	if got.Possession != const4assetus.PossessionLeasing {
		t.Errorf("possession = %q, want leasing", got.Possession)
	}
	if got.CountryID != "IE" {
		t.Errorf("countryID = %q, want IE", got.CountryID)
	}
	if got.ParentCategoryID != const4assetus.CategoryVehicles {
		t.Errorf("parentCategoryID = %q, want vehicles", got.ParentCategoryID)
	}
	if got.YearOfBuild == nil || *got.YearOfBuild != year {
		t.Errorf("yearOfBuild = %v, want %d", got.YearOfBuild, year)
	}
	if got.Geo == nil || got.Geo.Lat != 53.3 {
		t.Errorf("geo = %v, want lat 53.3", got.Geo)
	}
	if got.DateOfPurchase != "2020-05-01" {
		t.Errorf("dateOfPurchase = %q, want 2020-05-01", got.DateOfPurchase)
	}
	if !got.CanHaveExpense || got.FinancialDirection != "expense" {
		t.Errorf("financial flags not persisted: canHaveExpense=%v direction=%q", got.CanHaveExpense, got.FinancialDirection)
	}
	if got.GroupID != "g1" || len(got.MemberIDs) != 2 {
		t.Errorf("relationships not persisted: groupID=%q memberIDs=%v", got.GroupID, got.MemberIDs)
	}
	if got.ExtraType != extras4assetus.AssetExtraTypeVehicle {
		t.Errorf("extraType = %q, want vehicle", got.ExtraType)
	}
	extraData, err := got.GetExtraData()
	if err != nil {
		t.Fatalf("GetExtraData failed: %v", err)
	}
	vehicle, ok := extraData.(*extras4assetus.AssetVehicleExtra)
	if !ok {
		t.Fatalf("extra is %T, want *AssetVehicleExtra", extraData)
	}
	if vehicle.Vin != "WVWZZZ1JZXW000001" || vehicle.EngineFuel != const4assetus.FuelTypePetrol {
		t.Errorf("vehicle extra not round-tripped: vin=%q fuel=%q", vehicle.Vin, vehicle.EngineFuel)
	}
}

// AC: status settable on create — status=draft persists; empty status -> active.
func TestCreateAsset_StatusDraftAndDefault(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	_ = newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)

	draft, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		Name:         "Draft Asset",
		Category:     const4assetus.CategoryOther,
		Status:       const4assetus.StatusDraft,
	})
	if err != nil {
		t.Fatalf("CreateAsset (draft) failed: %v", err)
	}
	if draft.Asset.Status != const4assetus.StatusDraft {
		t.Errorf("status = %q, want draft", draft.Asset.Status)
	}

	active, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		Name:         "Active Asset",
		Category:     const4assetus.CategoryOther,
	})
	if err != nil {
		t.Fatalf("CreateAsset (default) failed: %v", err)
	}
	if active.Asset.Status != const4assetus.StatusActive {
		t.Errorf("status = %q, want active (default)", active.Asset.Status)
	}
}

// AC: backward-compatible — a create with only flat fields still works.
func TestCreateAsset_FlatOnlyStillWorks(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	_ = newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)

	resp, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		Name:         "Just A Name",
		Category:     const4assetus.CategoryOther,
	})
	if err != nil {
		t.Fatalf("CreateAsset (flat-only) failed: %v", err)
	}
	if resp.Asset.Status != const4assetus.StatusActive {
		t.Errorf("status = %q, want active", resp.Asset.Status)
	}
	if resp.Asset.ExtraType != "" {
		t.Errorf("expected no extra, got extraType=%q", resp.Asset.ExtraType)
	}
}

// AC: validation — out-of-category type and invalid possession are rejected.
func TestCreateAsset_RichValidationRejects(t *testing.T) {
	const spaceID coretypes.SpaceID = "family1"
	_ = newTestDBWithSpace(t, spaceID, coretypes.SpaceTypeFamily, testUserID)

	// A vehicle subtype on a non-vehicle category must be rejected.
	if _, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		Name:         "Bad Type",
		Category:     const4assetus.CategoryBooks,
		Type:         const4assetus.TypeVehicleCar,
	}); err == nil {
		t.Error("expected out-of-category type to be rejected")
	}

	// An invalid possession must be rejected.
	if _, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
		SpaceRequest: spaceRequest(spaceID),
		Name:         "Bad Possession",
		Category:     const4assetus.CategoryOther,
		Possession:   const4assetus.Possession("teleporting"),
	}); err == nil {
		t.Error("expected invalid possession to be rejected")
	}
}
