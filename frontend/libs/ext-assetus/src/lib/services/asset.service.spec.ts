import { HttpParams } from '@angular/common/http';
import { Injector, runInInjectionContext } from '@angular/core';
import { Firestore } from '@angular/fire/firestore';
import { SneatApiService } from '@sneat/api';
import { of } from 'rxjs';
import { AssetService } from './asset.service';

vi.mock('@angular/fire/firestore', () => ({
  Firestore: class {},
  collection: vi.fn(),
  collectionData: vi.fn(() => of([])),
}));

describe('AssetService', () => {
  let service: AssetService;
  let post: ReturnType<typeof vi.fn>;
  let get: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    post = vi.fn(() => of({}));
    get = vi.fn(() => of({}));
    const injector = Injector.create({
      providers: [
        { provide: Firestore, useValue: {} },
        { provide: SneatApiService, useValue: { post, get } },
      ],
    });
    service = runInInjectionContext(injector, () => new AssetService());
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('createAsset posts to create_asset with the request', () => {
    const request = {
      spaceID: 's1',
      name: 'Car',
      category: 'vehicles' as const,
      condition: 'good' as const,
    };
    service.createAsset(request).subscribe();
    expect(post).toHaveBeenCalledWith('assetus/create_asset', request);
  });

  it('getAsset gets asset with spaceID+assetID params', () => {
    service.getAsset('s1', 'a1').subscribe();
    expect(get).toHaveBeenCalledTimes(1);
    const [endpoint, params] = get.mock.calls[0];
    expect(endpoint).toBe('assetus/asset');
    expect((params as HttpParams).get('spaceID')).toBe('s1');
    expect((params as HttpParams).get('assetID')).toBe('a1');
  });

  it('updateAsset posts to update_asset with the request', () => {
    const request = {
      spaceID: 's1',
      assetID: 'a1',
      name: 'Car',
      category: 'vehicles' as const,
      condition: 'good' as const,
      visibility: 'private' as const,
    };
    service.updateAsset(request).subscribe();
    expect(post).toHaveBeenCalledWith('assetus/update_asset', request);
  });

  it('removeAsset posts to remove_asset with the request', () => {
    const request = { spaceID: 's1', assetID: 'a1', hardDelete: true };
    service.removeAsset(request).subscribe();
    expect(post).toHaveBeenCalledWith('assetus/remove_asset', request);
  });

  it('transferAsset posts to transfer_asset with the request', () => {
    const request = { spaceID: 's1', assetID: 'a1', toSpaceID: 's2' };
    service.transferAsset(request).subscribe();
    expect(post).toHaveBeenCalledWith('assetus/transfer_asset', request);
  });

  it('recordHistoryEvent posts to record_history_event with the request', () => {
    const request = {
      spaceID: 's1',
      assetID: 'a1',
      type: 'repaired' as const,
      note: 'fixed',
    };
    service.recordHistoryEvent(request).subscribe();
    expect(post).toHaveBeenCalledWith('assetus/record_history_event', request);
  });

  it('getHistory gets asset_history with spaceID+assetID params', () => {
    service.getHistory('s1', 'a1').subscribe();
    expect(get).toHaveBeenCalledTimes(1);
    const [endpoint, params] = get.mock.calls[0];
    expect(endpoint).toBe('assetus/asset_history');
    expect((params as HttpParams).get('spaceID')).toBe('s1');
    expect((params as HttpParams).get('assetID')).toBe('a1');
  });

  it('addVehicleRecord posts to create_vehicle_record with the request', () => {
    const request = {
      spaceID: 's1',
      assetID: 'a1',
      fuelVolume: 40,
      fuelVolumeUnit: 'l' as const,
      fuelCost: 60,
      currency: 'USD',
      mileage: 12345,
      mileageUnit: 'km' as const,
    };
    service.addVehicleRecord(request).subscribe();
    expect(post).toHaveBeenCalledWith('assetus/create_vehicle_record', request);
  });
});
