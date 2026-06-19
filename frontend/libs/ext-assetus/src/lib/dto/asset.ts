// Asset domain model — mirrors the Go backend
// (github.com/sneat-co/assetus/backend). Assetus is OWNERSHIP-ONLY: there is no
// sharing/availability/borrow/lend field and no ext.yardius.

import { CountryId, ITitledRecord, IWithSpaceIDs } from '@sneat/dto';

// --- Enums (string values must match the backend exactly) ---

// AssetCategory mirrors backend const4assetus.Category exactly: the MVP set
// unioned with the ported legacy categories (legacy sport_gear→sports_equipment,
// vehicle→vehicles, misc→other already reconciled; dwelling/document/debt
// retained).
export type AssetCategory =
  | 'books'
  | 'games'
  | 'toys'
  | 'sports_equipment'
  | 'tools'
  | 'electronics'
  | 'clothing'
  | 'vehicles'
  | 'camping_equipment'
  | 'other'
  | 'dwelling'
  | 'document'
  | 'debt';

export type AssetCondition =
  | 'new'
  | 'excellent'
  | 'good'
  | 'fair'
  | 'needs_repair'
  | 'broken';

// Ownership lifecycle — read-only-ish (the backend drives transitions).
// Mirrors backend const4assetus.Status: the MVP set unioned with the ported
// legacy 'draft' pre-active state.
export type AssetStatus =
  | 'draft'
  | 'active'
  | 'transferred'
  | 'archived'
  | 'disposed'
  | 'lost';

// AssetPossession mirrors backend const4assetus.Possession.
export type AssetPossession =
  | 'unknown'
  | 'undisclosed'
  | 'owning'
  | 'leasing'
  | 'renting';

// AssetType is the optional per-category subtype (backend const4assetus.Type).
// It is an open string so any category's subtype set is admitted; the named
// per-category unions below document the backend sets.
export type AssetVehicleType =
  | 'aircraft'
  | 'boat'
  | 'bus'
  | 'car'
  | 'helicopter'
  | 'motorcycle'
  | 'truck'
  | 'van';
export type AssetDwellingType =
  | 'apartment'
  | 'house'
  | 'office'
  | 'shop'
  | 'land'
  | 'garage'
  | 'warehouse';
export type AssetSportsEquipmentType =
  | 'bicycle'
  | 'kite'
  | 'kite_bar'
  | 'kite_board'
  | 'kite_hydrofoil'
  | 'prone_hydrofoil'
  | 'surf_board'
  | 'wetsuit'
  | 'wing'
  | 'wing_board'
  | 'wing_hydrofoil';
export type AssetDocumentType =
  | 'passport'
  | 'id_card'
  | 'driving_license'
  | 'marriage_cert'
  | 'birth_cert';
export type AssetType =
  | AssetVehicleType
  | AssetDwellingType
  | AssetSportsEquipmentType
  | AssetDocumentType
  | string;

// EngineType mirrors backend const4assetus.EngineType ('' = unknown).
export type EngineType =
  | ''
  | 'other'
  | 'combustion'
  | 'electric'
  | 'phev'
  | 'hybrid'
  | 'steam';

// FuelType mirrors backend const4assetus.FuelType ('' = unknown).
export type FuelType =
  | ''
  | 'other'
  | 'bio'
  | 'petrol'
  | 'diesel'
  | 'hydrogen';

// Fuel-volume unit, mileage unit and currency are OPEN strings on the backend
// (no typed enum). The named unions below are kept for ergonomics only and MUST
// serialize as plain strings.
export type FuelVolumeUnit = 'l' | 'g' | string;
export type MileageUnit = 'km' | 'mile' | string;

export type AssetVisibility =
  | 'private'
  | 'family'
  | 'friends'
  | 'friends_of_friends'
  | 'specific_space'
  | 'public';

// Derived from the owning space, read-only.
export type OwnerType =
  | 'individual'
  | 'family'
  | 'sports_club'
  | 'community'
  | 'school'
  | 'organisation';

export type HistoryEventType =
  | 'purchased'
  | 'repaired'
  | 'transferred'
  | 'sold'
  | 'donated'
  | 'lost';

// --- Value objects ---

// IMoney mirrors the backend MonetaryAmount / money.Amount ({currency,value}).
export interface IMoney {
  currency: string;
  value: number;
}

export interface IOwner {
  spaceID: string;
  spaceType: string;
  ownerType: OwnerType;
}

// IGeoPoint mirrors the backend GeoPoint ({lat,lng}).
export interface IGeoPoint {
  lat: number;
  lng: number;
}

// IAssetDates mirrors the backend embedded AssetDates — optional ISO
// 'YYYY-MM-DD' date strings.
export interface IAssetDates {
  dateOfBuild?: string;
  dateOfPurchase?: string;
  dateInsuredTill?: string;
  dateCertifiedTill?: string;
}

// --- Relationship sub-entities (mirror backend WithAssetRelationships) ---

// ISubAssetInfo mirrors the backend SubAssetInfo ({id,title,type,countryID,
// subType,expires}).
export interface ISubAssetInfo extends ITitledRecord {
  type: AssetCategory;
  countryID?: CountryId;
  subType?: string;
  expires?: string; // ISO 'YYYY-MM-DD'
}

// IAssetGroupCounts mirrors the backend AssetGroupCounts ({assets}).
export interface IAssetGroupCounts {
  assets?: number;
}

