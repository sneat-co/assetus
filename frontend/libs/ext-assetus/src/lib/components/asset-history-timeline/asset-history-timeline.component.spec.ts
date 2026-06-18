import { TestBed } from '@angular/core/testing';
import { ErrorLogger } from '@sneat/core';
import { of } from 'rxjs';
import { IHistoryEvent } from '../../dto';
import { AssetService } from '../../services';
import { AssetHistoryTimelineComponent } from './asset-history-timeline.component';

describe('AssetHistoryTimelineComponent', () => {
  const events: IHistoryEvent[] = [
    {
      id: 'e1',
      type: 'transferred',
      occurredAt: '2026-01-01T00:00:00Z',
      actorRef: 'u1',
      fromOwner: { spaceID: 's1', spaceType: 'family', ownerType: 'family' },
      toOwner: { spaceID: 's2', spaceType: 'private', ownerType: 'individual' },
    },
  ];
  const getHistory = vi.fn().mockReturnValue(of({ assetID: 'a1', events }));

  beforeEach(() => {
    getHistory.mockClear();
    TestBed.configureTestingModule({
      imports: [AssetHistoryTimelineComponent],
      providers: [
        { provide: AssetService, useValue: { getHistory } },
        {
          provide: ErrorLogger,
          useValue: { logError: () => undefined, logErrorHandler: () => () => undefined },
        },
      ],
    });
  });

  it('renders the transferred event with prior and new owner', () => {
    const fixture = TestBed.createComponent(AssetHistoryTimelineComponent);
    fixture.componentRef.setInput('spaceID', 's1');
    fixture.componentRef.setInput('assetID', 'a1');
    fixture.detectChanges();
    expect(getHistory).toHaveBeenCalledWith('s1', 'a1');
    const host = fixture.nativeElement as HTMLElement;
    expect(host.textContent).toContain('transferred');
    expect(host.textContent).toContain('s1');
    expect(host.textContent).toContain('s2');
  });
});
