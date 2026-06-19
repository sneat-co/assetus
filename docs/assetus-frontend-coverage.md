# Assetus frontend capability-coverage (FE Task 7)

Maps **every** legacy Assetus frontend capability — each `core` DTO/interface/enum,
each of the ~19 components, each page, and each service — to its disposition in the
MVP workspace `assetus/frontend/libs/ext-assetus`. No legacy frontend capability is
left without a home: each row is either **Ported** (to a concrete, verified path),
**Reconciled** with an existing MVP equivalent, an **intentional change**, or
assigned to the **liabilities / service-provider sibling Feature** with a reason.

Verifies: assetus-frontend-port#ac:liabilities-disposition-recorded,
assetus-frontend-port#ac:frontend-capability-coverage-complete,
assetus-frontend-port#ac:frontend-builds-and-tests-pass

Legacy sources:
- DTOs / uimodels / data: `sneat-libs/libs/extensions/assetus/core/src/lib/`
- Components: `sneat-libs/libs/extensions/assetus/components/src/lib/`
- Pages: `sneat-apps/libs/extensions/assetus/pages/src/lib/`

Ported paths below are relative to
`assetus/frontend/libs/ext-assetus/src/lib/` and were each verified to exist.
The component/page dispositions are folded from
`assetus/frontend/libs/ext-assetus/COMPONENTS-COVERAGE.md` (FE Task 5).

---

## 1. Core DTOs — `core/src/lib/dto/`

### 1.1 `assetus-types.ts` (enums)

| Legacy capability | Disposition |
|---|---|
| `AssetStatus` (`active`/`archived`/`draft`) | Ported (reconciled with MVP) -> `dto/asset.ts` `AssetStatus` (MVP superset `draft`,`active`,`transferred`,`archived`,`disposed`,`lost`; legacy `draft` retained). |
| `AssetCategory` (`undefined`/`dwelling`/`vehicle`/`document`/`debt`/`misc`) | Ported (intentional reconciliation) -> `dto/asset.ts` `AssetCategory`. Backend reconciliation baked in: legacy `vehicle`->`vehicles`, `misc`->`other`; legacy `undefined` dropped (no backend value); `dwelling`/`document`/`debt` retained. |
| `AssetVehicleType` | Ported -> `dto/asset.ts` `AssetVehicleType` (legacy `bicycle` moved under `AssetSportsEquipmentType` per backend `Type`). |
| `AssetRealEstateType` (`house`/`apartment`/`land`) | Ported -> `dto/asset.ts` `AssetDwellingType` (superset incl. `office`/`shop`/`garage`/`warehouse`). |
| `AssetType` (union) | Ported -> `dto/asset.ts` `AssetType` (open string union over the per-category subtypes). |
| `EngineType` / `EngineTypes` enum + const exports | Ported -> `dto/asset.ts` `EngineType` (`''`,`other`,`combustion`,`electric`,`phev`,`hybrid`,`steam`). Legacy const aliases (`EngineTypeElectric` etc.) intentionally not re-exported — the MVP uses the string-literal union, not const aliases. |
| `FuelType` / `FuelTypes` enum + const exports | Ported -> `dto/asset.ts` `FuelType` (`''`,`other`,`bio`,`petrol`,`diesel`,`hydrogen`). Const aliases intentionally dropped (string union). |
| `FuelVolumeUnit` + `FuelVolumeUnitTypes` | Ported -> `dto/asset.ts` `FuelVolumeUnit` (open string per backend; serializes as plain string). |
| `MileageUnit` + `MileageUnitTypes` | Ported -> `dto/asset.ts` `MileageUnit` (open string per backend). |
| `AssetPossession` + `AssetPossessions` + const exports | Ported -> `dto/asset.ts` `AssetPossession` (`unknown`,`undisclosed`,`owning`,`leasing`,`renting`). |

### 1.2 `dto-asset.ts` (interfaces)

