---
format: https://specscore.md/plan-specification
status: Implemented
---
# Plan: Assetus Frontend Port — Port Plan

**Status:** Implemented
**Source Feature:** assetus-frontend-port
**Date:** 2026-06-19
**Owner:** alex
**Supersedes:** —

## Summary

Port the legacy Assetus Angular frontend — `core` DTOs, the ~19 `components`, the services, and the pages — into the new MVP-shaped workspace at `assetus/frontend/libs/ext-assetus`, so the migrated DTOs reach parity with the unified backend model (`unified-assetus-data-model`) with **zero functionality lost** and **≥80% test coverage**. This is a **migration, not a rewrite**: each task relocates and adapts existing legacy Angular/TS files, naming the legacy source and the new merge target, and reconciling with the MVP equivalents already present rather than authoring from scratch.

Legacy source roots: `sneat-libs/libs/extensions/assetus/{core,components}/src/lib/` (DTOs + components) and `sneat-apps/libs/extensions/assetus/pages/src/lib/` (pages). New merge target root: `assetus/frontend/libs/ext-assetus/src/lib/`.

## Approach

Bottom-up dependency order: DTOs/enums → typed extras → relationship DTOs → services → components → pages → coverage table + green build → tests. The backend is now implemented, so each DTO/service task adapts the legacy file to mirror the **concrete backend json field names, enum string values and route paths** in the Feature's REQ:backend-wire-contract — NOT a guessed shape. Each task's mechanic is **copy-or-relocate the legacy file, then adapt it onto the MVP equivalent and the backend contract** — match the backend enums (`const4assetus`), add the backend's optional fields with their exact json names (`groupID`/`categoryID`/`countryID`, etc.), keep the MVP `IAssetDbo`/`AssetService`/history-timeline/MVP-pages core as the spine, and reconcile (not duplicate) the MVP files at `assetus/frontend/libs/ext-assetus/src/lib/`. Legacy validation (per-doc-type `standardDocTypesByID`, engine/fuel data, per-category `type`) and the full vehicle fuel-record payload are carried over to match the backend (`extras4assetus`/`dbo4assetus`). The standalone liabilities provider/plan UI is accounted for in the coverage table and may be deferred to the liabilities sibling Feature. The workspace test runner is **vitest** (`@analogjs/vitest-angular`, `@vitest/coverage-v8`), not jest; coverage is gathered via `npx nx test ext-assetus --coverage --watch=false`. The final tasks prove nothing was lost (frontend capability-coverage table), that the workspace lints/builds/tests green, and that coverage is ≥80% with a test per ported DTO/component.

Backend contract anchors (from REQ:backend-wire-contract): `AssetBase`/`AssetDbo`/`AssetBrief` (`backend/dbo4assetus/asset_dbo.go`), `WithAssetRelationships`/`AssetGroupInfo`/`SubAssetInfo` (`backend/dbo4assetus/relationships.go`), `VehicleRecordDbo`/`AddVehicleRecordRequest` (`backend/dbo4assetus/vehicle_record_dbo.go` + `backend/dto4assetus/add_vehicle_record.go`), typed extras (`backend/extras4assetus/`), enums (`backend/const4assetus/`), create/update requests (`backend/dto4assetus/{create,update}_asset.go`), routes (`backend/api4assetus/routes.go`). MVP frontend merge targets: `dto/asset.ts`, `services/{asset.service.ts,interfaces.ts,assetus-core-services.module.ts}`, `assetus-routing.ts`, `pages/{assets,asset,transfer}`, `pages/assets/new-asset-dialog.component.ts`, `components/asset-history-timeline`, `components/index.ts`, `dto/index.ts`, `index.ts`.

## Tasks

### Task 1: Port & union the DTO core + enums (`dto/asset.ts`)

**Verifies:** assetus-frontend-port#ac:dto-superset-roundtrip, assetus-frontend-port#ac:enum-union-no-value-dropped
**Depends-On:** —
**Status:** done

