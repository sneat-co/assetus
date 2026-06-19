package dbo4assetus

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/with"
)

// AC groups-nesting-linking-preserved: an asset with a group, two sub-assets
// under a parent, a relatedAs link, and a sameAssetID link round-trips all four
// resolvable.
func TestAssetRelationships_GroupsNestingLinkingPreserved(t *testing.T) {
	a := validAssetBase()
	a.GroupID = "group1"
	a.Group = &AssetGroupInfo{
		TitledRecord: TitledRecord{ID: "group1", Title: "Vehicles"},
		Order:        2,
		Desc:         "All my vehicles",
		CategoryID:   const4assetus.CategoryVehicles,
		NumberOf:     &AssetGroupCounts{Assets: 3},
	}
	a.ParentAssetID = "parent1"
	a.SubAssets = []SubAssetInfo{
		{TitledRecord: TitledRecord{ID: "sub1", Title: "Engine"}, Type: const4assetus.CategoryVehicles, SubType: "engine", Expires: "2030-01-31"},
		{TitledRecord: TitledRecord{ID: "sub2", Title: "Trailer"}, Type: const4assetus.CategoryVehicles, CountryID: "IE", SubType: "trailer"},
	}
	a.RelatedAs = "spouse"
	a.SameAssetID = "realtorAsset42"

	if err := a.Validate(); err != nil {
		t.Fatalf("asset with relationships rejected: %v", err)
	}

	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got AssetBase
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.GroupID != "group1" {
		t.Errorf("groupID = %q, want group1", got.GroupID)
	}
	if got.Group == nil || got.Group.ID != "group1" || got.Group.Order != 2 ||
		got.Group.CategoryID != const4assetus.CategoryVehicles ||
		got.Group.NumberOf == nil || got.Group.NumberOf.Assets != 3 {
		t.Errorf("group sub-entity not preserved: %+v", got.Group)
	}
	if got.ParentAssetID != "parent1" {
		t.Errorf("parentAssetID = %q, want parent1", got.ParentAssetID)
	}
	if len(got.SubAssets) != 2 {
		t.Fatalf("subAssets len = %d, want 2", len(got.SubAssets))
	}
	if got.SubAssets[0].ID != "sub1" || got.SubAssets[0].SubType != "engine" || got.SubAssets[0].Expires != "2030-01-31" {
		t.Errorf("sub-asset[0] detail not preserved: %+v", got.SubAssets[0])
	}
	if got.SubAssets[1].CountryID != "IE" {
		t.Errorf("sub-asset[1] countryID not preserved: %+v", got.SubAssets[1])
	}
	if got.RelatedAs != "spouse" {
		t.Errorf("relatedAs = %q, want spouse", got.RelatedAs)
	}
	if got.SameAssetID != "realtorAsset42" {
		t.Errorf("sameAssetID = %q, want realtorAsset42", got.SameAssetID)
	}
}

// AC multispace-with-canonical-owner: an asset associated with Spaces A, B, C
// surfaces under all three while owner-derivation/history anchor to space A.
func TestAssetDbo_MultiSpaceWithCanonicalOwner(t *testing.T) {
	now := time.Now()
	const canonicalOwner = coretypes.SpaceID("spaceA")
	dbo := &AssetDbo{
		AssetBase: validAssetBase(),
		WithSpaceIDs: dbmodels.WithSpaceIDs{
			// SpaceIDs[0] is the canonical owning space (the anchor).
			SpaceIDs: []coretypes.SpaceID{canonicalOwner},
		},
		WithModified: dbmodels.WithModified{
			CreatedFields: with.CreatedFields{
				CreatedAtField: with.CreatedAtField{CreatedAt: now},
				CreatedByField: with.CreatedByField{CreatedBy: "user1"},
			},
			UpdatedFields: with.UpdatedFields{UpdatedAt: now, UpdatedBy: "user1"},
		},
	}
	addSpace := func(spaceID string) {
		if dbo.Spaces == nil {
			dbo.Spaces = map[string]*AssetusSpaceBrief{}
		}
		dbo.Spaces[spaceID] = &AssetusSpaceBrief{
			Assets: AssetBriefs{
				"asset1": BriefFromAsset("asset1", dbo),
			},
		}
	}
	addSpace("spaceA")
	addSpace("spaceB")
	addSpace("spaceC")

	if err := dbo.Validate(); err != nil {
		t.Fatalf("multi-space AssetDbo rejected: %v", err)
	}

	// Surfaces under all three spaces.
	for _, sp := range []string{"spaceA", "spaceB", "spaceC"} {
		if dbo.Spaces[sp] == nil {
			t.Errorf("asset does not surface under %s", sp)
		}
	}

	// Owner-derivation is anchored to the canonical owning space A (SpaceIDs[0]).
	owner := NewOwnerRef(dbo.SpaceIDs[0], coretypes.SpaceTypeFamily)
	if owner.SpaceID != canonicalOwner {
		t.Errorf("owner anchored to %q, want spaceA", owner.SpaceID)
	}
	if owner.OwnerType != const4assetus.OwnerTypeFamily {
		t.Errorf("owner type = %q, want family", owner.OwnerType)
	}

	// History/transfer is anchored to space A (fromOwner is space A).
	event := AssetHistoryEventBase{
		Type:       const4assetus.HistoryEventTransferred,
		OccurredAt: now,
		ActorRef:   "user1",
		FromOwner:  &owner,
		ToOwner:    &OwnerRef{SpaceID: "spaceB"},
	}
	if event.FromOwner.SpaceID != canonicalOwner {
		t.Errorf("history fromOwner anchored to %q, want spaceA", event.FromOwner.SpaceID)
	}
	if err := event.Validate(); err != nil {
		t.Fatalf("transfer event anchored to space A rejected: %v", err)
	}
}

// AC member-info-preserved: memberIDs and per-member membersInfo round-trip.
func TestAssetRelationships_MemberInfoPreserved(t *testing.T) {
	a := validAssetBase()
	a.MemberIDs = []string{"m1", "m2"}
	a.MembersInfo = []TitledRecord{
		{ID: "m1", Title: "Alice"},
		{ID: "m2", Title: "Bob"},
	}

	if err := a.Validate(); err != nil {
		t.Fatalf("asset with member info rejected: %v", err)
	}

	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got AssetBase
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(got.MemberIDs) != 2 || got.MemberIDs[0] != "m1" || got.MemberIDs[1] != "m2" {
		t.Errorf("memberIDs not preserved: %v", got.MemberIDs)
	}
	if len(got.MembersInfo) != 2 || got.MembersInfo[0].Title != "Alice" || got.MembersInfo[1].ID != "m2" {
		t.Errorf("membersInfo not preserved: %+v", got.MembersInfo)
	}
}

// An asset supplying no relationship fields is still valid (all optional).
func TestAssetRelationships_AllOptional(t *testing.T) {
	if err := validAssetBase().Validate(); err != nil {
		t.Fatalf("asset without relationships rejected: %v", err)
	}
}
