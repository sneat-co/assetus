---
format: https://specscore.md/plan-specification
status: Approved
---
# Plan: Extension Library Architecture

**Status:** Approved
**Source:** idea:adopt-extension-library-architecture
**Date:** 2026-06-20
**Owner:** alexandertrakhimenok
**Supersedes:** —

## Summary

Splits the monolithic `ext-assetus` lib into the contract/shared/internal convention (defined in sneat-libs `spec/features/extension-library-architecture`): `@sneat/extension-assetus-contract` (IAsset types/DTOs/contexts/enums + `ASSET_SERVICE` token, runtime-light), `@sneat/extension-assetus-shared` (app-facing routing + pages + components, injecting the token), `@sneat/extension-assetus-internal` (AssetService + services module, providing the token). Adds nx tier tags + an `enforce-module-boundaries` tier matrix. Scoped to the assetus repo; the cross-repo Part B (publish `-contract`, rewire sneat-libs' type-only imports) is deferred.

## Approach

Same hard-cutover model proven on contactus/calendarius: move, never copy; one definition per symbol; old lib dropped as it empties; tier tasks sequential. Sequence: infra + scaffold three libs → contract cutover (types + `ASSET_SERVICE` token) → internal cutover (AssetService + provider factory) → shared cutover (routing/pages/components, inject token) → reroute the app + flip enforcement strict + verify → final sweep. The assetus app consumes the extension via one lazy-routing import, so it is rerouted onto `@sneat/extension-assetus-shared`.

**Repo-specific notes:** the assetus repo already has `@nx/enforce-module-boundaries` at `error` with `type:lib`/`scope:assetus` tags; Task 1 introduces the `type:contract`/`type:shared`/`type:internal` taxonomy + tier matrix (initially `warn` during migration, flipped to `error` in Task 5). `eslint.config.mjs` (ESM) is the config file. The shared tier must never import internal; the app-facing routing lives in shared and obtains `AssetService` via the contract token. `-internal` is omitted from `tsconfig.base.json` paths.

**Execution model (subagent-driven):** each task runs as a dispatched subagent; tasks editing shared workspace files (`tsconfig.base.json`, `eslint.config.mjs`) are serialized. Commit per task. The final task is a reconciling sweep.

**Deferred (Part B, cross-repo, publish-gated):** publish `@sneat/extension-assetus-contract` and rewire sneat-libs' three type-only `@sneat/extension-assetus` imports (`IAssetContext`/`IAssetBrief`/`IAssetDbo`) onto it. Out of scope here.

## Tasks

### Task 1: Convention infra + scaffold three libs

**Source:** idea:adopt-extension-library-architecture
**Depends-On:** —
**Status:** done

Add the `type:contract`/`type:shared`/`type:internal` tag taxonomy and a tier dependency matrix to `eslint.config.mjs`'s `enforce-module-boundaries` (initially `warn`). Scaffold empty `@sneat/extension-assetus-contract`, `@sneat/extension-assetus-shared`, and `@sneat/extension-assetus-internal` libs (mirroring `ext-assetus` config), with tier + `scope:assetus` tags. Add `paths` for `-contract` and `-shared` only; omit `-internal`.

**Notes:** Edits `tsconfig.base.json` + `eslint.config.mjs` — run solo.

### Task 2: Contract cutover — types + ASSET_SERVICE token

**Source:** idea:adopt-extension-library-architecture
**Depends-On:** 1
**Status:** done

Move pure types from `ext-assetus` (`dto`, `contexts`, `uimodels`, `data`, `constants`, `services/interfaces.ts`) into `extension-assetus-contract`, runtime-light. Define `ASSET_SERVICE` + `IAssetService` (interface derived from how the shared UI uses `AssetService`). Repoint references; delete moved symbols from `ext-assetus`.

**Notes:** Broad type-import repoint; run solo.

### Task 3: Internal cutover — AssetService + provider factory

**Source:** idea:adopt-extension-library-architecture
**Depends-On:** 2
**Status:** done

Move `AssetService` and `assetus-core-services.module.ts` into `extension-assetus-internal`; bind `ASSET_SERVICE` to `AssetService` via `provideAssetusInternal(): Provider[]`. No other lib imports `-internal`. Drop the emptied services dir from `ext-assetus`.

**Notes:** Run after Task 2.

### Task 4: Shared cutover — routing, pages, components

**Source:** idea:adopt-extension-library-architecture
**Depends-On:** 3
**Status:** pending

Move the app-facing routing, pages (9 page groups), the ~28 components, and `space-menu` into `extension-assetus-shared`, refactoring `AssetService` access to inject `ASSET_SERVICE` (zero `-internal` imports). Drop the now-empty `ext-assetus` lib (or repurpose it) once content has migrated.

**Notes:** Largest task; run after Task 3.

### Task 5: Reroute app + flip enforcement strict + verify

**Source:** idea:adopt-extension-library-architecture
**Depends-On:** 4
**Status:** pending

Reroute the assetus app's `@sneat/extension-assetus` import(s) onto `@sneat/extension-assetus-shared`, and wire `provideAssetusInternal()` at the app bootstrap. Flip the tier matrix in `eslint.config.mjs` from `warn` to `error`. Run full CI (build, lint, test) green; confirm a deliberate forbidden edge (shared → internal) fails lint.

**Notes:** Edits `eslint.config.mjs` + app — run solo.

### Task 6: Final lib removal + sweep

**Source:** idea:adopt-extension-library-architecture
**Depends-On:** 5
**Status:** pending

Confirm `ext-assetus` and its `project.json`/`tsconfig` `paths` entry are gone (remove residue) and the three new libs build/lint/test clean. Confirm no dangling `@sneat/extension-assetus` (old) imports remain in-repo.

## Open Questions

- Does the assetus app render extension pages via lazy routing only (one import), or are there deeper direct component imports? Confirm during Task 5 — affects how much of the UI must be in `shared` vs can stay `internal`.
- Does `ext-assetus` import `@sneat/contactus-services` (old published) anywhere that should instead be a contract token? It is one import; route it via a contract token if it's a runtime call, or leave as a published-package dep if type-only.

---
*This document follows the https://specscore.md/plan-specification*
