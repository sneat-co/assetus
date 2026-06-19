package dto4assetus

import (
	"testing"

	"github.com/sneat-co/assetus/backend/const4assetus"
)

func TestCreateAssetRequest_Validate(t *testing.T) {
	valid := CreateAssetRequest{Name: "Book", Category: const4assetus.CategoryBooks, Condition: const4assetus.ConditionGood}
	valid.SpaceID = "s1"
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid request rejected: %v", err)
	}

	noName := valid
	noName.Name = "  "
	if err := noName.Validate(); err == nil {
		t.Error("expected missing-name rejection")
	}

	badCat := valid
	badCat.Category = "not-a-category"
	if err := badCat.Validate(); err == nil {
		t.Error("expected invalid-category rejection")
	}

	badVis := valid
	badVis.Visibility = "everyone"
	if err := badVis.Validate(); err == nil {
		t.Error("expected invalid-visibility rejection")
	}

	noSpace := valid
	noSpace.SpaceID = ""
	if err := noSpace.Validate(); err == nil {
		t.Error("expected missing-spaceID rejection")
	}

	// Status is OPTIONAL on create: a valid status (draft) is accepted.
	draft := valid
	draft.Status = const4assetus.StatusDraft
	if err := draft.Validate(); err != nil {
		t.Errorf("draft status rejected: %v", err)
	}

	// An invalid status is rejected when supplied.
	badStatus := valid
	badStatus.Status = "borrowed"
	if err := badStatus.Validate(); err == nil {
		t.Error("expected invalid-status rejection")
	}
}

func TestUpdateAssetRequest_Validate(t *testing.T) {
	valid := UpdateAssetRequest{AssetID: "a1", Name: "Book", Category: const4assetus.CategoryBooks, Condition: const4assetus.ConditionGood, Visibility: const4assetus.VisibilityFamily}
	valid.SpaceID = "s1"
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid request rejected: %v", err)
	}
	noAsset := valid
	noAsset.AssetID = ""
	if err := noAsset.Validate(); err == nil {
		t.Error("expected missing-assetID rejection")
	}
	badCond := valid
	badCond.Condition = "mint"
	if err := badCond.Validate(); err == nil {
		t.Error("expected invalid-condition rejection")
	}
}

func TestRemoveAndGetAndHistoryRequests_Validate(t *testing.T) {
	r := RemoveAssetRequest{AssetID: "a1"}
	r.SpaceID = "s1"
	if err := r.Validate(); err != nil {
		t.Errorf("valid remove rejected: %v", err)
	}
	r.AssetID = ""
	if err := r.Validate(); err == nil {
		t.Error("expected missing-assetID rejection on remove")
	}

	g := GetAssetRequest{AssetID: "a1"}
	g.SpaceID = "s1"
	if err := g.Validate(); err != nil {
		t.Errorf("valid get rejected: %v", err)
	}

	h := GetHistoryRequest{AssetID: "a1"}
	h.SpaceID = "s1"
	if err := h.Validate(); err != nil {
		t.Errorf("valid get-history rejected: %v", err)
	}

	rec := RecordHistoryEventRequest{AssetID: "a1", Type: const4assetus.HistoryEventRepaired}
	rec.SpaceID = "s1"
	if err := rec.Validate(); err != nil {
		t.Errorf("valid record rejected: %v", err)
	}
	recTransfer := rec
	recTransfer.Type = const4assetus.HistoryEventTransferred
	if err := recTransfer.Validate(); err == nil {
		t.Error("expected Transferred type to be rejected by generic record endpoint")
	}
}

func TestTransferAssetRequest_Validate(t *testing.T) {
	v := TransferAssetRequest{AssetID: "a1", ToSpaceID: "s2"}
	v.SpaceID = "s1"
	if err := v.Validate(); err != nil {
		t.Fatalf("valid transfer rejected: %v", err)
	}
	same := v
	same.ToSpaceID = "s1"
	if err := same.Validate(); err == nil {
		t.Error("expected same-space transfer rejection")
	}
	noDest := v
	noDest.ToSpaceID = ""
	if err := noDest.Validate(); err == nil {
		t.Error("expected missing-destination rejection")
	}
}
