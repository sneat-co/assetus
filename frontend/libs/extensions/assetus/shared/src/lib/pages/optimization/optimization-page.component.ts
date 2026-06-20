import { Component } from '@angular/core';
import {
  IonBackButton,
  IonButton,
  IonButtons,
  IonCard,
  IonCardContent,
  IonCardHeader,
  IonCardTitle,
  IonContent,
  IonHeader,
  IonIcon,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';
import {
  SpaceBaseComponent,
  SpaceComponentBaseParams,
} from '@sneat/space-components';
import { ClassName } from '@sneat/ui';

// Savings/optimization suggestions page. Ported from the legacy
// `optimization` page; the legacy `CommuneBasePage`/`CommuneTopPage` plumbing
// is replaced with the MVP `SpaceBaseComponent` convention. Content is still
// the static demo card from the legacy page (no optimization service yet).
@Component({
  selector: 'assetus-optimization-page',
  templateUrl: './optimization-page.component.html',
  providers: [
    { provide: ClassName, useValue: 'OptimizationPageComponent' },
    SpaceComponentBaseParams,
  ],
  imports: [
    IonHeader,
    IonToolbar,
    IonButtons,
    IonBackButton,
    IonTitle,
    IonContent,
    IonCard,
    IonCardHeader,
    IonCardTitle,
    IonCardContent,
    IonButton,
    IonIcon,
  ],
})
export class OptimizationPageComponent extends SpaceBaseComponent {}
