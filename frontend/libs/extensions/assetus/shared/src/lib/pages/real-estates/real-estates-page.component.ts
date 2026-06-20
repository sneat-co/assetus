import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  IonBackButton,
  IonButton,
  IonButtons,
  IonContent,
  IonHeader,
  IonIcon,
  IonInput,
  IonItem,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';
import { IIdAndBrief } from '@sneat/core';
import { IAssetBrief } from '@sneat/extension-assetus-contract';
import { AssetsListComponent } from '../../components';
import {
  SpaceBaseComponent,
  SpaceComponentBaseParams,
} from '@sneat/space-components';
import { ClassName } from '@sneat/ui';

// Properties (real-estate) list page. Ported from the legacy `real-estates`
// page: the legacy `CommuneBasePageParams`/`AssetsBasePage`/`IAssetService`
// plumbing is replaced with the MVP `SpaceBaseComponent` convention, and the
// list renders through the ported `AssetsListComponent` filtered to the
// 'dwelling' category. Asset loading is delegated to that component.
@Component({
  selector: 'assetus-real-estates-page',
  templateUrl: './real-estates-page.component.html',
  providers: [
    { provide: ClassName, useValue: 'RealEstatesPageComponent' },
    SpaceComponentBaseParams,
  ],
  imports: [
    FormsModule,
    IonToolbar,
    IonButtons,
    IonHeader,
    IonBackButton,
    IonTitle,
    IonButton,
    IonIcon,
    IonContent,
    IonItem,
    IonInput,
    AssetsListComponent,
  ],
})
export class RealEstatesPageComponent extends SpaceBaseComponent {
  protected filter = '';
  protected assets?: IIdAndBrief<IAssetBrief>[];

  protected clearFilter(): void {
    this.filter = '';
  }

  protected goNew(): void {
    this.navigateForwardToSpacePage('new-asset', {
      state: { assetType: 'dwelling' },
    }).catch(
      this.errorLogger.logErrorHandler('Failed to navigate to new asset page'),
    );
  }
}
