package dbo4assetus

import (
	"encoding/json"
	"testing"

	"github.com/crediterra/money"
	"github.com/sneat-co/assetus/backend/const4assetus"
)

// AC optional-legacy-fields-roundtrip: a legacy asset with CountryID,
// dateInsuredTill, custom fields, tags, and a geo value round-trips unchanged,
// and an asset created without any of them is still valid (each is optional).
func TestAssetBase_OptionalLegacyFieldsRoundtrip(t *testing.T) {
	a := validAssetBase()
	a.CountryID = "IE"
	a.DateInsuredTill = "2030-01-31"
	a.WithCustomFields.FieldsStr = map[string]map[string]string{"color": {"en": "red"}}
	a.Tags = []string{"vintage", "rare"}
	a.Geo = &GeoPoint{Lat: 53.3498, Lng: -6.2603}

	if err := a.Validate(); err != nil {
		t.Fatalf("asset with all legacy fields rejected: %v", err)
	}

	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got AssetBase
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.CountryID != "IE" {
		t.Errorf("CountryID = %q, want IE", got.CountryID)
	}
	if got.DateInsuredTill != "2030-01-31" {
		t.Errorf("DateInsuredTill = %q, want 2030-01-31", got.DateInsuredTill)
	}
	if got.FieldsStr["color"]["en"] != "red" {
		t.Errorf("custom field color = %v, want red", got.FieldsStr["color"])
	}
	if len(got.Tags) != 2 || got.Tags[0] != "vintage" || got.Tags[1] != "rare" {
		t.Errorf("Tags = %v, want [vintage rare]", got.Tags)
	}
	if got.Geo == nil || got.Geo.Lat != 53.3498 || got.Geo.Lng != -6.2603 {
		t.Errorf("Geo = %v, want {53.3498 -6.2603}", got.Geo)
	}

	// Each of the five is optional: an MVP-shaped asset without any of them is valid.
	mvp := validAssetBase()
	if mvp.CountryID != "" || mvp.DateInsuredTill != "" || mvp.FieldsStr != nil || mvp.Tags != nil || mvp.Geo != nil {
		t.Fatal("MVP asset unexpectedly has legacy fields set")
	}
	if err := mvp.Validate(); err != nil {
		t.Fatalf("MVP asset (no legacy fields) rejected: %v", err)
	}
}

// AC leasing-asset-representable: a leasing asset stores Possession=leasing on
// the optional possession field with an independent Status, with no separate
// record type.
func TestAssetBase_LeasingAssetRepresentable(t *testing.T) {
	a := validAssetBase()
	a.Possession = const4assetus.PossessionLeasing
	a.Status = const4assetus.StatusActive // status independent of possession

	if err := a.Validate(); err != nil {
		t.Fatalf("leasing asset rejected: %v", err)
	}
	if a.Possession != const4assetus.PossessionLeasing {
		t.Errorf("Possession = %q, want leasing", a.Possession)
	}
	if a.Status != const4assetus.StatusActive {
		t.Errorf("Status = %q, want active", a.Status)
	}
	// Possession and Status round-trip independently.
	data, _ := json.Marshal(a)
	var got AssetBase
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Possession != const4assetus.PossessionLeasing || got.Status != const4assetus.StatusActive {
		t.Errorf("got possession=%q status=%q, want leasing/active", got.Possession, got.Status)
	}
}

// AC possession-defaults-owning: an MVP-shaped asset created with no possession
// value resolves to owning by default.
func TestAssetBase_PossessionDefaultsOwning(t *testing.T) {
	a := validAssetBase()
	if a.Possession != "" {
		t.Fatalf("expected empty possession on MVP asset, got %q", a.Possession)
	}
	if got := a.WithPossessionDefault(); got != const4assetus.PossessionOwning {
		t.Errorf("WithPossessionDefault() = %q, want owning", got)
	}
	// An explicitly-set possession is preserved (not overridden to owning).
	a.Possession = const4assetus.PossessionLeasing
	if got := a.WithPossessionDefault(); got != const4assetus.PossessionLeasing {
		t.Errorf("WithPossessionDefault() = %q, want leasing", got)
	}
}

// AC financial-fields-have-a-home: every legacy financial field has a recorded
// disposition as an optional unified field on the core asset, and they round-trip.
func TestAssetBase_FinancialFieldsHaveAHome(t *testing.T) {
	a := validAssetBase()
	a.Category = const4assetus.CategoryDebt // the debt category
	a.Totals = []money.Amount{money.NewAmount(money.CurrencyEUR, 12345)}
	a.CanHaveIncome = true
	a.CanHaveExpense = true
	a.FinancialDirection = "expense"
	a.Liabilities = []AssetLiabilityInfo{{ID: "liab1", ServiceTypes: []LiabilityServiceType{"electricity"}}}
	a.NotUsedServiceTypes = []LiabilityServiceType{"gas"}

	if err := a.Validate(); err != nil {
		t.Fatalf("asset with financial fields rejected: %v", err)
	}

	data, _ := json.Marshal(a)
	var got AssetBase
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Category != const4assetus.CategoryDebt {
		t.Errorf("Category = %q, want debt", got.Category)
	}
	if len(got.Totals) != 1 || got.Totals[0].Currency != money.CurrencyEUR || got.Totals[0].Value != 12345 {
		t.Errorf("Totals = %v, want one EUR 12345", got.Totals)
	}
	if !got.CanHaveIncome || !got.CanHaveExpense {
		t.Errorf("CanHaveIncome=%v CanHaveExpense=%v, want both true", got.CanHaveIncome, got.CanHaveExpense)
	}
	if got.FinancialDirection != "expense" {
		t.Errorf("FinancialDirection = %q, want expense", got.FinancialDirection)
	}
	if len(got.Liabilities) != 1 || got.Liabilities[0].ID != "liab1" {
		t.Errorf("Liabilities = %v, want one liab1", got.Liabilities)
	}
	if len(got.NotUsedServiceTypes) != 1 || got.NotUsedServiceTypes[0] != "gas" {
		t.Errorf("NotUsedServiceTypes = %v, want [gas]", got.NotUsedServiceTypes)
	}
}
