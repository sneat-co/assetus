# Legacy Frontend Assetus — Retirement Plan (handoff for a fresh session)

**Status:** Ready to execute · **Owner decision pending per step** · Tracks GitHub issue sneat-co/assetus#4
**Goal:** Migrate every consumer of the **legacy** frontend assetus packages onto the **new** unified `@sneat/extension-assetus` lib, then **delete** the legacy frontend assetus code — with zero functionality lost.

> **Backend is already done.** The legacy backend (`sneat-go-backend/pkg/extensions/assetus`) was deleted in sneat-go-backend#194; the monolith runs on the published standalone module `github.com/sneat-co/assetus/backend v0.1.0` via `assetusext.Extension()`. This plan is **frontend only**.

---

## 1. What "legacy frontend assetus" is

Two published libraries that live in **`sneat-libs/libs/extensions/assetus/`**:

| Package | Path | Contents |
|---|---|---|
| `@sneat/mod-assetus-core` | `sneat-libs/libs/extensions/assetus/core` | DTOs, enums, UI context wrappers, reference data (`carMakes`), `uimodels` (`AssetGroup`) |
| `@sneat/ext-assetus-components` | `sneat-libs/libs/extensions/assetus/components` | `AssetService`, the ~19 Angular components, shared bases |

Plus legacy **pages** in `sneat-apps/libs/extensions/assetus/pages/`.