Relocate the legacy enums (`sneat-libs/.../core/src/lib/dto/assetus-types.ts`) and the legacy `IAssetBrief`/`IAssetDbo`/`AssetLiabilityInfo` fields (`sneat-libs/.../core/src/lib/dto/dto-asset.ts`) into the MVP `assetus/frontend/libs/ext-assetus/src/lib/dto/asset.ts`, reconciling onto the MVP enums/fields already there and mirroring the **backend** json names + enum values (REQ:backend-wire-contract; `backend/dbo4assetus/asset_dbo.go`, `backend/const4assetus/`). Match `AssetCategory` to the backend `Category` set (`books`,`games`,`toys`,`sports_equipment`,`tools`,`electronics`,`clothing`,`vehicles`,`camping_equipment`,`other`,`dwelling`,`document`,`debt` — reconciliation already baked in), `AssetStatus` to `draft`/`active`/`transferred`/`archived`/`disposed`/`lost`, `AssetPossession` to `unknown`/`undisclosed`/`owning`/`leasing`/`renting`, `AssetType` to the per-category `Type` subtypes; keep `AssetCondition`/`AssetVisibility` unchanged; type fuel-volume unit / mileage unit / currency as open strings (no backend enum). Add the optional record fields with the **exact backend json names** (`isRequest`/`countryID`/`type`/`AssetDates` (`dateOfBuild`/`dateOfPurchase`/`dateInsuredTill`/`dateCertifiedTill`)/`yearOfBuild`/`possession`/`parentAssetID`/`parentCategoryID`/`geo`/`photos`/`memberIDs`/`membersInfo`/`groupID`+`group`/`subAssets`/`sameAssetID`/`relatedAs`/custom fields/`extraType`/`extra`/`totals`/`canHaveIncome`/`canHaveExpense`/`financialDirection`/`liabilities`/`notUsedServiceTypes`) as **optional** so a minimal MVP asset stays valid. Note the draft-vs-real corrections: `groupID` (not `groupId`) + `group` object, `categoryID`/`countryID` capital-ID, and there is no separate `title`/`desc` (`name` is canonical).

### Task 2: Port the polymorphic typed extras + vehicle fuel-record (`dto/`)

**Verifies:** assetus-frontend-port#ac:vehicle-extra-frontend-no-field-dropped, assetus-frontend-port#ac:document-extra-frontend-full-shape
**Depends-On:** 1
**Status:** done

Relocate the legacy typed-extra DTOs into `assetus/frontend/libs/ext-assetus/src/lib/dto/`, resolved by `extraType` ∈ {`vehicle`,`dwelling`,`document`}, mirroring the backend `extras4assetus` shapes (REQ:backend-wire-contract). Vehicle (`backend/extras4assetus/asset_vehicle.go`+`with_engine_data.go`+`with_make_model.go`): `make`/`model`/`regNumber`/`vin`, engine `engineType`/`engineFuel`/`engineCC`/`engineKW`/`engineNM`/`engineSerialNumber`, plain due-dates `nctExpires`/`taxExpires`/`nextServiceDue` (no task-link IDs on the backend — port plain dates). Vehicle record: both the persisted `VehicleRecordDbo{fuel:{volume,unit,amount,fuelCost,currency},mileage:{value,unit}}` (`backend/dbo4assetus/vehicle_record_dbo.go`) and the flat append request `AddVehicleRecordRequest{fuelVolume,fuelVolumeUnit,fuelCost,currency,mileage,mileageUnit}` (`backend/dto4assetus/add_vehicle_record.go`) — units/currency as plain strings. Dwelling (`backend/extras4assetus/asset_dwelling.go`): `address`/`rent_price`(`{value,currency}`)/`numberOfBedrooms`/`areaSqM`. Document (`backend/extras4assetus/asset_document.go`+`doc_type_schema.go`): `docType` (`passport`/`id_card`/`driving_license`/`marriage_cert`/`birth_cert`)/`number`/`batchNumber`/`countryID`/`issuedBy`/`issuedOn`/`effectiveFrom`/`expiresOn` (+`regNumber` alias), with the `standardDocTypesByID`/`DocTypeStandardFields` per-doc-type schema (passport/driving_license: number+validTill required; marriage_cert: number+issuedOn, exclude validTill, Members max 2). No field dropped; an asset with no extra stays valid.

