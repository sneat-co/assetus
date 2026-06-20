---
format: https://specscore.md/idea-specification
status: Approved
---

# Idea: Adopt extension library-architecture convention for assetus

**Status:** Approved
**Date:** 2026-06-20
**Owner:** alexandertrakhimenok
**Promotes To:** —
**Supersedes:** —
**Related Ideas:** —

## Problem Statement

How might we split the monolithic ext-assetus lib into the contract/shared/internal convention so external repos can depend on light asset contracts instead of the full bundle?

## Context

sneat-libs (contactus) can only import asset types from @sneat/extension-assetus as type-only to dodge the prebuilt-bundle peer-resolution wall. The convention (defined in sneat-libs spec/features/extension-library-architecture) splits an extension into a runtime-light contract, a public shared tier, and a private internal tier so a -contract lib can be depended on safely. assetus is a single-extension repo with one in-repo consumer (its app).

## Recommended Direction

Split ext-assetus into @sneat/extension-assetus-{contract,shared,internal}: contract holds IAsset types/DTOs/contexts/enums + ASSET_SERVICE token (runtime-light); shared holds the app-facing routing + pages + components (injects the token); internal holds AssetService + the services module (provides the token). Add nx tier tags + enforce-module-boundaries matrix. Defer the cross-repo Part B (publish -contract + rewire sneat-libs type-only imports).

## Alternatives Considered

<!-- 2–3 directions that lost, and why each lost. -->

## MVP Scope

In-repo split + nx enforcement, full CI green, app builds against the shared routing.

## Not Doing (and Why)

- Publishing the new packages — deferred; cross-repo release coordination
- Rewiring sneat-libs type-only imports to the new contract — deferred; publish-gated, same category as the deferred sneat-apps wiring

## Key Assumptions to Validate

| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true | placeholder dealbreaker assumption | describe how to validate |
| Should-be-true | … | … |
| Might-be-true | … | … |


## SpecScore Integration

- **New Features this would create:** TBD at design time
- **Existing Features affected:** none
- **Dependencies:** none

## Open Questions

None at this time.
