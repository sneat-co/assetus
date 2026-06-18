import { TestBed } from '@angular/core/testing';
import { Firestore as AngularFirestore } from '@angular/fire/firestore';
import { SneatApiService } from '@sneat/api';
import { of } from 'rxjs';
import { AssetService } from './asset.service';

describe('AssetService', () => {
  let service: AssetService;
  const post = vi.fn().mockReturnValue(of({ id: 'a1', asset: {} }));
  const get = vi.fn().mockReturnValue(of({ id: 'a1', asset: {}, owner: {} }));

  beforeEach(() => {
    post.mockClear();
    get.mockClear();
    TestBed.configureTestingModule({
      providers: [
        AssetService,
        { provide: SneatApiService, useValue: { post, get } },
        { provide: AngularFirestore, useValue: {} },
      ],
    });
    service = TestBed.inject(AssetService);
  });

  it('creates', () => {
    expect(service).toBeTruthy();
  });

  it('posts create_asset to the assetus path', () => {
    service
      .createAsset({
        spaceID: 's1',
        name: 'Drill',
        category: 'tools',
        condition: 'good',
      })
      .subscribe();
    expect(post).toHaveBeenCalledWith(
      'assetus/create_asset',
      expect.objectContaining({ name: 'Drill' }),
    );
  });

  it('posts update_asset', () => {
    service
      .updateAsset({
        spaceID: 's1',
        assetID: 'a1',
        name: 'Drill',
        category: 'tools',
        condition: 'fair',
        visibility: 'family',
      })
      .subscribe();
    expect(post).toHaveBeenCalledWith(
      'assetus/update_asset',
      expect.objectContaining({ assetID: 'a1' }),
    );
  });

  it('posts remove_asset with hardDelete flag', () => {
    service
      .removeAsset({ spaceID: 's1', assetID: 'a1', hardDelete: true })
      .subscribe();
    expect(post).toHaveBeenCalledWith(
      'assetus/remove_asset',
      expect.objectContaining({ hardDelete: true }),
    );
  });

  it('posts transfer_asset', () => {
    service
      .transferAsset({ spaceID: 's1', assetID: 'a1', toSpaceID: 's2' })
      .subscribe();
    expect(post).toHaveBeenCalledWith(
      'assetus/transfer_asset',
      expect.objectContaining({ toSpaceID: 's2' }),
    );
  });

  it('posts record_history_event', () => {
    service
      .recordHistoryEvent({ spaceID: 's1', assetID: 'a1', type: 'repaired' })
      .subscribe();
    expect(post).toHaveBeenCalledWith(
      'assetus/record_history_event',
      expect.objectContaining({ type: 'repaired' }),
    );
  });

  it('gets a single asset', () => {
    service.getAsset('s1', 'a1').subscribe();
    expect(get).toHaveBeenCalledWith('assetus/asset', expect.anything());
  });

  it('gets asset history', () => {
    service.getHistory('s1', 'a1').subscribe();
    expect(get).toHaveBeenCalledWith(
      'assetus/asset_history',
      expect.anything(),
    );
  });
});
