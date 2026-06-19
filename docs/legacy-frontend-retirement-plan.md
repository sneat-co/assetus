# Legacy Frontend Assetus — Retirement Plan (multi-agent handoff)

**Status:** Ready to execute · Tracks GitHub issue sneat-co/assetus#4
**Goal:** Migrate every consumer of the **legacy** frontend assetus packages onto the **new** unified `@sneat/extension-assetus` lib, then **delete** the legacy frontend assetus code — zero functionality lost.

> **Backend is already done.** Legacy backend (`sneat-go-backend/pkg/extensions/assetus`) deleted in sneat-go-backend#194; monolith runs on published `github.com/sneat-co/assetus/backend v0.1.0` via `assetusext.Extension()`. **This plan is frontend only.**

---

## 0. Non-negotiable requirements (apply to EVERY phase, stream, and PR)

These are **hard gates**. A stream is not done — and its PR must not merge — unless all hold. Do not relax them; surface a blocker instead.

1. **ZERO FUNCTIONALITY LOST.** Every capability, field, behaviour, and rendered output of each migrated consumer MUST be preserved exactly. Migration is a *repoint of imports*, not a feature change: do not drop, rename, or simplify any behaviour. For each migrated symbol (`IAssetContext`, `IAssetDocumentContext`, `AssetService`, `AssetGroup`, etc.), verify the new lib's equivalent is a **superset or exact match** of what the consumer used; if the new lib is missing anything the consumer relied on, **add it to the new lib (Phase 1)** rather than degrade the consumer. Any unavoidable behaviour change is a blocker for human decision, not a silent cut.
2. **GOOD TEST COVERAGE — no regression, ≥80%.** Every touched project MUST keep (or improve) its statement coverage, and MUST NOT drop below **80%**. `@sneat/extension-assetus` (ext-assetus) is at ~82% — keep it ≥80%. Any code newly **ported into the new lib in Phase 1** (context wrappers, `carMakes`, `AssetGroup`, service surface) MUST ship with tests. Migrated consumers MUST retain their existing tests passing; where a consumer had behaviour tests tied to the legacy import, port/adapt them — do not delete tests to make the build green.
3. **Evidence, not assertion.** Each PR must show: `nx run-many -t lint test build` green for the touched projects **with `--coverage`**, the measured coverage number(s), and a one-line behaviour-parity note per migrated consumer. "Build passes" alone is insufficient — coverage and parity must be demonstrated.
4. **Verify against the actual legacy code**, not assumptions — read the legacy symbol definitions in `sneat-libs/libs/extensions/assetus` and confirm the new lib matches, exactly as the original unification was reviewed.

---

## 0.5 Concurrency & git safety (HARD rules — parallel streams MUST NOT clobber each other)

Streams/subagents run concurrently and may share a filesystem/checkout. A stray `git switch` or a blanket `git add` in one stream can corrupt another's work, and an unrelated session can switch a branch underneath an agent. Enforce ALL of:

1. **Dedicated branch per stream, created at start, recorded.** Each stream creates and records its own branch off `origin/main`:
   `git switch -c retire-legacy-assetus/<repo> origin/main` (e.g. `retire-legacy-assetus/sneat-apps`).
   The orchestrator keeps the `{repo → branch}` map and passes each subagent its **exact expected branch name**.

