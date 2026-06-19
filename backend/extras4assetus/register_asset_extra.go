// Package extras4assetus carries the polymorphic typed extras for an asset:
// the vehicle, dwelling and document extras resolved by extraType. It builds on
// the flat core (dbo4assetus.AssetBase) and reuses const4assetus enums.
//
// An asset with no extra stays valid: the core/extra registry falls back to a
// no-op extra when no factory is registered for the (empty) extraType.
package extras4assetus

import (
	"github.com/sneat-co/sneat-core-modules/core/extra"
)

// Asset extra type identifiers, used as the extraType discriminator on
// extra.WithExtraField.
const (
	AssetExtraTypeVehicle  extra.Type = "vehicle"
	AssetExtraTypeDwelling extra.Type = "dwelling"
	AssetExtraTypeDocument extra.Type = "document"
)

// AssetExtra is the interface every typed asset extra implements. It is the
// core/extra.Data contract (RequiredFields/IndexedFields/GetBrief/Validate)
// plus per-extra validation hooks.
type AssetExtra = extra.Data

// assetExtraFactories is the registry of typed-extra factories keyed by
// extraType, mirroring the legacy RegisterAssetExtraFactory/NewAssetExtra
// mechanism. Each factory is ALSO registered with the shared core/extra
// registry so extra.WithExtraField.GetExtraData resolves the typed extra.
var assetExtraFactories = map[extra.Type]func() AssetExtra{}

// RegisterAssetExtraFactory registers a typed-extra factory by extraType.
func RegisterAssetExtraFactory(t extra.Type, f func() AssetExtra) {
	assetExtraFactories[t] = f
	extra.RegisterFactory(t, func() extra.Data { return f() })
}

// NewAssetExtra creates a typed extra for the given extraType, or nil if no
// factory is registered.
func NewAssetExtra(t extra.Type) AssetExtra {
	if f, ok := assetExtraFactories[t]; ok {
		return f()
	}
	return nil
}
