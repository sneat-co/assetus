import { ChangeDetectionStrategy, Component, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  IonBackButton,
  IonButtons,
  IonCard,
  IonContent,
  IonHeader,
  IonItem,
  IonLabel,
  IonList,
  IonRadio,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';
import {
  AssetAddDocumentComponent,
  AssetAddDwellingComponent,
  AssetAddVehicleComponent,
} from '../../components';
import {
  SpaceBaseComponent,
  SpaceComponentBaseParams,
} from '@sneat/space-components';
import { ClassName } from '@sneat/ui';
import { AssetCategory } from '../../dto';

interface INewAssetCategory {
  readonly id: AssetCategory;
  readonly title: string;
}

// "New asset" wizard: pick an asset kind, then render the matching
// asset-add-* component. Ported from the legacy `new-asset` page; the legacy
// `@sneat/mod-assetus-core` `IAssetCategory` is replaced with the MVP
// `AssetCategory` ids ('vehicles', 'dwelling', 'document').
@Component({
  imports: [
    FormsModule,
    AssetAddDocumentComponent,
    AssetAddVehicleComponent,
    AssetAddDwellingComponent,
    IonHeader,
    IonToolbar,
    IonButtons,
    IonBackButton,
    IonTitle,
    IonContent,
    IonCard,
    IonItem,
    IonLabel,
    IonList,
    IonRadio,
  ],
  providers: [
    { provide: ClassName, useValue: 'NewAssetPageComponent' },
    SpaceComponentBaseParams,
  ],
  selector: 'assetus-new-asset-page',
  templateUrl: './new-asset-page.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NewAssetPageComponent extends SpaceBaseComponent {
  protected readonly category = signal<INewAssetCategory | undefined>(undefined);

  public categories: INewAssetCategory[] = [
    { id: 'vehicles', title: 'Vehicle' },
    { id: 'dwelling', title: 'Real estate' },
  ];

  constructor() {
    super();
    const assetType = window.history.state?.['assetType'] as
      | AssetCategory
      | undefined;
    if (assetType) {
      this.category.set(this.categories.find((c) => c.id === assetType));
    }
  }

  public selectCategory(category: INewAssetCategory): void {
    this.category.set(category);
  }
}
