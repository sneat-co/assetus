import { TestBed } from '@angular/core/testing';
import { spacePageTestProviders } from '../../../testing/test-providers';
import { RealEstatePageComponent } from './real-estate-page.component';

// Render spec for the ported RealEstatePageComponent. It extends
// SpaceBaseComponent and reads its asset from navigation state.
describe('RealEstatePageComponent', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      imports: [RealEstatePageComponent],
      providers: spacePageTestProviders(),
    }),
  );

  it('creates and renders the Property fallback title', () => {
    const fixture = TestBed.createComponent(RealEstatePageComponent);
    fixture.detectChanges();
    const host = fixture.nativeElement as HTMLElement;
    expect(fixture.componentInstance).toBeTruthy();
    // With no asset in nav state, the static "Property" header is shown.
    expect(host.textContent).toContain('Property');
  });
});
