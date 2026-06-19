---
format: https://specscore.md/plan-specification
status: Implemented
---
# Plan: Unified Assetus Data Model — Port Plan

**Status:** Implemented
**Source Feature:** unified-assetus-data-model
**Date:** 2026-06-19
**Owner:** alex
**Supersedes:** —

## Summary

Port the legacy Assetus backend into the new MVP backend so the unified model is a strict superset of both, with **zero functionality lost**. This is a **migration, not a rewrite**: each task relocates and adapts existing legacy Go code from `sneat-go-backend/pkg/extensions/assetus/` onto the MVP ownership-core spine in `assetus/backend/`, preserving every legacy enum value, struct field, validation rule, child collection, facade, and route. Scope is backend only; the frontend Angular port and standalone liabilities provider/plan records are separate sibling Features.

## Approach

Bottom-up dependency order: enums → flat record → typed extras → relationships/multi-space/member-info → facades/API → coverage table + green build. Each task's mechanic is **copy-or-git-move the legacy package, then adapt it onto the MVP equivalent** (union the enums, add fields as optional, keep the MVP history/transfer/visibility core as the spine) — not author from scratch. Legacy validation (per-category `Type`, EngineType↔FuelType, per-doc-type schema) and the full vehicle-record/fuel child collection are carried over verbatim where practical. The financial dimension is folded onto the core record as optional fields per the owner decision. The final task proves nothing was lost (capability-coverage table) and that the merged module compiles and tests green.

Legacy source root: `sneat-go-backend/pkg/extensions/assetus/`. MVP merge target root: `assetus/backend/`.

## Tasks

### Task 1: Port & union the enums (`const4assetus`)

**Verifies:** unified-assetus-data-model#ac:category-superset, unified-assetus-data-model#ac:out-of-category-type-rejected, unified-assetus-data-model#ac:status-no-value-dropped
**Depends-On:** —
**Status:** done

Relocate legacy `const4assetus` (`AssetCategory`, `AssetStatus`, `AssetPossession`, per-category `AssetType`, `AssetDocumentType`, `EngineType`, `FuelType`, and the engine↔fuel compatibility validation) into `assetus/backend/const4assetus`, unioning with the MVP enums already there (`Condition`, `Visibility`, `Status` lifecycle states, `OwnerType`, `HistoryEventType`). Build the documented legacy→unified mapping table including `sport_gear`→`sports_equipment`, `vehicle`↔`vehicles`, `misc`→`other`, `debt` retained as a Category value, and `undefined`. Preserve per-category `Type` validation so out-of-category subtypes are rejected.

### Task 2: Port & merge the flat record superset (`dbo4assetus`, `briefs4assetus`)

**Verifies:** unified-assetus-data-model#ac:optional-legacy-fields-roundtrip, unified-assetus-data-model#ac:leasing-asset-representable, unified-assetus-data-model#ac:possession-defaults-owning, unified-assetus-data-model#ac:financial-fields-have-a-home
**Depends-On:** 1
**Status:** done

Adapt the legacy `AssetBrief`/`AssetBaseDbo`/`AssetDbo` (and `briefs4assetus`) fields onto the MVP flat `AssetBase`/`AssetDbo` in `assetus/backend/dbo4assetus`, keeping the MVP fields and adding the legacy ones as **optional**: `CountryID`, `AssetDates`, `yearOfBuild`, custom fields, tags, `geo`, `isRequest`, `parentCategoryID`, `Possession` (default `owning`), and the financial fields folded into core (`ITotalsHolder` totals, income/expense capability flags + direction, asset-side liability linkage `liabilities`/`notUsedServiceTypes`/`AssetLiabilityInfo`). Ensure a leased asset (possession independent of status) round-trips.

### Task 3: Port the polymorphic typed extras + vehicle records (`extras4assetus`)

**Verifies:** unified-assetus-data-model#ac:typed-extras-optional, unified-assetus-data-model#ac:vehicle-extra-no-field-dropped, unified-assetus-data-model#ac:document-extra-full-shape
**Depends-On:** 2
**Status:** done

Relocate legacy `extras4assetus` onto the flat core in `assetus/backend`: `WithExtraField` registration by `extraType`; `AssetVehicleExtra` (make/model/regNumber/VIN, engine incl. `engineSerialNumber`, EngineType↔FuelType validation); the full vehicle-record child collection (`VehicleRecordDbo` = `VehicleMileage{value,unit}` + `VehicleFuelRecord{volume,unit,amount}` + `fuelCost`/`currency` + `CreatedFields`); `AssetDwellingExtra` (address/rent/bedrooms/area); `AssetDocumentExtra` full shape (`docType`/`number`/`batchNumber`/`countryID`/`issuedBy`/`issuedOn`/`effectiveFrom`/`expiresOn`) with the `standardDocTypesByID` per-doc-type validation schema. Carry over the legacy indexed-field declarations (vehicle make/model/regNumber/vin; document expiresOn/effectiveFrom) per REQ:query-index-fields. An asset with no extra stays valid.

