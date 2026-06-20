import { TestBed } from '@angular/core/testing';
import { ErrorLogger } from '@sneat/core';
import { ModalController } from '@ionic/angular/standalone';
import { ASSET_SERVICE } from '@sneat/extension-assetus-contract';
import { NewAssetDialogComponent } from './new-asset-dialog.component';

// Guard against the production bug where the New Asset dialog would not open.
// The dialog is created by Ionic's ModalController in the APP-ROOT injector, so
// it never sees any component-scoped providers. It injects the ASSET_SERVICE
// token (bound to the concrete AssetService at app bootstrap via
// provideAssetusInternal); here a stub stands in for that root binding so the
// dialog and the token both resolve.
describe('NewAssetDialogComponent (root injector)', () => {
  const assetServiceStub = { createAsset: vi.fn() };

  beforeEach(() =>
    TestBed.configureTestingModule({
      imports: [NewAssetDialogComponent],
      providers: [
        { provide: ASSET_SERVICE, useValue: assetServiceStub },
        { provide: ErrorLogger, useValue: { logError: vi.fn() } },
        { provide: ModalController, useValue: { dismiss: vi.fn() } },
      ],
    }),
  );

  it('creates the dialog and resolves ASSET_SERVICE from root scope', () => {
    const fixture = TestBed.createComponent(NewAssetDialogComponent);
    fixture.componentInstance.spaceID = 's1';
    fixture.detectChanges();
    expect(fixture.componentInstance).toBeTruthy();
    expect(TestBed.inject(ASSET_SERVICE)).toBeTruthy();
  });
});
