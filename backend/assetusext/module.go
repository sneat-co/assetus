package assetusext

import (
	"github.com/sneat-co/assetus/backend/api4assetus"
	"github.com/sneat-co/assetus/backend/const4assetus"
	"github.com/sneat-co/sneat-go-core/extension"
)

// Extension returns the assetus extension config (routes + module ID).
func Extension() extension.Config {
	return extension.NewExtension(const4assetus.ExtensionID,
		extension.RegisterRoutes(api4assetus.RegisterHttpRoutes),
	)
}
