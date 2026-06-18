# Legacy Assetus → MVP: Gap Analysis & Migration Report

**Date:** 2026-06-18
**Status:** Reference (informs the follow-up port/retire decision)
**Scope:** What the pre-existing ("legacy") Assetus code did, what the new MVP
implements, and the concrete differences/conflicts that prevented a straight
port. Companion to the relocation-notice PRs (sneat-go-backend #193, sneat-libs
#2, sneat-apps #3411).

## TL;DR

The legacy code and the MVP are **two different products that happen to share a
name**:

- **Legacy Assetus** = a *rich, typed asset registry* for high-value/serialised
  goods — **vehicles, real-estate/dwellings, sport gear, documents** — with
  make/model/VIN/engine specs, mileage logs, insurance/certification dates,
  liabilities & service providers. It models *what kind of thing it is and its
  technical/financial attributes*.
- **MVP Assetus** = an *ownership system of record* — the smallest reusable slice
  that answers *"what do we own, in what condition, who can see it, and what is
  its ownership history?"* — deliberately category-agnostic, with **condition,
  ownership-lifecycle status, visibility, append-only history, and Space→Space
  transfer**. It is the substrate a later **Yardius** sharing layer reads.

The MVP Feature **explicitly defers** the legacy's rich-registry capabilities
(see its *Not Doing* section). Because the two models **conflict at the enum and
schema level** (different `Category` sets, different `Status` semantics,
`Possession` vs ownership-lifecycle, polymorphic "extras" vs a flat record,
multi-space vs single-owning-Space), the MVP was built **clean** rather than
ported. **No legacy code was deleted** — to honour "no functionality lost" it
remains in the source repos behind `RELOCATED.md` notices until its still-wanted
capabilities are ported or formally retired.

---

## 1. Legacy inventory

### 1.1 Backend — `sneat-go-backend/pkg/extensions/assetus/` (51 Go files)

Same persistence pattern as the MVP (Firestore via dalgo,
`dal4spaceus.RunModuleSpaceWorkerWithUserCtx` membership gating,
`/spaces/{spaceID}/ext/assetus/assets/{assetID}`), but a very different model.

**Closed enums (`const4assetus/`):**

| Enum | Legacy values |
|------|---------------|
| `AssetCategory` | `vehicle`, `dwelling`, `sport_gear`, `document` |
| `AssetStatus` | `active`, `archived`, `draft` (aliases of `dbmodels.Status*`) |
| `AssetPossession` | `unknown`, `undisclosed`, `owning`, `leasing`, `renting` |
| `AssetType` (subtype, validated per-category) | vehicles: `car/bus/van/truck/motorcycle/boat/aircraft/helicopter`; real-estate: `apartment/house/office/shop/land/garage/warehouse`; sport gear: `bicycle/kite/kite_bar/kite_board/kite_hydrofoil/prone_hydrofoil/surf_board/wetsuit/wing/wing_board/wing_hydrofoil`; documents: `passport/id_card/driving_license/marriage_cert/birth_cert` |
| `EngineType` / `FuelType` | combustion/hybrid/PHEV/electric/steam/… × petrol/diesel/hydrogen/… |

**Record model (`dbo4assetus/`, `briefs4assetus/`):**

- `AssetBrief`: `Title`, `Status`, `Category`, `Type`, `Possession`, **`CountryID`
  (required)**, a **polymorphic `WithExtraField`** (typed extension), and
  `WithOptionalRelatedAs` (a `linkage` relation hook).
- `AssetBaseDbo` = `AssetBrief` + `WithAssetSpaces` (a **multi-space** brief map —
  one asset can be briefed under several spaces) + `TagsField` +
  `WithCustomFields` + `AssetDates` (`dateOfBuild`, `dateOfPurchase`,
  `dateInsuredTill`, `dateCertifiedTill`).
- `AssetDbo` = `AssetBaseDbo` + `WithModified` + `WithUserIDs` + `WithSpaceIDs`.
- **Polymorphic extras** (`extras4assetus/`, registered by `extraType`):
  - `AssetVehicleExtra`: make/model/regNumber/VIN, `WithEngineData`
    (type/fuel/CC/kW/Nm), `VehicleMileage`.
  - `AssetDwellingExtra`: address, rent price, bedrooms, area (m²).
  - `AssetDocumentExtra`: regNumber, issuedOn/effectiveFrom/expiresOn.
- `VehicleMileage` records — a separate child collection with its own
  `add_vehicle_record` facade + `http_create_vehicle_record` endpoint.

**Facades / API:** `CreateAsset`, `GetAsset`, `UpdateAsset`, `DeleteAsset`
(hard delete only), `AddVehicleRecord`. No history, no transfer, no visibility.

### 1.2 Frontend — `sneat-libs/libs/extensions/assetus/{core,components}` + `sneat-apps/libs/extensions/assetus/pages`

- **`core`** DTOs: `IAssetBrief` / `IAssetDbo` (mirrors the backend rich model;
  adds `desc`, `memberIDs`, `parentAssetID`, `groupId`, `subAssets`,
  `liabilities`, `membersInfo`, vehicle NCT/tax/service due-date task links),
  `IAssetVehicleExtra` / `IAssetDwellingExtra`, `IAssetDtoGroup` (asset groups),
  plus a whole **Liabilities / Service-Provider** sub-domain
  (`DtoServiceProvider`, `DtoServicePlan`, `DtoServiceType`, `LiabilityServiceType`,
  `SettlementType`) — utility/insurance providers and plans linked to assets.
- **`components`** (~19): `vehicle-card`, `vehicle-engine`, `make-model-card`,
  `asset-reg-number-input`, `mileage-dialog`, `edit-dwelling-card`,
  `real-estate-location`, `asset-dates`, `period-segment`,
  `asset-possesion-card`, `asset-liabilities`, `asset-contacts-group`,
  `asset-add-{vehicle,dwelling,document,service}`, `assets-list`, `asset-card`, …
- **`pages`** (sneat-apps): `assets`, `asset`, `new-asset`, `real-estates`,
  `asset-group`, `liability/select-service-provider`.

---

## 2. MVP inventory (this repo)

**Enums (`const4assetus/`):** `Category` (`books/games/toys/sports_equipment/
tools/electronics/clothing/vehicles/camping_equipment/other`), `Condition`
(`new/excellent/good/fair/needs_repair/broken`), `Status` (ownership lifecycle:
`active/transferred/archived/disposed/lost`), `Visibility`
(`private/family/friends/friends_of_friends/specific_space/public`), `OwnerType`
(derived), `HistoryEventType` (`purchased/repaired/transferred/sold/donated/lost`).

**Record:** `AssetDbo` = flat `AssetBase` (name, description, category, condition,
status, visibility, optional acquisitionDate/purchasePrice/estimatedValue/
location/notes/tags, reserved `photos[]`) + `WithModified` + `WithSpaceIDs`
(single owning Space). Child collection `history/` (append-only). `OwnerRef`
derived at read.

**Facades:** create (membership-gated, visibility default/override), get (with
derived owner), update (history untouched), record-history, get-history, remove
(soft-archive default + hard-delete), transfer (Space→Space, relocates asset +
history, appends `Transferred`).

---

## 3. Side-by-side comparison

### 3.1 Capabilities

| Capability | Legacy | MVP |
|---|---|---|
| Membership-gated CRUD on a space module | ✅ | ✅ (same dalgo/spaceus pattern) |
| Tags, custom fields | ✅ (`TagsField`, `WithCustomFields`) | tags ✅; custom fields ❌ (dropped) |
| **Condition** (new/good/needs-repair/…) | ❌ | ✅ |
| **Ownership-lifecycle status** (active/transferred/archived/disposed/lost) | ❌ (only active/archived/draft) | ✅ |
| **Visibility** (private…public, inherits Space default) | ❌ | ✅ |
| **Append-only history** | ❌ | ✅ |
| **Ownership transfer (Space→Space) + history** | ❌ | ✅ |
| **Owner-type derivation** from Space type | ❌ | ✅ |
| Soft-archive vs hard-delete distinction | ❌ (delete = hard; archived is a status) | ✅ (soft-archive default + explicit hard-delete) |
| Typed sub-categories (`AssetType`) | ✅ (car/house/passport/…) | ❌ (deferred) |
| **Possession** (owning/leasing/renting) | ✅ | ❌ (out of MVP; owner == Space) |
| Polymorphic **extras** (vehicle/dwelling/document) | ✅ | ❌ (deferred) |
| Make/model/VIN/**engine** specs | ✅ | ❌ (deferred) |
| **Mileage** log (child collection) | ✅ | ❌ (deferred) |
| Insurance/certification **dates** | ✅ (`AssetDates`) | ❌ (only optional `acquisitionDate`) |
| **Liabilities / service providers / plans** | ✅ (whole sub-domain) | ❌ (out of scope; Debtus territory) |
| **Asset groups** | ✅ (`IAssetDtoGroup`) | ❌ |
| **Sub-assets / parent asset** | ✅ (`parentAssetID`, `subAssets`) | ❌ |
| **Multi-space** association | ✅ (`WithAssetSpaces`) | ❌ (single owning Space; multi-owner deferred to `linkage`) |
| Required `CountryID` | ✅ | ❌ (not modelled) |
| Reserved `photos[]`, `ext.yardius` | ❌ | ✅ (reserved, unpopulated) |

### 3.2 Enum conflicts (why a naive merge breaks)

| Concept | Legacy | MVP | Conflict |
|---|---|---|---|
| `Category` | `vehicle/dwelling/sport_gear/document` (a *kind-of-thing*) | `books/games/…/vehicles/…/other` (a *household-goods taxonomy*) | **Different axes & values.** Only `vehicles` overlaps by name. Legacy `dwelling`/`document` have no MVP equivalent; MVP `books/toys/clothing/…` have no legacy equivalent. |
| `Status` | `active/archived/draft` | `active/transferred/archived/disposed/lost` | **Same field name, different domain.** `draft` is gone; `transferred/disposed/lost` are new ownership-lifecycle states. |
| `Possession` | `owning/leasing/renting/…` | — | MVP drops possession entirely (owner == the Space). A leased/rented asset has no MVP representation. |
| `Condition` | — | `new/…/broken` | New required field with no legacy source value. |
| `Visibility` | — | `private…public` | New required field; legacy had none → would need a default on migration. |
| `Type` (subtype) | required, validated per category | — | MVP has no subtype; legacy data has no flat-model target. |

---

## 4. What was NOT migrated (deliberately deferred)

Per the MVP Feature's *Not Doing* section, none of the following were ported:

1. **Vehicles sub-domain** — make/model/VIN, engine (type/fuel/CC/kW/Nm),
   mileage logs, NCT/tax/service due-date reminders.
2. **Real-estate/dwelling sub-domain** — address, rent price, bedrooms, area.
3. **Documents sub-domain** — passport/licence/etc. with issue/expiry dates.
4. **Typed subcategories** (`AssetType`) and the per-category validation.
5. **Possession** semantics (owning/leasing/renting).
6. **Liabilities, service providers, service plans, settlement types** — the
   whole utilities/insurance sub-domain (this is closer to **Debtus**' remit).
7. **Asset groups**, **sub-assets / parent asset**, `sameAssetID` linking.
8. **Multi-space association** (`WithAssetSpaces`) and member-level info.
9. **Insurance/certification date tracking** (`AssetDates`) and the
   reminder/task-link fields.
10. **Custom fields**, required `CountryID`, `geo`/`relatedAs` hooks.
11. The ~19 legacy Angular components and the asset/real-estate/group pages.

## 5. What the MVP ADDS that legacy never had

Condition; ownership-lifecycle status (transferred/disposed/lost); visibility
(with Space-default inheritance + per-asset override); an append-only per-asset
**history**; **ownership transfer** with prior/new-owner preservation;
**owner-type derivation** from the Space type; the soft-archive vs hard-delete
distinction; and the reserved forward-compatibility hooks (`photos[]`,
`ext.yardius`) for the future Yardius sharing layer.

---

## 6. Migration path & risks (for the follow-up)

If/when the legacy rich capabilities are wanted in this repo:

1. **Treat them as additive extension modules on the MVP core**, not a schema
   merge. The MVP's flat `AssetBase` + a reintroduced polymorphic `extra`
   (vehicle/dwelling/document) is the natural shape — keep the MVP's
   condition/status/visibility/history/transfer as the shared core and hang the
   typed extras off it. Avoid reviving the legacy `Status`/`Category`/`Possession`
   enums as-is; map them onto the MVP sets instead.
2. **Data migration is non-trivial** (greenfield repo, *no backward compat*
   required by product, but existing legacy *data* in `sneat-eur3-1` would need a
   one-off transform): legacy `Category` → MVP `Category` + a subtype/extra;
   legacy `Status` (`draft`→? , `archived`→`archived`); supply a default
   `Condition` and `Visibility`; `Possession=leasing/renting` has no MVP home and
   needs a product decision.
3. **Backend coupling to resolve before deleting legacy:**
   `pkg/extensions/brandus/dbo4brands/make_test.go` imports
   `const4assetus.AssetCategory`, and `pkg/extensions/standard_extensions.go`
   registers `assetus.Extension()`. Both must change when the legacy package is
   removed; then `go build ./... && go test ./...` must pass.
4. **Frontend consumers to migrate:** `@sneat/ext-assetus-components` and
   `@sneat/mod-assetus-core` are consumed by `sneat-apps` (and indirectly the
   space app shell). Repoint or retire them, then `pnpm install && nx run-many -t
   lint test build`.
5. **Liabilities/service-providers**: decide whether they belong in Assetus at
   all — they read more like **Debtus** (obligations) or a dedicated `servius`
   module than an ownership record.

## 7. Recommendation

Keep the MVP as the canonical ownership core. Port the **vehicle/dwelling/document
typed-extras** (the genuinely asset-shaped capabilities) into this repo as
opt-in extension modules in a follow-up Feature, mapping their enums onto the MVP
sets. Route **liabilities/service-providers** to a separate Feature/module
decision. Only after the wanted capabilities live here (and consumers are
repointed) should the legacy directories be deleted and the relocation-notice
PRs be followed by deletion PRs.