### Task 4: Port relationships, multi-space & member info onto the record

**Verifies:** unified-assetus-data-model#ac:groups-nesting-linking-preserved, unified-assetus-data-model#ac:multispace-with-canonical-owner, unified-assetus-data-model#ac:member-info-preserved
**Depends-On:** 3
**Status:** done

Adapt the legacy relationship structures onto the unified record: `WithAssetSpaces` multi-space association preserved alongside the MVP single canonical owning Space (lifecycle/history anchored to the owner); asset groups as a sub-entity (`IAssetDtoGroup` fields `order`/`desc`/`categoryId`/`numberOf`, not just `groupId`); parent/sub-asset nesting with per-sub-asset detail (`ISubAssetInfo`: `type`/`countryId`/`subType`/`expires`); asset linking (`sameAssetID`/`relatedAs`/`linkage`); and per-asset member info (`memberIDs`/`membersInfo`).

### Task 5: Port & merge facades + API, preserving the MVP ownership core (`facade4assetus`, `api4assetus`, `dal4assetus`)

**Verifies:** unified-assetus-data-model#ac:history-transfer-intact, unified-assetus-data-model#ac:condition-optional-visibility-default
**Depends-On:** 4
**Status:** done

Merge the legacy facades/routes (`CreateAsset`/`GetAsset`/`UpdateAsset`/`DeleteAsset`/`AddVehicleRecord` + HTTP endpoints) into the MVP `facade4assetus`/`api4assetus`, keeping the MVP ownership core intact: append-only history, Space→Space transfer (relocates asset + history, appends `Transferred`), owner-type derivation, soft-archive-vs-hard-delete, and create-time visibility defaulting to the owning Space's default. Adapt `AddVehicleRecord` to the fuel-bearing request payload (`fuelVolume`/`fuelVolumeUnit`/`fuelCost`/`currency`/`mileage`/`mileageUnit`) and reconcile the `dal4assetus` paths/child collections.

### Task 6: Capability-coverage table + green build + retirement note

**Verifies:** unified-assetus-data-model#ac:capability-coverage-complete, unified-assetus-data-model#ac:backend-builds-and-tests-pass
**Depends-On:** 5
**Status:** done

Author the capability-coverage table mapping every legacy capability in `docs/legacy-gap-analysis.md` §3.1/§4 (corrected for the audit findings) to its unified home — concrete field, intentional change, or named sibling Feature. Make the merged `assetus/backend` module compile and pass tests (`go build ./...` and `go test ./...` from `backend/`), and write the retirement note listing which legacy directories (and consumers to repoint: `brandus/dbo4brands/make_test.go`, `standard_extensions.go`, `sneat-apps`) become deletable.

### Task 7: Backend test coverage (≥80% + a test per ported capability)

**Verifies:** unified-assetus-data-model#ac:backend-test-coverage
**Depends-On:** 6
**Status:** done

Add/port unit tests so every ported capability is exercised — each unioned enum and its mapping, each optional-field round-trip, each typed extra (including the vehicle fuel-record), the relationship/multi-space/member-info structures, and the facades (create/get/update/remove/transfer/add-vehicle-record). Carry over any reusable legacy tests from `sneat-go-backend/pkg/extensions/assetus/` rather than re-authoring. Verify overall backend statement coverage is ≥80% via `go test -cover ./...` from `backend/`. (Frontend test coverage is the equivalent requirement for the separate frontend-port sibling Feature — see Open Questions.)

## Open Questions

- **Frontend test coverage (sibling Feature).** REQ:test-coverage requires an equivalent ≥80% coverage bar for the migrated Angular DTOs/components, but the frontend port is a separate sibling Feature out of scope for this backend plan. This requirement is recorded so it is carried into that Feature when it is specified; it is not implemented here.
- **`yearOfBuild` vs `dateOfBuild`.** Legacy carries both (with a consolidation TODO); the port preserves both unless the owner opts to consolidate during Task 2.
- **Tasks-module reminder integration.** Vehicle service/tax/inspection due-dates are preserved as plain dates (Task 3); re-wiring them as live task reminders is a separate sibling Feature, not in this plan.

---
*This document follows the https://specscore.md/plan-specification*
