import { Route } from '@angular/router';
import {
  assetusRoutes,
  AssetusSpaceMenuComponent,
} from '@sneat/extension-assetus';
import { SpaceComponentBaseParams } from '@sneat/space-components';

// Thin, assetus-only space shell. It provides SpaceComponentBaseParams (which
// resolves the active space from the :spaceType/:spaceID route params) to all
// children, then mounts ONLY the assetus routes — unlike sneat-app's
// @sneat/space-pages, which bundles every extension. This keeps assetus.app
// decoupled while reusing the published @sneat/space-components context wiring.
export const assetusSpaceRoutes: Route[] = [
  {
    path: '',
    providers: [SpaceComponentBaseParams],
    children: [
      {
        // assetus-specific side menu (space selector + the space's Assets) instead
        // of the generic SpaceMenuComponent, which hardcodes every sneat-app
        // extension (Budget, Contacts, …) — none of which exist here.
        path: '',
        component: AssetusSpaceMenuComponent,
        outlet: 'menu',
      },
      {
        path: '',
        pathMatch: 'full',
        redirectTo: 'assets',
      },
      ...assetusRoutes,
    ],
  },
];