| Legacy capability | Disposition |
|---|---|
| `AssetLiabilityInfo` (`{id,serviceTypes,serviceProvider}`) | Ported (asset-side linkage) -> `dto/asset.ts` `IAssetLiabilityInfo` (`{id,serviceTypes}`). The `serviceProvider` sub-object is a service-provider concern -> liabilities sibling (see §4). |
| `ISubAssetInfo` | Ported -> `dto/asset.ts` `ISubAssetInfo` (`{id,title,type,countryID,subType,expires}`). |
| `IAssetBrief<ExtraType,Extra>` | Ported (reconciled) -> `dto/asset.ts` `IAssetBrief` (`{id,name,category,condition,status,visibility}`, backend canonical). The legacy generic `<ExtraType,Extra>` embedding is replaced by the polymorphic `extraType`/`extra` discriminator on the record (`dto/extras.ts`). |
| `IAssetusSpaceDbo` (`{assets}`) | Ported -> `dto/asset.ts` `IAssetusSpaceBrief` (`{assets:Record<id,IAssetBrief>}`), used by `IWithAssetSpaces.spaces`. |
| `AssetExtraType` | Ported -> `dto/extras.ts` `AssetExtraType` (`vehicle`/`dwelling`/`document`) + const exports. |
| `IAssetExtra` (index signature base) | Intentional change — dropped. Replaced by the concrete typed-extra interfaces resolved by `extraType`; the open `[key:string]:unknown` base is not needed. |
| `IAssetusSpaceContext` | Reconciled — superseded. The MVP reads assets via `AssetService.watchAssets()` (`services/asset.service.ts`) returning `IIdAndAssetDbo[]`, not a legacy `INavContext` space wrapper. |
| `IAssetDboBase<…>` | Ported (flattened) -> `dto/asset.ts` `IAssetDbo` carries all its fields (`parentAssetID`,`desc`->`description`,`memberIDs`,`parentCategoryID`,`sameAssetID`,`groupId`->`groupID`,`subAssets`,`membersInfo`,`liabilities`,`notUsedServiceTypes`). |
| `IWithAssetExtra<…>` (`extraType`/`extra`) | Ported -> the `extraType`/`extra` discriminator pattern (`dto/extras.ts` `AssetExtraType` + the typed extra interfaces). |
| `IAssetDbo<ExtraType,Extra>` (+`userIDs`) | Ported (reconciled) -> `dto/asset.ts` `IAssetDbo` (non-generic superset). Legacy `userIDs` intentionally dropped — not in the backend wire contract. |
| `IEngine` | Ported -> `dto/extras.ts` `IEngine` (+`engineSerialNumber`). |
| `IWithMakeAndModel` | Ported -> `dto/extras.ts` `IWithMakeAndModel`. |
| `IAssetVehicleExtra` | Ported -> `dto/extras.ts` `IAssetVehicleExtra` (`make`/`model`/`vin`/`regNumber`/engine + `nctExpires`/`taxExpires`/`nextServiceDue`). Legacy `*TaskId` task-link IDs degrade to plain due-dates (Open Question: tasks-module reminder integration -> sibling). |
| `IAssetDwellingExtra` | Ported -> `dto/extras.ts` `IAssetDwellingExtra` (`address`/`rent_price`/`numberOfBedrooms`/`areaSqM`). |
| `IAssetCategory` (titled category descriptor) | Reconciled — superseded. The MVP uses `dto/asset.ts` `assetCategories` + `categoryOptions` (`ILabeledOption`) for category metadata; the legacy per-category `iconName`/`order`/`canHaveIncome`/`canHaveExpense` descriptor is reproduced locally where needed (the ported `new-asset` page defines a local category shape — see COMPONENTS-COVERAGE.md). |
| `IAssetDtoGroupCounts` (`{assets}`) | Ported -> `dto/asset.ts` `IAssetGroupCounts` (`{assets}`). |
| `IAssetDtoGroup` | Ported -> `dto/asset.ts` `IAssetGroupInfo` (`{id,title,order,desc,categoryID,numberOf,totals}`) carried on `IAssetDbo.group`. |

### 1.3 `dto-document.ts`

