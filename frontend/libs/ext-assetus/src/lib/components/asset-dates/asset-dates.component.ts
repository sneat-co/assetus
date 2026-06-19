import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  IonDatetime,
  IonItem,
  IonItemDivider,
  IonItemGroup,
  IonLabel,
} from '@ionic/angular/standalone';
import {
  IAssetDbo,
  IAssetDboBase,
  IAssetVehicleExtra,
} from '@sneat/mod-assetus-core';

interface AssetDate {
  name: string;
  title: string;
  value?: string;
}

// Ported from @sneat/ext-assetus-components (legacy assetus components lib).
@Component({
  selector: 'sneat-asset-dates',
  templateUrl: './asset-dates.component.html',
  imports: [
    FormsModule,
    IonItemGroup,
    IonItemDivider,
    IonItem,
    IonLabel,
    IonDatetime,
  ],
})
export class AssetDatesComponent {
  private assetDto?: IAssetDboBase;

  @Input() set asset(v: IAssetDboBase) {
    this.assetDto = v;
    switch (v.category) {
      case 'vehicle': {
        const vehicle = v as IAssetDbo<'vehicle', IAssetVehicleExtra>;
        this.items = [
          {
            name: 'nctExpires',
            title: 'NCT expires',
            value: vehicle.extra?.nctExpires,
          },
          {
            name: 'taxExpires',
            title: 'Tax expires',
            value: vehicle.extra?.taxExpires,
          },
          {
            name: 'nextServiceDue',
            title: 'Next service due',
            value: vehicle.extra?.nextServiceDue,
          },
        ];
        break;
      }
      case 'dwelling': {
        // const property = v as IDwelling;
        this.items = [
          {
            name: 'leaseExpires',
            title: 'Lease expires',
            value: 'property.leaseExpires',
          },
        ];
        break;
      }
      default:
        this.items = [];
        break;
    }
  }

  @Output() changed = new EventEmitter<{ name: string; value: string }>();

  protected items?: AssetDate[];

  trackByName(i: number, v: AssetDate): string {
    return v.name;
  }

  onChange(name: string, $event: CustomEvent): void {
    const value = $event.detail.value;
    this.changed.emit({ name, value });
    let title: string;
    switch (name) {
      case 'nctExpires':
        title = 'NCT expires';
        break;
      case 'taxExpires':
        title = 'Tax expires';
        break;
      case 'nextServiceDue':
        title = 'Next service due';
        break;
      default:
        title = name;
        break;
    }
    throw new Error('not implemented yet, title=' + title);
  }
}
