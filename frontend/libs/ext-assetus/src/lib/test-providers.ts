import { Provider } from '@angular/core';
import { Firestore as AngularFirestore } from '@angular/fire/firestore';
import { NavController } from '@ionic/angular/standalone';
import { APP_INFO, ErrorLogger, LOGGER_FACTORY } from '@sneat/core';
import { SneatUserService } from '@sneat/auth-core';
import { SpaceService, SpaceNavService } from '@sneat/space-services';
import { of } from 'rxjs';
import { AssetService } from './services';

// Minimal mock provider set for unit-testing the ext-assetus components that
// extend SpaceBaseComponent without pulling in the full @sneat/app bootstrap.
// Each token is mocked just enough for the component to construct and render.
export function assetusTestProviders(
  assetServiceMock: Partial<AssetService> = {},
): Provider[] {
  const errorLogger = {
    logError: () => undefined,
    logErrorHandler: () => () => undefined,
  };
  return [
    { provide: ErrorLogger, useValue: errorLogger },
    {
      provide: LOGGER_FACTORY,
      useValue: { getLogger: () => ({ error: () => undefined }) },
    },
    { provide: APP_INFO, useValue: { appId: 'assetus', appTitle: 'Assetus' } },
    {
      provide: SneatUserService,
      useValue: { userState: of(undefined), currentUserID: undefined },
    },
    { provide: SpaceService, useValue: {} },
    {
      provide: SpaceNavService,
      useValue: { navigateForwardToSpacePage: () => Promise.resolve(true) },
    },
    { provide: NavController, useValue: { navigateForward: () => Promise.resolve(true) } },
    { provide: AngularFirestore, useValue: {} },
    {
      provide: AssetService,
      useValue: {
        watchAssets: () => of([]),
        getAsset: () => of({ id: 'a1', asset: {}, owner: {} }),
        getHistory: () => of({ assetID: 'a1', events: [] }),
        ...assetServiceMock,
      },
    },
  ];
}
