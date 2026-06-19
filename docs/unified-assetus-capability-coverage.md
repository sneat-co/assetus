# Unified Assetus — Capability-Coverage Table & Retirement Note

**Date:** 2026-06-19
**Status:** Authoritative coverage map for the legacy → unified port (Tasks 1–6).
**Companion to:** `docs/legacy-gap-analysis.md` (the §3.1 capabilities table and §4
deferred list this document maps from).

This document closes the loop on the legacy gap analysis. Every legacy capability
from `docs/legacy-gap-analysis.md` §3.1 (the side-by-side capabilities table) and
§4 (the "what was NOT migrated" deferred list) — **plus** the extra capabilities
surfaced by the code audits — has exactly one row below, mapped to one of:

- **Covered** — a concrete Go type/field/package in `assetus/backend` (cited and
  verified against the source).
- **Intentional change** — a deliberate model difference from legacy.
- **Sibling Feature: liabilities** — deferred to the liabilities/service-provider
  sub-module (a named follow-up Feature).
- **Sibling Feature: frontend port** — deferred to the Angular frontend-port
  Feature (the ~19 legacy components and the asset/real-estate/group pages).

No capability is left unmapped.

---

## 1. Coverage table

### 1.1 §3.1 capabilities (legacy gap-analysis side-by-side)

| Legacy capability | Disposition | Notes |
|---|---|---|
| Membership-gated CRUD on a space module | Covered — `facade4assetus` (create/get/update/remove) over `dal4assetus`, same `dal4spaceus` membership gating as legacy | Unchanged persistence pattern at `/spaces/{spaceID}/ext/assetus/assets/{assetID}`. |
| Tags | Covered — `dbo4assetus.AssetBase.Tags []string` | Validated non-empty in `AssetBase.Validate`. |
| Custom fields | Covered — `dbo4assetus.AssetBase` embeds `dbmodels.WithCustomFields` | Re-instated from legacy (gap analysis had marked it dropped in the MVP); now ported. |
| Condition (new/good/needs-repair/…) | Covered — `const4assetus.Condition` + `AssetBase.Condition` | MVP addition, retained. |
| Ownership-lifecycle Status (active/transferred/archived/disposed/lost) | Covered — `const4assetus.Status` (`status.go`) + `AssetBase.Status` | Union also retains legacy `draft` (`StatusDraft`). |
| Visibility (private…public, Space-default inheritance) | Covered — `const4assetus.Visibility` + `AssetBase.Visibility` | MVP addition, retained. |
| Append-only history | Covered — `dbo4assetus.AssetHistoryEvent` + `dal4assetus/asset_history.go` + facade record/get-history | Child collection `history/`. |
| Ownership transfer (Space→Space) + history | Covered — `facade4assetus` transfer | Relocates asset + history, appends `Transferred` event. |
| Owner-type derivation from Space type | Covered — `dbo4assetus.Owner` / `owner.go`, `const4assetus.OwnerType` | Derived at read. |
| Soft-archive vs hard-delete distinction | Covered — `facade4assetus` remove (soft-archive default + explicit hard-delete) | MVP behaviour, retained. |
| Typed sub-categories (`AssetType`) | Covered — `const4assetus.Type` (`asset_type.go`) + `AssetBase.Type`, validated per-Category via `const4assetus.ValidateType` | Full legacy value set ported (vehicle/dwelling/sports/document subtypes). |
| Possession (owning/leasing/renting) | Covered — `const4assetus.Possession` (`possession.go`) + `AssetBase.Possession`, default via `AssetBase.WithPossessionDefault` | Optional; defaults to `owning`. |
| Polymorphic extras (vehicle/dwelling/document) | Covered — `extras4assetus` (`AssetVehicleExtra`/`AssetDwellingExtra`/`AssetDocumentExtra`), registered via `RegisterAssetExtraFactory` | `extra.Data` interface; factory registration in each file's `init`. |
| Make/model/VIN/engine specs | Covered — `extras4assetus.AssetVehicleExtra` (`WithMakeModelRegNumberFields`, `Vin`) + `WithEngineData` (type/fuel/CC/kW/Nm) | See also engine rows below. |
| Mileage log (child collection) | Covered — `dbo4assetus.VehicleMileage` in `VehicleRecordDbo`, child collection `vehicleRecords/` (`dal4assetus/vehicle_record.go`), `facade4assetus.AddVehicleRecord` | Now a real facade (legacy was a TODO stub). |
| Insurance/certification dates (`AssetDates`) | Covered — `dbo4assetus.AssetDates` (dateOfBuild/dateOfPurchase/dateInsuredTill/dateCertifiedTill) embedded in `AssetBase` | Date-string validated. |
| Liabilities / service providers / plans | **Sibling Feature: liabilities** — asset-side linkage seam present (`dbo4assetus.AssetLiabilityInfo`, `LiabilityServiceType`, `AssetBase.Liabilities`, `NotUsedServiceTypes`); the provider/plan/settlement sub-domain itself is deferred | Only the asset→liability hook is in this repo; `DtoServiceProvider`/`DtoServicePlan`/`SettlementType` are the sibling Feature. |
| Asset groups (`IAssetDtoGroup`) | Covered — `dbo4assetus.AssetGroupInfo` (order/desc/categoryID/numberOf/totals) + `AssetGroupCounts`, linked via `WithAssetRelationships.GroupID`/`Group` | Full group shape, not just a groupId. |
| Sub-assets / parent asset | Covered — `WithAssetRelationships.ParentAssetID` + `SubAssets []SubAssetInfo` | Per-sub detail in `SubAssetInfo`. |
| Multi-space association (`WithAssetSpaces`) | Covered — `dbo4assetus.WithAssetSpaces` (spaceID→`AssetusSpaceBrief`), additive to owning `WithSpaceIDs` | Owning Space stays the single authoritative owner. |
| Required `CountryID` | Covered (as optional) — `AssetBase.CountryID geo.CountryAlpha2` | **Intentional change:** required-in-legacy → optional-in-unified; validated when present. |
| Member-level info | Covered — `WithAssetRelationships.MemberIDs` + `MembersInfo []TitledRecord` | Per-asset member IDs and per-member titled info. |
| Reserved `photos[]`, `ext.yardius` | Covered — `AssetBase.Photos []string` (reserved); `ext.yardius` intentionally absent | MVP forward-compat hook, retained. |