| Legacy capability | Disposition |
|---|---|
| `AssetDocumentType` | Ported (reconciled) -> `dto/asset.ts` `AssetDocumentType` (`passport`/`id_card`/`driving_license`/`marriage_cert`/`birth_cert`, backend taxonomy; legacy `marriage_certificate`->`marriage_cert`, `birth_certificate`->`birth_cert`). Legacy `unspecified`/`rent_lease`/`insurance_policy` dropped (not in backend `Type` set). |
| `IDocTypeField` | Ported -> `dto/extras.ts` `IDocTypeField`. |
| `IDocTypeStandardFields` | Ported -> `dto/extras.ts` `IDocTypeStandardFields`. |
| `DocTypeDef` (`{id,title,emoji,fields}`) | Ported -> `dto/extras.ts` `IDocTypeDef` (`{id,fields}`). Presentational `title`/`emoji` intentionally dropped from the schema DTO (validation-only). |
| `standardDocTypesByID` | Ported -> `dto/extras.ts` `standardDocTypesByID` (+ `docTypeSchema()` lookup). |
| `IAssetDocumentExtra` | Ported -> `dto/extras.ts` `IAssetDocumentExtra` (`docType`/`number`/`batchNumber`/`countryID`/`issuedBy`/`issuedOn`/`effectiveFrom`/`expiresOn` + `regNumber` alias). |
| `IAssetDocumentContext` | Reconciled — superseded (legacy `IAssetContext` nav-context type, see §1.6). |

### 1.4 `types.ts`

| Legacy capability | Disposition |
|---|---|
| `CurrencyCode` (`USD`/`EUR`) + `CurrencyUSD`/`CurrencyEUR`/`CurrencyList` | Intentional change — replaced by the backend's open-string currency. Money is `IMoney{currency:string,value}` (`dto/asset.ts`); the backend has no typed `CurrencyCode` enum, so the 2-value legacy union is intentionally not ported. |

### 1.5 `uimodels/` — `asset-group.ts`

| Legacy capability | Disposition |
|---|---|
| `AssetGroup` uimodel class (`totals`/`id`/`numberOf` over `IAssetGroupContext`) | Reconciled — superseded. The group sub-entity is the plain `IAssetGroupInfo` DTO (`dto/asset.ts`) carried on `IAssetDbo.group` (`numberOf`/`totals` are fields); the ported `asset-group-page.component.ts` reads the group directly from nav state rather than wrapping it in a `Totals`-bearing uimodel class. |

### 1.6 `contexts/` — `asset-context.ts`

| Legacy capability | Disposition |
|---|---|
| `IAssetContext<ExtraType,Extra>` (`ISpaceItemNavContext` wrapper) | Reconciled — superseded by `IIdAndAssetDbo` (`services/asset.service.ts`) + the flat `IAssetDbo`. The legacy `INavContext`/`ISpaceItemNavContext` brief+dbo nav wrapper is replaced by direct DBO access in the MVP pages. |
| `IAssetVehicleContext` (`IAssetContext<'vehicle',…>`) | Reconciled — superseded (same as above; vehicle data read via `IAssetDbo` + `extra` of `extraType='vehicle'`). |
| `IAssetDwellingContext` (`IAssetContext<'dwelling',…>`) | Reconciled — superseded (`IAssetDbo` + `extraType='dwelling'`). |
| `IAssetGroupContext` (`INavContext<IAssetDtoGroup>`) | Reconciled — superseded; the ported `asset-group-page` reads `IAssetGroupInfo` from nav state. |

### 1.7 `data/` — vehicle reference data

| Legacy capability | Disposition |
|---|---|
| `data/vehicles.ts` (`Engine`, `engines`, `IMake`/`IModel`) | Consumed-from-published — not re-ported. The ported `make-model-card`/`vehicle-engine` components import `carMakes`/`IMake`/`IModel`/`engines` from the published `@sneat/mod-assetus-core` (the legacy `core` lib resolves from node_modules). Static reference data, not a unified-DTO concern. |
| `data/car-makes-with-models.ts` (`carMakes`, `IMake`/`IModel`) | Consumed-from-published — not re-ported (same as above; imported by `components/make-model-card/make-model-card.component.ts` from `@sneat/mod-assetus-core`). |

