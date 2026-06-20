import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  IonButton,
  IonButtons,
  IonCard,
  IonDatetime,
  IonIcon,
  IonInput,
  IonItem,
  IonPopover,
} from '@ionic/angular/standalone';
import { ClassName, ISelectItem, SelectFromListComponent } from '@sneat/ui';
import { timestamp } from '@sneat/dto';
import {
  AssetPossession,
  AssetVehicleType,
  IAssetDboBase,
  IAssetExtra,
  IAssetDocumentExtra,
} from '@sneat/extension-assetus-contract';
import {
  IAssetContext,
  IAssetDocumentContext,
} from '@sneat/extension-assetus-contract';
import { SpaceComponentBaseParams } from '@sneat/space-components';
import { format, parseISO } from 'date-fns';
import { AddAssetBaseComponent } from '../add-asset-base.component';

// Ported from legacy ext-assetus-components (legacy assetus components lib).
@Component({
  selector: 'assetus-asset-add-document',
  templateUrl: './asset-add-document.component.html',
  providers: [
    {
      provide: ClassName,
      useValue: 'AssetAddVehicleComponent',
    },
    SpaceComponentBaseParams,
  ],
  imports: [
    IonCard,
    SelectFromListComponent,
    FormsModule,
    IonItem,
    IonInput,
    IonButtons,
    IonButton,
    IonIcon,
    IonPopover,
    IonDatetime,
  ],
})
export class AssetAddDocumentComponent
  extends AddAssetBaseComponent
  implements OnChanges
{
  @Input() public documentAsset?: IAssetDocumentContext;

  protected documentType?: AssetVehicleType;
  protected readonly documentTypes: ISelectItem[] = [
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
    const spaceChanges = changes['space'];
    if (spaceChanges && this.space) {
      const a: IAssetDocumentContext = this.documentAsset ?? ({
        id: '',
        space: this.space ?? { id: '' },
        dbo: {
          status: 'draft',
          // Live AssetCategory for a document (the legacy add-document component
          // seeded 'vehicle' here — a pre-existing inconsistency with its
          // 'document' extraType; aligned to the coherent live value).
          category: 'document',
          extraType: 'document',
          extra: {},
          type: this.documentType,
          title: '',
          possession: undefined as unknown as AssetPossession,
          createdAt: new Date().toISOString() as unknown as timestamp,
          createdBy: '-',
          updatedAt: new Date().toISOString() as unknown as timestamp,
          updatedBy: '-',
        },
      } as unknown as IAssetDocumentContext);
      this.documentAsset = { ...a, space: this.space };
    }
  }

  protected onAssetChanged(_asset: IAssetContext): void {
    // TODO: Implement asset changed logic
  }

  protected onVehicleTypeChanged(): void {
    if (this.documentAsset?.dbo) {
      this.documentAsset = {
        ...this.documentAsset,
        dbo: {
          ...this.documentAsset.dbo,
          type: this.documentType,
          extraType: 'document',
          extra: {},
        } as unknown as IAssetDocumentContext['dbo'],
      };
    }
  }

  protected formatDate(value?: string | string[] | null): string {
    return value && !Array.isArray(value)
      ? format(parseISO(value), 'dd MMMM yyyy')
      : '';
  }

  protected submitDocumentForm(): void {
    if (!this.space) {
      throw 'no team context';
    }
    if (!this.documentType) {
      throw 'no vehicleType';
    }
    const assetDto = this.documentAsset?.dbo;
    if (!assetDto) {
      throw new Error('no asset');
    }
    this.isSubmitting = true;
    let request: {
      asset: IAssetDboBase<'document', IAssetDocumentExtra & IAssetExtra>;
      spaceID: string;
    } = {
      asset: {
        ...assetDto,
        status: 'active',
        category: 'document',
      } as unknown as IAssetDboBase<'document', IAssetDocumentExtra & IAssetExtra>,
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