### Task 3: Port the relationship DTOs (groups, sub-assets, linking, multi-space, members)

**Verifies:** assetus-frontend-port#ac:relationships-frontend-preserved
**Depends-On:** 2
**Status:** done

Adapt the legacy relationship DTOs onto the unified frontend record, mirroring the backend `WithAssetRelationships` (`backend/dbo4assetus/relationships.go`, REQ:backend-wire-contract): the group linkage `groupID` (string) + group sub-entity `group`=`AssetGroupInfo` (`id`/`title`/`order`/`desc`/`categoryID`/`numberOf`=`{assets}`/`totals`) and its `AssetGroup` uimodel (`sneat-libs/.../core/src/lib/uimodels/asset-group.ts`); parent/sub-asset nesting `parentAssetID`+`subAssets` with per-sub-asset `SubAssetInfo` (`id`/`title`/`type`/`countryID`/`subType`/`expires`); asset linking `sameAssetID`/`relatedAs`; multi-space association alongside the canonical owning Space; and `memberIDs` (string[]) + `membersInfo` (`{id,title}[]`). Use the exact backend json names (capital-ID `groupID`/`categoryID`/`countryID`; `group` is an object, not a `groupId` string). Relocate into `assetus/frontend/libs/ext-assetus/src/lib/dto/` (+ a `uimodels`/`contexts` home as needed) and export from `dto/index.ts`/`index.ts`.

### Task 4: Port & merge the Assetus frontend services

**Verifies:** assetus-frontend-port#ac:services-frontend-ported
**Depends-On:** 3
**Status:** done

