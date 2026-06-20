import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { IonButton, IonCard } from '@ionic/angular/standalone';
import { ISelectItem, SelectFromListComponent } from '@sneat/ui';
import { timestamp } from '@sneat/dto';
import { SpaceComponentBaseParams } from '@sneat/space-components';
import {
  AssetVehicleType,
  EngineTypes,
  FuelTypes,
  IAssetDboBase,
  IAssetExtra,
  IAssetVehicleExtra,
} from '@sneat/extension-assetus-contract';
import {
  IAssetContext,
  IAssetVehicleContext,
} from '@sneat/extension-assetus-contract';
import { format, parseISO } from 'date-fns';
import { AddAssetBaseComponent } from '../add-asset-base.component';
import { VehicleCardComponent } from '../../vehicle-card/vehicle-card.component';
import { ClassName } from '@sneat/ui';

// Ported from legacy ext-assetus-components (legacy assetus components lib).
// Extends the published legacy AddAssetBaseComponent (its createAsset<Extra>
// flow differs from the MVP AssetService).
@Component({
  selector: 'assetus-asset-add-vehicle',
  templateUrl: './asset-add-vehicle.component.html',
  providers: [
    {
      provide: ClassName,
      useValue: 'AssetAddVehicleComponent',
    },
    SpaceComponentBaseParams,
  ],
  imports: [
    SelectFromListComponent,
    FormsModule,
    IonCard,
    VehicleCardComponent,
    IonButton,
  ],
})
export class AssetAddVehicleComponent
  extends AddAssetBaseComponent
  implements OnChanges
{
  @Input() public vehicleAsset?: IAssetVehicleContext;

  protected vehicleType?: AssetVehicleType;
  protected readonly vehicleTypes: ISelectItem[] = [
    { id: 'car', title: 'Car', iconName: 'car-outline' },
    { id: 'motorbike', title: 'Motorbike', iconName: 'bicycle-outline' },
    { id: 'boat', title: 'Boat', iconName: 'boat-outline' },
  ];

  protected countryIso2 = 'IE';
  protected regNumber = '';
  protected vin = '';
  protected yearOfBuild = '';
  protected engine = '';
  protected engines?: string[];

  protected nctExpires = ''; // ISO date string 'YYYY-MM-DD'
  protected taxExpires = ''; // ISO date string 'YYYY-MM-DD'
  protected nextServiceDue = ''; // ISO date string 'YYYY-MM-DD'

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['space'] && this.space) {
      const a: IAssetVehicleContext = this.vehicleAsset ?? ({
        id: '',
        space: this.space ?? { id: '' },
        dbo: {
          status: 'draft',
          category: 'vehicles',
          extraType: 'vehicle',
          extra: {
            make: '',
            model: '',
            engineFuel: FuelTypes.unknown,
            engineType: EngineTypes.unknown,
          },
          type: this.vehicleType,
          possession: 'owning',
          createdAt: new Date().toISOString() as unknown as timestamp,
          createdBy: '-',
          updatedAt: new Date().toISOString() as unknown as timestamp,
          updatedBy: '-',
        },
      } as unknown as IAssetVehicleContext);
      this.vehicleAsset = { ...a, space: this.space };
    }
  }

  protected onAssetChanged(asset: IAssetContext): void {
    this.vehicleAsset = asset as IAssetVehicleContext;
  }

  onVehicleTypeChanged(): void {
    if (this.vehicleAsset?.dbo) {
      this.vehicleAsset = {
        ...this.vehicleAsset,
        dbo: {
          ...this.vehicleAsset.dbo,
          type: this.vehicleType,
          extraType: 'vehicle',
          extra: {
            make: '',
            model: '',
            regNumber: '',
            engineType: '',
            engineFuel: '',
          },
        } as unknown as IAssetVehicleContext['dbo'],
      };
    }
  }

  formatDate(value?: string | string[] | null): string {
    return value && !Array.isArray(value)
      ? format(parseISO(value), 'dd MMMM yyyy')
      : '';
  }

  protected submitVehicleForm(): void {
    if (!this.space) {
      throw 'no team context';
    }
    if (!this.vehicleType) {
      throw 'no vehicleType';
    }
    const assetDto = this.vehicleAsset?.dbo;
    if (!assetDto) {
      throw new Error('no asset');
    }
    this.isSubmitting = true;
    let request: {
      asset: IAssetDboBase<'vehicle', IAssetVehicleExtra & IAssetExtra>;
      spaceID: string;
    } = {
      asset: {
        ...assetDto,
        status: 'active',
        category: 'vehicles',
      } as unknown as IAssetDboBase<'vehicle', IAssetVehicleExtra & IAssetExtra>,
      spaceID: this.space?.id,
    };
    if (this.yearOfBuild) {
      request = {
        ...request,
        asset: { ...request.asset, yearOfBuild: +this.yearOfBuild },
      };
    }

    this.createAssetAndGoToAssetPage(request, this.space);
  }
}
