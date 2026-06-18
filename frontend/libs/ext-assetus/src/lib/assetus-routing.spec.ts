import { assetusRoutes } from './assetus-routing';

describe('assetusRoutes', () => {
  it('exposes the assets overview route', () => {
    expect(assetusRoutes.some((r) => r.path === 'assets')).toBe(true);
  });

  it('exposes the asset detail route with an assetID param', () => {
    expect(assetusRoutes.some((r) => r.path === 'asset/:assetID')).toBe(true);
  });

  it('lazy-loads every route via loadComponent', () => {
    for (const route of assetusRoutes) {
      expect(typeof route.loadComponent).toBe('function');
    }
  });
});
