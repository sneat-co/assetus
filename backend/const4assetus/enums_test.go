package const4assetus

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/coretypes"
)

// AC: status-is-ownership-only — sharing/availability values must be rejected.
func TestValidateStatus(t *testing.T) {
	valid := []Status{StatusActive, StatusTransferred, StatusArchived, StatusDisposed, StatusLost}
	for _, s := range valid {
		if err := ValidateStatus(s); err != nil {
			t.Errorf("ValidateStatus(%q) unexpected error: %v", s, err)
		}
	}
	// Sharing/availability states are NOT Assetus statuses and must be rejected.
	for _, s := range []Status{"borrowed", "reserved", "lent", "sold-elsewhere", "bogus"} {
		if err := ValidateStatus(s); err == nil {
			t.Errorf("ValidateStatus(%q) expected rejection, got nil", s)
		}
	}
	if err := ValidateStatus(""); err == nil {
		t.Error("ValidateStatus(\"\") expected missing-field error, got nil")
	}
}

func TestValidateCategory(t *testing.T) {
	for _, c := range Categories {
		if err := ValidateCategory(c); err != nil {
			t.Errorf("ValidateCategory(%q) unexpected error: %v", c, err)
		}
	}
	if err := ValidateCategory("not-a-category"); err == nil {
		t.Error("ValidateCategory(invalid) expected error, got nil")
	}
	if err := ValidateCategory(""); err == nil {
		t.Error("ValidateCategory(empty) expected error, got nil")
	}
}

func TestValidateCondition(t *testing.T) {
	for _, c := range Conditions {
		if err := ValidateCondition(c); err != nil {
			t.Errorf("ValidateCondition(%q) unexpected error: %v", c, err)
		}
	}
	if err := ValidateCondition("mint"); err == nil {
		t.Error("ValidateCondition(invalid) expected error, got nil")
	}
	if err := ValidateCondition(""); err == nil {
		t.Error("ValidateCondition(empty) expected error, got nil")
	}
}

func TestValidateVisibility(t *testing.T) {
	for _, v := range Visibilities {
		if err := ValidateVisibility(v); err != nil {
			t.Errorf("ValidateVisibility(%q) unexpected error: %v", v, err)
		}
	}
	if err := ValidateVisibility("everyone-on-earth"); err == nil {
		t.Error("ValidateVisibility(invalid) expected error, got nil")
	}
	if err := ValidateVisibility(""); err == nil {
		t.Error("ValidateVisibility(empty) expected error, got nil")
	}
}

func TestValidateHistoryEventType(t *testing.T) {
	for _, e := range HistoryEventTypes {
		if err := ValidateHistoryEventType(e); err != nil {
			t.Errorf("ValidateHistoryEventType(%q) unexpected error: %v", e, err)
		}
	}
	if err := ValidateHistoryEventType("borrowed"); err == nil {
		t.Error("ValidateHistoryEventType(invalid) expected error, got nil")
	}
}

// AC: owner-type-derived-existing-spaces (the four existing space types).
func TestDeriveOwnerType_ExistingSpaceTypes(t *testing.T) {
	cases := map[coretypes.SpaceType]OwnerType{
		coretypes.SpaceTypePrivate: OwnerTypeIndividual,
		coretypes.SpaceTypeFamily:  OwnerTypeFamily,
		coretypes.SpaceTypeClub:    OwnerTypeSportsClub,
		coretypes.SpaceTypeCompany: OwnerTypeOrganisation,
	}
	for spaceType, want := range cases {
		if got := DeriveOwnerType(spaceType); got != want {
			t.Errorf("DeriveOwnerType(%q) = %q, want %q", spaceType, got, want)
		}
	}
}

// AC: owner-type-derived-new-spaces — these map once spaceus ships the types.
func TestDeriveOwnerType_NewSpaceTypes(t *testing.T) {
	if got := DeriveOwnerType("community"); got != OwnerTypeCommunity {
		t.Errorf("DeriveOwnerType(community) = %q, want %q", got, OwnerTypeCommunity)
	}
	if got := DeriveOwnerType("school"); got != OwnerTypeSchool {
		t.Errorf("DeriveOwnerType(school) = %q, want %q", got, OwnerTypeSchool)
	}
	if got := DeriveOwnerType("totally-unknown"); got != OwnerTypeUnknown {
		t.Errorf("DeriveOwnerType(unknown) = %q, want unknown", got)
	}
}

func TestDefaultVisibilityForSpaceType(t *testing.T) {
	if got := DefaultVisibilityForSpaceType(coretypes.SpaceTypeFamily); got != VisibilityFamily {
		t.Errorf("default for family = %q, want %q", got, VisibilityFamily)
	}
	if got := DefaultVisibilityForSpaceType(coretypes.SpaceTypePrivate); got != VisibilityPrivate {
		t.Errorf("default for private = %q, want %q", got, VisibilityPrivate)
	}
	if got := DefaultVisibilityForSpaceType(coretypes.SpaceTypeClub); got != VisibilitySpecificSpace {
		t.Errorf("default for club = %q, want %q", got, VisibilitySpecificSpace)
	}
}
