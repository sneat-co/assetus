package facade4assetus

import (
	"context"
	"testing"
	"time"

	"github.com/dal-go/dalgo/adapters/dalgo2memory"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

// testUserID is the default acting user in facade tests.
const testUserID = "user1"

// newTestDBWithSpace builds an in-memory dalgo DB seeded with a single Space
// record (members = userIDs), wires facade.GetSneatDB to it, and returns the DB.
// Membership-gated facades read this Space to enforce access and derive owner
// type / default visibility.
func newTestDBWithSpace(t *testing.T, spaceID coretypes.SpaceID, spaceType coretypes.SpaceType, userIDs ...string) dal.DB {
	t.Helper()
	db := dalgo2memory.NewDB()
	now := time.Now()
	space := dbo4spaceus.NewSpaceEntry(spaceID)
	space.Data.Type = spaceType
	space.Data.Title = "Test " + string(spaceType) + " space"
	space.Data.Status = dbmodels.StatusActive
	space.Data.CreatedAt = now
	space.Data.CreatedBy = "seed"
	space.Data.IncreaseVersion(now, "seed")
	space.Data.UserIDs = userIDs
	if err := space.Data.Validate(); err != nil {
		t.Fatalf("seed space invalid: %v", err)
	}
	ctx := context.Background()
	if err := db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		return tx.Insert(ctx, space.Record)
	}); err != nil {
		t.Fatalf("failed to seed space: %v", err)
	}
	facade.GetSneatDB = func(context.Context) (dal.DB, error) { return db, nil }
	return db
}

// userCtx returns a ContextWithUser for the given user.
func userCtx(userID string) facade.ContextWithUser {
	return facade.NewContextWithUserID(context.Background(), userID)
}

// spaceRequest builds a dto4spaceus.SpaceRequest for the given space.
func spaceRequest(spaceID coretypes.SpaceID) dto4spaceus.SpaceRequest {
	return dto4spaceus.SpaceRequest{SpaceID: spaceID}
}
