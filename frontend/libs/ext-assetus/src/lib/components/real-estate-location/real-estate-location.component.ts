import { Component, Input } from '@angular/core';
import {
  IonButton,
  IonButtons,
  IonIcon,
  IonInput,
  IonItem,
  IonLabel,
  IonTextarea,
} from '@ionic/angular/standalone';
import { IAssetContext } from '@sneat/mod-assetus-core';

// Ported from @sneat/ext-assetus-components (legacy assetus components lib).
@Component({
  selector: 'sneat-real-estate-location',
  templateUrl: './real-estate-location.component.html',
  imports: [
    IonItem,
    IonLabel,
    IonTextarea,
    IonButtons,
    IonButton,
    IonIcon,
    IonInput,
  ],
})
export class RealEstateLocationComponent {
  @Input() asset?: IAssetContext;
}
