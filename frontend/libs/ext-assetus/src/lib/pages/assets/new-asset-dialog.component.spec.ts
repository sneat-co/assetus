import { TestBed } from '@angular/core/testing';
import { Firestore } from '@angular/fire/firestore';
import { SneatApiService } from '@sneat/api';
import { ErrorLogger } from '@sneat/core';
import { ModalController } from '@ionic/angular/standalone';
import { AssetService } from '../../services';
import { NewAssetDialogComponent } from './new-asset-dialog.component';

// Guard against the production bug where the New Asset dialog would not open.
// The dialog is created by Ionic's ModalController in the APP-ROOT injector, so
// it never sees AssetusCoreServicesModule (which the page imports). It must be
// able to resolve AssetService from root scope alone. This TestBed deliberately
// does NOT import AssetusCoreServicesModule — it provides only root-equivalent
// doubles for AssetService's transitive deps (Firestore, SneatApiService), plus
// the tokens the dialog itself injects (ErrorLogger, ModalController). With
// AssetService not providedIn:'root' this fails with NullInjectorError; once it
// is providedIn:'root' the component and the service resolve.
describe('NewAssetDialogComponent (root injector)', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      imports: [NewAssetDialogComponent],
      providers: [
        { provide: Firestore, useValue: {} },
        { provide: SneatApiService, useValue: { post: vi.fn(), get: vi.fn() } },
        { provide: ErrorLogger, useValue: { logError: vi.fn() } },
        { provide: ModalController, useValue: { dismiss: vi.fn() } },
      ],
    }),
  );

  it('creates the dialog and resolves AssetService from root scope', () => {
    const fixture = TestBed.createComponent(NewAssetDialogComponent);
    fixture.componentInstance.spaceID = 's1';
    fixture.detectChanges();
    expect(fixture.componentInstance).toBeTruthy();
    expect(TestBed.inject(AssetService)).toBeTruthy();
  });
});
