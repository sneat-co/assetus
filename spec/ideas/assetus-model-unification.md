---
format: https://specscore.md/idea-specification
status: Specified
---

# Idea: Assetus Model Unification: Reconcile Legacy & MVP

**Status:** Specified
**Date:** 2026-06-18
**Owner:** alex
**Promotes To:** assetus-frontend-port, unified-assetus-data-model
**Supersedes:** —
**Related Ideas:** —

## Problem Statement

How might we reconcile the legacy rich-registry Assetus and the new MVP ownership-core into a single superset model — resolving each enum/schema conflict one by one — so the legacy code can be retired with zero functionality lost?

## Context

Both Assetus implementations were authored by the same owner: the legacy rich asset-registry (vehicles/dwellings/sport-gear/documents with typed extras, possession, liabilities) and the newer MVP ownership-core (condition/status/visibility/append-only history/Space-to-Space transfer). docs/legacy-gap-analysis.md framed them as 'two different products that happen to share a name' and recommended building the MVP clean while deferring legacy capabilities behind RELOCATED.md notices. The owner disputes that framing: the two were specified for slightly different use cases, not as incompatible products. There is NO live legacy data (greenfield), so no Firestore transform is required; this is a pure model/code reconciliation. Driver: honor 'no functionality lost' before any legacy directories are deleted.

## Recommended Direction

Treat unification as a conflict-by-conflict reconciliation that produces a SUPERSET model (legacy union MVP), not a rebuild and not a cautious defer-everything port. Keep the MVP ownership-core as the canonical spine (condition, ownership-lifecycle status, visibility, append-only history, Space-to-Space transfer, owner-type derivation). Fold legacy richness back in as additive/optional structure: restore the two-level Category+Type taxonomy (union the value sets; sport_gear maps to sports_equipment); union the Status enum (add draft); reintroduce optional Possession (default owning); reintroduce polymorphic typed extras (vehicle/dwelling/document) hanging off the flat core; fold CountryID, AssetDates, custom fields, the `geo` hook, and asset-to-asset linking (`sameAssetID`/`relatedAs`/`linkage`) back as optional; preserve asset groups (`groupId`/`IAssetDtoGroup`), sub-assets / parent-asset nesting (`parentAssetID`/`subAssets`), and member-level info (`memberIDs`/`membersInfo`); preserve multi-space association (`WithAssetSpaces`) so an asset can still surface under several spaces, while ADDING a single canonical owning Space for the ownership-lifecycle/history/transfer model to anchor to; keep Liabilities/service-providers in scope but as an isolated sub-module. Each enum conflict resolves by union or mapping, documented in a conflict-resolution table. Outcome: one Assetus that is a strict superset of both, after which the legacy directories can be deleted with zero functionality lost.

## Alternatives Considered

- **Keep them separate, defer legacy (the gap-analysis recommendation).** Build the
  MVP clean and only port "genuinely asset-shaped" extras later, routing
  liabilities elsewhere. **Rejected:** it indefinitely strands legacy capabilities
  behind `RELOCATED.md` and defers the "no functionality lost" obligation rather
  than discharging it. The owner authored both models for adjacent use cases — the
  "two different products" premise is overstated.
- **Schema-merge the two record models directly (revive legacy enums as-is).**
  Reconcile by widening the legacy `AssetDbo`/enums in place. **Rejected:** the
  legacy `Status` (`draft`-centric) and single-level `Category` are *narrower and
  differently-shaped* than the MVP's lifecycle/visibility/history model; merging
  onto the legacy base loses the MVP's ownership-core gains. The MVP spine is the
  better canonical base.
- **Graft MVP features onto the legacy codebase (legacy as base).** Add condition/
  visibility/history/transfer onto legacy. **Rejected:** legacy carries multi-space
  `WithAssetSpaces`, required `CountryID`, and `linkage`/`relatedAs` hooks that the
  MVP deliberately simplified away; starting from legacy reimports complexity the
  MVP already resolved. Superset-on-MVP-spine keeps the simpler core and adds only
  optional richness.

## MVP Scope

A unified Assetus data-model specification plus the reconciled backend Go packages (const4assetus enums + dbo4assetus record + extras4assetus typed extras), driven by a conflict-resolution table that addresses EVERY conflict in docs/legacy-gap-analysis.md section 3.2 and EVERY capability in sections 3.1 and 4 (including asset groups, sub-assets/parent nesting, asset linking, member-level info, multi-space association, possession, typed extras, AssetDates, custom fields, CountryID, and liabilities) — each maps to a unified field or is explicitly preserved. Success = each legacy capability is representable in the unified model, the unified backend builds and tests pass (go build ./... and go test ./...), and a written retirement note confirms which legacy directories become deletable. Frontend/Angular component port is sequenced as a follow-on, not lost.

