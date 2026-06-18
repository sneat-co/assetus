package facade4assetus

import (
	"testing"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/assetus/backend/dto4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// AC: owner-type-derived-existing-spaces — assets owned by private/family/club/
// company Spaces read back owner type Individual/Family/SportsClub/Organisation.
func TestGetAsset_OwnerTypeDerived_ExistingSpaceTypes(t *testing.T) {
	cases := []struct {
		spaceID   coretypes.SpaceID
		spaceType coretypes.SpaceType
		wantOwner const4assetus.OwnerType
	}{
		{"private1", coretypes.SpaceTypePrivate, const4assetus.OwnerTypeIndividual},
		{"family1", coretypes.SpaceTypeFamily, const4assetus.OwnerTypeFamily},
		{"club1", coretypes.SpaceTypeClub, const4assetus.OwnerTypeSportsClub},
		{"company1", coretypes.SpaceTypeCompany, const4assetus.OwnerTypeOrganisation},
	}
	for _, tc := range cases {
		t.Run(string(tc.spaceType), func(t *testing.T) {
			_ = newTestDBWithSpace(t, tc.spaceID, tc.spaceType, testUserID)
			created, err := CreateAsset(userCtx(testUserID), dto4assetus.CreateAssetRequest{
				SpaceRequest: spaceRequest(tc.spaceID),
				Name:         "Owned Item",
				Category:     const4assetus.CategoryOther,
				Condition:    const4assetus.ConditionGood,
			})
			if err != nil {
				t.Fatalf("CreateAsset failed: %v", err)
			}
			resp, err := GetAsset(userCtx(testUserID), dto4assetus.GetAssetRequest{
				SpaceRequest: spaceRequest(tc.spaceID),
				AssetID:      created.ID,
			})
			if err != nil {
				t.Fatalf("GetAsset failed: %v", err)
			}
			if resp.Owner.SpaceID != tc.spaceID {
				t.Errorf("owner spaceID = %q, want %q", resp.Owner.SpaceID, tc.spaceID)
			}
			if resp.Owner.OwnerType != tc.wantOwner {
				t.Errorf("owner type = %q, want %q", resp.Owner.OwnerType, tc.wantOwner)
			}
		})
	}
}

// AC: owner-type-derived-new-spaces — covered at the model level by
// const4assetus.DeriveOwnerType (community→Community, school→School) since the
// spaceus precondition for those space types has not shipped; an end-to-end read
// test will be added once spaceus exposes them (see backstage NEEDS-REVIEW).
