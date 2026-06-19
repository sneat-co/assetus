import { TestBed } from '@angular/core/testing';
import { spacePageTestProviders } from '../../../testing/test-providers';
import { OptimizationPageComponent } from './optimization-page.component';

// Render spec for the ported OptimizationPageComponent. The page extends
// SpaceBaseComponent, so the standard Sneat DI chain must resolve for it to
// create and render its static "Savings" card.
describe('OptimizationPageComponent', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      imports: [OptimizationPageComponent],
      providers: spacePageTestProviders(),
    }),
  );

  it('creates and renders the savings page', () => {
    const fixture = TestBed.createComponent(OptimizationPageComponent);
    fixture.detectChanges();
    const host = fixture.nativeElement as HTMLElement;
    expect(fixture.componentInstance).toBeTruthy();
    expect(host.querySelector('ion-title')?.textContent).toContain('Savings');
  });
});
