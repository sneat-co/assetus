import { TestBed } from '@angular/core/testing';
import { FormControl } from '@angular/forms';
import { spacePageTestProviders } from '../../../testing/test-providers';
import { VehicleCardComponent } from './vehicle-card.component';

// Render + logic spec for the ported VehicleCardComponent. It composes the
// make-model / engine / possession / reg-number children and needs the
// standard provider chain.
describe('VehicleCardComponent', () => {
  let fixture: ReturnType<typeof TestBed.createComponent<VehicleCardComponent>>;
  let component: VehicleCardComponent;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [VehicleCardComponent],
      providers: [...spacePageTestProviders()],
    });
    fixture = TestBed.createComponent(VehicleCardComponent);
    component = fixture.componentInstance;
  });

  it('creates and renders', () => {
    component.space = { id: 's1' } as never;
    component.vehicleAsset = {
      id: 'a1',
      space: { id: 's1' },
      dbo: { category: 'vehicle', countryID: 'IE', extra: { make: '', model: '' } },
    } as never;
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('emits an updated vehicle context when the make changes', () => {
    component.vehicleAsset = {
      id: 'a1',
      space: { id: 's1' },
      dbo: { category: 'vehicle', extra: { make: '', model: '' } },
    } as never;
    let emitted: { dbo: { extra: { make: string } } } | undefined;
    component.vehicleAssetChange.subscribe((v) => (emitted = v as never));

    component.makeChanged('Audi');

    expect(emitted?.dbo.extra.make).toBe('Audi');
  });

  const setVehicle = (): void => {
    component.vehicleAsset = {
      id: 'a1',
      space: { id: 's1' },
      dbo: { category: 'vehicle', extra: { make: '', model: '' } },
    } as never;
  };

  const onRegNumberChanged = (v: string): void =>
    (
      component as unknown as { onRegNumberChanged(v: string): void }
    ).onRegNumberChanged(v);
  const onRegNumberSkipped = (): void =>
    (
      component as unknown as { onRegNumberSkipped(): void }
    ).onRegNumberSkipped();
  const onAssetChanged = (a: unknown): void =>
    (component as unknown as { onAssetChanged(a: unknown): void }).onAssetChanged(
      a,
    );
  const modelChanged = (m: string): void =>
    (component as unknown as { modelChanged(m: string): void }).modelChanged(m);

  it('sets the regNumber form control from the dbo extra on changes', () => {
    component.vehicleAsset = {
      id: 'a1',
      space: { id: 's1' },
      dbo: {
        category: 'vehicles',
        extra: { make: '', model: '', regNumber: 'AA-11' },
      },
    } as never;
    component.ngOnChanges({ vehicleAsset: {} as never });
    expect(
      (component as unknown as { regNumber: { value: string } }).regNumber.value,
    ).toBe('AA-11');
  });

  it('does not overwrite a dirty regNumber control on changes', () => {
    const ctrl = (component as unknown as { regNumber: FormControl<string> })
      .regNumber;
    ctrl.setValue('USER-TYPED');
    ctrl.markAsDirty();
    component.vehicleAsset = {
      id: 'a1',
      space: { id: 's1' },
      brief: { extra: { regNumber: 'AA-11' } },
      dbo: { category: 'vehicle', extra: { make: '', model: '' } },
    } as never;
    component.ngOnChanges({ vehicleAsset: {} as never });
    expect(ctrl.value).toBe('USER-TYPED');
  });

  it('ignores ngOnChanges that do not touch vehicleAsset', () => {
    setVehicle();
    expect(() => component.ngOnChanges({ space: {} as never })).not.toThrow();
  });

  it('emits a new countryID when the country changes', () => {
    setVehicle();
    let emitted: { dbo: { countryID: string } } | undefined;
    component.vehicleAssetChange.subscribe((v) => (emitted = v as never));

    component.countryChanged('FR');

    expect(emitted?.dbo.countryID).toBe('FR');
  });

  it('does not emit on country change without a dbo', () => {
    component.vehicleAsset = undefined;
    const spy = vi.fn();
    component.vehicleAssetChange.subscribe(spy);
    component.countryChanged('FR');
    expect(spy).not.toHaveBeenCalled();
  });

  it('emits the new model on model change', () => {
    setVehicle();
    let emitted: { dbo: { extra: { model: string } } } | undefined;
    component.vehicleAssetChange.subscribe((v) => (emitted = v as never));

    modelChanged('A4');

    expect(emitted?.dbo.extra.model).toBe('A4');
  });

  it('emits the reg number on reg number change', () => {
    setVehicle();
    let emitted: { dbo: { extra: { regNumber: string } } } | undefined;
    component.vehicleAssetChange.subscribe((v) => (emitted = v as never));

    onRegNumberChanged('XY-99');

    expect(emitted?.dbo.extra.regNumber).toBe('XY-99');
    expect(
      (component as unknown as { regNumber: FormControl<string> }).regNumber
        .value,
    ).toBe('XY-99');
  });

  it('flags the reg number as skipped', () => {
    onRegNumberSkipped();
    expect(
      (component as unknown as { regNumberSkipped: boolean }).regNumberSkipped,
    ).toBe(true);
  });

  it('re-emits the asset passed to onAssetChanged', () => {
    let emitted: { id: string } | undefined;
    component.vehicleAssetChange.subscribe((v) => (emitted = v as never));

    onAssetChanged({ id: 'a2', dbo: { category: 'vehicle' } } as never);

    expect(emitted?.id).toBe('a2');
  });

  it('populates models and year fields via the input setters', () => {
    setVehicle();
    component.make = 'Audi';
    component.model = 'A4';
    component.yearBuild = 2020;

    expect(component.makeVal).toBe('Audi');
    expect(component.modelVal).toBe('A4');
    expect(component.yearBuildNumber).toBe(2020);
    expect(component.yearBuildVal).toBe('2020');
    expect(component.makes?.length).toBeGreaterThan(0);
  });

  it('clears models when the make setter receives an empty value', () => {
    component.make = '';
    expect(component.models).toBeUndefined();
    expect(component.modelVal).toBeUndefined();
  });
});
