// Asset domain model — mirrors the Go backend
// (github.com/sneat-co/assetus/backend). Assetus is OWNERSHIP-ONLY: there is no
// sharing/availability/borrow/lend field and no ext.yardius.

import { IWithSpaceIDs } from '@sneat/dto';

// --- Enums (string values must match the backend exactly) ---

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
  | 'other';

export type AssetCondition =
  | 'new'
  | 'excellent'
  | 'good'
  | 'fair'
  | 'needs_repair'
  | 'broken';

// Ownership lifecycle — read-only-ish (the backend drives transitions).
export type AssetStatus =
  | 'active'
  | 'transferred'
  | 'archived'
  | 'disposed'
  | 'lost';

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

export interface IMoney {
  currency: string;
  value: number;
}

export interface IOwner {
  spaceID: string;
  spaceType: string;
  ownerType: OwnerType;
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

export interface IAssetDbo extends IAssetBrief, IWithSpaceIDs {
  description?: string;
  acquisitionDate?: string; // ISO date
  purchasePrice?: IMoney;
  estimatedValue?: IMoney;
  location?: string;
  notes?: string;
  tags?: string[];
  // Backend serializes these as ISO datetime strings.
  createdAt?: string;
  updatedAt?: string;
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
