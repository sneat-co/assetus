import { Component, computed, inject, signal } from '@angular/core';
import {
  IonBackButton,
  IonButton,
  IonButtons,
  IonContent,
  IonHeader,
  IonIcon,
  IonItem,
  IonLabel,
  IonList,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';
import { SpaceServiceModule } from '@sneat/space-services';
import {
  SpaceBaseComponent,
  SpaceComponentBaseParams,
} from '@sneat/space-components';
import { ClassName } from '@sneat/ui';
import { Subscription } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import {
  categoryOptions,
  conditionOptions,
  IAssetDbo,
  visibilityOptions,
} from '@sneat/extension-assetus-contract';
import { AssetService, AssetusCoreServicesModule } from '../../services';
import { AssetHistoryTimelineComponent } from '../../components/asset-history-timeline/asset-history-timeline.component';

// Read-only asset detail page (route `asset/:assetID`). Reads the asset live
// from Firestore via AssetService.watchAssetByID(space, id) so it auto-refreshes
// on any change. Displays the asset fields read-only, hosts the history timeline
// and an Edit button that navigates to the manage page (`asset/:assetID/edit`).
@Component({
  selector: 'assetus-asset-page',
  templateUrl: './asset-page.component.html',
  imports: [
    AssetusCoreServicesModule,
    SpaceServiceModule,
    AssetHistoryTimelineComponent,
    IonHeader,
    IonToolbar,
    IonTitle,
    IonButtons,
    IonBackButton,
    IonButton,
    IonIcon,
    IonContent,
    IonList,
    IonItem,
    IonLabel,
  ],
  providers: [
    { provide: ClassName, useValue: 'AssetPageComponent' },
    SpaceComponentBaseParams,
  ],
})
export class AssetPageComponent extends SpaceBaseComponent {
  private readonly assetService = inject(AssetService);

  protected readonly $assetID = signal<string | undefined>(undefined);
  protected readonly $asset = signal<IAssetDbo | undefined>(undefined);

  // Human labels for the enum values shown on the read-only view.
  protected readonly $categoryLabel = computed(() =>
    this.labelFor(categoryOptions, this.$asset()?.category),
  );
  protected readonly $conditionLabel = computed(() =>
    this.labelFor(conditionOptions, this.$asset()?.condition),
  );
  protected readonly $visibilityLabel = computed(() =>
    this.labelFor(visibilityOptions, this.$asset()?.visibility),
  );

  private assetSub?: Subscription;

  constructor() {
    super();
    this.route.paramMap.pipe(takeUntil(this.destroyed$)).subscribe((params) => {
      const assetID = params.get('assetID') ?? undefined;
      this.$assetID.set(assetID);
      this.loadAsset(assetID);
    });
  }

  private labelFor(
    opts: readonly { value: string; label: string }[],
    value?: string,
  ): string {
    return opts.find((o) => o.value === value)?.label ?? value ?? '';
  }

  private loadAsset(assetID?: string): void {
    if (!assetID || !this.space?.id) {
      return;
    }
    // Live read from Firestore: re-applies the asset on every change so the page
    // auto-refreshes (edit, transfer, external edit) — consistent with the live
    // assets list. Re-subscribing replaces any prior subscription.
    this.assetSub?.unsubscribe();
    this.assetSub = this.assetService
      .watchAssetByID(this.space, assetID)
      .pipe(takeUntil(this.destroyed$))
      .subscribe({
        next: (ctx) => {
          if (ctx.dbo) {
            this.applyAsset(ctx.dbo);
          }
        },
        error: this.errorLogger.logErrorHandler('Failed to load asset'),
      });
  }

  private applyAsset(asset: IAssetDbo): void {
    this.$asset.set(asset);
  }

  protected edit(): void {
    const assetID = this.$assetID();
    if (!assetID) {
      return;
    }
    this.navigateForwardToSpacePage(`asset/${assetID}/edit`).catch(
      this.errorLogger.logErrorHandler('Failed to navigate to edit asset'),
    );
  }
}
