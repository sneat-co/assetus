import { TestBed } from '@angular/core/testing';
import { componentTestProviders } from '../../../testing/test-providers';
import { AssetPossessionCardComponent } from './asset-possession-card.component';

// Render + logic spec for the ported AssetPossessionCardComponent.
describe('AssetPossessionCardComponent', () => {
  let fixture: ReturnType<
    typeof TestBed.createComponent<AssetPossessionCardComponent>
  >;
  let component: AssetPossessionCardComponent;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [AssetPossessionCardComponent],
      providers: componentTestProviders(),
    });
    fixture = TestBed.createComponent(AssetPossessionCardComponent);
    component = fixture.componentInstance;
  });

  it('creates and renders', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('emits the asset with the new possession on change', () => {
    component.asset = {
      id: 'a1',
      space: { id: 's1' },
      dbo: { category: 'vehicle', possession: 'owning' },
    } as never;
    let emitted: { dbo: { possession: string } } | undefined;
    component.assetChange.subscribe((v) => (emitted = v as never));

    (component as unknown as { onPossessionChanged(p: string): void })
      .onPossessionChanged('renting');

    expect(emitted?.dbo.possession).toBe('renting');
  });

  it('does not emit when there is no asset dbo', () => {
    component.asset = undefined;
    const spy = vi.fn();
    component.assetChange.subscribe(spy);

    (component as unknown as { onPossessionChanged(p: string): void })
      .onPossessionChanged('renting');

    expect(spy).not.toHaveBeenCalled();
  });
});
