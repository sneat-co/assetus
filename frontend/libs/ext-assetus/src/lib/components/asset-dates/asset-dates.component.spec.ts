import { TestBed } from '@angular/core/testing';
import { AssetDatesComponent } from './asset-dates.component';

// Render + logic spec for the ported AssetDatesComponent.
describe('AssetDatesComponent', () => {
  let fixture: ReturnType<typeof TestBed.createComponent<AssetDatesComponent>>;
  let component: AssetDatesComponent;

  beforeEach(() => {
    TestBed.configureTestingModule({ imports: [AssetDatesComponent] });
    fixture = TestBed.createComponent(AssetDatesComponent);
    component = fixture.componentInstance;
  });

  it('creates and renders', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('builds the vehicle date rows from the vehicle extra', () => {
    component.asset = {
      category: 'vehicle',
      extra: { nctExpires: '2026-01-01', taxExpires: '2026-02-02' },
    } as never;
    const items = (component as unknown as { items: { name: string; value?: string }[] }).items;
    expect(items.map((i) => i.name)).toEqual([
      'nctExpires',
      'taxExpires',
      'nextServiceDue',
    ]);
    expect(items[0].value).toBe('2026-01-01');
  });

  it('builds a single lease row for a dwelling', () => {
    component.asset = { category: 'dwelling' } as never;
    const items = (component as unknown as { items: { name: string }[] }).items;
    expect(items.map((i) => i.name)).toEqual(['leaseExpires']);
  });

  it('emits the change then throws the not-implemented marker', () => {
    const spy = vi.fn();
    component.changed.subscribe(spy);
    expect(() =>
      component.onChange('nctExpires', { detail: { value: '2027-01-01' } } as CustomEvent),
    ).toThrowError(/not implemented/);
    expect(spy).toHaveBeenCalledWith({ name: 'nctExpires', value: '2027-01-01' });
  });

  it('builds an empty list for an unknown asset category', () => {
    component.asset = { category: 'other' } as never;
    const items = (component as unknown as { items: unknown[] }).items;
    expect(items).toEqual([]);
  });

  it('renders a date row per item once an asset is set', () => {
    component.asset = {
      category: 'vehicle',
      extra: { nctExpires: '2026-01-01' },
    } as never;
    fixture.detectChanges();
    const text: string = fixture.nativeElement.textContent;
    expect(text).toContain('NCT expires');
    expect(text).toContain('Tax expires');
    expect(text).toContain('Next service due');
  });

  it.each(['taxExpires', 'nextServiceDue', 'unknownDate'])(
    'emits the change then throws not-implemented for "%s"',
    (name) => {
      const spy = vi.fn();
      component.changed.subscribe(spy);
      expect(() =>
        component.onChange(name, {
          detail: { value: '2027-05-05' },
        } as CustomEvent),
      ).toThrowError(/not implemented/);
      expect(spy).toHaveBeenCalledWith({ name, value: '2027-05-05' });
    },
  );

  it('trackByName returns the row name', () => {
    expect(component.trackByName(0, { name: 'x', title: 'X' })).toBe('x');
  });
});