### 1.2 §4 deferred list ("what was NOT migrated")

| Legacy capability (deferred item) | Disposition | Notes |
|---|---|---|
| 1. Vehicles sub-domain (make/model/VIN, engine, mileage, NCT/tax/service due-dates) | Covered — `extras4assetus.AssetVehicleExtra` (+ `WithEngineData`, `WithMakeModelRegNumberFields`) with `NctExpires`/`TaxExpires`/`NextServiceDue`; mileage via `VehicleRecordDbo` | Reminders are stored as plain due-date strings (no task-link). |
| 2. Real-estate/dwelling sub-domain (address, rent price, bedrooms, area) | Covered — `extras4assetus.AssetDwellingExtra` (Address/RentPrice/NumberOfBedrooms/AreaSqM) | — |
| 3. Documents sub-domain (passport/licence/etc. with issue/expiry dates) | Covered — `extras4assetus.AssetDocumentExtra` (docType/number/issuedOn/effectiveFrom/expiresOn) | Per-doc-type validation — see §1.3. |
| 4. Typed subcategories (`AssetType`) + per-category validation | Covered — `const4assetus.Type` + `ValidateType(category, type)` | Per-category membership enforced. |
| 5. Possession semantics (owning/leasing/renting) | Covered — `const4assetus.Possession` | See §1.1 Possession row. |
| 6. Liabilities, service providers, service plans, settlement types | **Sibling Feature: liabilities** | Asset-side linkage seam only in this repo (see §1.1 liabilities row). |
| 7. Asset groups, sub-assets/parent asset, `sameAssetID` linking | Covered — `WithAssetRelationships` (`GroupID`/`Group`, `ParentAssetID`/`SubAssets`, `SameAssetID`) | `sameAssetID` = "same underlying asset" link. |
| 8. Multi-space association + member-level info | Covered — `WithAssetSpaces` + `WithAssetRelationships.MemberIDs`/`MembersInfo` | — |
| 9. Insurance/certification date tracking (`AssetDates`) + reminder/task-link fields | Covered (data) — `AssetDates`; vehicle reminders as due-date strings on `AssetVehicleExtra` | **Intentional change:** task-link/reminder wiring is not ported (plain date values only). |
| 10. Custom fields, required `CountryID`, `geo`/`relatedAs` hooks | Covered — `WithCustomFields`; `CountryID` (optional); `AssetBase.Geo *GeoPoint`; `WithOptionalRelatedAs` (in `WithAssetRelationships`) | `CountryID` optional, not required (intentional). |
| 11. ~19 legacy Angular components + asset/real-estate/group pages | **Sibling Feature: frontend port** | Backend data model homes all fields these UIs read; the components/pages themselves are the frontend Feature. |

