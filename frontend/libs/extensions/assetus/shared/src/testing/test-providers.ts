import { EnvironmentProviders, Provider } from '@angular/core';
import { provideRouter } from '@angular/router';
import { appEnvironmentConfig, getStandardSneatProviders } from '@sneat/app';
import { APP_INFO } from '@sneat/core';
import { provideErrorLogger } from '@sneat/logging';
import { SneatUserService } from '@sneat/auth-core';
import { SpaceService } from '@sneat/space-services';
import { BehaviorSubject, NEVER } from 'rxjs';

// Shared TestBed providers for the ported assetus pages/components that extend
// SpaceBaseComponent (or otherwise pull the standard Sneat DI chain). Mirrors
// the assetus-app component specs: the standard providers + a router + a stub
// user service so every injected dependency resolves without hitting Firebase.

const testEnvironmentConfig = appEnvironmentConfig({
  production: false,
  agents: {},
  firebaseConfig: {
    projectId: 'test',
    appId: 'test',
    apiKey: 'test',
    messagingSenderId: 'test',
  },
  signInMethod: 'redirect',
});

const userState$ = new BehaviorSubject<unknown>({
  status: 'authenticated',
  user: { uid: 'u1', isAnonymous: false, emailVerified: true, providerData: [] },
  record: { title: 'Test User', spaces: {} },
});

// Minimal providers for presentational components that pull @sneat/ui widgets
// (e.g. SelectFromListComponent), which inject the ErrorLogger token.
export function componentTestProviders(): (Provider | EnvironmentProviders)[] {
  return [provideErrorLogger()];
}

export function spacePageTestProviders(): (Provider | EnvironmentProviders)[] {
  return [
    ...getStandardSneatProviders(testEnvironmentConfig),
    provideRouter([]),
    {
      provide: APP_INFO,
      useValue: { appId: 'sneat', appTitle: 'Assetus' },
    },
    {
      provide: SneatUserService,
      useValue: { userState: userState$, currentUserID: 'u1' },
    },
    // SpaceService pulls a Firestore/Auth chain that is irrelevant to a static
    // page render; a stub whose watchSpace never emits keeps the page inert.
    {
      provide: SpaceService,
      useValue: { watchSpace: () => NEVER, onSpaceUpdated: () => undefined },
    },
  ];
}
