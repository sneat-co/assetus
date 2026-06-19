import { HttpParams } from '@angular/common/http';
import { Injector, runInInjectionContext } from '@angular/core';
import { Firestore, Timestamp } from '@angular/fire/firestore';
import { SneatApiService } from '@sneat/api';
import { of } from 'rxjs';
import { AssetService } from './asset.service';

// getHistory now reads Firestore directly via the free functions collection/
// query/orderBy/collectionData, so this spec mocks @angular/fire/firestore. The
// vitest.config for this lib runs the suite with isolated forks specifically so
// this module mock cannot leak into sibling specs that use the real Firestore
// (a leak made the suite flaky). The mock keeps the real Firestore DI token and
// Timestamp class so providers and `instanceof Timestamp` still behave; only the
// reads are stubbed. vi.mock is hoisted, so its spies come from vi.hoisted.
const fs = vi.hoisted(() => {
  const collection = vi.fn(() => ({ __type: 'collection' }));
  const query = vi.fn((c: unknown, ...constraints: unknown[]) => ({
    __type: 'query',
    collection: c,
    constraints,
  }));
  const orderBy = vi.fn((field: string, dir: string) => ({
    __type: 'orderBy',
    field,
    dir,
  }));
  // Swappable docs that collectionData emits per test.
  const state: { docs: unknown[] } = { docs: [] };
  const collectionData = vi.fn(() => of(state.docs));
  // Stand-in Firestore token + Timestamp class. The spec never imports the real
  // @angular/fire module (it can't JIT-compile under vitest); both the service
  // and this spec import Firestore/Timestamp from the same mocked module, so
  // their identities match.
  class Firestore {}
  class Timestamp {
    constructor(private readonly date: Date) {}
    static fromDate(d: Date): Timestamp {
      return new Timestamp(d);
    }
    toDate(): Date {
      return this.date;
    }
  }
  return { collection, query, orderBy, collectionData, state, Firestore, Timestamp };
});

vi.mock('@angular/fire/firestore', () => ({
  Firestore: fs.Firestore,
  Timestamp: fs.Timestamp,
  collection: fs.collection,
  query: fs.query,
  orderBy: fs.orderBy,
  collectionData: fs.collectionData,
  doc: vi.fn(),
  docData: vi.fn(),
}));

