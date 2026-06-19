# ext-assetus

This library was generated with [Nx](https://nx.dev).

## Running unit tests

Run `nx test ext-assetus` to execute the unit tests.

## Known tech debt — legacy decoupling (TODO)

The Assetus frontend port (PR #3) ported the legacy components/pages into this
library, but several still **import from the legacy published packages** rather
than from this lib's own unified DTOs/services. Until that is resolved, the legacy
frontend libraries (`@sneat/mod-assetus-core`, `@sneat/ext-assetus-components` —
sourced from `sneat-libs`) **cannot be deleted**, because this lib depends on them
at build time.

### Current legacy coupling
- **`@sneat/mod-assetus-core`** — imported by ~17 source files. Provides legacy
  DTOs/types still consumed by ported components/pages: `IAssetContext`,
  `IAssetDwellingContext`, `IAssetBrief`, `AssetCategory`, `AssetType`,
  `IAssetDtoGroup`, `carMakes`/`IMake`/`IModel`.
- **`@sneat/ext-assetus-components`** — imported by ~15 source files. Provides the
  legacy `AssetService`, legacy `IUpdateAssetRequest`, and shared base components.

### Decoupling work required
1. Replace every `@sneat/mod-assetus-core` type import with the unified equivalents
   already defined in `src/lib/dto/asset.ts` / `dto/extras.ts` (`IAssetDbo`,
   `AssetCategory`, `AssetType`, `IAssetGroupInfo`, the typed extras, etc.). The
   `IAssetContext`/`IAssetDwellingContext` UI-context wrappers and the
   `carMakes`/`IMake`/`IModel` reference data must be ported into this lib.
2. Replace `@sneat/ext-assetus-components` `AssetService`/`IUpdateAssetRequest`
   usages with this lib's own `services/asset.service.ts` and its request DTOs.
3. Remove both legacy packages from this project's dependencies; confirm
   `nx build ext-assetus` and `nx test ext-assetus` stay green.
4. Only then can the legacy `sneat-libs` assetus libraries be retired.

Tracked in the GitHub issue linked from PR #3.