### 1.3 Extra capabilities found by the code audits

| Legacy capability (audit-found) | Disposition | Notes |
|---|---|---|
| Vehicle **fuel** records: `VehicleFuelRecord` volume/unit/amount | Covered — `dbo4assetus.VehicleFuelRecord` (Volume/Unit/Amount) in `VehicleRecordDbo.Fuel` | — |
| Vehicle fuel records: fuelCost + currency | Covered — `VehicleFuelRecord.FuelCost` + `Currency` | Sourced from legacy `dal4assetus.Mileage` / `AddVehicleRecordRequest`. |
| `engineSerialNumber` | Covered — `extras4assetus.WithEngineData.EngineSerialNumber` | From legacy frontend `IEngine.engineSerialNumber`. |
| Full document extra: docType | Covered — `AssetDocumentExtra.DocType const4assetus.Type` | — |
| Full document extra: batchNumber | Covered — `AssetDocumentExtra.BatchNumber` | — |
| Full document extra: issuedBy | Covered — `AssetDocumentExtra.IssuedBy` | — |
| Full document extra: countryID | Covered — `AssetDocumentExtra.CountryID geo.CountryAlpha2` | Country-code validated when present. |
| Per-doc-type validation: `standardDocTypesByID` | Covered — `extras4assetus/doc_type_schema.go` (`standardDocTypesByID`, `DocTypeDef`, `DocTypeStandardFields`, `DocTypeField`), applied in `AssetDocumentExtra.validateDocTypeSchema` | Passport/driving-license require number+validity; birth/marriage cert require number+issuedOn, exclude validity. |
| Financial dimension: per-asset/group totals | Covered — `AssetBase.Totals []money.Amount`; group totals on `AssetGroupInfo.Totals` | OWNER DECISION (Task 2): optional fields on core asset, not a separate module. |
| Financial dimension: canHaveIncome / canHaveExpense | Covered — `AssetBase.CanHaveIncome` / `CanHaveExpense` | — |
| Financial dimension: income/expense direction | Covered — `AssetBase.FinancialDirection` (validated `""`/`income`/`expense`) | — |
| `isRequest` | Covered — `AssetBase.IsRequest bool` | — |
| `parentCategoryID` | Covered — `AssetBase.ParentCategoryID const4assetus.Category` | Validated as a Category when present. |
| Asset groups: order/desc/categoryId/numberOf | Covered — `AssetGroupInfo.Order`/`Desc`/`CategoryID`/`NumberOf` (`AssetGroupCounts`) | — |
| `ISubAssetInfo`: type/countryId/subType/expires | Covered — `dbo4assetus.SubAssetInfo` (Type/CountryID/SubType/Expires + embedded `TitledRecord`) | — |
| debt category | Covered — `const4assetus.CategoryDebt` | First-class Category value. |
| misc category | Covered — legacy `misc` → `const4assetus.CategoryOther` via `LegacyCategoryToUnified` | **Intentional change:** folded into `other`. |
| EngineType ↔ FuelType validation | Covered — `const4assetus.ValidateEngineFuel` (`engine_fuel.go`), called from `WithEngineData.Validate` | Legacy engine↔fuel compatibility matrix preserved. |
| Query/index fields | Covered — `extra.Data.IndexedFields()`: vehicle `make/model/make+model/regNumber/vin`; document `expiresOn/effectiveFrom` | Ported verbatim from legacy `IndexedFields` declarations. |
| Legacy enum value totality (Category/Status/Possession/Type/Engine/Fuel) | Covered — `const4assetus/legacy_mapping.go` (`LegacyCategoryToUnified`, `LegacyStatusToUnified`, `LegacyPossessionToUnified`, `LegacyTypeToUnified`, `LegacyEngineTypeToUnified`, `LegacyFuelTypeToUnified`) | Every legacy value maps to exactly one unified counterpart; none dropped. |

