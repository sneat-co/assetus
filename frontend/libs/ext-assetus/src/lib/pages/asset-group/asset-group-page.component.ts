import { Component } from '@angular/core';
import {
  IonBackButton,
  IonButton,
  IonButtons,
  IonContent,
  IonHeader,
  IonIcon,
  IonLabel,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';
import { Period } from '@sneat/dto';
import { IAssetContext, IAssetDtoGroup } from '@sneat/mod-assetus-core';
import { AssetCardComponent, PeriodSegmentComponent } from '../../components';
import {
  SpaceBaseComponent,
  SpaceComponentBaseParams,
} from '@sneat/space-components';
import { ClassName } from '@sneat/ui';

// Asset-group detail page. Ported from the legacy `asset-group` page onto the
// MVP `SpaceBaseComponent` convention. The legacy page injected the never-
// implemented `IAssetGroupService`/`IAssetService` stubs; those are dropped.
// The group + its assets are read from navigation state, and rendered via the
// ported `PeriodSegmentComponent` and `AssetCardComponent`.
@Component({
  selector: 'assetus-asset-group-page',
  templateUrl: './asset-group-page.component.html',
  providers: [
    { provide: ClassName, useValue: 'AssetGroupPageComponent' },
    SpaceComponentBaseParams,
  ],
  imports: [
    PeriodSegmentComponent,
    AssetCardComponent,
    IonHeader,
    IonToolbar,
    IonButtons,
    IonBackButton,
    IonTitle,
    IonButton,
    IonIcon,
    IonLabel,
    IonContent,
  ],
})
export class AssetGroupPageComponent extends SpaceBaseComponent {
  protected period: Period = 'month';
  protected assetGroup?: IAssetDtoGroup;
  protected assets?: IAssetContext[];

  constructor() {
    super();
    const state = window.history.state;
    this.assetGroup = state?.['assetGroupDto'] as IAssetDtoGroup | undefined;
    this.assets = state?.['assets'] as IAssetContext[] | undefined;
  }

  protected periodChanged(period: Period): void {
    this.period = period;
  }
}
