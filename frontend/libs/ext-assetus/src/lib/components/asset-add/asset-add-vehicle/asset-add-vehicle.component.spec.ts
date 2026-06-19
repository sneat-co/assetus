import { TestBed } from '@angular/core/testing';
import { WritableSignal } from '@angular/core';
import { EMPTY } from 'rxjs';
import { AssetService } from '@sneat/ext-assetus-components';
import { spacePageTestProviders } from '../../../../testing/test-providers';
import { AssetAddVehicleComponent } from './asset-add-vehicle.component';

// Render + logic spec for the ported AssetAddVehicleComponent, which extends the
// legacy AddAssetBaseComponent (SpaceBaseComponent + legacy AssetService).
describe('AssetAddVehicleComponent', () => {
  let fixture: ReturnType<
    typeof TestBed.createComponent<AssetAddVehicleComponent>
  >;
  let component: AssetAddVehicleComponent;
  let createAsset: ReturnType<typeof vi.fn>;

  // Seeds the read-only `space` getter by updating the protected $spaceRef
  // signal that backs it (the getter derives `space` from $spaceRef). We never
  // assign component.space directly.
  const seedSpace = (id = 's1'): void => {
    (
      component as unknown as { $spaceRef: WritableSignal<{ id: string }> }
    ).$spaceRef.set({ id });
  };

  const vehicleAssetInput = () =>
    ({
      id: 'a1',
      space: { id: 's1' },
      dbo: { category: 'vehicle', extra: { make: '', model: '' } },
    }) as never;

  beforeEach(() => {
    createAsset = vi.fn().mockReturnValue(EMPTY);
    TestBed.configureTestingModule({
      imports: [AssetAddVehicleComponent],
      providers: [
        ...spacePageTestProviders(),
        { provide: AssetService, useValue: { createAsset } },
      ],
    });
    fixture = TestBed.createComponent(AssetAddVehicleComponent);
    component = fixture.componentInstance;
  });

  it('creates and renders with a populated vehicle asset', () => {
    component.vehicleAsset = vehicleAssetInput();
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('maps the vehicle type onto the current vehicle asset', () => {
    component.vehicleAsset = vehicleAssetInput();
    (component as unknown as { vehicleType: string }).vehicleType = 'car';
    component.onVehicleTypeChanged();
    expect(
      (component as unknown as { vehicleAsset: { dbo: { type: string } } })
        .vehicleAsset.dbo.type,
    ).toBe('car');
  });

  it('onVehicleTypeChanged is a no-op when there is no asset dbo', () => {
    component.vehicleAsset = undefined;
    expect(() => component.onVehicleTypeChanged()).not.toThrow();
    expect(component.vehicleAsset).toBeUndefined();
  });

  it('onAssetChanged replaces the current vehicle asset', () => {
    const replacement = { id: 'a2', space: { id: 's1' }, dbo: {} } as never;
    (
      component as unknown as { onAssetChanged(a: unknown): void }
    ).onAssetChanged(replacement);
    expect(component.vehicleAsset).toBe(replacement);
  });

  it('ngOnChanges seeds a draft asset when space arrives', () => {
    seedSpace();
    component.ngOnChanges({ space: {} as never });
    expect(component.vehicleAsset?.dbo?.category).toBe('vehicle');
    expect(component.vehicleAsset?.space?.id).toBe('s1');
  });

  it('ngOnChanges keeps the provided asset but binds the space', () => {
    seedSpace();
    component.vehicleAsset = vehicleAssetInput();
    component.ngOnChanges({ space: {} as never });
    const asset = component.vehicleAsset as unknown as {
      id: string;
      space?: { id: string };
    };
    expect(asset?.id).toBe('a1');
    expect(asset?.space?.id).toBe('s1');
  });

  it('ngOnChanges ignores change sets without a space key', () => {
    component.ngOnChanges({ vehicleAsset: {} as never });
    expect(component.vehicleAsset).toBeUndefined();
  });

  it('formatDate returns an empty string for a missing value', () => {
    expect(component.formatDate(undefined)).toBe('');
  });

  it('formatDate returns an empty string for an array value', () => {
    expect(component.formatDate(['2026-03-15'])).toBe('');
  });

  it('formatDate formats an ISO date', () => {
    expect(component.formatDate('2026-03-15')).toContain('2026');
  });

  it('submitVehicleForm throws when there is no vehicle type', () => {
    seedSpace();
    component.vehicleAsset = vehicleAssetInput();
    expect(() =>
      (
        component as unknown as { submitVehicleForm(): void }
      ).submitVehicleForm(),
    ).toThrow();
  });

  it('submitVehicleForm throws when there is no asset', () => {
    seedSpace();
    component.vehicleAsset = undefined;
    (component as unknown as { vehicleType: string }).vehicleType = 'car';
    expect(() =>
      (
        component as unknown as { submitVehicleForm(): void }
      ).submitVehicleForm(),
    ).toThrow();
  });

  it('submitVehicleForm creates the asset and sets isSubmitting', () => {
    seedSpace();
    component.vehicleAsset = vehicleAssetInput();
    (component as unknown as { vehicleType: string }).vehicleType = 'car';
    (component as unknown as { yearOfBuild: string }).yearOfBuild = '2020';
    (
      component as unknown as { submitVehicleForm(): void }
    ).submitVehicleForm();
    expect(createAsset).toHaveBeenCalledOnce();
    expect(
      (component as unknown as { isSubmitting: boolean }).isSubmitting,
    ).toBe(true);
    const request = createAsset.mock.calls[0][1] as {
      asset: { status: string; yearOfBuild: number };
    };
    expect(request.asset.status).toBe('active');
    expect(request.asset.yearOfBuild).toBe(2020);
  });

  it('submitVehicleForm omits yearOfBuild when it is blank', () => {
    seedSpace();
    component.vehicleAsset = vehicleAssetInput();
    (component as unknown as { vehicleType: string }).vehicleType = 'car';
    (
      component as unknown as { submitVehicleForm(): void }
    ).submitVehicleForm();
    const request = createAsset.mock.calls[0][1] as {
      asset: { yearOfBuild?: number };
    };
    expect(request.asset.yearOfBuild).toBeUndefined();
  });
});