**Coverage summary:** every legacy capability has a row; every "Covered" row cites a
verified `assetus/backend` symbol; the only non-Covered dispositions are the two
named sibling Features (liabilities sub-module; frontend port) and the explicitly
intentional changes (optional `CountryID`, `misc`→`other`, no reminder/task-link
wiring). **No capability is left unmapped.**

---

## 2. Retirement note

Once the consumers below are repointed and `go build ./... && go test ./...` pass,
the following legacy backend directories become **deletable**:

- `sneat-go-backend/pkg/extensions/assetus/` — the entire legacy backend extension
  (const4assetus/dbo4assetus/extras4assetus/facade/api). Its still-wanted
  capabilities now live in `assetus/backend` (this repo) per §1.
- `sneat-libs/libs/extensions/assetus/{core,components}` — **NOT yet deletable.**
  Pending the **frontend-port sibling Feature**. The backend data model already
  homes every field these DTOs/components read, but the Angular code itself is
  retired only after that Feature lands.
- `sneat-apps/libs/extensions/assetus/pages` (assets/asset/new-asset/real-estates/
  asset-group/liability) — **NOT yet deletable.** Same frontend-port sibling
  Feature gate.

### 2.1 Consumers to repoint FIRST (backend)

These import or register the legacy package and must be changed before the legacy
backend directory is removed; afterwards `go build ./... && go test ./...` must pass:

1. `pkg/extensions/brandus/dbo4brands/make_test.go` — imports
   `const4assetus.AssetCategory`. Repoint to the unified
   `const4assetus.Category` (this repo) or drop the dependency.
2. `pkg/extensions/standard_extensions.go` — registers `assetus.Extension()`.
   Remove (or repoint to the unified extension) when the legacy package is deleted.
3. `sneat-apps` — consumes `@sneat/ext-assetus-components` and
   `@sneat/mod-assetus-core`; repoint or retire, then
   `pnpm install && nx run-many -t lint test build`. (Bundled with the
   frontend-port sibling Feature.)

### 2.2 Not-yet-deletable (pending sibling Features)

| Item | Blocked by |
|---|---|
| `sneat-libs/.../assetus/{core,components}` | Sibling Feature: frontend port |
| `sneat-apps/.../assetus/pages` | Sibling Feature: frontend port |
| Liabilities/service-provider sub-domain (`DtoServiceProvider`, `DtoServicePlan`, `DtoServiceType`, `SettlementType`) | Sibling Feature: liabilities sub-module |

---

## 3. Build & test verification

From `assetus/backend`:

- `go build ./...` → exit 0.
- `go test ./...` → exit 0.

Both confirmed green for this Task 6 deliverable.
