import { assetusRoutes } from './assetus-routing';

describe('assetusRoutes', () => {
  it('exposes the assets overview route', () => {
    expect(assetusRoutes.some((r) => r.path === 'assets')).toBe(true);
  });

  it('exposes the asset detail route with an assetID param', () => {
    expect(assetusRoutes.some((r) => r.path === 'asset/:assetID')).toBe(true);
  });

  it('exposes the ported asset-shaped pages', () => {
    for (const path of [
      'new-asset',
      'real-estates',
      'real-estate/:assetID',
      'asset-group',
      'optimization',
    ]) {
      expect(assetusRoutes.some((r) => r.path === path)).toBe(true);
    }
  });

  it('gives every route a title', () => {
    for (const route of assetusRoutes) {
      expect((route.data as { title?: string } | undefined)?.title).toBeTruthy();
    }
  });

  it('lazy-loads every route via loadComponent', () => {
    for (const route of assetusRoutes) {
      expect(typeof route.loadComponent).toBe('function');
    }
  });
});
