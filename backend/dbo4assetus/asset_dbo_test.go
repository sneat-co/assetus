package dbo4assetus

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/with"
)

func validAssetBase() AssetBase {
	return AssetBase{
		Name:       "Harry Potter Collection",
		Category:   const4assetus.CategoryBooks,
		Condition:  const4assetus.ConditionGood,
		Status:     const4assetus.StatusActive,
		Visibility: const4assetus.VisibilityFamily,
	}
}

func TestAssetBaseValidate_OK(t *testing.T) {
	if err := validAssetBase().Validate(); err != nil {
		t.Fatalf("valid asset base rejected: %v", err)
	}
}

func TestAssetBaseValidate_RequiredName(t *testing.T) {
	a := validAssetBase()
	a.Name = "   "
	if err := a.Validate(); err == nil {
		t.Error("expected missing-name error, got nil")
	}
}

// AC: status-is-ownership-only — the write boundary rejects sharing states.
func TestAssetBaseValidate_RejectsSharingStatus(t *testing.T) {
	for _, bad := range []const4assetus.Status{"borrowed", "reserved"} {
		a := validAssetBase()
		a.Status = bad
		if err := a.Validate(); err == nil {
			t.Errorf("AssetBase.Validate() accepted sharing status %q, want rejection", bad)
		}
	}
}

func TestAssetDboValidate_OK(t *testing.T) {
	now := time.Now()
	dbo := &AssetDbo{
		AssetBase: validAssetBase(),
		WithSpaceIDs: dbmodels.WithSpaceIDs{
			SpaceIDs: []coretypes.SpaceID{"space1"},
		},
		WithModified: dbmodels.WithModified{
			CreatedFields: with.CreatedFields{
				CreatedAtField: with.CreatedAtField{CreatedAt: now},
				CreatedByField: with.CreatedByField{CreatedBy: "user1"},
			},
			UpdatedFields: with.UpdatedFields{UpdatedAt: now, UpdatedBy: "user1"},
		},
	}
	if err := dbo.Validate(); err != nil {
		t.Fatalf("valid AssetDbo rejected: %v", err)
	}
}

// AC: no-sharing-state-in-assetus — no borrow/lend/sell/swap/reserve/availability
// field and no ext.yardius are present on the persisted model.
func TestAssetModel_HasNoSharingOrYardiusFields(t *testing.T) {
	banned := []string{"borrow", "lend", "lent", "sell", "sold", "swap", "reserve", "reserved", "available", "availability", "yardius", "loan", "share", "sharing"}
	for _, typ := range []reflect.Type{
		reflect.TypeOf(AssetBase{}),
		reflect.TypeOf(AssetDbo{}),
		reflect.TypeOf(AssetBrief{}),
		reflect.TypeOf(AssetusSpaceDbo{}),
	} {
		for _, name := range fieldNames(typ) {
			lower := strings.ToLower(name)
			for _, b := range banned {
				if strings.Contains(lower, b) {
					t.Errorf("%s has forbidden sharing/availability field %q (contains %q)", typ.Name(), name, b)
				}
			}
		}
	}
}

func fieldNames(t reflect.Type) []string {
	var names []string
	if t.Kind() != reflect.Struct {
		return names
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		names = append(names, f.Name)
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			names = append(names, fieldNames(f.Type)...)
		}
	}
	return names
}

func TestNewOwnerRef_DerivesType(t *testing.T) {
	ref := NewOwnerRef("space1", coretypes.SpaceTypeFamily)
	if ref.SpaceID != "space1" {
		t.Errorf("SpaceID = %q, want space1", ref.SpaceID)
	}
	if ref.OwnerType != const4assetus.OwnerTypeFamily {
		t.Errorf("OwnerType = %q, want %q", ref.OwnerType, const4assetus.OwnerTypeFamily)
	}
}

func TestAssetHistoryEventValidate_TransferRequiresOwners(t *testing.T) {
	base := AssetHistoryEventBase{
		Type:       const4assetus.HistoryEventTransferred,
		OccurredAt: time.Now(),
		ActorRef:   "user1",
	}
	if err := base.Validate(); err == nil {
		t.Error("Transferred event without from/to owners should be rejected")
	}
	from := NewOwnerRef("spaceA", coretypes.SpaceTypeFamily)
	to := NewOwnerRef("spaceB", coretypes.SpaceTypeClub)
	base.FromOwner = &from
	base.ToOwner = &to
	if err := base.Validate(); err != nil {
		t.Errorf("valid Transferred event rejected: %v", err)
	}
}

func TestAssetHistoryEventValidate_NonTransfer(t *testing.T) {
	base := AssetHistoryEventBase{
		Type:       const4assetus.HistoryEventPurchased,
		OccurredAt: time.Now(),
		ActorRef:   "user1",
	}
	if err := base.Validate(); err != nil {
		t.Errorf("valid Purchased event rejected: %v", err)
	}
}
