---
format: https://specscore.md/feature-specification
status: Approved
---

# Feature: Assetus Frontend Port

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/assetus/spec/features/assetus-frontend-port?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/assetus/spec/features/assetus-frontend-port?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/assetus/spec/features/assetus-frontend-port?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/assetus/spec/features/assetus-frontend-port?op=request-change) |
**Status:** Approved
**Source Ideas:** assetus-model-unification
**Grade:** B

## Summary

Port the legacy Assetus Angular frontend — the `core` DTOs, the ~19 `components`, and the pages — into `assetus/frontend` (the new MVP-shaped Angular workspace), adapting them onto the existing MVP frontend so the migrated DTOs reach **parity with the unified backend model** (`unified-assetus-data-model`) and **no legacy frontend functionality is lost**. This is the frontend sibling of the backend `unified-assetus-data-model` Feature, recorded there as the deferred frontend-port and the equivalent ≥80% coverage bar.

## Problem

The unified backend Feature (`unified-assetus-data-model`) reconciles the legacy rich asset-registry and the MVP ownership-core into a single superset Go model. The frontend has the same split: the legacy Angular frontend (`sneat-libs/libs/extensions/assetus/{core,components}` + `sneat-apps/libs/extensions/assetus/pages`) carries the rich registry DTOs, ~19 components and the asset/real-estate/group/liability pages, while the new `assetus/frontend/libs/ext-assetus` carries only the flat MVP DTO (`IAssetBrief`/`IAssetDbo`), a history-timeline component, and a handful of MVP pages.

Until the legacy Angular code is ported into `assetus/frontend` and its DTOs mirror the unified backend contract, the legacy frontend directories cannot be retired without losing functionality, and the frontend cannot consume the unified backend. The owner's hard condition is **no functionality lost**: every legacy DTO, component and page must be ported or explicitly deferred with a reason, and the migrated code must carry the same ≥80% coverage bar the backend Feature records.

This is a **migration/port, not a rewrite**: tasks relocate and adapt the existing legacy Angular/TS files into `assetus/frontend`, naming the legacy source and the new merge target — not authoring from scratch.

## Behavior

### Backend Contract

#### REQ: backend-wire-contract

The backend is now implemented, so the frontend DTOs MUST mirror the **exact json field names, enum string values and HTTP route paths** below (extracted from `assetus/backend`), not a vague "parity". This is the canonical contract the rest of this Feature's REQs reference.

**Asset record — `AssetBase`/`AssetDbo`/`AssetBrief` (`backend/dbo4assetus/asset_dbo.go`).** json field names:
`name`, `description`, `category`, `condition`, `status`, `visibility`, `acquisitionDate`, `purchasePrice`, `estimatedValue` (each `{currency,value}` — `MonetaryAmount`), `location`, `notes`, `tags`, `photos`; optional ported fields `countryID`, `type`, `possession`, `parentCategoryID`, `yearOfBuild` (int), `isRequest` (bool), `geo` (`{lat,lng}` — `GeoPoint`); embedded `AssetDates` → `dateOfBuild`, `dateOfPurchase`, `dateInsuredTill`, `dateCertifiedTill` (all ISO `YYYY-MM-DD` strings); custom fields via `WithCustomFields` (`fieldsStr`/`fieldsInt`/`fieldsDate`/`fieldsAmount`); financial fields `totals` (`money.Amount[]`), `canHaveIncome` (bool), `canHaveExpense` (bool), `financialDirection` (`"income"`|`"expense"`), `liabilities` (`AssetLiabilityInfo[]` — `{id,serviceTypes}`), `notUsedServiceTypes` (string[]). `AssetBrief` is exactly `{id,name,category,condition,status,visibility}`. NOTE: there is **no** `title`/`desc`/`name`-vs-`title` split on the backend record — `name` is canonical; there is no top-level `extraType`/`extra` field on `AssetBase` (the typed extra is carried via the shared `core/extra.WithExtraField`, json `extraType`/`extra`).

