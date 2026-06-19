import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { AssetCardComponent } from './asset-card.component';

// Render + logic spec for the ported AssetCardComponent. It uses RouterModule
// for its [routerLink], so a router must be provided.
describe('AssetCardComponent', () => {
  let fixture: ReturnType<typeof TestBed.createComponent<AssetCardComponent>>;
  let component: AssetCardComponent;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [AssetCardComponent],
      providers: [provideRouter([])],
    });
    fixture = TestBed.createComponent(AssetCardComponent);
    component = fixture.componentInstance;
  });

  it('creates and renders the asset title', () => {
    component.asset = {
      id: 'a1',
      space: { id: 's1' },
      dbo: { category: 'vehicle', title: 'My Car' },
    } as never;
    fixture.detectChanges();
    const host = fixture.nativeElement as HTMLElement;
    expect(host.querySelector('ion-card-title')?.textContent).toContain('My Car');
  });

  it('switches the default segment to income when incomes outnumber expenses', () => {
    component.asset = {
      id: 'a1',
      space: { id: 's1' },
      dbo: {
        category: 'vehicle',
        totals: { incomes: { count: 5 }, expenses: { count: 1 } },
      },
    } as never;
    component.ngOnChanges({ asset: {} as never });
    expect((component as unknown as { segment: string }).segment).toBe('income');
  });

  it('segmentChanged updates the active segment', () => {
    component.segmentChanged({ detail: { value: 'income' } } as CustomEvent);
    expect((component as unknown as { segment: string }).segment).toBe('income');
  });
});
