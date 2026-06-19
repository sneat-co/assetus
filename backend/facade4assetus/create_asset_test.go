package facade4assetus

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dal4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
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
