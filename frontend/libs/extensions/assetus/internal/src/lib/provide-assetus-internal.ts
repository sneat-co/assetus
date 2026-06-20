import { Provider } from '@angular/core';
import { ASSET_SERVICE } from '@sneat/extension-assetus-contract';
import { AssetService } from './services';

// Registers the concrete AssetService and binds it to the ASSET_SERVICE token so
// consumers depend only on the IAssetService contract. Wired in at app bootstrap
// (consumers do not import this factory directly).
export function provideAssetusInternal(): Provider[] {
  return [
    AssetService,
    { provide: ASSET_SERVICE, useExisting: AssetService },
  ];
}