**Relationships — `WithAssetRelationships` (`backend/dbo4assetus/relationships.go`).** json: `groupID` (string), `group` (`AssetGroupInfo`), `parentAssetID`, `subAssets` (`SubAssetInfo[]`), `sameAssetID`, `relatedAs` (via `WithOptionalRelatedAs`), `memberIDs` (string[]), `membersInfo` (`TitledRecord[]` — `{id,title}`). `AssetGroupInfo` = `{id,title,order,desc,categoryID,numberOf,totals}` where `numberOf` is `AssetGroupCounts{assets:int}`. `SubAssetInfo` = `{id,title,type,countryID,subType,expires}`.

**Vehicle record — `VehicleRecordDbo` (`backend/dbo4assetus/vehicle_record_dbo.go`).** Persisted shape: `{fuel,mileage}` where `fuel` = `VehicleFuelRecord{volume,unit,amount,fuelCost,currency}` and `mileage` = `VehicleMileage{value,unit}`. The **request** to append one is `AddVehicleRecordRequest` (`backend/dto4assetus/add_vehicle_record.go`) with the flat payload `{spaceID,assetID,fuelVolume,fuelVolumeUnit,fuelCost,currency,mileage,mileageUnit}`.

**Typed extras (`backend/extras4assetus/`), discriminated by `extraType` ∈ {`vehicle`,`dwelling`,`document`}.**
- Vehicle (`asset_vehicle.go` + `with_engine_data.go` + `with_make_model.go`): `make`, `model`, `regNumber`, `vin`, engine `engineType`/`engineFuel`/`engineCC`/`engineKW`/`engineNM`/`engineSerialNumber`, plus plain due-dates `nctExpires`/`taxExpires`/`nextServiceDue`. `IndexedFields`: `make`,`model`,`make+model`,`regNumber`,`vin`.
- Dwelling (`asset_dwelling.go`): `address` (`dbmodels.Address`), `rent_price` (`{value,currency}`), `numberOfBedrooms` (int), `areaSqM` (int).
- Document (`asset_document.go` + `doc_type_schema.go`): `regNumber` (number alias), `docType` (a `const4assetus.Type` document subtype), `number`, `batchNumber`, `countryID`, `issuedBy`, `issuedOn`, `effectiveFrom`, `expiresOn`. `IndexedFields`: `expiresOn`,`effectiveFrom`. Per-doc-type schema `standardDocTypesByID` keyed by docType: e.g. `passport` & `driving_license` require `number`+`validTill`(=`expiresOn`) and allow `Members` max 1; `birth_cert` & `marriage_cert` require `number`+`issuedOn`, **exclude** `validTill`, allow `Members` max 1 (birth) / 2 (marriage); `other` requires a title.

**Enum string values (`backend/const4assetus/`).**
- `Category` (`category.go`): `books`,`games`,`toys`,`sports_equipment`,`tools`,`electronics`,`clothing`,`vehicles`,`camping_equipment`,`other`,`dwelling`,`document`,`debt`. (Legacy reconciliation already applied in the backend set: `sport_gear`→`sports_equipment`, `vehicle`→`vehicles`, `misc`→`other`.)
- `Status` (`status.go`): `draft`,`active`,`transferred`,`archived`,`disposed`,`lost`.
- `Possession` (`possession.go`): `unknown`,`undisclosed`,`owning`,`leasing`,`renting`.
- `Type` (`asset_type.go`, per-category subtype): vehicles → `aircraft`,`boat`,`bus`,`car`,`helicopter`,`motorcycle`,`truck`,`van`; dwelling → `apartment`,`house`,`office`,`shop`,`land`,`garage`,`warehouse`; sports_equipment → `bicycle`,`kite`,`kite_bar`,`kite_board`,`kite_hydrofoil`,`prone_hydrofoil`,`surf_board`,`wetsuit`,`wing`,`wing_board`,`wing_hydrofoil`; document → `passport`,`id_card`,`driving_license`,`marriage_cert`,`birth_cert`. (The document `Type` subtype set IS the `AssetDocumentType` taxonomy.)
- `EngineType` (`engine_fuel.go`): `""`(unknown),`other`,`combustion`,`electric`,`phev`,`hybrid`,`steam`.
- `FuelType` (`engine_fuel.go`): `""`(unknown),`other`,`bio`,`petrol`,`diesel`,`hydrogen`.
- `Condition`: `new`,`excellent`,`good`,`fair`,`needs_repair`,`broken`. `Visibility`: `private`,`family`,`friends`,`friends_of_friends`,`specific_space`,`public`.
- There is **no typed `FuelVolumeUnit`/`MileageUnit`/`CurrencyCode` enum on the backend** — fuel volume unit, mileage unit and currency are carried as open `string` on the vehicle-record request/record. The frontend MAY keep a named union for ergonomics but MUST serialize plain strings.

