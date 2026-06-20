# assetus

Assetus — the Sneat **ownership system of record**. A Sneat **Space** owns
physical **Assets** carrying typed metadata, condition, an ownership-lifecycle
status, visibility, and an append-only per-asset history. The owning Space can
transfer an asset to another Space. Assetus answers *"what do we own?"*; a later
**Yardius** sharing layer reads a pinned contract and owns all sharing
transactions — Assetus never embeds borrow/lend/sell logic.

**License:** [AGPL-3.0](LICENSE)

## Repository layout

This repo hosts two independent toolchains in subdirectories — neither
`package.json` nor `go.mod` lives at the repo root (mirrors the
[`listus`](https://github.com/sneat-co/listus) layout):

| Directory | Stack | Description |
|-----------|-------|-------------|
| [`backend/`](backend) | Go 1.26 | The `assetus` space module (Firestore via dalgo under `/spaces/{spaceID}/ext/assetus/...`) |
| [`frontend/`](frontend) | Nx · Angular · Ionic · pnpm | The `assetus-app` standalone app and the `@sneat/extension-assetus-*` libraries (see [Frontend](#frontend)) |

## Backend

The backend is a Sneat **space module** persisted on **Firestore via dalgo**
under `/spaces/{spaceID}/ext/assetus/...`, the same convention as `contactus`,
`eventus`, and `listus`. Membership of the owning Space (`spaceus`) is the access
boundary.

```bash
cd backend
go build ./...
go test ./...
```

### Data model

- **`Asset`** aggregate at `/spaces/{spaceID}/ext/assetus/assets/{assetID}` —
  `name`, `description?`, `category`, `condition`, `status`, `visibility`, and
  optional `acquisitionDate` / `purchasePrice` / `estimatedValue` / `location` /
  `notes` / `tags[]`. `photos[]` and `ext.yardius` are **reserved but not
  populated** in the MVP.
- **`AssetHistoryEvent`** append-only child collection at
  `…/assets/{assetID}/history/{eventID}`.
- **Closed enums** validated on write: `Category`, `Condition`, `Status`
  (`Active`/`Transferred`/`Archived`/`Disposed`/`Lost` only — sharing/availability
  states such as `Borrowed`/`Reserved` are rejected), `Visibility`.
- **Owner == owning Space**; the owner **type** is derived from the Space type
  (`private`→Individual, `family`→Family, `club`→SportsClub, `company`→Organisation;
  `community`→Community and `school`→School once those `spaceus` space types ship).

Spec: [`sneat-co/backstage` → `spec/features/assetus-mvp`](https://github.com/sneat-co/backstage/tree/main/spec/features/assetus-mvp).

## Frontend

```bash
cd frontend
pnpm install
npx nx build assetus-app
npx nx run-many -t lint test
```

### Library structure (extension library-architecture convention)

The assetus frontend follows the **extension library-architecture** convention —
an extension is split into three libraries by *runtime weight* and *visibility*,
so other repos can depend on a light **contract** instead of the full bundle, and
cross-extension calls go through dependency-inverted `InjectionToken`s rather than
direct implementation imports. The convention is defined in
[`sneat-co/sneat-libs` → `spec/features/extension-library-architecture`](https://github.com/sneat-co/sneat-libs/tree/main/spec/features/extension-library-architecture/README.md).

| Lib | nx tags | Holds | May depend on |
|-----|---------|-------|---------------|
| [`@sneat/extension-assetus-contract`](frontend/libs/extensions/assetus/contract) | `type:contract` | Asset DTOs/types/enums + the `ASSET_SERVICE` `InjectionToken` (`IAssetService`). Runtime-light — no components/services. | other contracts + foundational `@sneat/*` |
| [`@sneat/extension-assetus-shared`](frontend/libs/extensions/assetus/shared) | `type:shared` | The app-facing UI: routing, pages, components, space-menu. Obtains services via the `ASSET_SERVICE` token. | `-contract` + foundational — **never `-internal`** |
| [`@sneat/extension-assetus-internal`](frontend/libs/extensions/assetus/internal) | `type:internal` | `AssetService` + `provideAssetusInternal()`. Private implementation. | `-contract` / `-shared` + foundational |

The boundary matrix is enforced by `@nx/enforce-module-boundaries` in
`frontend/eslint.config.mjs` (a `type:shared → type:internal` import fails lint).
`-internal` is consumed only by the composition-root **app**, which wires
`provideAssetusInternal()` at bootstrap (`frontend/apps/assetus-app/src/main.ts`)
to bind `ASSET_SERVICE` to the concrete `AssetService`.

The local plan that performed this split:
[`spec/plans/extension-library-architecture.md`](spec/plans/extension-library-architecture.md).
