import { TestBed } from '@angular/core/testing';
import { AssetService } from '@sneat/ext-assetus-components';
import { SpaceNavService } from '@sneat/space-services';
import { spacePageTestProviders } from '../../../testing/test-providers';
import { RealEstatesPageComponent } from './real-estates-page.component';

// Render spec for the ported RealEstatesPageComponent. It extends
// SpaceBaseComponent and embeds AssetsListComponent (legacy AssetService +
// SpaceNavService + ModalController).
describe('RealEstatesPageComponent', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      imports: [RealEstatesPageComponent],
      providers: [
        ...spacePageTestProviders(),
        { provide: AssetService, useValue: { deleteAsset: vi.fn() } },
        {
          provide: SpaceNavService,
          useValue: { navigateForwardToSpacePage: vi.fn(() => Promise.resolve(true)) },
        },
      ],
    }),
  );

  it('creates and renders the Properties title', () => {
    const fixture = TestBed.createComponent(RealEstatesPageComponent);
    fixture.detectChanges();
    const host = fixture.nativeElement as HTMLElement;
    expect(host.querySelector('ion-title')?.textContent).toContain('Properties');
  });

  it('clearFilter resets the filter and removes the clear button', () => {
    const fixture = TestBed.createComponent(RealEstatesPageComponent);
    const component = fixture.componentInstance as unknown as {
      filter: string;
      clearFilter(): void;
    };
    component.filter = 'house';
    fixture.detectChanges();
    expect(
      (fixture.nativeElement as HTMLElement).querySelector(
        'ion-button[fill="clear"]',
      ),
    ).toBeTruthy();

    component.clearFilter();
    expect(component.filter).toBe('');
  });

  it('goNew navigates forward to the new-asset page', () => {
    const fixture = TestBed.createComponent(RealEstatesPageComponent);
    fixture.detectChanges();
    expect(() =>
      (fixture.componentInstance as unknown as { goNew(): void }).goNew(),
    ).not.toThrow();
  });
});
