# Assetus components coverage (FE Task 5)

Disposition of EVERY legacy UI component from
`sneat-libs/libs/extensions/assetus/components/src/lib/` against the MVP
`ext-assetus` lib. None silently dropped.

Verifies: assetus-frontend-port#ac:all-components-accounted

## Approach notes

The legacy UI components are tightly coupled to the legacy
`@sneat/mod-assetus-core` context/DTO model (`IAssetContext`,
`IAssetVehicleContext`, `IAssetDwellingContext`, `IAssetDbo<category, extra>`
with an embedded `.extra`, `dbo.totals.{incomes,expenses}`) and to the legacy
local `AssetService` (`createAsset<Extra>(space, request)`, `deleteAsset`,
`updateAsset` with `assetCategory`/`regNumber`, `addVehicleRecord`). That legacy
model is a DIFFERENT shape from the MVP unified `IAssetDbo` (no `.extra`;
`totals: IMoney[]`; flat `AssetService`).

Per the port-oriented discipline ("adapt EXISTING components; keep the build
GREEN; a green build with an honest coverage table beats a broken build"), the
components were ported into `ext-assetus` with their source preserved, rewiring
only the local relative imports:

- legacy `../services` (legacy `AssetService` / request DTOs) ->
  `@sneat/ext-assetus-components` (the legacy lib is published in node_modules),
- legacy `../add-asset-base-component` -> `@sneat/ext-assetus-components`,
- sibling components -> the locally ported paths under
  `src/lib/components/...`.

The legacy peer libs they consume (`@sneat/mod-assetus-core`, `@sneat/ui`,
`@sneat/components`, `@sneat/space-models`, `@sneat/space-components`,
`@sneat/dto`, `@sneat/core`) all resolve from node_modules, so the ports build
under `strictTemplates`. Components whose ONLY blocker is the legacy
service/DTO model still build because that model is consumed from the published
`@sneat/ext-assetus-components` package rather than rewritten.

A full unified-DTO rewrite of these context-bound components is intentionally
out of scope for this task (it would be a rewrite, not a port, and would
fan out into the asset pages/services). It is tracked implicitly by the MVP's
own unified-DTO pages (asset-page, assets-page) which already consume the
unified `IAssetDbo` directly.

## Ported -> new path under `libs/ext-assetus/src/lib/components/`

| Legacy component | Ported to | Notes |
|---|---|---|
| period-segment | `period-segment/period-segment.component.ts` | Pure; consumes `@sneat/dto` `Period`. Smoke + logic spec added. |
| vehicle-engine | `vehicle-engine/vehicle-engine.component.ts` | |
| make-model-card | `make-model-card/make-model-card.component.ts` | |
| asset-reg-number-input | `asset-reg-number-input/asset-reg-number-input.component.ts` | `../services` -> `@sneat/ext-assetus-components`. |
| mileage-dialog | `mileage-dialog/mileage-dialog.component.ts` | `../services` -> `@sneat/ext-assetus-components`. |
| edit-dwelling-card | `edit-dwelling-card/edit-dwelling-card.component.ts` | Exports `AddDwellingCardComponent`. |
| real-estate-location | `real-estate-location/real-estate-location.component.ts` | |
| asset-dates | `asset-dates/asset-dates.component.ts` | Added missing Ionic imports for strictTemplates. |
| asset-possesion-card | `asset-possession-card/asset-possession-card.component.ts` | |
| asset-card | `asset-card/asset-card.component.ts` | |
| assets-list | `assets-list/assets-list.component.ts` | `../services` -> published lib; `toSorted` -> non-mutating `[...].sort()` (ES2022 target). |
| vehicle-card | `vehicle-card/vehicle-card.component.ts` | Composes ported make-model / engine / possession / reg-number. |
| asset-add/asset-add-vehicle | `asset-add/asset-add-vehicle/asset-add-vehicle.component.ts` | Extends published `AddAssetBaseComponent`. |
| asset-add/asset-add-dwelling | `asset-add/asset-add-dwelling/asset-add-dwelling.component.ts` | Extends published `AddAssetBaseComponent`. |
| asset-add/asset-add-document | `asset-add/asset-add-document/asset-add-document.component.ts` | Extends published `AddAssetBaseComponent`. |

`assets-list/asset-list-item.component.ts` — the legacy `assets-list` folder also
contains an unused `asset-list-item.component.ts` that is not referenced by
`assets-list.component` (the list renders items inline) nor exported by the
legacy index. Not ported (dead/unused legacy code); the list rendering it
duplicated is covered by the ported `assets-list`.

## Reconciled with existing MVP equivalent (not re-ported)

| Legacy component | MVP equivalent | Notes |
|---|---|---|
| asset-card / assets-list (navigation + edit) | `pages/asset/asset-page.component.ts`, `pages/assets/assets-page.component.ts` | The MVP pages already implement asset detail/edit/remove/transfer/history against the unified `IAssetDbo`. The legacy `asset-card`/`assets-list` are still ported above as presentational widgets, but the page-level behaviour is the MVP's. |

## Deferred -> liabilities / service-provider sibling Feature (reason recorded)

| Legacy component | Reason |
|---|---|
| asset-liabilities | Manages asset liabilities + service-types (`MultiSelectorComponent`, `notUsedServiceTypes`, legacy `AssetService` liability calls). Belongs to the liabilities sibling Feature. |
| asset-add-service | "Add asset service" page bound to the service-provider flow (`AssetBasePage` + `AssetComponentBaseParams` service plumbing). Service-provider/liabilities sibling. |
| asset-contacts-group | Bound to `@sneat/contactus-core` (`IContact2Asset`, `ContactType`) — asset<->service-provider contact grouping. Service-provider-bound; defer to the liabilities/service sibling. |

## Shared bases / non-UI (accounted, not separately ported)

| Legacy file | Disposition |
|---|---|
| asset-add/add-asset-base-component.ts | Consumed from published `@sneat/ext-assetus-components` by the ported `asset-add-*` subclasses (legacy `createAsset<Extra>` flow). Not re-ported. |
| asset-base-page.ts | Service-provider page base used by `asset-add-service`. Deferred with the liabilities/service sibling. |
| asset-component-base-params.ts | Injectable params wrapper around the legacy `AssetService`; consumed from the published lib by the deferred bases. Not re-ported. |
| services/ (asset-service.ts, asset-service.dto.ts, assetus-space.service.ts, assetus-services.module.ts) | NOT UI components. The legacy service is consumed from published `@sneat/ext-assetus-components`; the MVP has its own `services/asset.service.ts` (FE Tasks 1-4). |
| testing/test-utils.ts | NOT a UI component (test harness). |

## Pages

Disposition of EVERY legacy page from
`sneat-libs/libs/extensions/assetus/pages/src/lib/` against the MVP
`ext-assetus` lib. None silently dropped.

Verifies: assetus-frontend-port#ac:all-pages-accounted

Ported pages follow the MVP convention: standalone Ionic components extending
`@sneat/space-components` `SpaceBaseComponent` with `SpaceComponentBaseParams`
+ a `ClassName` provider, plain Ion-header titles (the legacy
`SpacePageTitleComponent` is not exported by the installed
`@sneat/space-components`), and routes wired in `assetus-routing.ts`.

### Ported -> new path under `libs/ext-assetus/src/lib/pages/`

| Legacy page | Ported to | Route | Notes |
|---|---|---|---|
| new-asset | `pages/new-asset/new-asset-page.component.ts` | `new-asset` | Category picker -> ported `AssetAddVehicle/Dwelling/Document` components. Legacy `@sneat/mod-assetus-core` `IAssetCategory` replaced with a local shape over the MVP `AssetCategory` ids (`vehicles`/`dwelling`/`document`). |
| real-estates | `pages/real-estates/real-estates-page.component.ts` | `real-estates` | Properties list via the ported `AssetsListComponent` filtered to `assetType="dwelling"`. Legacy `AssetsBasePage`/`CommuneBasePageParams`/`IAssetService` plumbing replaced with `SpaceBaseComponent`; list does its own loading. |
| real-estate | `pages/real-estate/real-estate-page.component.ts` | `real-estate/:assetID` | Property detail rendering the ported `RealEstateLocationComponent` + a static finance section. Asset context read from nav state. Landlord/tenant `asset-contacts-group` deferred (see below). |
| asset-group | `pages/asset-group/asset-group-page.component.ts` | `asset-group` | Group detail via ported `PeriodSegmentComponent` + `AssetCardComponent`. Legacy injected the never-implemented `IAssetGroupService`/`IAssetService` stubs (the legacy file would not compile); those are dropped — group + assets are read from nav state. |
| optimization | `pages/optimization/optimization-page.component.ts` | `optimization` | Savings suggestions page. Static demo card preserved; legacy `CommuneBasePage`/`CommuneTopPage` replaced with `SpaceBaseComponent`. |

### Reconciled with existing MVP page (not re-ported)

| Legacy page | MVP page | Notes |
|---|---|---|
| assets | `pages/assets/assets-page.component.ts` | MVP already implements the per-space assets list (live `watchAssets`, New-asset dialog, Active/Archived filter) against the unified `IAssetDbo`. The legacy `assets` page (+ its `AssetsBasePage`) is superseded. |
| asset | `pages/asset/asset-page.component.ts` | MVP already implements asset detail/edit/archive/hard-delete/transfer/history against the unified `IAssetDbo`. The legacy `asset` page (+ `AssetBasePage`) is superseded. |

### Shared base pages (accounted, not separately ported)

| Legacy file | Disposition |
|---|---|
| assets-base.page.ts | Behaviour folded into the MVP `assets-page` (reconciled above). The `goNew(category)` navigation is reproduced inline by `real-estates-page`. Not re-ported as a base. |
| asset-base.page.ts | Legacy `CommuneBasePage`-derived base for `asset`/`real-estate`. Behaviour folded into the MVP `asset-page` and the ported `real-estate-page`. Also the base for the deferred service-provider pages (see below). Not re-ported as a base. |

### Deferred -> liabilities / service-provider sibling Feature (reason recorded)

| Legacy page | Reason |
|---|---|
| liability/liability-new | Creates a liability (debt) record; belongs to the liabilities sibling Feature. Its service-provider internals are intentionally not ported. |
| liability/select-service-provider | Service-provider picker built on the legacy `AssetBasePage` + `LiabilityServiceType` service-provider flow. Deferred to the liabilities/service-provider sibling. |
| real-estate landlord/tenant contacts | The legacy `real-estate` page rendered `asset-contacts-group` (landlords/tenants) bound to `@sneat/contactus-core` `IContact2Asset`. That component is already deferred to the liabilities/service sibling (see Components section); the contact groups are therefore omitted from the ported `real-estate-page`. |
