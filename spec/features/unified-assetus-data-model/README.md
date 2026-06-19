---
format: https://specscore.md/feature-specification
status: Approved
---

# Feature: Unified Assetus Data Model

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/sneat-co/assetus/spec/features/unified-assetus-data-model?op=explore) | [Edit](https://specscore.studio/app/github.com/sneat-co/assetus/spec/features/unified-assetus-data-model?op=edit) | [Ask question](https://specscore.studio/app/github.com/sneat-co/assetus/spec/features/unified-assetus-data-model?op=ask) | [Request change](https://specscore.studio/app/github.com/sneat-co/assetus/spec/features/unified-assetus-data-model?op=request-change) |
**Status:** Approved
**Source Ideas:** assetus-model-unification
**Grade:** B

## Summary

Reconcile the legacy rich-registry Assetus and the MVP ownership-core into a single superset backend data model (enums, record, typed extras, relationships, multi-space, member info), so legacy capabilities are all representable and the legacy directories can be retired with zero functionality lost.

## Problem

There are two Assetus implementations authored by the same owner for adjacent use cases: a **legacy** rich asset-registry (typed vehicle/dwelling/document extras, possession, multi-space association, asset groups, sub-assets, liabilities) and a newer **MVP** ownership-core (condition, ownership-lifecycle status, visibility, append-only history, Space→Space transfer, owner derivation). `docs/legacy-gap-analysis.md` declared them irreconcilable and recommended deferring the legacy capabilities behind `RELOCATED.md` notices.

The approved Idea `assetus-model-unification` rejects that premise: with **no live legacy data** (greenfield), almost every conflict resolves by union or a two-level taxonomy. The legacy code currently sits stranded behind relocation notices; the owner's hard condition is **no functionality lost** before any legacy directory is deleted. This Feature defines the single superset backend model that makes every legacy capability representable, so the legacy directories become deletable.

## Behavior

### Enums

#### REQ: category-two-level

The unified model MUST restore the legacy **two-level** taxonomy: a `Category` plus an optional per-category `Type` subtype. `Category` MUST be the **union** of the legacy and MVP category sets — legacy `vehicle`, `dwelling`, `document` and MVP `books`, `games`, `toys`, `sports_equipment`, `tools`, `electronics`, `clothing`, `vehicles`, `camping_equipment`, `other` — with the legacy `sport_gear` value mapped onto `sports_equipment` and legacy `vehicle` reconciled with MVP `vehicles` to a single canonical value. The per-category `Type` (e.g. vehicles: `car`/`van`/`motorcycle`/`boat`/…; real-estate: `apartment`/`house`/…; documents: `passport`/`driving_license`/…; sports: `bicycle`/`kite`/`surf_board`/…) MUST be preserved and validated per category where the legacy model validated it, and MUST be optional (no subtype required) for categories that had none.

#### REQ: status-union

The unified `Status` enum MUST be the **union** of both lifecycles: `draft`, `active`, `transferred`, `archived`, `disposed`, `lost`. `draft` (legacy editing state) and `active`/`archived` (shared) MUST be preserved; the MVP ownership-lifecycle states `transferred`, `disposed`, `lost` MUST be preserved. No legacy or MVP `Status` value may be dropped.

#### REQ: possession-optional

The legacy `Possession` enum (`unknown`, `undisclosed`, `owning`, `leasing`, `renting`) MUST be reintroduced as an **optional** field defaulting to `owning`. Possession describes the legal relationship to the asset and MUST be orthogonal to `Status` (ownership-lifecycle) and to the canonical owning Space — i.e. a `leasing`/`renting` asset MUST remain representable without a separate record type.

#### REQ: condition-visibility-retained

The MVP `Condition` (`new`/`excellent`/`good`/`fair`/`needs_repair`/`broken`) and `Visibility` (`private`/`family`/`friends`/`friends_of_friends`/`specific_space`/`public`) enums MUST be retained unchanged. `Condition` MUST be optional (legacy records had no condition value). `Visibility` MUST default to the owning Space's default visibility when unset, preserving MVP semantics.

#### REQ: enum-mapping-total

Every legacy enum value (`AssetCategory`, `AssetStatus`, `AssetPossession`, `AssetType`, `AssetDocumentType`, `EngineType`, `FuelType`) MUST map to exactly one unified value, documented in a mapping table in the implementation. No legacy enum value may be left without a unified target. The mapping MUST explicitly account for the legacy `AssetCategory` values the gap-analysis omitted — `debt` and `misc` (and `undefined`) — mapping `misc`→`other` and retaining `debt` as a unified `Category` value (per the owner decision in REQ:financial-dimension-disposition to keep the financial dimension in core). No legacy enum value may be silently absent.

### Record

#### REQ: flat-core-superset

The unified `AssetDbo` MUST be a **superset** of the MVP flat `AssetBase` plus the still-wanted legacy fields as **optional**: `CountryID` (optional, not required as in legacy), `AssetDates` (`dateOfBuild`/`dateOfPurchase`/`dateInsuredTill`/`dateCertifiedTill`), `yearOfBuild` (legacy keeps it alongside `dateOfBuild`; preserve unless explicitly consolidated), custom fields, tags, the `geo` hook, `isRequest` (the legacy boolean distinguishing a wanted/request record from an owned asset), `parentCategoryID` (the legacy parent-**category** pointer, distinct from `parentAssetID`), and the financial fields folded into core per REQ:financial-dimension-disposition (`ITotalsHolder` totals, income/expense capability flags + direction, and the asset-side liability linkage `liabilities`/`notUsedServiceTypes`/`AssetLiabilityInfo`). Fields the MVP already carries (name, description, category, condition, status, visibility, acquisitionDate, purchasePrice, estimatedValue, location, notes, photos) MUST be retained.

#### REQ: ownership-core-preserved

The MVP ownership-core MUST be preserved unchanged: the append-only `history/` child collection, Space→Space `transfer` (relocating asset + history and appending a `Transferred` event), owner-type derivation from the owning Space, and the soft-archive-vs-hard-delete distinction. Reintroducing legacy structure MUST NOT regress any of these behaviors.

### Typed extras

#### REQ: polymorphic-extras

The legacy polymorphic **typed extras** MUST be reintroduced as an **optional** extension hanging off the flat core, registered by `extraType`, covering at least `vehicle`, `dwelling`, and `document`. An asset with no extra MUST remain valid (the flat core stands alone).

#### REQ: vehicle-extra-complete

The `vehicle` extra MUST carry every legacy vehicle attribute with no field dropped: make, model, registration number, VIN, engine data (`EngineType`/`FuelType`/CC/kW/Nm **and `engineSerialNumber`**), and the vehicle-record child collection. The vehicle-record child collection MUST preserve the **full** legacy `VehicleRecordDbo` shape — not just the odometer reading: the mileage (`VehicleMileage{value, unit}`) **and** the fuel-purchase data (`VehicleFuelRecord{volume, unit, amount}` plus the request-level `fuelCost`/`currency`), along with each record's creation metadata (`CreatedFields`). The `AddVehicleRecord` facade + request MUST preserve the fuel-bearing payload (`fuelVolume`/`fuelVolumeUnit`/`fuelCost`/`currency`/`mileage`/`mileageUnit`). The legacy **EngineType↔FuelType compatibility validation** (which fuel types are legal per engine type) MUST be preserved, not relaxed. The vehicle extra MUST also preserve the legacy **service/tax/inspection (e.g. NCT) due-date** fields; where the legacy tasks-module task-link integration is unavailable, those MUST degrade to plain due-date values stored on the extra rather than being dropped.

The `dwelling` extra MUST carry address, rent price, bedrooms, and area. The `document` extra MUST carry the **full** legacy field set with no field dropped: `docType` (the `AssetDocumentType` taxonomy — `passport`/`driving_license`/`birth_certificate`/`marriage_certificate`/`rent_lease`/`insurance_policy`/…), registration `number`, `batchNumber`, document `countryID`, `issuedBy`, `issuedOn`, `effectiveFrom`, and `expiresOn`. The per-document-type **field-validation schema** (legacy `standardDocTypesByID` / `IDocTypeStandardFields` — e.g. passport requires number+validTill, marriage certificate allows two members) MUST be preserved as the document-extra validation rules.

#### REQ: query-index-fields

The legacy extras' indexed-field declarations MUST be preserved so the corresponding queries remain possible: vehicle by `make`, `model`, `make+model`, `regNumber`, `vin`; document by `expiresOn`, `effectiveFrom`. Reintroducing the typed extras MUST NOT drop these index declarations.

### Relationships

#### REQ: groups-and-nesting

Asset **grouping** and **parent/sub-asset nesting** MUST be preserved with no field dropped. The asset **group** is a legacy sub-entity (`IAssetDtoGroup`), not merely a `groupId` string: its fields (`order`, `desc`, `categoryId`, `numberOf`/counts) MUST be preserved as a group record. Parent/sub-asset nesting (`parentAssetID`, `subAssets`) MUST preserve the per-sub-asset detail the legacy `ISubAssetInfo` carries (`type`, `countryId`, `subType`, `expires`), not collapse sub-assets to bare IDs.

#### REQ: asset-linking

Asset-to-asset **linking** (`sameAssetID` / `relatedAs` / the legacy `linkage` relation hook) MUST be preserved so two asset records can be linked as related or as the same underlying asset.

### Multi-space and ownership

#### REQ: canonical-owner-with-multispace

The unified model MUST add a **single canonical owning Space** that anchors the ownership-lifecycle, history, transfer, and owner derivation, AND MUST preserve the legacy **multi-space association** (`WithAssetSpaces`) so the same asset can still surface under several spaces. Adding the canonical owner MUST NOT remove an asset's ability to appear under multiple spaces; only the lifecycle/history is anchored to the one owning Space.

### Member info

#### REQ: member-level-info

Per-asset **member-level information** (`memberIDs` / `membersInfo`) MUST be preserved so an asset can record which members it concerns and their per-member info.

### Financial dimension

#### REQ: financial-dimension-disposition

The legacy model carries a financial dimension that MUST NOT be silently dropped: per-asset and per-group **totals** (`ITotalsHolder`), the per-category `canHaveIncome`/`canHaveExpense` capability flags, income/expense **direction**, the legacy `debt` category, and the asset-side liability linkage (`liabilities`, `notUsedServiceTypes`, `AssetLiabilityInfo` held on the asset itself). **Owner decision: these live as optional fields on the unified core asset/group**, not in a separate sibling Feature. Concretely: `ITotalsHolder` totals are optional fields on the unified asset and group; `canHaveIncome`/`canHaveExpense` and income/expense direction are retained on the category/asset; the legacy `debt` value is retained as a unified `Category` value (per REQ:enum-mapping-total); and the asset-side liability linkage fields are optional fields on the unified `AssetDbo`. (Standalone liability provider/plan records remain the separate liabilities sibling Feature; only the asset-side financial fields are folded into core here.)

### Testing

#### REQ: test-coverage

The ported backend MUST have comprehensive automated test coverage: a unit test exercising **every ported capability** — each unioned enum, each optional field round-trip, each typed extra (including the vehicle fuel-record), the relationship/multi-space/member-info structures, and the facades (create/get/update/remove/transfer/add-vehicle-record) — with overall backend statement coverage of **at least 80%** (`go test -cover ./...`). The frontend port (its own sibling Feature) MUST carry an **equivalent coverage bar** (≥80% plus a test per ported DTO/component) for the migrated Angular code; that requirement is recorded here so it is not lost when the frontend-port Feature is specified.

### Retirement readiness

#### REQ: capability-coverage-total

Every legacy capability enumerated in `docs/legacy-gap-analysis.md` §3.1 and §4 MUST be representable in the unified model, documented in a capability-coverage table that maps each legacy capability to its unified field/structure (or to an explicitly-named sibling Feature for liabilities and the frontend port). No legacy capability may be left with no unified home. This requirement is the testable encoding of the "no functionality lost" condition.

#### REQ: build-green

The unified backend (`backend/.../const4assetus`, `dbo4assetus`, `extras4assetus`, and facades) MUST compile and its tests MUST pass: `go build ./...` and `go test ./...` succeed from the repo's `backend/` module.

## Acceptance Criteria

### AC: category-superset (verifies REQ:category-two-level, REQ:enum-mapping-total)

**Given** a legacy asset with `Category=sport_gear` and `Type=kite_board`, and an MVP asset with `Category=books`,
**When** both are expressed in the unified model,
**Then** the first maps to the unified `Category=sports_equipment` with `Type=kite_board` preserved, the second to `Category=books` with no subtype, and the legacy↔unified mapping table contains an entry for every legacy `AssetCategory`, `AssetType`, `EngineType`, and `FuelType` value (no engine/fuel enum value dropped).

### AC: out-of-category-type-rejected (verifies REQ:category-two-level)

**Given** an asset with `Category=document` and `Type=kite_board` (a sports subtype not valid for documents),
**When** it is validated in the unified model,
**Then** validation rejects the asset because `Type` is not in the set permitted for `Category=document`, preserving the legacy per-category `Type` validation.

### AC: status-no-value-dropped (verifies REQ:status-union)

**Given** the legacy `AssetStatus` set `{active, archived, draft}` and the MVP `Status` set `{active, transferred, archived, disposed, lost}`,
**When** the unified `Status` enum is defined,
**Then** it contains exactly the union `{draft, active, transferred, archived, disposed, lost}` and no value from either source set is absent.

### AC: leasing-asset-representable (verifies REQ:possession-optional)

**Given** a legacy asset with `Possession=leasing`,
**When** it is expressed in the unified model,
**Then** its possession is stored as `leasing` on the optional possession field, its `Status` is independently set (e.g. `active`), and no separate record type is required to represent a leased asset.

### AC: possession-defaults-owning (verifies REQ:possession-optional)

**Given** an MVP-shaped asset created with no possession value supplied,
**When** it is persisted in the unified model,
**Then** its possession resolves to `owning` by default.

### AC: condition-optional-visibility-default (verifies REQ:condition-visibility-retained)

**Given** a legacy asset that carries no condition value, created in a Space whose default visibility is `family`,
**When** it is expressed in the unified model with no condition and no visibility supplied,
**Then** condition remains unset (valid) and visibility resolves to `family` (the owning Space default).

### AC: optional-legacy-fields-roundtrip (verifies REQ:flat-core-superset)

**Given** a legacy asset with `CountryID`, `dateInsuredTill`, custom fields, tags, and a `geo` value,
**When** it is stored and re-read via the unified `AssetDbo`,
**Then** all five are present and unchanged, and a unified asset created without any of them is still valid (each is optional).

### AC: history-transfer-intact (verifies REQ:ownership-core-preserved)

**Given** a unified asset owned by Space A with one or more history events,
**When** it is transferred to Space B,
**Then** the asset and its history relocate to Space B, a `Transferred` event is appended, owner-type is re-derived from Space B, and the append-only history is never mutated in place.

### AC: typed-extras-optional (verifies REQ:polymorphic-extras)

**Given** the unified model with registered `vehicle`, `dwelling`, and `document` extras,
**When** one asset is created with a `vehicle` extra and another with no extra at all,
**Then** both are valid — the first carries its typed extra resolved by `extraType`, and the second is a valid flat-core asset with no extra.

### AC: vehicle-extra-no-field-dropped (verifies REQ:vehicle-extra-complete)

**Given** a legacy vehicle asset with make, model, registration number, VIN, engine (type/fuel/CC/kW/Nm and engine serial number), a service/tax/inspection (NCT) due-date, and three vehicle records each carrying both mileage and fuel data (fuel volume, fuel unit, fuel cost, currency),
**When** it is expressed in the unified model,
**Then** every one of those attributes — including the engine serial number and the service/tax/inspection due-date (as a plain due-date value) — is stored on the `vehicle` extra, and the three records are preserved in the vehicle-record child collection with their mileage AND fuel volume/unit/cost/currency intact, with no attribute dropped.

### AC: document-extra-full-shape (verifies REQ:vehicle-extra-complete)

**Given** a legacy document asset of `docType=passport` with number, batch number, country, issued-by, issued-on, effective-from, and expires-on values,
**When** it is expressed in the unified model,
**Then** all of `docType`, `number`, `batchNumber`, `countryID`, `issuedBy`, `issuedOn`, `effectiveFrom`, and `expiresOn` are stored on the `document` extra, and the per-doc-type validation schema for `passport` (e.g. number + validity required) is applied — none of these fields is dropped.

### AC: financial-fields-have-a-home (verifies REQ:financial-dimension-disposition)

**Given** the legacy financial dimension (per-asset/group totals, `canHaveIncome`/`canHaveExpense`, income/expense direction, the `debt` category, and asset-side liability linkage),
**When** the capability-coverage table is checked,
**Then** every one of those financial fields has a recorded disposition — either an optional unified field or an assignment to the liabilities/financial sibling Feature — and none is left unassigned.

### AC: groups-nesting-linking-preserved (verifies REQ:groups-and-nesting, REQ:asset-linking)

**Given** a legacy asset that belongs to a group, has two sub-assets under a parent, is linked to one asset via `relatedAs`, and is linked to another as the same underlying asset via `sameAssetID`,
**When** it is expressed in the unified model,
**Then** the group membership, the parent/sub-asset nesting, the `relatedAs` link, and the `sameAssetID` ("same underlying asset") link are all preserved and resolvable.

### AC: multispace-with-canonical-owner (verifies REQ:canonical-owner-with-multispace)

**Given** a legacy asset associated with Spaces A, B, and C,
**When** it is expressed in the unified model with Space A as the canonical owner,
**Then** the asset still surfaces under Spaces A, B, and C, while its history, transfer, and owner-derivation are anchored to Space A.

### AC: member-info-preserved (verifies REQ:member-level-info)

**Given** a legacy asset carrying `memberIDs` and per-member `membersInfo`,
**When** it is expressed in the unified model,
**Then** the member IDs and per-member info are preserved and readable.

### AC: capability-coverage-complete (verifies REQ:capability-coverage-total)

**Given** the capability list in `docs/legacy-gap-analysis.md` §3.1 and §4,
**When** the unified model's capability-coverage table is checked against it,
**Then** every legacy capability has a row mapping it either to a concrete unified field/structure or to a named sibling Feature (liabilities sub-module; frontend port), and no capability is left unmapped.

### AC: backend-builds-and-tests-pass (verifies REQ:build-green)

**Given** the unified backend packages in the `backend/` module,
**When** `go build ./...` and `go test ./...` are run from `backend/`,
**Then** both commands exit zero.

### AC: backend-test-coverage (verifies REQ:test-coverage)

**Given** the merged `assetus/backend` module,
**When** `go test -cover ./...` is run from `backend/`,
**Then** every package's tests pass and overall statement coverage is at least 80%, with a test present for each ported capability (enum union, optional-field round-trip, typed extras incl. vehicle fuel-record, relationships/multi-space/member-info, and the facades).

## Architecture

- **Spine:** the MVP ownership-core (`dbo4assetus.AssetDbo` flat `AssetBase` + history + transfer + owner derivation) remains canonical and is extended, not replaced.
- **Enums (`const4assetus`):** unified `Category`(+`Type`), `Status`, `Possession`, `Condition`, `Visibility`, `OwnerType`, `HistoryEventType`, `EngineType`, `FuelType`, plus a documented legacy→unified mapping table.
- **Record (`dbo4assetus`):** flat superset core + optional legacy fields (`CountryID`, `AssetDates`, custom fields, `geo`) + relationship fields (`groupId`, `parentAssetID`, `subAssets`, `sameAssetID`/`relatedAs`) + multi-space `WithAssetSpaces` alongside the canonical owning Space + `memberIDs`/`membersInfo`.
- **Extras (`extras4assetus`):** polymorphic `vehicle`/`dwelling`/`document` extras registered by `extraType`, with `VehicleMileage` as a child collection.
- **Sibling Features (out of this Feature, not lost):** *liabilities & service-providers sub-module* — owns standalone provider/plan/settlement records only; the asset-side financial fields are folded into core here per the owner decision in REQ:financial-dimension-disposition. *Legacy frontend/Angular port*.
- **Consumers to repoint before any legacy deletion (deletion is out of scope here):** `brandus/dbo4brands/make_test.go`, `standard_extensions.go`, `sneat-apps`.

## Testing Strategy

All ACs are testable (Go unit tests + build/test commands), but Rehearse stub files are **deferred to the plan/implement phase** rather than scaffolded at spec time — the unified package layout isn't fixed yet, so stub paths would be premature. The representability ACs become Go unit tests that construct a legacy-shaped value and assert every field has a unified home; `backend-builds-and-tests-pass` is verified by the build/test commands themselves; `capability-coverage-complete` is verified against the coverage table as a checklist test.

## Not Doing

- **Live data migration / Firestore transform** — no legacy production data exists (greenfield); reconciliation is model/code only.
- **Liabilities & service-providers** — preserved as a capability but specified in a separate sibling Feature (a whole bounded context); tracked in the coverage table, not built here.
- **Frontend/Angular component & page port** — sequenced as a separate sibling Feature; tracked in the coverage table, not built here.
- **Deleting the legacy directories** — a separate follow-up once consumers are repointed.
- **Net-new capabilities beyond legacy ∪ MVP** — unification only.

## Assumption Carryover

- *Must-be-true* (from Idea): no live legacy data — **carried, unblocking** the data-migration cut.
- *Must-be-true*: every legacy capability is representable — **promoted to REQ:capability-coverage-total** and AC `capability-coverage-complete`.
- *Must-be-true*: single owning Space can coexist with multi-space association — **resolved** by REQ:canonical-owner-with-multispace (preserve both).
- *Should-be-true*: typed extras don't regress history/transfer/visibility — **promoted to REQ:ownership-core-preserved** and AC `history-transfer-intact`.
- *Open question* "liabilities home" — **answered for this Feature**: out of scope, sibling Feature; coverage table records it so nothing is lost.

## Open Questions

- **Financial dimension home (kept in core for now — revisit).** Per-asset/group totals (`ITotalsHolder`), `canHaveIncome`/`canHaveExpense`, income/expense direction, the `debt` category, and asset-side liability linkage currently live as **optional fields on the unified core asset** per REQ:financial-dimension-disposition (owner decision). This is retained as-is, but flagged as an open question: should the financial dimension eventually move out of the ownership core into the **liabilities/financial sibling Feature** (closer to Debtus), to keep the asset record focused on ownership? No functionality is at risk either way — this is a bounded-context placement decision to revisit when the liabilities sibling Feature is specified.
- **Tasks-module reminder integration (sibling Feature, not a cut).** The legacy vehicle NCT/tax/service due-date links pointed into a tasks module. This Feature preserves the **due-date values** on the `vehicle` extra (REQ:vehicle-extra-complete), so no data is lost; whether to re-wire them as live task reminders into a tasks module is deferred to a separate sibling Feature once that integration target is confirmed to exist.
- **`yearOfBuild` vs `dateOfBuild`.** Legacy carries both (with a TODO to consolidate). Preserve both, or consolidate to `dateOfBuild` on migration? Spec preserves both by default.

*(Resolved during review: `draft` lives in the unified `Status` enum per REQ:status-union; per-category `Type` validation is preserved per REQ:category-two-level and AC:out-of-category-type-rejected — these are no longer open.)*
