import {
  assetCategories,
  assetConditions,
  assetVisibilities,
  categoryOptions,
  conditionOptions,
  defaultVisibilityForSpaceType,
  deriveOwnerType,
  visibilityOptions,
} from './asset';

describe('defaultVisibilityForSpaceType', () => {
  it('maps private space to private visibility', () => {
    expect(defaultVisibilityForSpaceType('private')).toBe('private');
  });
  it('maps family space to family visibility', () => {
    expect(defaultVisibilityForSpaceType('family')).toBe('family');
  });
  it('falls back to specific_space for other space types', () => {
    expect(defaultVisibilityForSpaceType('sports_club')).toBe('specific_space');
    expect(defaultVisibilityForSpaceType(undefined)).toBe('specific_space');
  });
});

describe('deriveOwnerType', () => {
  it('maps private to individual', () => {
    expect(deriveOwnerType('private')).toBe('individual');
  });
  it('maps family to family', () => {
    expect(deriveOwnerType('family')).toBe('family');
  });
  it('passes through known org-like space types', () => {
    expect(deriveOwnerType('sports_club')).toBe('sports_club');
    expect(deriveOwnerType('community')).toBe('community');
    expect(deriveOwnerType('school')).toBe('school');
  });
  it('defaults unknown space types to organisation', () => {
    expect(deriveOwnerType('whatever')).toBe('organisation');
  });
});

describe('select option lists', () => {
  it('produces an option per category/condition/visibility value', () => {
    expect(categoryOptions).toHaveLength(assetCategories.length);
    expect(conditionOptions).toHaveLength(assetConditions.length);
    expect(visibilityOptions).toHaveLength(assetVisibilities.length);
  });
  it('titleizes multi-word values for labels', () => {
    const sports = categoryOptions.find((o) => o.value === 'sports_equipment');
    expect(sports?.label).toBe('Sports Equipment');
    const repair = conditionOptions.find((o) => o.value === 'needs_repair');
    expect(repair?.label).toBe('Needs Repair');
  });
});
