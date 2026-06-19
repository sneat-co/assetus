package extras4assetus

import (
	"github.com/sneat-co/sneat-core-modules/core/extra"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

func init() {
	RegisterAssetExtraFactory(AssetExtraTypeDwelling, func() AssetExtra {
		return new(AssetDwellingExtra)
	})
}

var _ extra.Data = (*AssetDwellingExtra)(nil)

// AssetDwellingExtra is the typed extra for dwelling assets, ported from the
// legacy extras4assetus.AssetDwellingExtra: address, rent price, bedrooms, area.
type AssetDwellingExtra struct {
	Address   *dbmodels.Address `json:"address,omitempty" firestore:"address,omitempty"`
	RentPrice struct {
		Value    float64 `json:"value,omitempty" firestore:"value,omitempty"`
		Currency string  `json:"currency,omitempty" firestore:"currency,omitempty"`
	} `json:"rent_price,omitempty" firestore:"rent_price,omitempty"`
	NumberOfBedrooms int `json:"numberOfBedrooms,omitempty" firestore:"numberOfBedrooms,omitempty"`
	AreaSqM          int `json:"areaSqM,omitempty" firestore:"areaSqM,omitempty"`
}

// RequiredFields implements extra.Data.
func (v *AssetDwellingExtra) RequiredFields() []string {
	return nil
}

// IndexedFields implements extra.Data.
func (v *AssetDwellingExtra) IndexedFields() []string {
	return nil
}

// GetBrief implements extra.Data.
func (v *AssetDwellingExtra) GetBrief() extra.Data {
	return &AssetDwellingExtra{
		NumberOfBedrooms: v.NumberOfBedrooms,
		AreaSqM:          v.AreaSqM,
	}
}

// Validate returns an error if the dwelling extra is not valid.
func (v *AssetDwellingExtra) Validate() error {
	if v.Address != nil {
		if err := v.Address.Validate(); err != nil {
			return err
		}
	}
	if v.NumberOfBedrooms < 0 {
		return validation.NewErrBadRecordFieldValue("numberOfBedrooms", "negative value")
	}
	if v.AreaSqM < 0 {
		return validation.NewErrBadRecordFieldValue("areaSqM", "negative value")
	}
	if v.RentPrice.Value < 0 {
		return validation.NewErrBadRecordFieldValue("rent_price.value", "negative value")
	}
	return nil
}
