import { Period } from '@sneat/dto';
import { PeriodSegmentComponent } from './period-segment.component';

// Smoke + logic spec for the ported PeriodSegmentComponent. The component has no
// injected dependencies, so it can be exercised without a TestBed (matching the
// pure-spec style already used in this lib).
describe('PeriodSegmentComponent', () => {
  let component: PeriodSegmentComponent;

  beforeEach(() => {
    component = new PeriodSegmentComponent();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('segmentChanged updates period and emits the change', () => {
    const newPeriod: Period = 'month';
    const changedSpy = vi.fn();
    component.changed.subscribe(changedSpy);

    component.segmentChanged({ detail: { value: newPeriod } } as CustomEvent);

    expect(component.period).toBe(newPeriod);
    expect(changedSpy).toHaveBeenCalledWith(newPeriod);
  });
});
