// Main entry point for assetus.app
import { bootstrapApplication } from '@angular/platform-browser';
import { provideRouter } from '@angular/router';
import {
  getStandardSneatProviders,
  provideAppInfo,
  provideRolesByType,
} from '@sneat/app';
import { SneatApp } from '@sneat/core';
import { authRoutes } from '@sneat/auth-ui';
import { App } from './app/app';
import { appRoutes } from './app/app.routes';
import { assetusAppEnvironmentConfig } from './environments/environment';
import { registerIonicons } from './register-ionicons';

bootstrapApplication(App, {
  providers: [
    ...getStandardSneatProviders(assetusAppEnvironmentConfig),
    // 'assetus' is not yet a member of the @sneat/core SneatApp union; cast
    // until it is registered upstream (tracked follow-up).
    provideAppInfo({ appId: 'assetus' as SneatApp, appTitle: 'Assetus.app' }),
    provideRouter([...appRoutes, ...authRoutes]),
    provideRolesByType(undefined),
  ],
}).catch((err) => console.error(err));

registerIonicons();