// IAssetGroupInfo mirrors the backend AssetGroupInfo
// ({id,title,order,desc,categoryID,numberOf,totals}).
export interface IAssetGroupInfo extends ITitledRecord {
  order?: number;
  desc?: string;
  categoryID?: AssetCategory;
  numberOf?: IAssetGroupCounts;
  totals?: IMoney[];
}

// --- Multi-space association (mirror backend WithAssetSpaces) ---

// IAssetusSpaceBrief mirrors the backend AssetusSpaceBrief ({assets}): the
// per-space projection of an asset's briefs (backend AssetBriefs:
// map[assetID]*AssetBrief).
export interface IAssetusSpaceBrief {
  assets?: Record<string, IAssetBrief>;
}

// IWithAssetSpaces mirrors the backend WithAssetSpaces ({spaces}): the
// multi-space association mapping spaceID -> per-space asset briefs, so a single
// asset record can be associated with multiple spaces. This is ADDITIVE to the
// single owning space carried via IWithSpaceIDs.spaceIDs.
export interface IWithAssetSpaces {
  spaces?: Record<string, IAssetusSpaceBrief>;
}

// --- Financial / liability sub-entities ---

// IAssetLiabilityInfo mirrors the backend AssetLiabilityInfo ({id,serviceTypes}).
export interface IAssetLiabilityInfo {
  id: string;
  serviceTypes?: string[];
}

// --- Asset ---

export interface IAssetBrief {
  id: string;
  name: string;
  category: AssetCategory;
  condition: AssetCondition;
  status: AssetStatus;
  visibility: AssetVisibility;
}

export interface IAssetDbo extends IAssetBrief, IWithSpaceIDs, IWithAssetSpaces {
  description?: string;
  acquisitionDate?: string; // ISO date
  purchasePrice?: IMoney;
  estimatedValue?: IMoney;
  location?: string;
  notes?: string;
  tags?: string[];
  photos?: string[];
  // Backend serializes these as ISO datetime strings.
  createdAt?: string;
  updatedAt?: string;

  // --- Ported legacy optional fields (backend AssetBase json names) ---
  isRequest?: boolean;
  countryID?: CountryId;
  type?: AssetType; // per-category subtype
  possession?: AssetPossession;
  parentCategoryID?: AssetCategory;
  yearOfBuild?: number;
  geo?: IGeoPoint;

  // Embedded AssetDates (flattened, like the backend).
  dateOfBuild?: string;
  dateOfPurchase?: string;
  dateInsuredTill?: string;
  dateCertifiedTill?: string;

  // --- Ported financial fields (backend AssetBase json names) ---
  totals?: IMoney[];
  canHaveIncome?: boolean;
  canHaveExpense?: boolean;
  financialDirection?: 'income' | 'expense';
  liabilities?: IAssetLiabilityInfo[];
  notUsedServiceTypes?: string[];

  // --- Ported relationship fields (backend WithAssetRelationships json names) ---
  groupID?: string;
  group?: IAssetGroupInfo;
  parentAssetID?: string;
  subAssets?: ISubAssetInfo[];
  sameAssetID?: string;
  // relatedAs mirrors the backend dbmodels.WithOptionalRelatedAs.RelatedAs — a
  // plain optional string naming the relationship role (json 'relatedAs').
  relatedAs?: string;
  memberIDs?: string[];
  membersInfo?: ITitledRecord[];
}

// --- History (append-only) ---

export interface IHistoryEvent {
  id: string;
  type: HistoryEventType;
  occurredAt: string; // ISO datetime
  actorRef: string;
  note?: string;
  fromOwner?: IOwner;
  toOwner?: IOwner;
}

// --- Select option helpers ---

export interface ILabeledOption<T extends string> {
  value: T;
  label: string;
}

const titleize = (s: string): string =>
  s
    .split('_')
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(' ');

function options<T extends string>(values: readonly T[]): ILabeledOption<T>[] {
  return values.map((value) => ({ value, label: titleize(value) }));
}

export const assetCategories: readonly AssetCategory[] = [
  'books',
  'games',
  'toys',
  'sports_equipment',
  'tools',
  'electronics',
  'clothing',
  'vehicles',
  'camping_equipment',
  'other',
  'dwelling',
  'document',
  'debt',
];

export const assetConditions: readonly AssetCondition[] = [
  'new',
  'excellent',
  'good',
  'fair',
  'needs_repair',
  'broken',
];

export const assetVisibilities: readonly AssetVisibility[] = [
  'private',
  'family',
  'friends',
  'friends_of_friends',
  'specific_space',
  'public',
];

export const categoryOptions = options(assetCategories);
export const conditionOptions = options(assetConditions);
export const visibilityOptions = options(assetVisibilities);

// --- Derivations matching the backend ---

// The default visibility a new asset inherits from its owning space type:
// private→private, family→family, everything else→specific_space.
export function defaultVisibilityForSpaceType(
  spaceType: string | undefined,
): AssetVisibility {
  switch (spaceType) {
    case 'private':
      return 'private';
    case 'family':
      return 'family';
    default:
      return 'specific_space';
  }
}

// Derives the OwnerType from the owning space type. Mirrors the backend's
// read-only derivation.
export function deriveOwnerType(spaceType: string | undefined): OwnerType {
  switch (spaceType) {
    case 'private':
      return 'individual';
    case 'family':
      return 'family';
    case 'sports_club':
      return 'sports_club';
    case 'community':
      return 'community';
    case 'school':
      return 'school';
    default:
      return 'organisation';
  }
}
