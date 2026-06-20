import { TestBed } from '@angular/core/testing';
import { ASSET_SERVICE } from '@sneat/extension-assetus-contract';
import { spacePageTestProviders } from '../../../testing/test-providers';
import { NewAssetPageComponent } from './new-asset-page.component';

// Render spec for the ported NewAssetPageComponent. It extends
// SpaceBaseComponent and embeds the asset-add-* components (AssetService),
// so it needs the standard chain plus a stub service.
describe('NewAssetPageComponent', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      imports: [NewAssetPageComponent],
      providers: [
        ...spacePageTestProviders(),
        { provide: ASSET_SERVICE, useValue: { createAsset: vi.fn() } },
      ],
    }),
  );

  it('creates and renders the new-asset title with category options', () => {
    const fixture = TestBed.createComponent(NewAssetPageComponent);
    fixture.detectChanges();
    const host = fixture.nativeElement as HTMLElement;
    expect(host.querySelector('ion-title')?.textContent).toContain('New asset');
    expect(fixture.componentInstance.categories.length).toBeGreaterThan(0);
  });

  it('selectCategory sets the active category signal', () => {
    const fixture = TestBed.createComponent(NewAssetPageComponent);
    const category = fixture.componentInstance.categories[0];
    fixture.componentInstance.selectCategory(category);
    expect(
      (fixture.componentInstance as unknown as {
        category: { (): unknown };
      }).category(),
    ).toBe(category);
  });

  it('renders the asset-add host once a category is selected', () => {
    const fixture = TestBed.createComponent(NewAssetPageComponent);
    const vehicle = fixture.componentInstance.categories.find(
      (c) => c.id === 'vehicles',
    );
    fixture.componentInstance.selectCategory(vehicle as never);
    fixture.detectChanges();
    const host = fixture.nativeElement as HTMLElement;
    // The category picker list is gone; the vehicle add host is rendered.
    expect(host.querySelector('assetus-asset-add-vehicle')).toBeTruthy();
    expect(host.textContent).not.toContain('Select asset kind');
  });

  it('renders the dwelling add host for the real-estate category', () => {
    const fixture = TestBed.createComponent(NewAssetPageComponent);
    const dwelling = fixture.componentInstance.categories.find(
      (c) => c.id === 'dwelling',
    );
    fixture.componentInstance.selectCategory(dwelling as never);
    fixture.detectChanges();
    const host = fixture.nativeElement as HTMLElement;
    expect(host.querySelector('assetus-asset-add-dwelling')).toBeTruthy();
  });
});
