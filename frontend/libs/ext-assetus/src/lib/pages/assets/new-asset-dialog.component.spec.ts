import { TestBed } from '@angular/core/testing';
import { ModalController } from '@ionic/angular/standalone';
import { ErrorLogger } from '@sneat/core';
import { of } from 'rxjs';
import { AssetService } from '../../services';
import { NewAssetDialogComponent } from './new-asset-dialog.component';

describe('NewAssetDialogComponent', () => {
  const createAsset = vi.fn().mockReturnValue(of({ id: 'a1', asset: {} }));
  const dismiss = vi.fn().mockResolvedValue(undefined);

  beforeEach(() => {
    createAsset.mockClear();
    dismiss.mockClear();
    TestBed.configureTestingModule({
      imports: [NewAssetDialogComponent],
      providers: [
        { provide: AssetService, useValue: { createAsset } },
        { provide: ModalController, useValue: { dismiss } },
        {
          provide: ErrorLogger,
          useValue: { logError: () => undefined, logErrorHandler: () => () => undefined },
        },
      ],
    });
  });

  it('defaults visibility to the space default on init', () => {
    const fixture = TestBed.createComponent(NewAssetDialogComponent);
    fixture.componentInstance.spaceID = 's1';
    fixture.componentInstance.spaceType = 'family';
    fixture.detectChanges();
    expect(fixture.nativeElement.querySelector('ion-title')?.textContent).toContain(
      'New asset',
    );
  });

  it('calls createAsset with the entered name and inherited visibility', () => {
    const fixture = TestBed.createComponent(NewAssetDialogComponent);
    const cmp = fixture.componentInstance as unknown as {
      spaceID: string;
      spaceType?: string;
      name: string;
      create: () => void;
      ngOnInit: () => void;
    };
    cmp.spaceID = 's1';
    cmp.spaceType = 'private';
    cmp.ngOnInit();
    cmp.name = 'Tent';
    cmp.create();
    expect(createAsset).toHaveBeenCalledWith(
      expect.objectContaining({ spaceID: 's1', name: 'Tent', visibility: 'private' }),
    );
  });
});
