import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { spacePageTestProviders } from '../../../testing/test-providers';
import { AssetGroupPageComponent } from './asset-group-page.component';

// Render spec for the ported AssetGroupPageComponent. It extends
// SpaceBaseComponent and embeds the AssetCardComponent (RouterModule) +
// PeriodSegmentComponent.
describe('AssetGroupPageComponent', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      imports: [AssetGroupPageComponent],
      providers: [...spacePageTestProviders(), provideRouter([])],
    }),
  );

  it('creates and defaults the period to month', () => {
    const fixture = TestBed.createComponent(AssetGroupPageComponent);
    fixture.detectChanges();
    expect(fixture.componentInstance).toBeTruthy();
    expect(
      (fixture.componentInstance as unknown as { period: string }).period,
    ).toBe('month');
  });

  it('periodChanged updates the active period', () => {
    const fixture = TestBed.createComponent(AssetGroupPageComponent);
    (fixture.componentInstance as unknown as {
      periodChanged(p: string): void;
    }).periodChanged('year');
    expect(
      (fixture.componentInstance as unknown as { period: string }).period,
    ).toBe('year');
  });
});
