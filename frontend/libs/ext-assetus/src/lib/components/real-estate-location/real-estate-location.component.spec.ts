import { TestBed } from '@angular/core/testing';
import { RealEstateLocationComponent } from './real-estate-location.component';

// Smoke/render spec for the ported RealEstateLocationComponent. It has no
// injected dependencies, so a bare TestBed render proves the standalone
// component + its Ionic template compile and instantiate.
describe('RealEstateLocationComponent', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      imports: [RealEstateLocationComponent],
    }),
  );

  it('creates and renders', () => {
    const fixture = TestBed.createComponent(RealEstateLocationComponent);
    fixture.detectChanges();
    expect(fixture.componentInstance).toBeTruthy();
  });
});