**HTTP routes (`backend/api4assetus/routes.go`), all under the `assetus/` API prefix.**
`POST /v0/assetus/create_asset`, `GET /v0/assetus/asset`, `POST /v0/assetus/update_asset`, `POST /v0/assetus/remove_asset`, `POST /v0/assetus/transfer_asset`, `POST /v0/assetus/record_history_event`, `GET /v0/assetus/asset_history`, `POST /v0/assetus/create_vehicle_record`. (The vehicle-record route is `create_vehicle_record`, not `add_vehicle_record`.)

**Known divergence requiring a human decision (DTO vs API request surface).** The persisted `AssetDbo` carries the full rich field set above, but the implemented `CreateAssetRequest`/`UpdateAssetRequest` (`backend/dto4assetus/{create,update}_asset.go`) currently accept ONLY the flat MVP set (`name`,`description`,`category`,`condition`,`visibility`,`acquisitionDate`,`purchasePrice`,`estimatedValue`,`location`,`notes`,`tags`) — and create does not accept `status` at all. The rich optional fields (`type`,`possession`,`countryID`,`AssetDates`,relationships,financials,typed extras) are therefore **not yet settable through the create/update HTTP API**, even though the record can hold them (they would be written by other paths, e.g. Firestore writes / future request fields). The frontend create/update **request DTOs** MUST mirror the current backend request surface; the rich fields belong on the **record DTO** (`IAssetDbo`) read back from Firestore/`asset_history`. Whether to widen the backend create/update requests to accept the rich fields is owned by `unified-assetus-data-model` and is flagged in Open Questions.

### DTO parity

#### REQ: dto-core-superset

The frontend `IAssetBrief`/`IAssetDbo` in `assetus/frontend/libs/ext-assetus/src/lib/dto/asset.ts` MUST become a **superset** whose json field names mirror the backend `AssetBrief`/`AssetBase`/`AssetDbo` exactly (see REQ:backend-wire-contract). `IAssetBrief` is `{id,name,category,condition,status,visibility}`. The MVP `IAssetDbo` fields already present (`description`, `acquisitionDate`, `purchasePrice`, `estimatedValue`, `location`, `notes`, `tags`, `createdAt`/`updatedAt`, `spaceID`s) MUST be retained, and the still-wanted legacy fields MUST be added as **optional**, using the backend json names: `isRequest`, `countryID`, `type` (per-category subtype), `yearOfBuild` (int) plus the `AssetDates` group `dateOfBuild`/`dateOfPurchase`/`dateInsuredTill`/`dateCertifiedTill`, `possession`, `parentAssetID`, `parentCategoryID`, `geo` (`{lat,lng}`), `photos`, `memberIDs`, `membersInfo` (`{id,title}[]`), `groupID` (string) + `group` (the group sub-entity), `subAssets`, `sameAssetID`, `relatedAs`, custom fields (`fieldsStr`/`fieldsInt`/`fieldsDate`/`fieldsAmount`), the typed extra via `extraType`/`extra`, and the financial dimension as core fields (`totals`, `canHaveIncome`/`canHaveExpense`, `financialDirection`, `liabilities`/`notUsedServiceTypes`). NOTE the corrections vs the earlier draft: it is `groupID` (not `groupId`) plus a `group` object, `categoryID`/`countryID` (not `…Id`); there is **no** separate `title`/`desc` on the record (`name` is canonical). An asset created with none of the optional legacy fields MUST stay valid.

#### REQ: enum-union-frontend

