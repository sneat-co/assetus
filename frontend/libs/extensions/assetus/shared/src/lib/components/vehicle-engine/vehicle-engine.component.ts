import { Component, EventEmitter, Input, Output } from '@angular/core';
import { IonInput, IonItem, IonItemGroup } from '@ionic/angular/standalone';
import { ISelectItem, SelectFromListComponent } from '@sneat/ui';
import {
  EngineType,
  EngineTypeCombustion,
  EngineTypeElectric,
  EngineTypeHybrid,
  EngineTypePHEV,
  EngineTypes,
  FuelType,
  FuelTypes,
} from '@sneat/extension-assetus-contract';
import { IAssetVehicleContext } from '@sneat/extension-assetus-contract';

// Ported from legacy ext-assetus-components (legacy assetus components lib).
@Component({
  selector: 'assetus-vehicle-engine',
  templateUrl: './vehicle-engine.component.html',
  imports: [SelectFromListComponent, IonItemGroup, IonItem, IonInput],
})
export class VehicleEngineComponent {
  @Input() public vehicleAsset?: IAssetVehicleContext;
  @Output() public readonly vehicleAssetChange =
    new EventEmitter<IAssetVehicleContext>();

  protected get hasBattery(): boolean {
    const et = this.vehicleAsset?.dbo?.extra?.engineType;
    return (
      et === EngineTypeElectric ||
      et === EngineTypePHEV ||
      et === EngineTypeHybrid
    );
  }

  protected get hasCombustion(): boolean {
    const et = this.vehicleAsset?.dbo?.extra?.engineType;
    return (
      et === EngineTypeCombustion ||
      et === EngineTypePHEV ||
      et === EngineTypeHybrid
    );
  }

  readonly engineTypes: ISelectItem[] = [
    { id: FuelTypes.petrol, title: 'Petrol', emoji: '🔥' },
    { id: FuelTypes.diesel, title: 'Diesel', emoji: '🔥' },
    { id: EngineTypes.electric, title: 'Electric', emoji: '🔌' },
    { id: 'phev_petrol', title: 'PHEV Petrol', emoji: '🔌🔥' },
    { id: 'phev_diesel', title: 'PHEV Diesel', emoji: '🔌🔥' },
    { id: 'hybrid_petrol', title: 'Hybrid Petrol', emoji: '🔋🔥' },
    { id: 'hybrid_diesel', title: 'Hybrid Diesel', emoji: '🔋🔥' },
    { id: FuelTypes.hydrogen, title: 'Hydrogen', emoji: '💧' },
    { id: EngineTypes.steam, title: 'Steam', emoji: '🚂' },
    { id: EngineTypes.other, title: 'Other', labelColor: 'medium' },
    { id: '', title: 'Unknown', labelColor: 'medium' },
  ];

  protected onEngineTypeChanged(v: string): void {
    let engineType: EngineType = EngineTypes.unknown;
    let engineFuel: FuelType = FuelTypes.unknown;
    if (this.vehicleAsset?.dbo) {
      switch (v) {
        case FuelTypes.diesel:
          engineType = EngineTypes.combustion;
          engineFuel = FuelTypes.diesel;
          break;
        case FuelTypes.petrol:
          engineType = EngineTypes.combustion;
          engineFuel = FuelTypes.petrol;
          break;
        case 'phev_diesel':
          engineType = EngineTypes.phev;
          engineFuel = FuelTypes.diesel;
          break;
        case 'phev_petrol':
          engineType = EngineTypes.phev;
          engineFuel = FuelTypes.petrol;
          break;
        case 'hybrid_diesel':
          engineType = EngineTypes.hybrid;
          engineFuel = FuelTypes.diesel;
          break;
        case 'hybrid_petrol':
          engineType = EngineTypes.hybrid;
          engineFuel = FuelTypes.petrol;
          break;
        case EngineTypes.steam:
          engineType = EngineTypes.steam;
          engineFuel = FuelTypes.unknown;
          break;
        case 'other':
          engineType = EngineTypes.other;
          engineFuel = FuelTypes.other;
          break;
      }
      this.vehicleAsset = {
        ...this.vehicleAsset,
        dbo: {
          ...this.vehicleAsset.dbo,
          extra: {
            ...(this.vehicleAsset.dbo?.extra || { make: '', model: '' }),
            engineType: engineType,
            engineFuel: engineFuel,
          },
        },
      };
      this.vehicleAssetChange.emit(this.vehicleAsset);
    }
  }
}
