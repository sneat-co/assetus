import { TestBed } from '@angular/core/testing';
import { ModalController, provideIonicAngular } from '@ionic/angular/standalone';
import { ASSET_SERVICE } from '@sneat/extension-assetus-contract';
import { SpaceNavService } from '@sneat/space-services';
import { of } from 'rxjs';
import { componentTestProviders } from '../../../testing/test-providers';
import { AssetsListComponent } from './assets-list.component';

// Render + logic spec for the ported AssetsListComponent. It injects the lib's
// AssetService, the SpaceNavService, the ModalController and the ErrorLogger.
describe('AssetsListComponent', () => {
  let navigate: ReturnType<typeof vi.fn>;
  let removeAsset: ReturnType<typeof vi.fn>;
  let fixture: ReturnType<typeof TestBed.createComponent<AssetsListComponent>>;
  let component: AssetsListComponent;

  const assets = [
    { id: 'a1', brief: { category: 'dwelling', name: 'House' } },
    { id: 'a2', brief: { category: 'vehicles', name: 'Car' } },
  ];

  // Rich fixture that exercises every template branch: a vehicle with
  // make/model + regNumber badge, and a dwelling with address + yearOfBuild.
  // Shapes mirror the live IAssetDbo (name, the 'vehicles' category, and the
  // typed extra on the dbo).
  const richAssets = [
    {
      id: 'v1',
      brief: {
        category: 'vehicles',
        extraType: 'vehicle',
        name: 'My Car',
        extra: { make: 'Toyota', model: 'Corolla', regNumber: 'ABC123' },
      },
    },
    {
      id: 'd1',
      brief: {
        category: 'dwelling',
        extraType: 'dwelling',
        name: 'Home',
        yearOfBuild: 1990,
        extra: { address: '1 Main St' },
      },
    },
  ];

  beforeEach(() => {
    navigate = vi.fn(() => Promise.resolve(true));
    removeAsset = vi.fn(() => of(undefined));
    TestBed.configureTestingModule({
      imports: [AssetsListComponent],
      providers: [
        ...componentTestProviders(),
        provideIonicAngular(),
        { provide: ASSET_SERVICE, useValue: { removeAsset } },
        {
          provide: SpaceNavService,
          useValue: { navigateForwardToSpacePage: navigate },
        },
      ],
    });
    fixture = TestBed.createComponent(AssetsListComponent);
    component = fixture.componentInstance;
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllGlobals();
  });

  it('creates and renders', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('filters the list to the requested asset type', () => {
    component.allAssets = assets as never;
    component.assetType = 'dwelling' as never;
    component.ngOnChanges({ allAssets: {} as never });
    const filtered = (
      component as unknown as { assets: { id: string }[] }
    ).assets;
    expect(filtered.map((a) => a.id)).toEqual(['a1']);
  });

  it('leaves assets undefined when no source list is provided', () => {
    component.allAssets = undefined;
    component.ngOnChanges({ allAssets: {} as never });
    expect(
      (component as unknown as { assets?: unknown[] }).assets,
    ).toBeUndefined();
  });

  it('navigates to the asset page on select', () => {
    component.space = { id: 's1' } as never;
    (component as unknown as { goAsset(a: unknown): void }).goAsset(assets[0]);
    expect(navigate).toHaveBeenCalled();
  });

  it('keeps the full list when no filter or assetType is set', () => {
    component.allAssets = assets as never;
    component.ngOnChanges({ allAssets: {} as never });
    expect(
      (component as unknown as { assets: { id: string }[] }).assets.map(
        (a) => a.id,
      ),
    ).toEqual(['a1', 'a2']);
  });

  it('renders the badges, make/model and action buttons for each asset', () => {
    component.allAssets = richAssets as never;
    component.space = { id: 's1' } as never;
    component.ngOnChanges({ allAssets: {} as never });
    fixture.detectChanges();

    const host = fixture.nativeElement as HTMLElement;
    expect(host.querySelectorAll('ion-item').length).toBe(2);
    // Vehicle make/model + reg-number badge.
    expect(host.textContent).toContain('Toyota');
    expect(host.textContent).toContain('ABC123');
    // Dwelling address + year-of-build badge.
    expect(host.textContent).toContain('1 Main St');
    expect(host.textContent).toContain('1990');
    // The "Record miles" button only renders for the vehicle row.
    expect(host.textContent).toContain('Record miles');
  });

  it('renders the loading placeholder when assets are undefined', () => {
    fixture.detectChanges();
    const host = fixture.nativeElement as HTMLElement;
    expect(host.textContent).toContain('Loading...');
  });

  it('renders the empty placeholder when the list is empty', () => {
    component.allAssets = [] as never;
    component.ngOnChanges({ allAssets: {} as never });
    fixture.detectChanges();
    const host = fixture.nativeElement as HTMLElement;
    expect(host.textContent).toContain('No items created yet');
  });

  it('deletes an asset after the user confirms', () => {
    vi.useFakeTimers();
    vi.stubGlobal('confirm', () => true);
    component.allAssets = assets as never;
    component.space = { id: 's1' } as never;
    component.ngOnChanges({ allAssets: {} as never });

    const event = {
      stopPropagation: vi.fn(),
      preventDefault: vi.fn(),
    } as unknown as Event;
    (
      component as unknown as { delete(e: Event, a: unknown): void }
    ).delete(event, assets[0]);
    vi.runAllTimers();

    expect(removeAsset).toHaveBeenCalledWith({ spaceID: 's1', assetID: 'a1' });
    expect(
      (component as unknown as { assets: { id: string }[] }).assets.map(
        (a) => a.id,
      ),
    ).toEqual(['a2']);
  });

  it('does not delete when the user cancels the confirm', () => {
    vi.useFakeTimers();
    vi.stubGlobal('confirm', () => false);
    component.allAssets = assets as never;
    component.space = { id: 's1' } as never;
    component.ngOnChanges({ allAssets: {} as never });

    const event = {
      stopPropagation: vi.fn(),
      preventDefault: vi.fn(),
    } as unknown as Event;
    (
      component as unknown as { delete(e: Event, a: unknown): void }
    ).delete(event, assets[0]);
    vi.runAllTimers();

    expect(removeAsset).not.toHaveBeenCalled();
  });

  it('opens the mileage dialog modal', async () => {
    const present = vi.fn(() => Promise.resolve());
    // The component declares its own ModalController provider, so resolve it
    // through the component's injector rather than the root TestBed injector.
    const modalCtrl = fixture.debugElement.injector.get(ModalController);
    vi.spyOn(modalCtrl, 'create').mockResolvedValue({
      present,
    } as never);

    const event = {
      stopPropagation: vi.fn(),
      preventDefault: vi.fn(),
    } as unknown as Event;
    await (
      component as unknown as {
        addNewMilesAndFuel(e: Event, a: unknown): Promise<void>;
      }
    ).addNewMilesAndFuel(event, assets[1]);

    expect(present).toHaveBeenCalled();
  });
});
