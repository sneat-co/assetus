import {
  ChangeDetectionStrategy,
  Component,
  Input,
  OnChanges,
  inject,
  signal,
} from '@angular/core';
import {
  IonIcon,
  IonItem,
  IonLabel,
  IonList,
  IonListHeader,
  IonNote,
  IonText,
} from '@ionic/angular/standalone';
import { ErrorLogger, IErrorLogger } from '@sneat/core';
import {
  ASSET_SERVICE,
  IAssetService,
  IHistoryEvent,
} from '@sneat/extension-assetus-contract';

// Append-only asset history timeline (Task 11). Renders every recorded event,
// including the Transferred entry with its prior and new owner.
@Component({
  selector: 'assetus-asset-history-timeline',
  templateUrl: './asset-history-timeline.component.html',
  imports: [
    IonList,
    IonListHeader,
    IonItem,
    IonLabel,
    IonNote,
    IonIcon,
    IonText,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AssetHistoryTimelineComponent implements OnChanges {
  private readonly assetService: IAssetService = inject(ASSET_SERVICE);
  private readonly errorLogger = inject<IErrorLogger>(ErrorLogger);

  @Input() spaceID?: string;
  @Input() assetID?: string;

  protected readonly $events = signal<IHistoryEvent[]>([]);

  ngOnChanges(): void {
    if (!this.spaceID || !this.assetID) {
      this.$events.set([]);
      return;
    }
    this.assetService.getHistory(this.spaceID, this.assetID).subscribe({
      next: (res) => this.$events.set(res.events ?? []),
      error: this.errorLogger.logErrorHandler('Failed to load asset history'),
    });
  }

  protected iconFor(type: IHistoryEvent['type']): string {
    switch (type) {
      case 'purchased':
        return 'cart';
      case 'repaired':
        return 'build';
      case 'transferred':
        return 'swap-horizontal';
      case 'sold':
        return 'cash';
      case 'donated':
        return 'gift';
      case 'lost':
        return 'help-circle';
      default:
        return 'ellipse';
    }
  }
}