2. **Do NOT use git worktrees here — deliberate decision.** These are **nx/pnpm** workspaces: a fresh worktree has no `node_modules` (gitignored) and would force a `pnpm install` per worktree (slow, even with pnpm's store). And it buys nothing, because **this plan parallelizes BY REPO** — Streams A/B/C are already in *three different repo checkouts* (`assetus`, `sneat-libs`, `sneat-apps`), each with its own warm `node_modules` and independent branch. There is no shared working tree across streams to isolate. So: **one existing checkout per repo, one branch per repo (rule 1), guarded by rule 4.** Reuse the warm `node_modules`; do not reinstall.
   - **The only same-repo overlap** is Phase 1 (publish, `assetus`) vs Stream A (self-decouple, `assetus`). These are **naturally serial** — Stream A can't repoint to the new lib until Phase 1 has published it — so **serialize them in the one `assetus` checkout**. Do NOT run them concurrently.
   - **Worktree is a last resort only**, if you truly must have two concurrent writers in one repo: `git worktree add ../wt-retire-<repo> -b retire-legacy-assetus/<repo> origin/main` then `pnpm install` inside it (budget the reinstall). Prefer serializing over this.

3. **Track edited files explicitly.** Maintain the exact list of paths this stream/subagent edits. **NEVER** stage with `git add -A`, `git add .`, or `git add -u`.

4. **Pre-commit guard — run before EVERY commit (orchestrator and every subagent):**
   ```bash
   # 1. Assert we are on our own branch (abort if a parallel session switched it)
   test "$(git rev-parse --abbrev-ref HEAD)" = "<expected-branch>" || { echo "WRONG BRANCH — abort"; exit 1; }
   # 2. Stage ONLY this stream's tracked files, by explicit path:
   git add <file1> <file2> ...        # never -A / . / -u
   # 3. Confirm nothing leaked in from another stream:
   git diff --cached --name-only | sort > /tmp/staged.$$
   printf '%s\n' <file1> <file2> ... | sort > /tmp/expected.$$
   diff /tmp/staged.$$ /tmp/expected.$$ || { echo "UNEXPECTED STAGED FILES — abort"; git restore --staged .; exit 1; }
   # 4. Only now commit
   git commit -m "..."
   ```

5. **Never switch branches mid-work** in a shared (non-worktree) checkout, and **never force-push another stream's branch**. If you must read another branch, use `git show <branch>:<path>` — do not check it out.

6. **One writer per repo at a time** (the §1 model). The orchestrator MUST NOT dispatch two agents that write to the same repo concurrently — **serialize them** (preferred), or, only if unavoidable, worktree-isolate one (rule 2, last resort). Different repos run freely in parallel since they are separate checkouts.

---

## 1. Execution model (read first)

This is a **pipeline with barriers**, not a flat fan-out. There is a serial gate at the front and a serial barrier at the end; only the middle parallelizes.

```
PHASE 1  (SERIAL, blocking) ── single agent ──────────────────────────────
   Port the new lib's missing exports + PUBLISH @sneat/extension-assetus.
   Nothing downstream can start until the new lib exposes the needed
   symbols AND is published/consumable by sneat-apps + sneat-libs.
        │
        ▼
PHASE 2  (PARALLEL fan-out) ── 3 agents, one per repo ───────────────────
   A: assetus     — self-decouple ext-assetus (issue #4)
   B: sneat-libs  — contactus/shared + contactus/internal
   C: sneat-apps  — docus + budgetus
   Streams are independent (separate repos / nx workspaces). Within a
   stream, consumers run sequentially to avoid intra-workspace nx-cache
   contention. Each stream ends with its own green build + PR.
        │
        ▼
PHASE 3  (SERIAL barrier) ── single agent ───────────────────────────────
   Confirm zero remaining importers → delete legacy dirs → full build →
   deletion PRs → close #4.
```

**Why by repo, not by consumer:** two agents running `nx build/test` in the *same* workspace contend on the nx cache/daemon. Repos are separate nx workspaces, so the safe parallel unit is the repo. `docus`+`budgetus` share the sneat-apps workspace → one agent does both, sequentially.

**Why Phase 1 can't be parallelized in:** every consumer migration repoints imports to `@sneat/extension-assetus`; that target must already export the symbols and be installable. So publish is a hard prerequisite barrier.

**Workflow shape (if driven via the Workflow tool):**
`phase('publish') → agent(A1)` then `phase('migrate') → parallel([B, C])` (A's self-decouple can run in the publish phase or as a third parallel stream once exports exist) `→ phase('delete') → agent(barrier)`.

---

## 2. What "legacy frontend assetus" is

Two published libraries in **`sneat-libs/libs/extensions/assetus/`**:

| Package | Path | Contents |
|---|---|---|
| `@sneat/mod-assetus-core` | `…/assetus/core` | DTOs, enums, UI context wrappers, `carMakes` data, `uimodels` (`AssetGroup`) |
| `@sneat/ext-assetus-components` | `…/assetus/components` | `AssetService`, ~19 Angular components, shared bases |

Plus legacy **pages** in `sneat-apps/libs/extensions/assetus/pages/`.
The **new** replacement is **`@sneat/extension-assetus`** at `assetus/frontend/libs/ext-assetus`.

## 3. The 5 consumers (the blast radius) — exact symbols

| Stream | Consumer | Repo | Imports | Symbols |
|---|---|---|---|---|
| A | **ext-assetus** (new lib) | assetus | core ×17, components ×24 | `IAssetContext`, `IAssetDwellingContext`, `IAssetDocumentContext`, `IAssetBrief`, `AssetCategory`, `AssetType`, `IAssetDtoGroup`, `carMakes`/`IMake`/`IModel`, `AssetService`, legacy `IUpdateAssetRequest`, base components |
| B | **contactus/shared** | sneat-libs | core ×1, components ×2 | `IAssetContext`, `AssetService` |
| B | **contactus/internal** | sneat-libs | core ×1 | `IAssetContext` |
| C | **docus** | sneat-apps | core ×5, components ×7 | `IAssetDocumentContext`, `AssetService` |
| C | **budgetus** | sneat-apps | core ×1 | `AssetGroup` |

Only A is from the unification (#4); B and C are pre-existing, independent.

## 4. PHASE 1 — Port exports + publish (SERIAL prerequisite)

Audit `@sneat/extension-assetus`'s public `index.ts`; add/port every symbol §3 lists. Known gaps to port from legacy `core`:
- **Context wrappers** `IAssetContext` / `IAssetDocumentContext` / `IAssetDwellingContext` (`core/src/lib/contexts/`) — UI glue wrapping DTO + space context.
- **`carMakes` / `IMake` / `IModel`** (`core/src/lib/data/car-makes-with-models.ts`).
- **`AssetGroup`** uimodel (`core/src/lib/uimodels/asset-group.ts`) — reconcile with new `IAssetGroupInfo`.
- **`AssetService`** — new lib already has its own (8 routes, widened DTOs); confirm method/DTO parity with the legacy service docus/contactus used.
- Enums/DTOs (`AssetCategory`, `AssetType`, `IAssetBrief`) — present; confirm names match.

Gate: `nx run-many -t lint test build --projects=ext-assetus` green, then **publish** a new `@sneat/extension-assetus` version (same release mechanism as other `@sneat/extension-*` libs) so sneat-apps + sneat-libs can install it.

## 5. PHASE 2 — Per-repo migration (PARALLEL, 3 agents)

Each stream: repoint imports → remove legacy deps → green build → one PR. Run B and C concurrently; A can run alongside once Phase 1 exports exist.

### Stream A — assetus repo (self-decouple ext-assetus, issue #4)
Replace all 41 imports of `@sneat/mod-assetus-core` / `@sneat/ext-assetus-components` inside `assetus/frontend/libs/ext-assetus/src` with the lib's own unified types/services. Remove both from `ext-assetus/package.json` peerDependencies (added in PR #3 only for lint). Verify `nx run-many -t lint test build --projects=ext-assetus assetus-app` green.

### Stream B — sneat-libs (contactus)
Repoint `IAssetContext` (shared+internal) and `AssetService` (shared) to `@sneat/extension-assetus`. **Watch the assetus↔contactus circular-dependency risk** — assetus core historically imported contactus and vice versa; if repointing creates a cycle, break it by depending only on DTO-level exports, not components. Verify `nx run-many -t lint test build` for contactus + dependents.

### Stream C — sneat-apps (docus + budgetus, sequential within the agent)
1. `budgetus`: repoint `AssetGroup`. Verify.
2. `docus`: repoint `IAssetDocumentContext` + `AssetService` (12 imports). Verify.
Run `nx run-many -t lint test build --projects=docus budgetus` (+ dependents) green.

## 6. PHASE 3 — Delete (SERIAL barrier — only after all of Phase 2 merged)

1. **Importer check:** `grep -rl '@sneat/mod-assetus-core\|@sneat/ext-assetus-components' sneat-apps/libs sneat-libs/libs assetus/frontend/libs | grep -v node_modules` returns **nothing** outside the legacy dirs themselves.
2. **Migrate/retire legacy pages** `sneat-apps/libs/extensions/assetus/pages/` (superseded by the new lib's pages; confirm routing uses the new pages).
3. **Delete** `sneat-libs/libs/extensions/assetus/{core,components}` and `sneat-apps/libs/extensions/assetus`, plus their `project.json`/`tsconfig` path mappings and workspace references.
4. Full `nx run-many -t lint test build` across affected projects green.
5. Deletion PRs (one per repo); close issue #4.

## 7. Ready-to-use subagent prompts

> Each Phase-2 agent works in ONE repo only, on its **own dedicated branch `retire-legacy-assetus/<repo>` created off `origin/main`** in that repo's **existing checkout (NOT a worktree — §0.5 rule 2; reuse the warm node_modules)**, stages **only its tracked files by explicit path** through the §0.5 pre-commit guard, opens its own PR, and must end with the repo's `nx run-many -t lint test build --coverage` green for the touched projects. Do not let two agents run nx — or write — in the same workspace concurrently. **§0 + §0.5 apply to every prompt below: zero functionality lost, ≥80% / no-coverage-regression, behaviour-parity + coverage evidence in the PR, dedicated per-repo branch (no worktree — reuse warm node_modules), explicit-path staging, and the pre-commit branch+staged-files guard.** Each prompt below implicitly carries these — repeat them to the subagent, including its exact expected branch name.

**Phase 1 (publish) — assetus:**
> "In /Users/.../assetus/frontend, port the missing exports into `@sneat/extension-assetus` (libs/ext-assetus): context wrappers IAssetContext/IAssetDocumentContext/IAssetDwellingContext (from legacy sneat-libs/.../assetus/core/src/lib/contexts), carMakes/IMake/IModel (…/data/car-makes-with-models.ts), AssetGroup uimodel (…/uimodels/asset-group.ts) — reconcile with the new IAssetGroupInfo. Export everything the 5 consumers in docs/legacy-frontend-retirement-plan.md §3 need from the lib's index.ts. Verify `npx nx run-many -t lint test build --projects=ext-assetus` is green. Then publish a new lib version per the repo's release process. Stage; open a PR. Report the new version + the exported symbol list."

**Stream A — assetus (self-decouple):**
> "In assetus/frontend/libs/ext-assetus/src, replace every import from `@sneat/mod-assetus-core` and `@sneat/ext-assetus-components` with the lib's own unified types/services (dto/asset.ts, dto/extras.ts, services/asset.service.ts). Remove both packages from ext-assetus/package.json peerDependencies. `npx nx run-many -t lint test build --projects=ext-assetus assetus-app` must be green. Stage; open a PR closing issue #4's ext-assetus part."

**Stream B — sneat-libs (contactus):**
> "In sneat-libs, repoint contactus/shared and contactus/internal off `@sneat/mod-assetus-core` (IAssetContext) and `@sneat/ext-assetus-components` (AssetService) onto `@sneat/extension-assetus`. Beware an assetus↔contactus dependency cycle — depend on DTO-level exports only, not components. `npx nx run-many -t lint test build` for contactus + dependents must be green. Stage; open a PR."

**Stream C — sneat-apps (docus + budgetus):**
> "In sneat-apps, repoint two extensions off legacy assetus onto `@sneat/extension-assetus`, sequentially: (1) budgetus — AssetGroup; (2) docus — IAssetDocumentContext + AssetService (12 imports). `npx nx run-many -t lint test build --projects=docus budgetus` (+ dependents) must be green. Stage; open a PR."

**Phase 3 — delete (barrier):**
> "Only after the Phase-2 PRs are merged. Run the importer grep in docs/legacy-frontend-retirement-plan.md §6.1 — it must be empty outside the legacy dirs. Retire the legacy sneat-apps assetus pages, then delete sneat-libs/libs/extensions/assetus/{core,components} and sneat-apps/libs/extensions/assetus plus their project.json/tsconfig/workspace references. Full `nx run-many -t lint test build` green. Open deletion PRs per repo; close #4."

## 8. Definition of done

- **Zero functionality lost** — every migrated consumer behaves exactly as before (§0.1); any behaviour change was surfaced and approved, not silently cut.
- **Coverage held** — every touched project ≥80% statement coverage with no regression; newly-ported lib code ships with tests (§0.2).
- All 5 consumers import only `@sneat/extension-assetus` (or no assetus).
- `sneat-libs/libs/extensions/assetus` and `sneat-apps/libs/extensions/assetus` **deleted**.
- All affected projects build/lint/test green in CI, with coverage evidence in each PR.
- Issue sneat-co/assetus#4 closed.