The frontend enum types in `asset.ts` MUST exactly match the backend `const4assetus` string sets (see REQ:backend-wire-contract for the canonical lists): `AssetCategory` = `books`,`games`,`toys`,`sports_equipment`,`tools`,`electronics`,`clothing`,`vehicles`,`camping_equipment`,`other`,`dwelling`,`document`,`debt` (the legacy reconciliation `sport_gear`→`sports_equipment`, `vehicle`→`vehicles`, `misc`→`other` is already baked into this set; `debt`/`dwelling`/`document` retained); `AssetStatus` = `draft`,`active`,`transferred`,`archived`,`disposed`,`lost`; `AssetPossession` = `unknown`,`undisclosed`,`owning`,`leasing`,`renting`; `AssetType` = the per-category subtype set in the backend `Type`; `EngineType` = `other`,`combustion`,`electric`,`phev`,`hybrid`,`steam` (plus empty = unknown); `FuelType` = `other`,`bio`,`petrol`,`diesel`,`hydrogen` (plus empty = unknown); `AssetCondition` and `AssetVisibility` MUST be retained unchanged. Fuel-volume unit, mileage unit and currency are **open `string`** on the backend (NOT typed enums) — the frontend MAY keep a named union for ergonomics but MUST serialize plain strings. No backend enum value may be missing from the frontend set.

#### REQ: typed-extras-frontend

The legacy polymorphic typed-extra DTOs MUST be ported onto the frontend core as an optional `extra` resolved by `extraType` ∈ {`vehicle`,`dwelling`,`document`}, with json field names matching REQ:backend-wire-contract and no field dropped:
- `IAssetVehicleExtra` (backend `extras4assetus.AssetVehicleExtra`) — `make`, `model`, `regNumber`, `vin`, engine data `engineType`/`engineFuel`/`engineCC`/`engineKW`/`engineNM` **and `engineSerialNumber`**, plus the plain service/tax/inspection due-dates `nctExpires`/`taxExpires`/`nextServiceDue` (ISO `YYYY-MM-DD` strings; the legacy task-link IDs degrade to plain due-dates — the backend carries only the dates).
- The vehicle-record append DTO — backend persists `VehicleRecordDbo{fuel,mileage}` (`fuel`=`{volume,unit,amount,fuelCost,currency}`, `mileage`=`{value,unit}`), and the **request** `AddVehicleRecordRequest` is the flat payload `{spaceID,assetID,fuelVolume,fuelVolumeUnit,fuelCost,currency,mileage,mileageUnit}`. Both shapes MUST be ported with no field dropped, and `fuelVolumeUnit`/`mileageUnit`/`currency` serialized as plain strings.
- `IAssetDwellingExtra` (backend `AssetDwellingExtra`) — `address`, `rent_price` (`{value,currency}`), `numberOfBedrooms`, `areaSqM`.
- `IAssetDocumentExtra` (backend `AssetDocumentExtra`) — `docType` (a document `AssetType` subtype: `passport`/`id_card`/`driving_license`/`marriage_cert`/`birth_cert`), `number`, `batchNumber`, `countryID`, `issuedBy`, `issuedOn`, `effectiveFrom`, `expiresOn` (plus the `regNumber` number-alias), and the per-doc-type validation schema (`standardDocTypesByID`/`DocTypeStandardFields`, e.g. `passport`/`driving_license` require `number`+`validTill`=`expiresOn`; `marriage_cert` requires `number`+`issuedOn`, excludes `validTill`, allows Members max 2) MUST be ported as the document-extra validation rules.

#### REQ: relationships-frontend

The legacy relationship DTOs MUST be ported with json names matching the backend `WithAssetRelationships` (REQ:backend-wire-contract) and no field dropped: the asset **group** linkage `groupID` (string) **plus** the group sub-entity `group` = `AssetGroupInfo` `{id,title,order,desc,categoryID,numberOf,totals}` (where `numberOf`=`{assets}`), and its `AssetGroup` uimodel (`sneat-libs/.../core/src/lib/uimodels/asset-group.ts`); parent/sub-asset nesting `parentAssetID`+`subAssets` preserving the per-sub-asset `SubAssetInfo` detail `{id,title,type,countryID,subType,expires}`; asset-to-asset linking (`sameAssetID`/`relatedAs`); multi-space association alongside the canonical owning Space; and per-asset member info `memberIDs` (string[]) + `membersInfo` (`{id,title}[]`). NOTE the corrections vs the earlier draft: `groupID`/`categoryID`/`countryID` (capital ID), and `group` is an object, not merely a `groupId` string.

