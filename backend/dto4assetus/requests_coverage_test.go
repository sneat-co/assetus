package dto4assetus

import (
	"testing"

	"github.com/sneat-co/assetus/backend/const4assetus"
)

func TestAddVehicleRecordRequest_Validate(t *testing.T) {
	valid := AddVehicleRecordRequest{AssetID: "a1", Mileage: 1000, MileageUnit: "km"}
	valid.SpaceID = "s1"
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid request rejected: %v", err)
	}

	noSpace := valid
	noSpace.SpaceID = ""
	if err := noSpace.Validate(); err == nil {
		t.Error("expected missing-spaceID rejection")
	}

	noAsset := valid
	noAsset.AssetID = ""
	if err := noAsset.Validate(); err == nil {
		t.Error("expected missing-assetID rejection")
	}
}

func TestGetAssetRequest_Validate_MissingAssetID(t *testing.T) {
	g := GetAssetRequest{AssetID: ""}
	g.SpaceID = "s1"
	if err := g.Validate(); err == nil {
		t.Error("expected missing-assetID rejection on get")
	}

	noSpace := GetAssetRequest{AssetID: "a1"}
	if err := noSpace.Validate(); err == nil {
		t.Error("expected missing-spaceID rejection on get")
	}
}

func TestGetHistoryRequest_Validate_ErrorBranches(t *testing.T) {
	noAsset := GetHistoryRequest{AssetID: ""}
	noAsset.SpaceID = "s1"
	if err := noAsset.Validate(); err == nil {
		t.Error("expected missing-assetID rejection on get-history")
	}

	noSpace := GetHistoryRequest{AssetID: "a1"}
	if err := noSpace.Validate(); err == nil {
		t.Error("expected missing-spaceID rejection on get-history")
	}
}

func TestRecordHistoryEventRequest_Validate_ErrorBranches(t *testing.T) {
	noSpace := RecordHistoryEventRequest{AssetID: "a1", Type: const4assetus.HistoryEventRepaired}
	if err := noSpace.Validate(); err == nil {
		t.Error("expected missing-spaceID rejection on record")
	}

	noAsset := RecordHistoryEventRequest{AssetID: "", Type: const4assetus.HistoryEventRepaired}
	noAsset.SpaceID = "s1"
	if err := noAsset.Validate(); err == nil {
		t.Error("expected missing-assetID rejection on record")
	}

	badType := RecordHistoryEventRequest{AssetID: "a1", Type: "not-a-type"}
	badType.SpaceID = "s1"
	if err := badType.Validate(); err == nil {
		t.Error("expected invalid-type rejection on record")
	}
}

func TestRemoveAssetRequest_Validate_MissingSpace(t *testing.T) {
	noSpace := RemoveAssetRequest{AssetID: "a1"}
	if err := noSpace.Validate(); err == nil {
		t.Error("expected missing-spaceID rejection on remove")
	}
}

func TestTransferAssetRequest_Validate_ErrorBranches(t *testing.T) {
	noSpace := TransferAssetRequest{AssetID: "a1", ToSpaceID: "s2"}
	if err := noSpace.Validate(); err == nil {
		t.Error("expected missing-spaceID rejection on transfer")
	}

	noAsset := TransferAssetRequest{AssetID: "", ToSpaceID: "s2"}
	noAsset.SpaceID = "s1"
	if err := noAsset.Validate(); err == nil {
		t.Error("expected missing-assetID rejection on transfer")
	}
}

func TestUpdateAssetRequest_Validate_ErrorBranches(t *testing.T) {
	valid := UpdateAssetRequest{
		AssetID:    "a1",
		Name:       "Book",
		Category:   const4assetus.CategoryBooks,
		Condition:  const4assetus.ConditionGood,
		Visibility: const4assetus.VisibilityFamily,
	}
	valid.SpaceID = "s1"

	noSpace := valid
	noSpace.SpaceID = ""
	if err := noSpace.Validate(); err == nil {
		t.Error("expected missing-spaceID rejection on update")
	}

	noName := valid
	noName.Name = "  "
	if err := noName.Validate(); err == nil {
		t.Error("expected missing-name rejection on update")
	}

	badCat := valid
	badCat.Category = "not-a-category"
	if err := badCat.Validate(); err == nil {
		t.Error("expected invalid-category rejection on update")
	}

	badVis := valid
	badVis.Visibility = "everyone"
	if err := badVis.Validate(); err == nil {
		t.Error("expected invalid-visibility rejection on update")
	}
}
