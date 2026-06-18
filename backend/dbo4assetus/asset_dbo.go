package dbo4assetus

import (
	"fmt"
	"strings"
	"time"

	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

// AssetsCollection is the Firestore collection name for assets under the
// assetus space module: /spaces/{spaceID}/ext/assetus/assets/{assetID}.
const AssetsCollection = "assets"

// MonetaryAmount is an optional monetary metadata value (purchase price /
// estimated value). Reserved-simple: a currency code plus a decimal value.
type MonetaryAmount struct {
	Currency string  `json:"currency,omitempty" firestore:"currency,omitempty"`
	Value    float64 `json:"value,omitempty" firestore:"value,omitempty"`
}

// AssetBase carries the user-editable fields of an asset. It is embedded by
// AssetDbo (the persisted record) and reused by the create/edit facades.
type AssetBase struct {
	Name        string                   `json:"name" firestore:"name"`
	Description string                   `json:"description,omitempty" firestore:"description,omitempty"`
	Category    const4assetus.Category   `json:"category" firestore:"category"`
	Condition   const4assetus.Condition  `json:"condition" firestore:"condition"`
	Status      const4assetus.Status     `json:"status" firestore:"status"`
	Visibility  const4assetus.Visibility `json:"visibility" firestore:"visibility"`

	// Optional metadata.
	AcquisitionDate *time.Time      `json:"acquisitionDate,omitempty" firestore:"acquisitionDate,omitempty"`
	PurchasePrice   *MonetaryAmount `json:"purchasePrice,omitempty" firestore:"purchasePrice,omitempty"`
	EstimatedValue  *MonetaryAmount `json:"estimatedValue,omitempty" firestore:"estimatedValue,omitempty"`
	Location        string          `json:"location,omitempty" firestore:"location,omitempty"`
	Notes           string          `json:"notes,omitempty" firestore:"notes,omitempty"`
	Tags            []string        `json:"tags,omitempty" firestore:"tags,omitempty"`

	// Photos is reserved for a fast-follow and is NOT implemented in the MVP.
	// It is declared so the schema does not need reshaping when photo
	// upload/storage lands.
	Photos []string `json:"photos,omitempty" firestore:"photos,omitempty"`
}

// Validate enforces the closed-enum write boundary plus required fields. This is
// the single enforcement point reused by the create and update facades — out-of-set
// values (e.g. status "borrowed"/"reserved") are rejected before persistence.
func (v AssetBase) Validate() error {
	if strings.TrimSpace(v.Name) == "" {
		return validation.NewErrRecordIsMissingRequiredField("name")
	}
	if err := const4assetus.ValidateCategory(v.Category); err != nil {
		return err
	}
	if err := const4assetus.ValidateCondition(v.Condition); err != nil {
		return err
	}
	if err := const4assetus.ValidateStatus(v.Status); err != nil {
		return err
	}
	if err := const4assetus.ValidateVisibility(v.Visibility); err != nil {
		return err
	}
	for i, tag := range v.Tags {
		if strings.TrimSpace(tag) == "" {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("tags[%d]", i), "tag must not be empty")
		}
	}
	return nil
}

// AssetDbo is the persisted asset aggregate at
// /spaces/{spaceID}/ext/assetus/assets/{assetID}.
//
// The owner is the owning Space (carried via WithSpaceIDs); there is no separate
// owner entity and no sharing/availability fields. ext.yardius is intentionally
// absent in the MVP.
type AssetDbo struct {
	AssetBase
	dbmodels.WithModified
	dbmodels.WithSpaceIDs
}

// Validate returns an error if the record is not valid.
func (v *AssetDbo) Validate() error {
	if err := v.WithSpaceIDs.Validate(); err != nil {
		return err
	}
	if err := v.WithModified.Validate(); err != nil {
		return err
	}
	if err := v.AssetBase.Validate(); err != nil {
		return err
	}
	return nil
}

// AssetBrief is a denormalized summary of an asset for listing on the module
// space entry. It deliberately carries no sharing/availability fields.
type AssetBrief struct {
	ID         string                   `json:"id" firestore:"id"`
	Name       string                   `json:"name" firestore:"name"`
	Category   const4assetus.Category   `json:"category" firestore:"category"`
	Condition  const4assetus.Condition  `json:"condition" firestore:"condition"`
	Status     const4assetus.Status     `json:"status" firestore:"status"`
	Visibility const4assetus.Visibility `json:"visibility" firestore:"visibility"`
}

// Validate returns an error if the brief is not valid.
func (v AssetBrief) Validate() error {
	if strings.TrimSpace(v.ID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if strings.TrimSpace(v.Name) == "" {
		return validation.NewErrRecordIsMissingRequiredField("name")
	}
	if err := const4assetus.ValidateCategory(v.Category); err != nil {
		return err
	}
	if err := const4assetus.ValidateCondition(v.Condition); err != nil {
		return err
	}
	if err := const4assetus.ValidateStatus(v.Status); err != nil {
		return err
	}
	if err := const4assetus.ValidateVisibility(v.Visibility); err != nil {
		return err
	}
	return nil
}

// BriefFromAsset builds an AssetBrief from a persisted asset.
func BriefFromAsset(id string, asset *AssetDbo) *AssetBrief {
	return &AssetBrief{
		ID:         id,
		Name:       asset.Name,
		Category:   asset.Category,
		Condition:  asset.Condition,
		Status:     asset.Status,
		Visibility: asset.Visibility,
	}
}