### Components

#### REQ: components-ported

Every legacy component under `sneat-libs/libs/extensions/assetus/components/src/lib/` MUST be **ported or explicitly deferred with a reason** into `assetus/frontend/libs/ext-assetus/src/lib/components/` (merging alongside the existing `asset-history-timeline`), with no legacy component silently dropped. The set to account for: `vehicle-card`, `vehicle-engine`, `make-model-card`, `asset-reg-number-input`, `mileage-dialog`, `edit-dwelling-card`, `real-estate-location`, `asset-dates`, `period-segment`, `asset-possesion-card` (`asset-possession-card`), `asset-liabilities`, `asset-contacts-group`, `asset-add-vehicle`, `asset-add-dwelling`, `asset-add-document`, `asset-add-service`, `assets-list` (+ `asset-list-item`), `asset-card`, and the shared bases (`add-asset-base-component`, `asset-base-page`, `asset-component-base-params`). Ported components MUST consume the unified frontend DTOs and the new component exports MUST be re-exported from the library index.

#### REQ: services-ported

The legacy Assetus frontend services MUST be ported/merged into `assetus/frontend/libs/ext-assetus/src/lib/services/`: the asset service and its request DTOs (`asset-service.dto.ts` incl. `ICreateAssetRequest`/`IUpdateAssetRequest`/`IAddVehicleRecordRequest` with the fuel-bearing payload), the assetus space service, and the services module — adapted onto the existing MVP `AssetService`/core-services module so a single asset service exposes both the MVP ownership operations (transfer, history) and the ported create/update/add-vehicle-record operations. The service MUST call the **actual implemented backend routes** (REQ:backend-wire-contract): `create_asset`, `asset` (GET), `update_asset`, `remove_asset`, `transfer_asset`, `record_history_event`, `asset_history` (GET) and `create_vehicle_record` (the vehicle-record route is `create_vehicle_record`, NOT `add_vehicle_record`), under the shared `assetus/` API prefix. The MVP `AssetService` already wires the first seven; the add-vehicle-record method MUST be added pointing at `create_vehicle_record`. The create/update **request** DTOs MUST mirror the current backend request surface (the flat MVP set — see the divergence note in REQ:backend-wire-contract), even though the record DTO read back is the full superset.

### Pages

#### REQ: pages-ported

Every legacy page under `sneat-apps/libs/extensions/assetus/pages/src/lib/` MUST be **ported or explicitly deferred with a reason** into `assetus/frontend/libs/ext-assetus/src/lib/pages/` and wired into the routing module, with no legacy page silently dropped. The set to account for: `assets`, `asset`, `new-asset`, `real-estates` (+ `real-estate`), `asset-group`, `optimization`, the `liability/*` pages (`liability-new`, `select-service-provider`), and the shared `asset-base.page`/`assets-base.page`. The existing MVP pages (`assets`, `asset`, `transfer`, `new-asset-dialog`) MUST be reconciled with the ported pages rather than duplicated.

### Liabilities scope note

#### REQ: liabilities-sibling-scope

The standalone **liabilities / service-provider** provider/plan UI (`sneat-libs/.../core/src/lib/dto/{dto-liability,dto-service-provider}.ts`, the `asset-add-service`/`select-service-provider`/`liability-new` flows) MUST be **accounted for, not dropped**: the asset-side liability linkage fields (`liabilities`/`notUsedServiceTypes`/`AssetLiabilityInfo`) are folded into the unified core DTO here (REQ:dto-core-superset), while the standalone provider/plan/settlement records and their dedicated UI MAY be scoped to the **liabilities sibling Feature**. The capability-coverage table MUST record this disposition for every liabilities DTO/component/page so nothing is lost.

### Retirement readiness

#### REQ: frontend-capability-coverage

