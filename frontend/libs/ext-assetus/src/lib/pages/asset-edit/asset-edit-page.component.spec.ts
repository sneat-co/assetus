import { TestBed } from '@angular/core/testing';
import { Observable, of } from 'rxjs';
import { spacePageTestProviders } from '../../../testing/test-providers';
import { IAssetContext } from '@sneat/extension-assetus-contract';
import { IAssetDbo } from '@sneat/extension-assetus-contract';
import { AssetService } from '../../services';
import { AssetEditPageComponent } from './asset-edit-page.component';

// Render + behaviour spec for AssetEditPageComponent (the manage page). It loads
// the asset ONCE for prefill via watchAssetByID(space, id).pipe(take(1)), edits
// it via updateAsset, and hosts archive/hard-delete via removeAsset. A stub
// AssetService is provided at the component scope (AssetusCoreServicesModule
// provides AssetService there), mirroring the view-page spec.
describe('AssetEditPageComponent', () => {
  const dbo: IAssetDbo = {
    name: 'My Car',
    description: 'Daily driver',
    category: 'vehicles',
    condition: 'fair',
    visibility: 'family',
  } as IAssetDbo;

  let watchAssetByID: ReturnType<typeof vi.fn>;
  let updateAsset: ReturnType<typeof vi.fn>;
  let removeAsset: ReturnType<typeof vi.fn>;
  let serviceStub: {
    watchAssetByID: typeof watchAssetByID;
    updateAsset: typeof updateAsset;
    removeAsset: typeof removeAsset;
  };

  beforeEach(() => {
    watchAssetByID = vi.fn(
      (space: { id?: string }, id: string): Observable<IAssetContext> =>
        of({ id, space, dbo } as IAssetContext),
    );
    updateAsset = vi.fn(() => of({ asset: dbo }));
    removeAsset = vi.fn(() => of(undefined));
    serviceStub = { watchAssetByID, updateAsset, removeAsset };
    TestBed.configureTestingModule({
      imports: [AssetEditPageComponent],
      providers: [...spacePageTestProviders()],
    });
    TestBed.overrideComponent(AssetEditPageComponent, {
      add: { providers: [{ provide: AssetService, useValue: serviceStub }] },
    });
  });

  type EditPage = {
    $spaceRef: { set(ref: { id: string }): void };
    space: { id?: string };
    $assetID: { set(id: string): void };
    loadAsset(assetID: string): void;
    save(): void;
    archive(): void;
    remove(hardDelete: boolean): void;
    name: string;
    description: string;
    category: string;
    condition: string;
    visibility: string;
  };

  it('creates', () => {
    const fixture = TestBed.createComponent(AssetEditPageComponent);
    expect(fixture.componentInstance).toBeTruthy();
  });

  it('loadAsset reads once via watchAssetByID and prefills the form', () => {
    const fixture = TestBed.createComponent(AssetEditPageComponent);
    const cmp = fixture.componentInstance as unknown as EditPage;
    cmp.$spaceRef.set({ id: 's1' });

    cmp.loadAsset('a1');

    expect(watchAssetByID).toHaveBeenCalledTimes(1);
    const [spaceArg, idArg] = watchAssetByID.mock.calls[0];
    expect(spaceArg.id).toBe('s1');
    expect(idArg).toBe('a1');
    expect(cmp.name).toBe('My Car');
    expect(cmp.description).toBe('Daily driver');
    expect(cmp.category).toBe('vehicles');
    expect(cmp.condition).toBe('fair');
    expect(cmp.visibility).toBe('family');
  });

  it('save calls updateAsset with the form values', () => {
    const fixture = TestBed.createComponent(AssetEditPageComponent);
    const cmp = fixture.componentInstance as unknown as EditPage;
    cmp.$spaceRef.set({ id: 's1' });
    cmp.$assetID.set('a1');
    cmp.name = 'Updated';
    cmp.description = 'New desc';
    cmp.category = 'tools';
    cmp.condition = 'good';
    cmp.visibility = 'private';

    cmp.save();

    expect(updateAsset).toHaveBeenCalledTimes(1);
    expect(updateAsset).toHaveBeenCalledWith({
      spaceID: 's1',
      assetID: 'a1',
      name: 'Updated',
      description: 'New desc',
      category: 'tools',
      condition: 'good',
      visibility: 'private',
    });
  });

  it('archive calls removeAsset with hardDelete:false', () => {
    const fixture = TestBed.createComponent(AssetEditPageComponent);
    const cmp = fixture.componentInstance as unknown as EditPage;
    cmp.$spaceRef.set({ id: 's1' });
    cmp.$assetID.set('a1');

    cmp.archive();

    expect(removeAsset).toHaveBeenCalledTimes(1);
    expect(removeAsset).toHaveBeenCalledWith({
      spaceID: 's1',
      assetID: 'a1',
      hardDelete: false,
    });
  });

  it('remove(true) calls removeAsset with hardDelete:true', () => {
    const fixture = TestBed.createComponent(AssetEditPageComponent);
    const cmp = fixture.componentInstance as unknown as EditPage;
    cmp.$spaceRef.set({ id: 's1' });
    cmp.$assetID.set('a1');

    cmp.remove(true);

    expect(removeAsset).toHaveBeenCalledTimes(1);
    expect(removeAsset).toHaveBeenCalledWith({
      spaceID: 's1',
      assetID: 'a1',
      hardDelete: true,
    });
  });
});