The **new** replacement is **`@sneat/extension-assetus`** at `assetus/frontend/libs/ext-assetus` (the superset port from sneat-co/assetus#3).

## 2. Why it can't be deleted yet — the 5 consumers

Deleting the legacy packages now breaks the build of every importer below. Exact symbols each one uses:

| Consumer | Repo | Imports from legacy | Symbols used |
|---|---|---|---|
| **ext-assetus** (the NEW lib) | assetus | `mod-assetus-core` (17 files), `ext-assetus-components` (24 files) | `IAssetContext`, `IAssetDwellingContext`, `IAssetDocumentContext`, `IAssetBrief`, `AssetCategory`, `AssetType`, `IAssetDtoGroup`, `carMakes`/`IMake`/`IModel`, `AssetService`, legacy `IUpdateAssetRequest`, shared base components |
| **docus** | sneat-apps | `mod-assetus-core` (5), `ext-assetus-components` (7) | `IAssetDocumentContext`, `AssetService` |
| **budgetus** | sneat-apps | `mod-assetus-core` (1) | `AssetGroup` |
| **contactus/shared** | sneat-libs | `mod-assetus-core` (1), `ext-assetus-components` (2) | `IAssetContext`, `AssetService` |
| **contactus/internal** | sneat-libs | `mod-assetus-core` (1) | `IAssetContext` |

Only **ext-assetus** is from the unification (issue #4). The other four are **pre-existing, independent** dependencies — that is the real blast radius.

## 3. Prerequisite — the new lib must export the needed surface

Before any external consumer can migrate, `@sneat/extension-assetus` must (a) **export** every symbol the consumers need and (b) be **publishable/consumable** by sneat-apps + sneat-libs (it must be published to the registry the same way `@sneat/extension-listus`/others are, or wired via the workspace).

Audit the new lib's public `index.ts` and add/port whatever is missing. Known gaps to port from legacy `core`:
- **UI context wrappers**: `IAssetContext`, `IAssetDocumentContext`, `IAssetDwellingContext` (legacy `core/src/lib/contexts/`). The new lib has the DTOs but may not have these `*Context` wrappers — port them.
- **Reference data**: `carMakes`, `IMake`, `IModel` (legacy `core/src/lib/data/car-makes-with-models.ts`). Port into the new lib.
- **`AssetGroup`** uimodel (legacy `core/src/lib/uimodels/asset-group.ts`) — reconcile with the new `IAssetGroupInfo`.
- **`AssetService`** — the new lib already has its own `services/asset.service.ts` (8 routes, widened request DTOs). Consumers must move to it; confirm method/DTO parity with the legacy `AssetService` they used.
- Enums/DTOs (`AssetCategory`, `AssetType`, `IAssetBrief`) — already in the new lib; confirm names match.

## 4. Migration order (dependency-safe)

Execute top-to-bottom. Each step ends green (`nx run-many -t lint test build` for the touched projects) and is its own commit/PR per repo.

1. **[assetus repo] Complete the new lib's exports** (§3). Port the missing context wrappers + `carMakes` + `AssetGroup`; ensure `index.ts` exports everything the 5 consumers need. Verify `nx run-many -t lint test build --projects=ext-assetus` green. **Publish** a new version of `@sneat/extension-assetus`.
2. **[assetus repo] Self-decouple ext-assetus** (issue #4): replace all 17+24 imports of `@sneat/mod-assetus-core` / `@sneat/ext-assetus-components` inside `ext-assetus/src` with the lib's own unified types/services. Remove those two from `ext-assetus/package.json` peerDependencies (added in PR #3 only to satisfy lint). Verify green. **This removes the unification's own coupling.**
3. **[sneat-libs] Migrate `contactus/shared` + `contactus/internal`**: repoint `IAssetContext` + `AssetService` to `@sneat/extension-assetus`. Verify `nx run-many -t lint test build` for contactus + dependents. *(Do contactus before the sneat-apps consumers, since shared libs sit lower in the graph.)*
4. **[sneat-apps] Migrate `budgetus`**: repoint `AssetGroup`. Verify.
5. **[sneat-apps] Migrate `docus`**: repoint `IAssetDocumentContext` + `AssetService` (12 imports). Verify. *(Largest external consumer.)*
6. **[sneat-apps] Migrate the legacy `pages`** (`sneat-apps/libs/extensions/assetus/pages/`): these are superseded by the new lib's pages — confirm routing now uses the new lib's pages, then they can be removed with the legacy libs.
7. **[verify] No remaining importers**: `grep -rl '@sneat/mod-assetus-core\|@sneat/ext-assetus-components' sneat-apps/libs sneat-libs/libs assetus/frontend/libs | grep -v node_modules` returns **nothing** (outside the legacy dirs themselves).
8. **[delete]** Remove `sneat-libs/libs/extensions/assetus/{core,components}` and `sneat-apps/libs/extensions/assetus/pages`, plus their project.json/tsconfig path mappings + workspace references. Run full `nx run-many -t lint test build` across affected projects. Follow the relocation-notice PRs with deletion PRs.

## 5. Verification gate (every step)

- `nx run-many -t lint test build --projects=<touched>` exits 0.
- No new `@sneat/mod-assetus-core` / `@sneat/ext-assetus-components` imports introduced.
- Behaviour parity: the migrated consumer renders/works as before (component specs or manual check).
- After step 8: full affected-project build green AND the grep in step 7 is empty.

## 6. Risks / notes

- **Cross-repo publish coupling**: external consumers (sneat-apps, sneat-libs) can only migrate once `@sneat/extension-assetus` is published with the needed exports. Sequence the publish (step 1) first.
- **`AssetService` parity**: the new service has widened request DTOs and the 8 real routes; confirm docus/contactus usages map cleanly (esp. any method signatures they relied on).
- **Context wrappers** (`IAsset*Context`) are UI-model glue, not pure DTOs — port them carefully (they wrap the DTO + space context).
- This touches **three repos** (assetus, sneat-apps, sneat-libs) → expect **one PR per repo**, merged in the §4 order.
- Legacy `core` also feeds `contactus` — assetus and contactus had a bidirectional relationship; watch for circular-dependency surprises when repointing.

## 7. Definition of done

- All 5 consumers import only `@sneat/extension-assetus` (or no assetus at all).
- `sneat-libs/libs/extensions/assetus` and `sneat-apps/libs/extensions/assetus` are **deleted**.
- All affected projects build/lint/test green in CI.
- Issue sneat-co/assetus#4 closed.
