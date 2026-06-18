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
];
