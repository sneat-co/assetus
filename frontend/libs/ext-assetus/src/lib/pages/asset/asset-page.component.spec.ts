import { TestBed } from '@angular/core/testing';
import { Observable, of } from 'rxjs';
import { spacePageTestProviders } from '../../../testing/test-providers';
import { IAssetContext } from '@sneat/extension-assetus-contract';
import { IAssetDbo } from '@sneat/extension-assetus-contract';
import { ASSET_SERVICE } from '@sneat/extension-assetus-contract';
import { AssetPageComponent } from './asset-page.component';

// Render spec for AssetPageComponent. This page is now read-only: it reads the
// asset live from Firestore via AssetService.watchAssetByID(space, id) and
// applies ctx.dbo to its $asset signal (no editable form, no save). A stub
// AssetService emits an IAssetContext so we can assert the live load.
describe('AssetPageComponent', () => {
  const dbo: IAssetDbo = {
    name: 'My Car',
    description: 'Daily driver',
    category: 'vehicles',
    condition: 'fair',
    visibility: 'family',
  } as IAssetDbo;

  let watchAssetByID: ReturnType<typeof vi.fn>;
  let serviceStub: { watchAssetByID: typeof watchAssetByID };

  beforeEach(() => {
    watchAssetByID = vi.fn(
      (space: { id?: string }, id: string): Observable<IAssetContext> =>
        of({ id, space, dbo } as IAssetContext),
    );
    serviceStub = { watchAssetByID };
    TestBed.configureTestingModule({
      imports: [AssetPageComponent],
      providers: [...spacePageTestProviders()],
    });
    // The page injects the ASSET_SERVICE token; provide a stub at the component
    // level so the page resolves this fake instead of the real implementation.
    TestBed.overrideComponent(AssetPageComponent, {
      add: { providers: [{ provide: ASSET_SERVICE, useValue: serviceStub }] },
    });
  });

  it('creates', () => {
    const fixture = TestBed.createComponent(AssetPageComponent);
    expect(fixture.componentInstance).toBeTruthy();
  });

  it('loadAsset subscribes to watchAssetByID and applies ctx.dbo to $asset', () => {
    const fixture = TestBed.createComponent(AssetPageComponent);
    const cmp = fixture.componentInstance as unknown as {
      $spaceRef: { set(ref: { id: string }): void };
      space: { id?: string };
      loadAsset(assetID: string): void;
      $asset(): IAssetDbo | undefined;
      $categoryLabel(): string;
      $conditionLabel(): string;
      $visibilityLabel(): string;
    };
    cmp.$spaceRef.set({ id: 's1' });

    cmp.loadAsset('a1');

    expect(watchAssetByID).toHaveBeenCalledTimes(1);
    const [spaceArg, idArg] = watchAssetByID.mock.calls[0];
    expect(spaceArg.id).toBe('s1');
    expect(idArg).toBe('a1');
    // The page applied ctx.dbo to its signal and derives human labels.
    expect(cmp.$asset()).toBe(dbo);
    expect(cmp.$categoryLabel()).toBe('Vehicles');
    expect(cmp.$conditionLabel()).toBe('Fair');
    expect(cmp.$visibilityLabel()).toBe('Family');
  });

  it('is read-only: exposes no editable save behaviour', () => {
    const fixture = TestBed.createComponent(AssetPageComponent);
    const cmp = fixture.componentInstance as unknown as Record<string, unknown>;
    expect(cmp['save']).toBeUndefined();
    expect(cmp['archive']).toBeUndefined();
    expect(cmp['transfer']).toBeUndefined();
    expect(cmp['remove']).toBeUndefined();
  });
});