---

## 2. Components — `components/src/lib/`

Folded from `libs/ext-assetus/COMPONENTS-COVERAGE.md` (FE Task 5); every ported
path below was re-verified to exist under `components/`.

### 2.1 Ported -> `components/`

| Legacy component | Ported to |
|---|---|
| `period-segment` | `components/period-segment/period-segment.component.ts` |
| `vehicle-engine` | `components/vehicle-engine/vehicle-engine.component.ts` |
| `make-model-card` | `components/make-model-card/make-model-card.component.ts` |
| `asset-reg-number-input` | `components/asset-reg-number-input/asset-reg-number-input.component.ts` |
| `mileage-dialog` | `components/mileage-dialog/mileage-dialog.component.ts` |
| `edit-dwelling-card` | `components/edit-dwelling-card/edit-dwelling-card.component.ts` |
| `real-estate-location` | `components/real-estate-location/real-estate-location.component.ts` |
| `asset-dates` | `components/asset-dates/asset-dates.component.ts` |
| `asset-possesion-card` | `components/asset-possession-card/asset-possession-card.component.ts` |
| `asset-card` | `components/asset-card/asset-card.component.ts` |
| `assets-list` | `components/assets-list/assets-list.component.ts` |
| `vehicle-card` | `components/vehicle-card/vehicle-card.component.ts` |
| `asset-add/asset-add-vehicle` | `components/asset-add/asset-add-vehicle/asset-add-vehicle.component.ts` |
| `asset-add/asset-add-dwelling` | `components/asset-add/asset-add-dwelling/asset-add-dwelling.component.ts` |
| `asset-add/asset-add-document` | `components/asset-add/asset-add-document/asset-add-document.component.ts` |

### 2.2 Reconciled with existing MVP equivalent

| Legacy component | MVP equivalent |
|---|---|
| `asset-card` / `assets-list` (page-level navigation + edit) | `pages/asset/asset-page.component.ts`, `pages/assets/assets-page.component.ts` — page-level detail/edit/remove/transfer/history is the MVP's against the unified `IAssetDbo`. (The widgets are still ported above as presentational components.) |

### 2.3 Deferred -> liabilities / service-provider sibling Feature

| Legacy component | Reason |
|---|---|
| `asset-liabilities` | Manages asset liabilities + service-types (`MultiSelectorComponent`, `notUsedServiceTypes`, legacy liability service calls). Belongs to the liabilities sibling. |
| `asset-add-service` | "Add asset service" page bound to the service-provider flow (`AssetBasePage` + service plumbing). Service-provider/liabilities sibling. |
| `asset-contacts-group` | Bound to `@sneat/contactus-core` (`IContact2Asset`/`ContactType`) — asset<->service-provider contact grouping. Service-provider-bound; liabilities sibling. |

### 2.4 Shared bases / non-UI (accounted, not separately ported)

| Legacy file | Disposition |
|---|---|
| `asset-add/add-asset-base-component.ts` | Consumed from published `@sneat/ext-assetus-components` by the ported `asset-add-*` subclasses. Not re-ported. |
| `asset-base-page.ts` | Service-provider page base used by `asset-add-service`. Deferred with the liabilities/service sibling. |
| `asset-component-base-params.ts` | Injectable params wrapper around the legacy `AssetService`; consumed from the published lib by the deferred bases. Not re-ported. |
| `assets-list/asset-list-item.component.ts` | Dead/unused legacy code (not referenced by `assets-list.component`, not exported by the legacy index). Not ported; list rendering covered by the ported `assets-list`. |
| `testing/test-utils.ts` | Test harness (not a UI component). Not ported. |

---

## 3. Services — `components/src/lib/services/`

