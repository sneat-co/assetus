import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { RouterModule } from '@angular/router';
import {
  IonCard,
  IonCardContent,
  IonCardHeader,
  IonCardTitle,
  IonItem,
  IonLabel,
} from '@ionic/angular/standalone';
import { Period } from '@sneat/dto';
import { IAssetContext } from '../../contexts';

// Ported from legacy ext-assetus-components (legacy assetus components lib).
@Component({
  selector: 'assetus-asset-card',
  templateUrl: './asset-card.component.html',
  imports: [
    RouterModule,
    IonCard,
    IonCardHeader,
    IonCardTitle,
    IonCardContent,
    IonItem,
    IonLabel,
  ],
})
export class AssetCardComponent implements OnChanges {
  @Input() period?: Period;
  @Input({ required: true }) asset?: IAssetContext;

  protected segment: 'expenses' | 'income' = 'expenses';

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['asset'] && this.asset) {
      // Default the segment from the live backend `financialDirection` field
      // (IAssetDbo). The legacy nested {incomes,expenses} totals shape this
      // previously inspected no longer exists on the live assetus backend
      // (totals is now IMoney[]).
      if (this.asset?.dbo?.financialDirection === 'income') {
        this.segment = 'income';
      }
    }
  }

  segmentChanged(ev: CustomEvent): void {
    this.segment = ev.detail.value;
  }
}
