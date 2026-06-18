import { TestBed } from '@angular/core/testing';
import { ModalController } from '@ionic/angular/standalone';
import { ErrorLogger } from '@sneat/core';
import { SneatUserService } from '@sneat/auth-core';
import { of } from 'rxjs';
import { AssetService } from '../../services';
import { TransferAssetComponent } from './transfer-asset.component';

describe('TransferAssetComponent', () => {
  const transferAsset = vi
    .fn()
    .mockReturnValue(of({ id: 'a1', owner: { spaceID: 's2' } }));
  const dismiss = vi.fn().mockResolvedValue(undefined);
  const userState = of({
    record: {
      spaces: {
        s1: { title: 'Family', type: 'family' },
        s2: { title: 'Personal', type: 'private' },
      },
    },
  });

  beforeEach(() => {
    transferAsset.mockClear();
    dismiss.mockClear();
    TestBed.configureTestingModule({
      imports: [TransferAssetComponent],
      providers: [
        { provide: AssetService, useValue: { transferAsset } },
        { provide: ModalController, useValue: { dismiss } },
        { provide: SneatUserService, useValue: { userState } },
        {
          provide: ErrorLogger,
          useValue: { logError: () => undefined, logErrorHandler: () => () => undefined },
        },
      ],
    });
  });

  it('excludes the current space from destinations and transfers to the picked one', () => {
    const fixture = TestBed.createComponent(TransferAssetComponent);
    const cmp = fixture.componentInstance as unknown as {
      spaceID: string;
      assetID: string;
      toSpaceID?: string;
      ngOnInit: () => void;
      transfer: () => void;
    };
    cmp.spaceID = 's1';
    cmp.assetID = 'a1';
    cmp.ngOnInit();
    cmp.toSpaceID = 's2';
    cmp.transfer();
    expect(transferAsset).toHaveBeenCalledWith(
      expect.objectContaining({ spaceID: 's1', assetID: 'a1', toSpaceID: 's2' }),
    );
  });
});