Merge the legacy frontend services (`sneat-libs/.../components/src/lib/services/{asset-service.dto.ts,asset-service.ts,assetus-space.service.ts,assetus-services.module.ts}`) into `assetus/frontend/libs/ext-assetus/src/lib/services/`, adapting onto the existing MVP `asset.service.ts`/`interfaces.ts`/`assetus-core-services.module.ts` so one asset service exposes both the MVP ownership operations and the ported create/update/add-vehicle-record operations. Wire the methods to the **actual implemented routes** (`backend/api4assetus/routes.go`, REQ:backend-wire-contract): the MVP service already covers `create_asset`/`asset`/`update_asset`/`remove_asset`/`transfer_asset`/`record_history_event`/`asset_history`; ADD the add-vehicle-record method pointing at `create_vehicle_record` (NOT `add_vehicle_record`). Mirror the current backend request surface: `ICreateAssetRequest`/`IUpdateAssetRequest` stay the flat MVP set (create takes no `status`), `IAddVehicleRecordRequest` carries the fuel-bearing payload (`fuelVolume`/`fuelVolumeUnit`/`fuelCost`/`currency`/`mileage`/`mileageUnit`). Carry the assetus space service; reconcile, do not duplicate. (See the create/update request-surface divergence note in the Feature's Open Questions.)

### Task 5: Port the ~19 components, reconciling with the MVP component

**Verifies:** assetus-frontend-port#ac:all-components-accounted
**Depends-On:** 4
**Status:** done

Relocate every legacy component from `sneat-libs/.../components/src/lib/` into `assetus/frontend/libs/ext-assetus/src/lib/components/` (alongside `asset-history-timeline`), adapting each onto the unified DTOs/services: `vehicle-card`, `vehicle-engine`, `make-model-card`, `asset-reg-number-input`, `mileage-dialog`, `edit-dwelling-card`, `real-estate-location`, `asset-dates`, `period-segment`, `asset-possesion-card` (→ `asset-possession-card`), `asset-liabilities`, `asset-contacts-group`, `asset-add-{vehicle,dwelling,document,service}`, `assets-list`+`asset-list-item`, `asset-card`, plus the shared bases (`add-asset-base-component`, `asset-base-page`, `asset-component-base-params`) and the `car-makes-with-models`/`vehicles` data. Re-export ported components from `components/index.ts`; any component intentionally not ported gets an explicit deferral reason recorded for the coverage table (Task 7). The standalone liabilities `asset-add-service` UI may be marked for the liabilities sibling Feature.

### Task 6: Port the pages and wire routing, reconciling with the MVP pages

**Verifies:** assetus-frontend-port#ac:all-pages-accounted
**Depends-On:** 5
**Status:** done

Relocate every legacy page from `sneat-apps/.../pages/src/lib/` into `assetus/frontend/libs/ext-assetus/src/lib/pages/` and wire into `assetus-routing.ts`: `assets`, `asset`, `new-asset`, `real-estates`+`real-estate`, `asset-group`, `optimization`, the `liability/*` pages (`liability-new`, `select-service-provider`), and the shared `asset-base.page`/`assets-base.page`. Reconcile the existing MVP pages (`assets`, `asset`, `transfer`, `new-asset-dialog`) into the ported set rather than duplicating; any page intentionally not ported (e.g. the standalone liabilities pages → liabilities sibling Feature) gets an explicit deferral reason recorded for the coverage table (Task 7).

### Task 7: Frontend capability-coverage table + liabilities disposition + green build

**Verifies:** assetus-frontend-port#ac:liabilities-disposition-recorded, assetus-frontend-port#ac:frontend-capability-coverage-complete, assetus-frontend-port#ac:frontend-builds-and-tests-pass
**Depends-On:** 6
**Status:** done

Author the frontend capability-coverage table mapping every legacy frontend capability (each `core` DTO incl. `dto-liability.ts`/`dto-service-provider.ts`, each of the ~19 components, each page, each service) to its ported new location, an intentional change, or the named liabilities sibling Feature. Record the liabilities disposition explicitly: asset-side linkage fields folded into the core DTO (Task 1); standalone provider/plan/settlement DTOs/components/pages either ported or assigned to the sibling Feature — none unassigned. Make the workspace lint/build/test green from `assetus/frontend` (`nx run-many -t lint build test` or the `ext-assetus` targets exit zero).

### Task 8: Frontend test coverage (≥80% + a test per ported DTO/component)

**Verifies:** assetus-frontend-port#ac:frontend-test-coverage
**Depends-On:** 7
**Status:** done

Add/port unit tests so every ported DTO has a round-trip/shape test asserting the **backend json field names + enum string values** (incl. `draft`/`debt`/`dwelling`/`document` categories, `phev`/`hybrid`/`steam` engine types, `bio`/`hydrogen` fuels, the `groupID`/`categoryID`/`countryID` capital-ID names, the vehicle fuel-record extra both persisted and request shape, the document per-doc-type validation schema, the relationship DTOs) and every ported component has at least a render/smoke test. Carry over any reusable legacy specs rather than re-authoring. Verify overall statement coverage is ≥80% via `npx nx test ext-assetus --coverage --watch=false` using the workspace's vitest runner (`@analogjs/vitest-angular`, `@vitest/coverage-v8`), NOT jest.

## Open Questions

- **Rich fields not yet on the create/update HTTP requests (HUMAN DECISION).** The implemented `CreateAssetRequest`/`UpdateAssetRequest` (`backend/dto4assetus/{create,update}_asset.go`) accept only the flat MVP set (create takes no `status`); the rich optional fields live on `AssetDbo` but are not settable through the create/update API. Tasks 1/4 mirror the current request surface on the frontend request DTOs and the full superset on `IAssetDbo`. Whether to widen the backend requests is owned by `unified-assetus-data-model`; until resolved, ported create/edit components can only round-trip the flat set through the API.
- **Tasks-module reminder integration (sibling).** The backend `AssetVehicleExtra` carries only plain due-dates (`nctExpires`/`taxExpires`/`nextServiceDue`); re-wiring live task reminders is deferred to the sibling Feature noted in the backend Feature.
- **`yearOfBuild` vs `dateOfBuild` — resolved.** The backend keeps both (`yearOfBuild` int + `dateOfBuild` string on `AssetDates`); Task 1 mirrors both.
- **Standalone liabilities UI home.** Whether the provider/plan/settlement DTOs, the `asset-add-service`/`select-service-provider`/`liability-new` flows are built here or in the liabilities sibling Feature is recorded in the coverage table (Task 7); this plan folds only the asset-side linkage fields into core.

---
*This document follows the https://specscore.md/plan-specification*
