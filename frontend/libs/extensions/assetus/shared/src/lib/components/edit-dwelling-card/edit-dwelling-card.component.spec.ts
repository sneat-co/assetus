import { TestBed } from '@angular/core/testing';
import { spacePageTestProviders } from '../../../testing/test-providers';
import { AddDwellingCardComponent } from './edit-dwelling-card.component';

// Render + logic spec for the ported AddDwellingCardComponent. It composes the
// possession child and the country selector, so it needs the standard provider
// chain.
describe('AddDwellingCardComponent', () => {
  let fixture: ReturnType<
    typeof TestBed.createComponent<AddDwellingCardComponent>
  >;
  let component: AddDwellingCardComponent;

  const dwellingAsset = (
    dbo: Record<string, unknown> = {
      category: 'dwelling',
      extra: { rent_price: { value: 0, currency: 'USD' } },
    },
  ) =>
    ({
      id: 'a1',
      space: { id: 's1' },
      dbo,
    }) as never;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [AddDwellingCardComponent],
      providers: spacePageTestProviders(),
    });
    fixture = TestBed.createComponent(AddDwellingCardComponent);
    component = fixture.componentInstance;
  });

  it('creates and renders', () => {
    component.space = { id: 's1' } as never;
    component.dwellingAsset = dwellingAsset();
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('hydrates local fields from the dwelling asset on changes', () => {
    component.dwellingAsset = dwellingAsset({
      category: 'dwelling',
      name: 'Flat 1',
      extra: {
        address: { lines: '1 Main St' },
        rent_price: { value: 1200, currency: 'EUR' },
        numberOfBedrooms: 2,
        areaSqM: 75,
      },
    });
    component.ngOnChanges({ dwellingAsset: {} as never });
    const c = component as unknown as {
      title: string;
      address: string;
      rent_price_amount?: number;
      rent_price_currency: string;
      number_of_bedrooms?: number;
      areaSqM?: number;
    };
    expect(c.title).toBe('Flat 1');
    expect(c.address).toBe('1 Main St');
    expect(c.rent_price_amount).toBe(1200);
    expect(c.rent_price_currency).toBe('EUR');
    expect(c.number_of_bedrooms).toBe(2);
    expect(c.areaSqM).toBe(75);
  });

  it('falls back to defaults when the asset has no fields', () => {
    component.dwellingAsset = dwellingAsset({ category: 'dwelling' });
    component.ngOnChanges({ dwellingAsset: {} as never });
    const c = component as unknown as {
      title: string;
      address: string;
      rent_price_currency: string;
    };
    expect(c.title).toBe('');
    expect(c.address).toBe('');
    expect(c.rent_price_currency).toBe('USD');
  });

  it('seeds the country from the space and emits when space changes', () => {
    component.space = { id: 's1', dbo: { countryID: 'IE' } } as never;
    component.dwellingAsset = dwellingAsset({
      category: 'dwelling',
      extra: {},
    });
    let emitted: { dbo: { countryID: string } } | undefined;
    component.dwellingAssetChange.subscribe((v) => (emitted = v as never));
    component.ngOnChanges({ space: {} as never });
    expect(emitted?.dbo.countryID).toBe('IE');
    expect(
      (component.dwellingAsset as unknown as { dbo: { countryID: string } }).dbo
        .countryID,
    ).toBe('IE');
  });

  it('does not seed the country when the asset already has one', () => {
    component.space = { id: 's1', dbo: { countryID: 'IE' } } as never;
    component.dwellingAsset = dwellingAsset({
      category: 'dwelling',
      countryID: 'US',
      extra: {},
    });
    const emit = vi.fn();
    component.dwellingAssetChange.subscribe(emit);
    component.ngOnChanges({ space: {} as never });
    expect(emit).not.toHaveBeenCalled();
  });

  it('emits an updated dwelling when a brief field changes', () => {
    component.dwellingAsset = dwellingAsset({ category: 'dwelling', extra: {} });
    let emitted: { dbo: { title: string } } | undefined;
    component.dwellingAssetChange.subscribe((v) => (emitted = v as never));
    (
      component as unknown as { onBriefFieldChanged(f: string, v: string): void }
    ).onBriefFieldChanged('title', 'New title');
    expect(emitted?.dbo.title).toBe('New title');
  });

  it('onBriefFieldChanged is a no-op when there is no asset dbo', () => {
    component.dwellingAsset = undefined;
    const emit = vi.fn();
    component.dwellingAssetChange.subscribe(emit);
    (
      component as unknown as { onBriefFieldChanged(f: string, v: string): void }
    ).onBriefFieldChanged('title', 'x');
    expect(emit).not.toHaveBeenCalled();
  });

  it('emits an updated dwelling when an extra field changes', () => {
    component.dwellingAsset = dwellingAsset({
      category: 'dwelling',
      extra: { numberOfBedrooms: 1 },
    });
    let emitted: { dbo: { extra: { numberOfBedrooms: string } } } | undefined;
    component.dwellingAssetChange.subscribe((v) => (emitted = v as never));
    (
      component as unknown as { onExtraFieldChanged(f: string, v: string): void }
    ).onExtraFieldChanged('numberOfBedrooms', '4');
    expect(emitted?.dbo.extra.numberOfBedrooms).toBe('4');
  });

  it('emits an updated dwelling when a rent-price field changes', () => {
    component.dwellingAsset = dwellingAsset({
      category: 'dwelling',
      extra: { rent_price: { value: 0, currency: 'USD' } },
    });
    let emitted:
      | { dbo: { extra: { rent_price: { value: string } } } }
      | undefined;
    component.dwellingAssetChange.subscribe((v) => (emitted = v as never));
    (
      component as unknown as {
        onRentPriceFieldChanged(f: string, v: string): void;
      }
    ).onRentPriceFieldChanged('value', '999');
    expect(emitted?.dbo.extra.rent_price.value).toBe('999');
  });

  it('onRentPriceFieldChanged is a no-op without a rent price', () => {
    component.dwellingAsset = dwellingAsset({
      category: 'dwelling',
      extra: {},
    });
    const emit = vi.fn();
    component.dwellingAssetChange.subscribe(emit);
    (
      component as unknown as {
        onRentPriceFieldChanged(f: string, v: string): void;
      }
    ).onRentPriceFieldChanged('value', '999');
    expect(emit).not.toHaveBeenCalled();
  });

  it('emits an updated dwelling when the country changes', () => {
    component.dwellingAsset = dwellingAsset({ category: 'dwelling', extra: {} });
    let emitted: { dbo: { countryID: string } } | undefined;
    component.dwellingAssetChange.subscribe((v) => (emitted = v as never));
    (
      component as unknown as { onCountryChanged(v: string): void }
    ).onCountryChanged('FR');
    expect(emitted?.dbo.countryID).toBe('FR');
  });

  it('onAssetChanged replaces and re-emits the dwelling asset', () => {
    const replacement = { id: 'a2', space: { id: 's1' }, dbo: {} } as never;
    let emitted: unknown;
    component.dwellingAssetChange.subscribe((v) => (emitted = v));
    (
      component as unknown as { onAssetChanged(a: unknown): void }
    ).onAssetChanged(replacement);
    expect(component.dwellingAsset).toBe(replacement);
    expect(emitted).toBe(replacement);
  });
});
