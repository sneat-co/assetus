package const4assetus

import (
	"fmt"
	"slices"

	"github.com/strongo/validation"
)

// Category is a closed, write-validated set of asset categories. It is
// extensible by adding new members to the set below; values outside the set are
// rejected at the write boundary.
type Category string

const (
	CategoryBooks            Category = "books"
	CategoryGames            Category = "games"
	CategoryToys             Category = "toys"
	CategorySportsEquipment  Category = "sports_equipment"
	CategoryTools            Category = "tools"
	CategoryElectronics      Category = "electronics"
	CategoryClothing         Category = "clothing"
	CategoryVehicles         Category = "vehicles"
	CategoryCampingEquipment Category = "camping_equipment"
	CategoryOther            Category = "other"
)

// Categories is the closed set of valid categories.
var Categories = []Category{
	CategoryBooks,
	CategoryGames,
	CategoryToys,
	CategorySportsEquipment,
	CategoryTools,
	CategoryElectronics,
	CategoryClothing,
	CategoryVehicles,
	CategoryCampingEquipment,
	CategoryOther,
}

// IsValidCategory reports whether v is a member of the closed category set.
func IsValidCategory(v Category) bool {
	return slices.Contains(Categories, v)
}

// ValidateCategory returns an error if v is not a valid category. A category is
// required on every asset.
func ValidateCategory(v Category) error {
	if v == "" {
		return validation.NewErrRecordIsMissingRequiredField("category")
	}
	if !IsValidCategory(v) {
		return validation.NewErrBadRecordFieldValue("category",
			fmt.Sprintf("unknown category %q, expected one of %v", v, Categories))
	}
	return nil
}
