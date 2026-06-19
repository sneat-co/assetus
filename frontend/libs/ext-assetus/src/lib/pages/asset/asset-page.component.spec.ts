import { TestBed } from '@angular/core/testing';
import { Observable, of } from 'rxjs';
import { spacePageTestProviders } from '../../../testing/test-providers';
import { IAssetContext } from '../../contexts';
import { IAssetDbo } from '../../dto';
import { AssetService } from '../../services';
import { AssetPageComponent } from './asset-page.component';

// Render spec for AssetPageComponent. The page now reads the asset live from
// Firestore via AssetService.watchAssetByID(space, id) (replacing the one-shot
// API getAsset). A stub AssetService emits an IAssetContext so we can assert the
// page applies ctx.dbo to its form fields.
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
    // The page imports AssetusCoreServicesModule, which provides AssetService at
    // the component scope; override it at the component level so the page injects
    // this stub (a root-level override would be shadowed by that module).
    TestBed.overrideComponent(AssetPageComponent, {
      add: { providers: [{ provide: AssetService, useValue: serviceStub }] },
    });
  });

  it('creates', () => {
    const fixture = TestBed.createComponent(AssetPageComponent);
    expect(fixture.componentInstance).toBeTruthy();
  });

  it('loadAsset subscribes to watchAssetByID and applies ctx.dbo to the form', () => {
    const fixture = TestBed.createComponent(AssetPageComponent);
    const cmp = fixture.componentInstance as unknown as {
      $spaceRef: { set(ref: { id: string }): void };
      space: { id?: string };
      loadAsset(assetID: string): void;
      name: string;
      description: string;
      category: string;
      condition: string;
      visibility: string;
      $asset(): IAssetDbo | undefined;
    };
    cmp.$spaceRef.set({ id: 's1' });

    cmp.loadAsset('a1');

    expect(watchAssetByID).toHaveBeenCalledTimes(1);
    const [spaceArg, idArg] = watchAssetByID.mock.calls[0];
    expect(spaceArg.id).toBe('s1');
    expect(idArg).toBe('a1');
    // The page applied ctx.dbo to its signal and form fields.
    expect(cmp.$asset()).toBe(dbo);
    expect(cmp.name).toBe('My Car');
    expect(cmp.description).toBe('Daily driver');
    expect(cmp.category).toBe('vehicles');
    expect(cmp.condition).toBe('fair');
    expect(cmp.visibility).toBe('family');
  });
});
