import { TestBed } from '@angular/core/testing';
import { WritableSignal } from '@angular/core';
import { EMPTY } from 'rxjs';
import { ASSET_SERVICE } from '@sneat/extension-assetus-contract';
import { spacePageTestProviders } from '../../../../testing/test-providers';
import { AssetAddDwellingComponent } from './asset-add-dwelling.component';

// Render + logic spec for the ported AssetAddDwellingComponent, which extends
// the legacy AddAssetBaseComponent.
describe('AssetAddDwellingComponent', () => {
  let fixture: ReturnType<
    typeof TestBed.createComponent<AssetAddDwellingComponent>
  >;
  let component: AssetAddDwellingComponent;
  let createAsset: ReturnType<typeof vi.fn>;

  // Seeds the read-only `space` getter via the protected $spaceRef signal.
  const seedSpace = (id = 's1'): void => {
    (
      component as unknown as { $spaceRef: WritableSignal<{ id: string }> }
    ).$spaceRef.set({ id });
  };

  const dwellingAssetInput = (extra: Record<string, unknown> = {}) =>
    ({
      id: 'a1',
      space: { id: 's1' },
      dbo: { category: 'dwelling', extra },
    }) as never;

  beforeEach(() => {
    createAsset = vi.fn().mockReturnValue(EMPTY);
    TestBed.configureTestingModule({
      imports: [AssetAddDwellingComponent],
      providers: [
        ...spacePageTestProviders(),
        { provide: ASSET_SERVICE, useValue: { createAsset } },
      ],
    });
    fixture = TestBed.createComponent(AssetAddDwellingComponent);
    component = fixture.componentInstance;
  });

  it('creates and renders with a populated dwelling asset', () => {
    component.dwellingAsset = dwellingAssetInput({
      rent_price: { value: 0, currency: 'USD' },
    });
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('maps the dwelling type onto the current dwelling asset', () => {
    component.dwellingAsset = dwellingAssetInput();
    (component as unknown as { dwellingType: string }).dwellingType = 'house';
    (
      component as unknown as { onDwellingTypeChanged(): void }
    ).onDwellingTypeChanged();
    expect(
      (component as unknown as { dwellingAsset: { dbo: { type: string } } })
        .dwellingAsset.dbo.type,
    ).toBe('house');
  });

  it('onDwellingTypeChanged is a no-op when there is no asset dbo', () => {
    component.dwellingAsset = undefined;
    expect(() =>
      (
        component as unknown as { onDwellingTypeChanged(): void }
      ).onDwellingTypeChanged(),
    ).not.toThrow();
    expect(component.dwellingAsset).toBeUndefined();
  });

  it('onAssetChanged replaces the current dwelling asset', () => {
    const replacement = { id: 'a2', space: { id: 's1' }, dbo: {} } as never;
    (
      component as unknown as { onAssetChanged(a: unknown): void }
    ).onAssetChanged(replacement);
    expect(component.dwellingAsset).toBe(replacement);
  });

  it('ngOnChanges seeds a draft dwelling when space arrives', () => {
    seedSpace();
    component.ngOnChanges({ space: {} as never });
    expect(component.dwellingAsset?.dbo?.category).toBe('dwelling');
    expect(component.dwellingAsset?.space?.id).toBe('s1');
  });

  it('ngOnChanges ignores change sets without a space key', () => {
    component.ngOnChanges({ dwellingAsset: {} as never });
    expect(component.dwellingAsset).toBeUndefined();
  });

  it('submitDwellingForm throws when there is no dwelling asset', () => {
    seedSpace();
    component.dwellingAsset = undefined;
    expect(() =>
      (
        component as unknown as { submitDwellingForm(): void }
      ).submitDwellingForm(),
    ).toThrow();
  });

  it('submitDwellingForm coerces numeric extras and creates the asset', () => {
    seedSpace();
    component.dwellingAsset = dwellingAssetInput({
      numberOfBedrooms: '3',
      areaSqM: '85',
    });
    (
      component as unknown as { submitDwellingForm(): void }
    ).submitDwellingForm();
    expect(createAsset).toHaveBeenCalledOnce();
    expect(
      (component as unknown as { isSubmitting: boolean }).isSubmitting,
    ).toBe(true);
    const request = createAsset.mock.calls[0][0] as {
      status: string;
      extra: { numberOfBedrooms: number; areaSqM: number };
    };
    expect(request.status).toBe('active');
    expect(request.extra.numberOfBedrooms).toBe(3);
    expect(request.extra.areaSqM).toBe(85);
  });

  it('submitDwellingForm creates the asset without numeric extras', () => {
    seedSpace();
    component.dwellingAsset = dwellingAssetInput({});
    (
      component as unknown as { submitDwellingForm(): void }
    ).submitDwellingForm();
    expect(createAsset).toHaveBeenCalledOnce();
  });
});