## Not Doing (and Why)

- Live data migration / Firestore transform — no legacy production data exists (greenfield); reconciliation is model/code only
- Changing legacy behavior — multi-space association, asset groups, sub-assets, linking, and member-info are all preserved as-is; this Idea adds a canonical owning Space alongside them, it does not remove the ability for an asset to appear under several spaces
- Frontend/Angular component and page port — the legacy components and pages are sequenced as a follow-on Feature after the model+backend reconciles, not dropped
- Net-new capabilities beyond legacy union MVP — this is unification only; no features outside the union of the two existing models
- Deleting the legacy directories in this Idea — deletion is a separate follow-up once consumers (brandus make_test, standard_extensions registration, sneat-apps) are repointed

## Key Assumptions to Validate

| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true | There is no live legacy Assetus data in production (sneat-eur3-1) — so no data transform is in scope. | Owner confirmed (code-only). Spot-check the legacy Firestore `assets` collection is empty/unused before any deletion follow-up. |
| Must-be-true | Every legacy capability in `legacy-gap-analysis.md` §3.2/§4 has a representation in the unified model (no conflict is truly irreconcilable). | Walk the conflict-resolution table; each row must map to a unified field/enum value. Any row with no resolution is a blocker surfaced to the owner. |
| Must-be-true | The MVP ownership-core (single owning Space + history + transfer) can absorb legacy multi-space "briefing" via visibility/sharing without losing intent. | Design review of `WithAssetSpaces` use sites: confirm each was a *sharing/visibility* need, not genuine co-ownership. |
| Should-be-true | Reintroducing polymorphic typed `extra` onto the flat `AssetBase` does not regress the MVP's history/transfer/visibility behavior. | Unified backend builds; existing MVP facade tests still pass after extras are added. |
| Should-be-true | Liabilities/service-providers can live as an isolated Assetus sub-module without entangling the ownership core. | Module-boundary review: liabilities reference assets by ID only, no reverse coupling into `AssetDbo`. |
| Might-be-true | Frontend Angular components can be ported largely as-is once the unified DTOs land. | Defer until the model+backend Feature is done; re-evaluate against `@sneat/ext-assetus-components`. |


## SpecScore Integration

- **New Features this would create:**
  - *Unified Assetus data model* — the reconciled enums (`const4assetus`), record
    (`dbo4assetus`), and polymorphic typed extras (`extras4assetus`), carrying the
    conflict-resolution table as acceptance criteria.
  - *Liabilities & service-providers sub-module* (candidate split) — isolated
    obligation/provider/plan domain referencing assets by ID.
  - *Legacy frontend port* (follow-on) — migrate the ~19 Angular components and
    pages onto the unified DTOs.
- **Existing Features affected:** The current MVP Assetus Feature (its *Not Doing*
  section is the inverse of this Idea's scope — this Idea reverses several of those
  deferrals). Downstream consumers to repoint before legacy deletion:
  `brandus/dbo4brands/make_test.go`, `standard_extensions.go`, and `sneat-apps`.
- **Dependencies:** None blocking. Sequencing only — the unified model Feature must
  land before the frontend-port and legacy-deletion follow-ups.

## Open Questions

- **Liabilities home.** Does the liabilities/service-providers sub-domain belong
  inside Assetus at all, or in a separate `servius`/Debtus module? In MVP scope per
  the owner's "no functionality lost," but flagged as a candidate to split.
- **Unified `Category` shape.** Single flat union enum, or the restored two-level
  `Category` + `Type`? Leaning two-level (legacy already had it; richer), but this
  is the main taxonomy decision for the design Feature.
- **Possession vs ownership semantics.** With `Possession` reintroduced
  (own/lease/rent), how does a `leasing`/`renting` asset interact with
  ownership-lifecycle `Status` and Space→Space transfer? Needs an explicit rule.
- **`draft` Status placement.** Is `draft` a `Status` value (union approach) or a
  separate edit-state flag, given the MVP treats Status as ownership-lifecycle?
- **Canonical owner vs multi-space.** Legacy let one asset be briefed under several
  spaces with no single owner; the MVP's history/transfer model needs one canonical
  owning Space. The unified model preserves multi-space association AND adds a
  canonical owner — confirm these coexist cleanly (which Space "owns" history when
  an asset is shared across several?).
- **Reminder/task-link fields.** Legacy vehicle NCT/tax/service due-date task links
  point into a tasks module — confirm that integration target still exists, or the
  links become plain dates.