| Legacy capability | Disposition |
|---|---|
| `asset-service.ts` `AssetService` (`createAsset<Extra>`, `updateAsset`, `deleteAsset`, `addVehicleRecord`) | Ported (reconciled into one service) -> `services/asset.service.ts` MVP `AssetService`: `createAsset`/`updateAsset`/`removeAsset`/`addVehicleRecord` plus the MVP ownership ops `transferAsset`/`recordHistoryEvent`/`getHistory`/`watchAssets`. `addVehicleRecord` targets `create_vehicle_record` (not legacy `add_vehicle_record`/`assets/…`). |
| `asset-service.dto.ts` `ICreateAssetRequest`/`IUpdateAssetRequest`/`IAddVehicleRecordRequest` | Ported (reconciled to backend request surface) -> `services/interfaces.ts` `ICreateAssetRequest`/`IUpdateAssetRequest`/`ICreateVehicleRecordRequest` (extends `IAddVehicleRecordRequest` from `dto/extras.ts`, fuel-bearing payload). Mirrors the current flat backend request surface. |
| `assetus-space.service.ts` `AssetusSpaceService` | Reconciled — superseded. Per-space asset reads are `AssetService.watchAssets(spaceID)` (`services/asset.service.ts`); a separate space service is not needed in the MVP. |
| `assetus-services.module.ts` (NgModule) | Reconciled — superseded by `services/assetus-core-services.module.ts` (MVP services module). |

---

## 4. Liabilities / service-provider DTOs — `dto/dto-liability.ts` + `dto/dto-service-provider.ts`

**Disposition split (AC:liabilities-disposition-recorded).** The **asset-side
liability linkage** is folded into the unified core DTO here: `IAssetDbo.liabilities`
(`IAssetLiabilityInfo[]`) and `IAssetDbo.notUsedServiceTypes` in `dto/asset.ts` (§1.2).
The **standalone provider / plan / settlement / service-type** records and their
contact model are NOT asset-side linkage — each is assigned individually to the
**liabilities sibling Feature** below. None lumped, none dropped.

### 4.1 Standalone exports assigned to the liabilities sibling Feature

| Legacy export (source) | Disposition: liabilities sibling Feature — reason |
|---|---|
| `DtoServiceProvider` (`dto-service-provider.ts`) | Liabilities sibling — a standalone service-provider record (`countryId`/`status`/`assetCategoryId`/`serviceTypes`/`contact`); not an asset field. Its dedicated provider UI (`asset-add-service`/`select-service-provider`) is deferred with it. |
| `DtoServicePlan` (`dto-service-provider.ts`) | Liabilities sibling — a tariff/plan record (`eab`/`pricePerUnit`/`offers`/`unit`/`settlementType`) under a provider; no asset-side home. |
| `DtoServiceType` (`dto-service-provider.ts`) | Liabilities sibling — a catalogued service-type definition (`serviceCategoryId` + `LinkedToAssetCategories`); a provider/catalog concern, distinct from the per-asset `serviceTypes` string list folded into `IAssetLiabilityInfo`. |
| `LiabilityServiceType` (`dto-liability.ts`) | Liabilities sibling — the typed service-type union (`DwellingServiceType`/`VehicleServiceType`/etc.). The asset-side `serviceTypes` is folded into core as `string[]` (`IAssetLiabilityInfo.serviceTypes` / `notUsedServiceTypes`), but the full typed taxonomy belongs to the liabilities/service-provider catalog. |
| `SettlementType` (`dto-liability.ts`, `rural`/`urban`) | Liabilities sibling — a service-plan attribute (consumed only by `DtoServicePlan.settlementType`); no asset-side use. |
| `ServiceProviderContacts` (`dto-service-provider.ts`) | Liabilities sibling — the provider contact model (`usUrl` + `groups` of `ServiceProviderContactsGroup`/`ServiceProviderContact`), surfaced by the deferred `asset-contacts-group` component. Provider-bound, not an asset field. |

### 4.2 Other liabilities/service-provider DTOs (same sibling, recorded for completeness)

