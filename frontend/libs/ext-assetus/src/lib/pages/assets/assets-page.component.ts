import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  IonBadge,
  IonButton,
  IonButtons,
  IonContent,
  IonHeader,
  IonIcon,
  IonItem,
  IonLabel,
  IonList,
  IonMenuButton,
  IonNote,
  IonSegment,
  IonSegmentButton,
  IonTitle,
  IonToolbar,
  ModalController,
} from '@ionic/angular/standalone';
import { ContactusServicesModule } from '@sneat/contactus-services';
import { SpaceServiceModule } from '@sneat/space-services';
import {
  SpaceBaseComponent,
  SpaceComponentBaseParams,
} from '@sneat/space-components';
import { ClassName } from '@sneat/ui';
import { takeUntil } from 'rxjs/operators';
import { Subscription } from 'rxjs';
import {
  ASSET_SERVICE,
  IAssetDbo,
  IAssetService,
} from '@sneat/extension-assetus-contract';
import { NewAssetDialogComponent } from './new-asset-dialog.component';

interface IAssetRow {
  id: string;
  dbo: IAssetDbo;
}

type AssetsFilter = 'active' | 'archived';

// Per-space assets list page (Task 8). Reads assets live via
// AssetService.watchAssets() and offers a "New asset" dialog. An Active/Archived
// segment filters out soft-archived assets (Task 10).
@Component({
  selector: 'assetus-assets-page',
  templateUrl: './assets-page.component.html',
  imports: [
    FormsModule,
    ContactusServicesModule,
    SpaceServiceModule,
    IonHeader,
    IonToolbar,
    IonTitle,
    IonButtons,
    IonMenuButton,
    IonButton,
    IonIcon,
    IonContent,
    IonList,
    IonItem,
    IonLabel,
    IonNote,
    IonBadge,
    IonSegment,
    IonSegmentButton,
  ],
  providers: [
    { provide: ClassName, useValue: 'AssetsPageComponent' },
    SpaceComponentBaseParams,
  ],
})
export class AssetsPageComponent extends SpaceBaseComponent {
  private readonly assetService: IAssetService = inject(ASSET_SERVICE);
  private readonly modalCtrl = inject(ModalController);

  protected readonly $assets = signal<IAssetRow[]>([]);
  protected readonly $filter = signal<AssetsFilter>('active');

  private assetsSub?: Subscription;

  constructor() {
    super();
    // Re-subscribe to the space's assets whenever the active space changes.
    this.spaceIDChanged$
      .pipe(takeUntil(this.destroyed$))
      .subscribe((spaceID) => this.watchAssets(spaceID));
  }

  protected get visibleAssets(): IAssetRow[] {
    const filter = this.$filter();
    return this.$assets().filter((a) =>
      filter === 'archived'
        ? a.dbo.status === 'archived'
        : a.dbo.status !== 'archived',
    );
  }

  protected setFilter(event: Event): void {
    const value = (event as CustomEvent).detail.value as AssetsFilter;
    this.$filter.set(value);
  }

  private watchAssets(spaceID?: string): void {
    this.assetsSub?.unsubscribe();
    if (!spaceID) {
      this.$assets.set([]);
      return;
    }
    this.assetsSub = this.assetService
      .watchAssets(spaceID)
      .pipe(takeUntil(this.destroyed$))
      .subscribe({
        next: (rows) => this.$assets.set(rows),
        error: this.errorLogger.logErrorHandler('Failed to watch assets'),
      });
  }

  protected goAsset(row: IAssetRow): void {
    this.navigateForwardToSpacePage(`asset/${row.id}`, {
      state: { asset: row.dbo, space: this.space },
    }).catch(this.errorLogger.logErrorHandler('Failed to navigate to asset'));
  }

  protected async newAsset(): Promise<void> {
    if (!this.space) {
      return;
    }
    const modal = await this.modalCtrl.create({
      component: NewAssetDialogComponent,
      componentProps: {
        spaceID: this.space.id,
        spaceType: this.space.type,
      },
    });
    await modal.present();
    await modal.onDidDismiss();
  }
}
