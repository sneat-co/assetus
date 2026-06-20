import { TestBed } from '@angular/core/testing';
import { componentTestProviders } from '../../../testing/test-providers';
import { VehicleEngineComponent } from './vehicle-engine.component';

// Render + logic spec for the ported VehicleEngineComponent. No injected
// dependencies, so a bare TestBed render exercises the standalone template,
// and the engine-type mapping is asserted on the emitted vehicle context.
describe('VehicleEngineComponent', () => {
  let fixture: ReturnType<typeof TestBed.createComponent<VehicleEngineComponent>>;
  let component: VehicleEngineComponent;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [VehicleEngineComponent],
      providers: componentTestProviders(),
    });
    fixture = TestBed.createComponent(VehicleEngineComponent);
    component = fixture.componentInstance;
  });

  const callOnEngineTypeChanged = (v: string): void =>
    (
      component as unknown as { onEngineTypeChanged(v: string): void }
    ).onEngineTypeChanged(v);

  const readHasBattery = (): boolean =>
    (component as unknown as { hasBattery: boolean }).hasBattery;

  const readHasCombustion = (): boolean =>
    (component as unknown as { hasCombustion: boolean }).hasCombustion;

  it('creates and renders the select branch when no engine type is set', () => {
    component.vehicleAsset = {
      id: 'a1',
      space: { id: 's1' },
      dbo: {
        category: 'vehicle',
        extraType: 'vehicle',
        extra: { make: 'Audi', model: 'A4' },
      },
    } as never;
    fixture.detectChanges();
    expect(component).toBeTruthy();
    const html: string = fixture.nativeElement.innerHTML;
    expect(html).toContain('sneat-select-from-list');
  });

  it('reports combustion-only for a combustion engine', () => {
    component.vehicleAsset = {
      id: 'a1',
      space: { id: 's1' },
      dbo: {
        category: 'vehicle',
        extraType: 'vehicle',
        extra: {
          make: 'Audi',
          model: 'A4',
          engineType: 'combustion',
          engineFuel: 'petrol',
        },
      },
    } as never;
    expect(readHasCombustion()).toBe(true);
    expect(readHasBattery()).toBe(false);
  });

  it('reports battery-only for an electric engine', () => {
    component.vehicleAsset = {
      id: 'a1',
      space: { id: 's1' },
      dbo: {
        category: 'vehicle',
        extraType: 'vehicle',
        extra: {
          make: 'Tesla',
          model: 'M3',
          engineType: 'electric',
          engineFuel: 'electric',
        },
      },
    } as never;
    expect(readHasBattery()).toBe(true);
    expect(readHasCombustion()).toBe(false);
  });

  it('reports both combustion and battery for a hybrid engine', () => {
    component.vehicleAsset = {
      id: 'a1',
      space: { id: 's1' },
      dbo: {
        category: 'vehicle',
        extraType: 'vehicle',
        extra: {
          make: 'Toyota',
          model: 'Prius',
          engineType: 'hybrid',
          engineFuel: 'petrol',
        },
      },
    } as never;
    fixture.detectChanges();
    expect(readHasBattery()).toBe(true);
    expect(readHasCombustion()).toBe(true);
  });

  it('maps "phev_diesel" to a phev engine type with diesel fuel and emits', () => {
    component.vehicleAsset = {
      id: 'a1',
      space: { id: 's1' },
      dbo: {
        category: 'vehicle',
        extraType: 'vehicle',
        extra: { make: 'Audi', model: 'A4' },
      },
    } as never;
    const emitted: unknown[] = [];
    component.vehicleAssetChange.subscribe((v) => emitted.push(v));

    (component as unknown as { onEngineTypeChanged(v: string): void })
      .onEngineTypeChanged('phev_diesel');

    expect(emitted).toHaveLength(1);
    const extra = (emitted[0] as { dbo: { extra: { engineType: string; engineFuel: string } } })
      .dbo.extra;
    expect(extra.engineType).toBe('phev');
    expect(extra.engineFuel).toBe('diesel');
  });

  it.each([
    ['petrol', 'combustion', 'petrol'],
    ['diesel', 'combustion', 'diesel'],
    ['phev_petrol', 'phev', 'petrol'],
    ['hybrid_diesel', 'hybrid', 'diesel'],
    ['hybrid_petrol', 'hybrid', 'petrol'],
    ['steam', 'steam', ''],
    ['other', 'other', 'other'],
  ])(
    'maps "%s" to engineType "%s" and engineFuel "%s"',
    (input, expectedType, expectedFuel) => {
      component.vehicleAsset = {
        id: 'a1',
        space: { id: 's1' },
        dbo: {
          category: 'vehicle',
          extraType: 'vehicle',
          extra: { make: 'Audi', model: 'A4' },
        },
      } as never;
      const emitted: unknown[] = [];
      component.vehicleAssetChange.subscribe((v) => emitted.push(v));

      callOnEngineTypeChanged(input);

      expect(emitted).toHaveLength(1);
      const extra = (
        emitted[0] as {
          dbo: { extra: { engineType: string; engineFuel: string } };
        }
      ).dbo.extra;
      expect(extra.engineType).toBe(expectedType);
      expect(extra.engineFuel).toBe(expectedFuel);
    },
  );

  it('falls back to unknown engine type for an unrecognised value but still emits', () => {
    component.vehicleAsset = {
      id: 'a1',
      space: { id: 's1' },
      dbo: {
        category: 'vehicle',
        extraType: 'vehicle',
        extra: { make: 'Audi', model: 'A4' },
      },
    } as never;
    const emitted: unknown[] = [];
    component.vehicleAssetChange.subscribe((v) => emitted.push(v));

    callOnEngineTypeChanged('hydrogen');

    expect(emitted).toHaveLength(1);
    const extra = (
      emitted[0] as {
        dbo: { extra: { engineType: string; engineFuel: string } };
      }
    ).dbo.extra;
    // Hydrogen is not handled by the switch, so it falls through to the
    // unknown defaults (both represented as empty strings in the enums).
    expect(extra.engineType).toBe('');
    expect(extra.engineFuel).toBe('');
  });

  it('does not emit when there is no vehicle dbo', () => {
    component.vehicleAsset = undefined;
    const spy = vi.fn();
    component.vehicleAssetChange.subscribe(spy);

    (component as unknown as { onEngineTypeChanged(v: string): void })
      .onEngineTypeChanged('petrol');

    expect(spy).not.toHaveBeenCalled();
  });
});
