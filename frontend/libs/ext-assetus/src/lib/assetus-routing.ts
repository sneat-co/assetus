import { Route } from '@angular/router';

export const assetusRoutes: Route[] = [
  {
    path: 'assets',
    data: { title: 'Assets' },
    loadComponent: () =>
      import('./pages/assets/assets-page.component').then(
        (m) => m.AssetsPageComponent,
      ),
  },
  {
    path: 'asset/:assetID',
    data: { title: 'Asset' },
    loadComponent: () =>
      import('./pages/asset/asset-page.component').then(
        (m) => m.AssetPageComponent,
      ),
  },
  {
    path: 'new-asset',
    data: { title: 'New asset' },
    loadComponent: () =>
      import('./pages/new-asset/new-asset-page.component').then(
        (m) => m.NewAssetPageComponent,
      ),
  },
  {
    path: 'real-estates',
    data: { title: 'Properties' },
    loadComponent: () =>
      import('./pages/real-estates/real-estates-page.component').then(
        (m) => m.RealEstatesPageComponent,
      ),
  },
  {
    path: 'real-estate/:assetID',
    data: { title: 'Property' },
    loadComponent: () =>
      import('./pages/real-estate/real-estate-page.component').then(
        (m) => m.RealEstatePageComponent,
      ),
  },
  {
    path: 'asset-group',
    data: { title: 'Asset group' },
    loadComponent: () =>
      import('./pages/asset-group/asset-group-page.component').then(
        (m) => m.AssetGroupPageComponent,
      ),
  },
  {
    path: 'optimization',
    data: { title: 'Savings' },
    loadComponent: () =>
      import('./pages/optimization/optimization-page.component').then(
        (m) => m.OptimizationPageComponent,
      ),
  },
];
