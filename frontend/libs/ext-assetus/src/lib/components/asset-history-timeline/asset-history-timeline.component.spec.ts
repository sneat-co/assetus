import { TestBed } from '@angular/core/testing';
import { of } from 'rxjs';
import { AssetService } from '../../services';
import { componentTestProviders } from '../../../testing/test-providers';
import { AssetHistoryTimelineComponent } from './asset-history-timeline.component';

// Render + logic spec for the AssetHistoryTimelineComponent (Task 11). It
// injects the MVP AssetService and the ErrorLogger.
describe('AssetHistoryTimelineComponent', () => {
  let getHistory: ReturnType<typeof vi.fn>;
  let fixture: ReturnType<
    typeof TestBed.createComponent<AssetHistoryTimelineComponent>
  >;
  let component: AssetHistoryTimelineComponent;

  beforeEach(() => {
    getHistory = vi.fn(() =>
      of({ events: [{ type: 'transferred', at: '2026-01-01' }] }),
    );
    TestBed.configureTestingModule({
      imports: [AssetHistoryTimelineComponent],
      providers: [
        ...componentTestProviders(),
        { provide: AssetService, useValue: { getHistory } },
      ],
    });
    fixture = TestBed.createComponent(AssetHistoryTimelineComponent);
    component = fixture.componentInstance;
  });

  it('creates and renders', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('loads history events when both ids are present', () => {
    component.spaceID = 's1';
    component.assetID = 'a1';
    component.ngOnChanges();
    expect(getHistory).toHaveBeenCalledWith('s1', 'a1');
    expect(
      (component as unknown as { $events: { (): unknown[] } }).$events(),
    ).toHaveLength(1);
  });

  it('clears events and skips the call when ids are missing', () => {
    component.spaceID = undefined;
    component.ngOnChanges();
    expect(getHistory).not.toHaveBeenCalled();
    expect(
      (component as unknown as { $events: { (): unknown[] } }).$events(),
    ).toEqual([]);
  });

  it('maps history event types to icon names', () => {
    const iconFor = (component as unknown as {
      iconFor(t: string): string;
    }).iconFor.bind(component);
    expect(iconFor('purchased')).toBe('cart');
    expect(iconFor('repaired')).toBe('build');
    expect(iconFor('transferred')).toBe('swap-horizontal');
    expect(iconFor('sold')).toBe('cash');
    expect(iconFor('donated')).toBe('gift');
    expect(iconFor('lost')).toBe('help-circle');
    expect(iconFor('unknown')).toBe('ellipse');
  });

  it('renders an item per event including the transferred owner line', () => {
    getHistory.mockReturnValueOnce(
      of({
        events: [
          {
            id: 'e1',
            type: 'transferred',
            occurredAt: '2026-01-01',
            actorRef: 'u1',
            note: 'handed over',
            fromOwner: { spaceID: 's1', ownerType: 'personal' },
            toOwner: { spaceID: 's2', ownerType: 'family' },
          },
          {
            id: 'e2',
            type: 'purchased',
            occurredAt: '2026-02-02',
            actorRef: 'u1',
          },
        ],
      }),
    );
    component.spaceID = 's1';
    component.assetID = 'a1';
    component.ngOnChanges();
    fixture.detectChanges();

    const host = fixture.nativeElement as HTMLElement;
    expect(host.querySelectorAll('ion-item').length).toBe(2);
    // The transferred owner line renders both spaceIDs.
    expect(host.textContent).toContain('s1');
    expect(host.textContent).toContain('s2');
    expect(host.textContent).toContain('handed over');
  });

  it('renders the empty placeholder when there are no events', () => {
    getHistory.mockReturnValueOnce(of({ events: [] }));
    component.spaceID = 's1';
    component.assetID = 'a1';
    component.ngOnChanges();
    fixture.detectChanges();

    const host = fixture.nativeElement as HTMLElement;
    expect(host.textContent).toContain('No history events yet.');
  });
});
