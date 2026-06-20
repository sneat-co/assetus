import { Component, signal } from '@angular/core';
import {
  IonBackButton,
  IonBadge,
  IonButtons,
  IonContent,
  IonHeader,
  IonItem,
  IonItemDivider,
  IonItemGroup,
  IonLabel,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';
import { IAssetContext } from '@sneat/extension-assetus-contract';
import { RealEstateLocationComponent } from '../../components';
import {
  SpaceBaseComponent,
  SpaceComponentBaseParams,
} from '@sneat/space-components';
import { ClassName } from '@sneat/ui';

// Real-estate (property) detail page. Ported from the legacy `real-estate`
// page onto the MVP `SpaceBaseComponent` convention. Shows the property
// location (via the ported `RealEstateLocationComponent`) plus a static
// finance section. The legacy landlords/tenants contact groups
// (`asset-contacts-group`) are deferred to the liabilities/service-provider
// sibling Feature (see COMPONENTS-COVERAGE.md), so they are not rendered here.
@Component({
  selector: 'assetus-real-estate-page',
  templateUrl: './real-estate-page.component.html',
  providers: [
    { provide: ClassName, useValue: 'RealEstatePageComponent' },
    SpaceComponentBaseParams,
  ],
  imports: [
    IonHeader,
    IonToolbar,
    IonButtons,
    IonBackButton,
    IonTitle,
    IonContent,
    IonItemGroup,
    IonItemDivider,
    IonItem,
    IonLabel,
    IonBadge,
    RealEstateLocationComponent,
  ],
})
export class RealEstatePageComponent extends SpaceBaseComponent {
  protected readonly $asset = signal<IAssetContext | undefined>(
    window.history.state?.['asset'] as IAssetContext | undefined,
  );
}