Every legacy frontend capability — each `core` DTO, each of the ~19 components, each page, and each service — MUST be representable in `assetus/frontend`, documented in a frontend capability-coverage table that maps each legacy file/capability to its ported new location, an intentional change, or an explicitly-named sibling Feature (liabilities). No legacy frontend capability may be left with no home. This is the testable encoding of "no functionality lost".

#### REQ: build-green-frontend

The migrated frontend MUST lint, build and test green from `assetus/frontend`: `nx run-many -t lint build test` (or the workspace's `ext-assetus` lint/build/test targets) exits zero.

### Testing

#### REQ: test-coverage-frontend

The ported frontend MUST carry the **equivalent coverage bar** recorded in the backend Feature: **≥80% statement coverage** plus **a test per ported DTO and per ported component**, runnable via the frontend's configured test runner (the workspace uses **vitest** via `@analogjs/vitest-angular`, NOT jest; coverage runs via `npx nx test ext-assetus --coverage --watch=false`, reporting v8 coverage). Each ported DTO MUST have a round-trip/shape test (incl. the backend enum string values, typed extras with the vehicle fuel-record, the document per-doc-type validation schema, relationships) and each ported component MUST have at least a render/smoke test.

## Acceptance Criteria

### AC: dto-superset-roundtrip (verifies REQ:dto-core-superset)

**Given** the ported `IAssetDbo` in `assetus/frontend/libs/ext-assetus/src/lib/dto/asset.ts`,
**When** an object is constructed carrying the MVP fields plus the legacy optionals using the backend json names (`isRequest`, `countryID`, `type`, `possession`, `AssetDates` (`dateOfBuild`/`dateOfPurchase`/`dateInsuredTill`/`dateCertifiedTill`), `yearOfBuild`, `parentAssetID`, `parentCategoryID`, `groupID`+`group`, `subAssets`, `sameAssetID`, `memberIDs`/`membersInfo`, `totals`, `canHaveIncome`/`canHaveExpense`/`financialDirection`, `liabilities`/`notUsedServiceTypes`),
**Then** it type-checks and round-trips with every field preserved, and a second object carrying only the required MVP fields (no legacy optionals) is also valid — proving each legacy field is optional.

### AC: enum-union-no-value-dropped (verifies REQ:enum-union-frontend)

**Given** the backend `const4assetus` enum string sets (`Category`, `Status`, `Possession`, `Type`, `EngineType`, `FuelType`) and the MVP enum sets,
**When** the unified frontend enum types are defined in `asset.ts`,
**Then** `AssetCategory` admits `document`/`debt`/`dwelling` and the MVP values (e.g. `books`), `AssetStatus` admits `draft` and `disposed`/`lost`/`transferred`, `AssetPossession` admits `owning`/`leasing`/`renting`/`unknown`/`undisclosed`, `EngineType` admits `phev`/`hybrid`/`steam`, `FuelType` admits `bio`/`hydrogen`, no backend enum value is absent (with `sport_gear`→`sports_equipment`, `vehicle`→`vehicles`, `misc`→`other` already reconciled in the backend `Category` set), and fuel-volume unit / mileage unit / currency are typed as open strings (no backend enum).

### AC: vehicle-extra-frontend-no-field-dropped (verifies REQ:typed-extras-frontend)

**Given** a ported `IAssetVehicleExtra` and an `IAddVehicleRecordRequest`,
**When** a vehicle asset is expressed with `make`, `model`, `regNumber`, `vin`, engine (`engineType`/`engineFuel`/`engineCC`/`engineKW`/`engineNM`/`engineSerialNumber`), the due-dates `nctExpires`/`taxExpires`/`nextServiceDue`, and a vehicle record carrying mileage AND fuel data — both the persisted shape (`fuel`=`{volume,unit,amount,fuelCost,currency}`, `mileage`=`{value,unit}`) and the flat request payload (`fuelVolume`/`fuelVolumeUnit`/`fuelCost`/`currency`/`mileage`/`mileageUnit`),
**Then** every one of those attributes — including `engineSerialNumber` and the NCT/tax/service due-dates — is present on the ported DTOs with no field dropped.

### AC: document-extra-frontend-full-shape (verifies REQ:typed-extras-frontend)

**Given** the ported `IAssetDocumentExtra` and `standardDocTypesByID`,
**When** a `docType=passport` document is expressed with `number`, `batchNumber`, `countryID`, `issuedBy`, `issuedOn`, `expiresOn` (and `effectiveFrom`),
**Then** all of those fields are present on the ported document extra with the backend json names and the per-doc-type validation schema for `passport` (`number` + `validTill`=`expiresOn` required) is applied — none of these fields is dropped.

### AC: relationships-frontend-preserved (verifies REQ:relationships-frontend)

**Given** the ported relationship DTOs (`group`=`AssetGroupInfo`/`AssetGroup`, `subAssets`=`SubAssetInfo[]`, the linking and multi-space/member fields),
**When** an asset that belongs to a group (`groupID`+`group`), has two sub-assets under a `parentAssetID`, is linked via `relatedAs`/`sameAssetID`, is associated with multiple spaces, and carries member info is expressed,
**Then** the group (with `id`/`title`/`order`/`desc`/`categoryID`/`numberOf`=`{assets}`/`totals`), the per-sub-asset `SubAssetInfo` detail (`id`/`title`/`type`/`countryID`/`subType`/`expires`), the links, the multi-space association, and `memberIDs`/`membersInfo` (`{id,title}[]`) are all preserved and resolvable with the backend json names.

### AC: all-components-accounted (verifies REQ:components-ported)

**Given** the legacy component inventory under `sneat-libs/.../components/src/lib/`,
**When** the frontend capability-coverage table is checked,
**Then** each of the ~19 components (`vehicle-card`, `vehicle-engine`, `make-model-card`, `asset-reg-number-input`, `mileage-dialog`, `edit-dwelling-card`, `real-estate-location`, `asset-dates`, `period-segment`, `asset-possession-card`, `asset-liabilities`, `asset-contacts-group`, `asset-add-{vehicle,dwelling,document,service}`, `assets-list`+`asset-list-item`, `asset-card`) and the shared bases has a row marking it ported (to a concrete new path) or explicitly deferred with a reason, with none silently dropped, and every ported component is re-exported from the library index.

### AC: services-frontend-ported (verifies REQ:services-ported)

**Given** the legacy Assetus frontend services and request DTOs,
**When** the ported `assetus/frontend/libs/ext-assetus/src/lib/services/` is checked,
**Then** a single asset service exposes both the MVP ownership operations (transfer via `transfer_asset`, history via `record_history_event`/`asset_history`) and the ported create/update/add-vehicle-record operations, the add-vehicle-record call targets the `create_vehicle_record` route (not `add_vehicle_record`), and `IAddVehicleRecordRequest` carries the fuel-bearing payload (`fuelVolume`/`fuelVolumeUnit`/`fuelCost`/`currency`/`mileage`/`mileageUnit`).

### AC: all-pages-accounted (verifies REQ:pages-ported)

**Given** the legacy page inventory under `sneat-apps/.../pages/src/lib/`,
**When** the frontend capability-coverage table and the routing module are checked,
**Then** each legacy page (`assets`, `asset`, `new-asset`, `real-estates`+`real-estate`, `asset-group`, `optimization`, `liability/*`, and the shared base pages) has a row marking it ported (and wired into routing) or explicitly deferred with a reason, the MVP pages are reconciled rather than duplicated, and none is silently dropped.

### AC: liabilities-disposition-recorded (verifies REQ:liabilities-sibling-scope)

**Given** the legacy liabilities/service-provider DTOs, components and pages,
**When** the frontend capability-coverage table is checked,
**Then** the asset-side liability linkage fields are recorded as folded into the unified core DTO, and every standalone provider/plan/settlement DTO/component/page has a recorded disposition (ported here or assigned to the named liabilities sibling Feature) — none is left unassigned.

### AC: frontend-capability-coverage-complete (verifies REQ:frontend-capability-coverage)

**Given** the full legacy frontend inventory (every `core` DTO, the ~19 components, every page, every service),
**When** the frontend capability-coverage table is checked against it,
**Then** every legacy capability has a row mapping it to a concrete ported location, an intentional change, or the named liabilities sibling Feature, and no capability is left unmapped.

### AC: frontend-builds-and-tests-pass (verifies REQ:build-green-frontend)

**Given** the migrated `assetus/frontend` workspace,
**When** `nx run-many -t lint build test` (or the `ext-assetus` lint/build/test targets) is run from `assetus/frontend`,
**Then** lint, build and test all exit zero.

### AC: frontend-test-coverage (verifies REQ:test-coverage-frontend)

**Given** the migrated `ext-assetus` library,
**When** `npx nx test ext-assetus --coverage --watch=false` is run via the workspace's vitest runner (`@analogjs/vitest-angular`),
**Then** the tests pass, overall statement coverage is at least 80%, and a test is present for each ported DTO (incl. enum-union values, the vehicle fuel-record extra, the document per-doc-type validation, and the relationship DTOs) and for each ported component (at least a render/smoke test).

## Architecture

- **Spine:** the MVP frontend `assetus/frontend/libs/ext-assetus` (flat `IAssetDbo`, `asset-history-timeline`, MVP pages, core `AssetService`) remains canonical and is extended, not replaced.
- **DTOs (`dto/asset.ts`):** flat MVP superset + ported optional legacy fields + unioned enums + ported typed-extra interfaces (vehicle/dwelling/document) + relationship DTOs (`IAssetDtoGroup`/`ISubAssetInfo`/linking) — mirroring the unified backend contract.
- **Components / services / pages:** legacy Angular files relocated and adapted under `ext-assetus/src/lib/{components,services,pages}`, consuming the unified DTOs and re-exported from the library index; MVP equivalents reconciled rather than duplicated.
- **Test runner:** the workspace uses **vitest** (`@analogjs/vitest-angular`, `@vitest/coverage-v8`), not jest; coverage is gathered via `npx nx test ext-assetus --coverage --watch=false`.
- **Sibling Features (out of this Feature, not lost):** the unified backend model (`unified-assetus-data-model`); the standalone *liabilities & service-providers* UI (provider/plan/settlement records and their dedicated pages) — tracked in the coverage table.

## Not Doing

- **Net-new frontend capability beyond legacy ∪ MVP** — migration/parity only.
- **Standalone liabilities provider/plan/settlement UI** — accounted for in the coverage table; the dedicated provider/plan records and pages may be built in the liabilities sibling Feature. Only the asset-side liability linkage fields are folded into the core DTO here.
- **Deleting the legacy frontend directories** — a follow-up once consumers are repointed and parity is proven.
- **Backend changes** — owned by `unified-assetus-data-model`; this Feature only consumes the unified contract.

## Open Questions

- **Rich fields not yet on the create/update HTTP requests (HUMAN DECISION).** The implemented `CreateAssetRequest`/`UpdateAssetRequest` accept only the flat MVP set; the rich optional fields (`type`/`possession`/`countryID`/`AssetDates`/relationships/financials/typed extras) live on the persisted `AssetDbo` but are not settable through the create/update API, and create does not accept `status`. This Feature mirrors the **current** backend request surface on the frontend request DTOs and the full superset on the record DTO. Whether to widen the backend create/update requests to carry the rich fields (so the ported UI can write them) is owned by `unified-assetus-data-model` and must be resolved by the owner; until then the ported create/edit components can only round-trip the flat set through the API.
- **Tasks-module reminder integration (sibling, not a cut).** The legacy vehicle NCT/tax/service due-dates carried task-link IDs into a tasks module. The backend `AssetVehicleExtra` carries only the plain due-date values (`nctExpires`/`taxExpires`/`nextServiceDue`); re-wiring live task reminders is deferred to the same sibling Feature noted in the backend Feature.
- **`yearOfBuild` vs `dateOfBuild` — resolved by the backend.** The backend keeps BOTH: `yearOfBuild` (int) on `AssetBase` and `dateOfBuild` (string) on `AssetDates`. The frontend DTO mirrors both.
- **Liabilities UI home.** Whether the standalone provider/plan/settlement UI is built here or in the liabilities sibling Feature is recorded in the coverage table; this Feature folds only the asset-side linkage fields (`liabilities`/`notUsedServiceTypes`) into core.

---
*This document follows the https://specscore.md/feature-specification*
