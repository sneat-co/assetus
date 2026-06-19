package dbo4assetus

import (
	"fmt"
	"strings"

	"github.com/crediterra/money"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/sneat-go-core/geo"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

// TitledRecord is the unified port of the legacy frontend ITitledRecord
// (an id plus an optional title). It backs the per-member info entries and is
// embedded by SubAssetInfo. Ported here rather than re-using a shared mixin so
// the assetus model stays self-contained.
type TitledRecord struct {
	ID    string `json:"id" firestore:"id"`
	Title string `json:"title,omitempty" firestore:"title,omitempty"`
}

// Validate returns an error if the record is missing its required ID.
func (v TitledRecord) Validate() error {
	if strings.TrimSpace(v.ID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	return nil
}

// SubAssetInfo is the per-sub-asset detail of a parent/sub-asset nesting,
// ported from the legacy frontend ISubAssetInfo (type/countryId/subType/expires
// plus the embedded ITitledRecord id/title).
type SubAssetInfo struct {
	TitledRecord
	Type      const4assetus.Category `json:"type" firestore:"type"`
	CountryID geo.CountryAlpha2      `json:"countryID,omitempty" firestore:"countryID,omitempty"`
	SubType   string                 `json:"subType,omitempty" firestore:"subType,omitempty"`
	// Expires is an optional ISO-date (YYYY-MM-DD) string.
	Expires string `json:"expires,omitempty" firestore:"expires,omitempty"`
}

// Validate returns an error if the sub-asset detail is not valid.
func (v SubAssetInfo) Validate() error {
	if err := v.TitledRecord.Validate(); err != nil {
		return err
	}
	if err := const4assetus.ValidateCategory(v.Type); err != nil {
		return err
	}
	if v.CountryID != "" && !geo.IsValidCountryAlpha2(v.CountryID) {
		return validation.NewErrBadRecordFieldValue("countryID",
			fmt.Sprintf("invalid country alpha-2 code: %q", v.CountryID))
	}
	if v.Expires != "" {
		if _, err := validate.DateString(v.Expires); err != nil {
			return validation.NewErrBadRecordFieldValue("expires", err.Error())
		}
	}
	return nil
}

// AssetGroupCounts is the optional per-group counts, ported from the legacy
// frontend IAssetDtoGroupCounts.
type AssetGroupCounts struct {
	Assets int `json:"assets,omitempty" firestore:"assets,omitempty"`
}

// AssetGroupInfo is the asset group as a sub-entity, ported from the legacy
// frontend IAssetDtoGroup (order/desc/categoryId/numberOf plus the embedded
// ITitledRecord id/title). It carries the full group shape, not just a groupId.
type AssetGroupInfo struct {
	TitledRecord
	Order      int                    `json:"order,omitempty" firestore:"order,omitempty"`
	Desc       string                 `json:"desc,omitempty" firestore:"desc,omitempty"`
	CategoryID const4assetus.Category `json:"categoryID,omitempty" firestore:"categoryID,omitempty"`
	NumberOf   *AssetGroupCounts      `json:"numberOf,omitempty" firestore:"numberOf,omitempty"`
	// Totals is the legacy ITotalsHolder dimension of a group (optional).
	Totals []money.Amount `json:"totals,omitempty" firestore:"totals,omitempty"`
}

// Validate returns an error if the group sub-entity is not valid.
func (v AssetGroupInfo) Validate() error {
	if err := v.TitledRecord.Validate(); err != nil {
		return err
	}
	if v.CategoryID != "" {
		if err := const4assetus.ValidateCategory(v.CategoryID); err != nil {
			return err
		}
	}
	for i, amount := range v.Totals {
		if err := amount.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("totals[%d]", i), err.Error())
		}
	}
	return nil
}

// WithAssetRelationships carries the optional relationship dimension of an asset,
// ported from the legacy AssetBrief/AssetDbo and frontend DTO. Every field is
// optional; an MVP-shaped asset supplying none of them is still valid.
//
//   - GroupID + Group:    group membership (linkage + the group sub-entity).
//   - ParentAssetID +
//     SubAssets:          parent/sub-asset nesting with per-sub-asset detail.
//   - SameAssetID:        link to the same underlying asset (realtor's/tenant's
//     asset ID) — "same underlying asset".
//   - WithOptionalRelatedAs.RelatedAs:
//     link to a related asset.
//   - MemberIDs + MembersInfo:
//     per-asset member IDs and per-member info.
type WithAssetRelationships struct {
	// GroupID links the asset to its group; Group carries the group sub-entity
	// detail (order/desc/categoryId/numberOf) when denormalized onto the asset.
	GroupID string          `json:"groupID,omitempty" firestore:"groupID,omitempty"`
	Group   *AssetGroupInfo `json:"group,omitempty" firestore:"group,omitempty"`

	// ParentAssetID + SubAssets express the parent/sub-asset nesting.
	ParentAssetID string         `json:"parentAssetID,omitempty" firestore:"parentAssetID,omitempty"`
	SubAssets     []SubAssetInfo `json:"subAssets,omitempty" firestore:"subAssets,omitempty"`

	// SameAssetID links to the same underlying asset (e.g. a realtor's or
	// tenant's asset ID for the same physical asset).
	SameAssetID string `json:"sameAssetID,omitempty" firestore:"sameAssetID,omitempty"`

	// RelatedAs links the asset to a related asset (ported legacy relatedAs).
	dbmodels.WithOptionalRelatedAs

	// MemberIDs + MembersInfo carry the per-asset member IDs and per-member info.
	MemberIDs   []string       `json:"memberIDs,omitempty" firestore:"memberIDs,omitempty"`
	MembersInfo []TitledRecord `json:"membersInfo,omitempty" firestore:"membersInfo,omitempty"`
}

// Validate returns an error if any present relationship field is not valid.
// All fields are optional: an empty WithAssetRelationships is valid.
func (v WithAssetRelationships) Validate() error {
	if err := v.WithOptionalRelatedAs.Validate(); err != nil {
		return err
	}
	if v.Group != nil {
		if err := v.Group.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("group", err.Error())
		}
	}
	for i, sub := range v.SubAssets {
		if err := sub.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("subAssets[%d]", i), err.Error())
		}
	}
	for i, id := range v.MemberIDs {
		if strings.TrimSpace(id) == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("memberIDs[%d]", i), "member ID must not be empty")
		}
	}
	for i, m := range v.MembersInfo {
		if err := m.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("membersInfo[%d]", i), err.Error())
		}
	}
	return nil
}