describe('AssetService', () => {
  let service: AssetService;
  let post: ReturnType<typeof vi.fn>;
  let get: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fs.collection.mockClear();
    fs.query.mockClear();
    fs.orderBy.mockClear();
    fs.collectionData.mockClear();
    fs.state.docs = [];
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

  it('createAsset forwards the rich superset fields (incl. vehicle extra + status=draft)', () => {
    const request = {
      spaceID: 's1',
      name: 'Car',
      category: 'vehicles' as const,
      condition: 'good' as const,
      visibility: 'private' as const,
      status: 'draft' as const,
      // rich superset fields
      type: 'car' as const,
      possession: 'owning' as const,
      countryID: 'IE',
      yearOfBuild: 2020,
      isRequest: false,
      geo: { lat: 53.3, lng: -6.2 },
      dateOfPurchase: '2020-01-15',
      totals: [{ currency: 'EUR', value: 15000 }],
      canHaveExpense: true,
      financialDirection: 'expense',
      groupID: 'g1',
      subAssets: [{ id: 'sa1', title: 'Spare', type: 'vehicles' as const }],
      // typed extra
      extraType: 'vehicle' as const,
      extra: { make: 'Toyota', model: 'Corolla', regNumber: '12-D-3456' },
    };
    service.createAsset(request).subscribe();
    expect(post).toHaveBeenCalledWith('assetus/create_asset', request);
    const [, sent] = post.mock.calls[0];
    expect(sent.extraType).toBe('vehicle');
    expect(sent.extra).toEqual({
      make: 'Toyota',
      model: 'Corolla',
      regNumber: '12-D-3456',
    });
    expect(sent.status).toBe('draft');
    expect(sent.type).toBe('car');
    expect(sent.totals).toEqual([{ currency: 'EUR', value: 15000 }]);
  });

  it('createAsset still works with only the flat MVP fields', () => {
    const request = {
      spaceID: 's1',
      name: 'Book',
      category: 'books' as const,
      condition: 'new' as const,
    };
    service.createAsset(request).subscribe();
    expect(post).toHaveBeenCalledWith('assetus/create_asset', request);
  });

  it('updateAsset forwards the rich superset fields but carries NO status', () => {
    const request = {
      spaceID: 's1',
      assetID: 'a1',
      name: 'Car',
      category: 'vehicles' as const,
      condition: 'good' as const,
      visibility: 'private' as const,
      // rich superset fields
      type: 'car' as const,
      possession: 'leasing' as const,
      parentCategoryID: 'vehicles' as const,
      liabilities: [{ id: 'l1', serviceTypes: ['insurance'] }],
      notUsedServiceTypes: ['tax'],
      extraType: 'vehicle' as const,
      extra: { vin: 'WVWZZZ1JZXW000001' },
    };
    service.updateAsset(request).subscribe();
    expect(post).toHaveBeenCalledWith('assetus/update_asset', request);
    const [, sent] = post.mock.calls[0];
    expect(sent.extraType).toBe('vehicle');
    expect(sent.extra).toEqual({ vin: 'WVWZZZ1JZXW000001' });
    expect('status' in sent).toBe(false);
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

  it('getHistory reads the asset history subcollection from Firestore ordered by occurredAt asc (no API call)', () => {
    service.getHistory('s1', 'a1').subscribe();

    // Reads spaces/s1/ext/assetus/assets/a1/history, not the backend.
    expect(get).not.toHaveBeenCalled();
    expect(fs.collection).toHaveBeenCalledWith(
      {}, // injected stub Firestore
      'spaces',
      's1',
      'ext',
      'assetus',
      'assets',
      'a1',
      'history',
    );
    expect(fs.orderBy).toHaveBeenCalledWith('occurredAt', 'asc');
    // The query is built from the history collection + the orderBy constraint.
    const [collArg] = fs.query.mock.calls[0];
    expect(collArg).toEqual({ __type: 'collection' });
    expect(fs.collectionData).toHaveBeenCalledWith(
      expect.objectContaining({ __type: 'query' }),
      { idField: 'id' },
    );
  });

  it('getHistory returns { assetID, events } and converts a Firestore Timestamp occurredAt to an ISO string', () => {
    fs.state.docs = [
      {
        id: 'h1',
        type: 'repaired',
        occurredAt: Timestamp.fromDate(new Date('2024-03-04T05:06:07.000Z')),
        actorRef: 'u1',
        note: 'fixed',
      },
    ];

    let result:
      | { assetID: string; events: { occurredAt: string }[] }
      | undefined;
    service.getHistory('s1', 'a1').subscribe((r) => (result = r as never));

    expect(result?.assetID).toBe('a1');
    expect(result?.events).toHaveLength(1);
    // Timestamp -> ISO string (matches IHistoryEvent.occurredAt: string).
    expect(result?.events[0].occurredAt).toBe('2024-03-04T05:06:07.000Z');
    expect(typeof result?.events[0].occurredAt).toBe('string');
  });

  it('getHistory passes through an occurredAt that is already an ISO string', () => {
    fs.state.docs = [
      {
        id: 'h1',
        type: 'repaired',
        occurredAt: '2024-01-02T03:04:05.000Z',
        actorRef: 'u1',
      },
    ];

    let result: { events: { occurredAt: string }[] } | undefined;
    service.getHistory('s1', 'a1').subscribe((r) => (result = r as never));

    expect(result?.events[0].occurredAt).toBe('2024-01-02T03:04:05.000Z');
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

  // watchAssetByID's required-spaceID guard is reachable without touching
  // Firestore; its happy path (a Firestore read) is left to integration cover.
  it('watchAssetByID throws when the space has no id', () => {
    expect(() =>
      service.watchAssetByID({ id: '' } as never, 'a1'),
    ).toThrowError('spaceID is required');
  });
});