| Legacy export (source) | Disposition |
|---|---|
| `DtoLiability` (`dto-liability.ts`) | Liabilities sibling — the standalone liability (debt) record (`type`/`direction`/`amount`/`period`/`serviceProvider`); created by the deferred `liability/liability-new` page. Asset-side linkage is the `IAssetLiabilityInfo` fold only. |
| `LiabilityType` / `LiabilityDirection` / `ServiceCategory` (`dto-liability.ts`) | Liabilities sibling — supporting unions of `DtoLiability` / `DtoServiceType`. |
| `DwellingServiceType` / `VehicleServiceType` / `EntertainmentServiceType` / `EducationServiceType` / `DwellingServiceType` / `DwellingTaxServiceType` (`dto-liability.ts`) | Liabilities sibling — the per-domain members of `LiabilityServiceType`. |
| `ServiceProviderStatus` (`dto-service-provider.ts`) | Liabilities sibling — the `DtoServiceProvider.status` union. |
| `ServiceProviderContact` / `ServiceProviderContactsGroup` (`dto-service-provider.ts`) | Liabilities sibling — sub-shapes of `ServiceProviderContacts`. |
| `ServicePlanOffer` (`dto-service-provider.ts`) | Liabilities sibling — a `DtoServicePlan.offers[]` sub-shape. |
| `LinkedToAssetCategories` (`dto-service-provider.ts`) | Liabilities sibling — the `DtoServiceType` base (`assetCategoryIds`). |

---

## 5. Pages — `sneat-apps/.../pages/src/lib/`

Folded from `libs/ext-assetus/COMPONENTS-COVERAGE.md`; ported paths re-verified.

### 5.1 Ported -> `pages/` (wired in `assetus-routing.ts`)

| Legacy page | Ported to | Route |
|---|---|---|
| `new-asset` | `pages/new-asset/new-asset-page.component.ts` | `new-asset` |
| `real-estates` | `pages/real-estates/real-estates-page.component.ts` | `real-estates` |
| `real-estate` | `pages/real-estate/real-estate-page.component.ts` | `real-estate/:assetID` |
| `asset-group` | `pages/asset-group/asset-group-page.component.ts` | `asset-group` |
| `optimization` | `pages/optimization/optimization-page.component.ts` | `optimization` |

### 5.2 Reconciled with existing MVP page

| Legacy page | MVP page |
|---|---|
| `assets` (+`AssetsBasePage`) | `pages/assets/assets-page.component.ts` (live `watchAssets`, New-asset dialog, Active/Archived filter against unified `IAssetDbo`). |
| `asset` (+`AssetBasePage`) | `pages/asset/asset-page.component.ts` (detail/edit/archive/delete/transfer/history against unified `IAssetDbo`). |

### 5.3 Shared base pages (accounted, not separately ported)

| Legacy file | Disposition |
|---|---|
| `assets-base.page.ts` | Folded into the MVP `assets-page`; the `goNew(category)` navigation reproduced inline by `real-estates-page`. Not re-ported as a base. |
| `asset-base.page.ts` | Folded into the MVP `asset-page` and the ported `real-estate-page`; also the base for the deferred service-provider pages. Not re-ported as a base. |

### 5.4 Deferred -> liabilities / service-provider sibling Feature

| Legacy page | Reason |
|---|---|
| `liability/liability-new` | Creates a liability (debt) record; liabilities sibling. |
| `liability/select-service-provider` | Service-provider picker on the legacy `AssetBasePage` + service-provider flow; liabilities sibling. |
| `real-estate` landlord/tenant contacts | Rendered `asset-contacts-group` (deferred component, §2.3); the contact groups are omitted from the ported `real-estate-page`. Liabilities/service sibling. |

---

## 6. Build & test (AC:frontend-builds-and-tests-pass)

From `assetus/frontend`:

- `npx nx build ext-assetus` -> **success** (`Built @sneat/extension-assetus`).
- `npx nx test ext-assetus --watch=false` -> **success** (6 test files, **49 tests passed**).

Both pass with no fixes required. (The Nx Cloud "complete setup" / 401 FREE-plan
notice is unrelated to the build/test outcome and is ignored per task scope.)
