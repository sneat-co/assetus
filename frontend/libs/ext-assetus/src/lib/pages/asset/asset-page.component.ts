import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  AlertController,
  IonBackButton,
  IonButton,
  IonButtons,
  IonContent,
  IonHeader,
  IonIcon,
  IonInput,
  IonItem,
  IonLabel,
  IonList,
  IonSelect,
  IonSelectOption,
  IonTextarea,
  IonTitle,
  IonToolbar,
  ModalController,
} from '@ionic/angular/standalone';
import { SpaceServiceModule } from '@sneat/space-services';
import {
  SpaceBaseComponent,
  SpaceComponentBaseParams,
} from '@sneat/space-components';
import { ClassName } from '@sneat/ui';
import { takeUntil } from 'rxjs/operators';
import {
  AssetCategory,
  AssetCondition,
  AssetVisibility,
  categoryOptions,
  conditionOptions,
  IAssetDbo,
  visibilityOptions,
} from '../../dto';
import { AssetService, AssetusCoreServicesModule } from '../../services';
import { AssetHistoryTimelineComponent } from '../../components/asset-history-timeline/asset-history-timeline.component';
import { TransferAssetComponent } from '../transfer/transfer-asset.component';

// Asset detail + edit page (Task 9). Loads the asset, lets a member edit
// name/description/category/condition/visibility and save via updateAsset().
// Hosts the remove UI (Task 10: soft-archive + confirmed hard-delete), the
// transfer entry point and the history timeline (Task 11).
@Component({
  selector: 'assetus-asset-page',
  templateUrl: './asset-page.component.html',
  imports: [
    FormsModule,
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
    IonInput,
    IonTextarea,
    IonSelect,
    IonSelectOption,
  ],
  providers: [
    { provide: ClassName, useValue: 'AssetPageComponent' },
    SpaceComponentBaseParams,
  ],
})
export class AssetPageComponent extends SpaceBaseComponent {
  private readonly assetService = inject(AssetService);
  private readonly modalCtrl = inject(ModalController);
  private readonly alertCtrl = inject(AlertController);

  protected readonly categoryOptions = categoryOptions;
  protected readonly conditionOptions = conditionOptions;
  protected readonly visibilityOptions = visibilityOptions;

  protected readonly $assetID = signal<string | undefined>(undefined);
  protected readonly $asset = signal<IAssetDbo | undefined>(undefined);
  protected saving = false;

  // Editable form fields.
  protected name = '';
  protected description = '';
  protected category: AssetCategory = 'other';
  protected condition: AssetCondition = 'good';
  protected visibility: AssetVisibility = 'private';

  constructor() {
    super();
    this.route.paramMap.pipe(takeUntil(this.destroyed$)).subscribe((params) => {
      const assetID = params.get('assetID') ?? undefined;
      this.$assetID.set(assetID);
      this.loadAsset(assetID);
    });
  }

  private loadAsset(assetID?: string): void {
    if (!assetID || !this.space?.id) {
      return;
    }
    this.assetService
      .getAsset(this.space.id, assetID)
      .pipe(takeUntil(this.destroyed$))
      .subscribe({
        next: (res) => this.applyAsset(res.asset),
        error: this.errorLogger.logErrorHandler('Failed to load asset'),
      });
  }

  private applyAsset(asset: IAssetDbo): void {
    this.$asset.set(asset);
    this.name = asset.name;
    this.description = asset.description ?? '';
    this.category = asset.category;
    this.condition = asset.condition;
    this.visibility = asset.visibility;
  }

  protected save(): void {
    const assetID = this.$assetID();
    if (!assetID || !this.space?.id) {
      return;
    }
    this.saving = true;
    this.assetService
      .updateAsset({
        spaceID: this.space.id,
        assetID,
        name: this.name.trim(),
        description: this.description.trim() || undefined,
        category: this.category,
        condition: this.condition,
        visibility: this.visibility,
      })
      .subscribe({
        next: (res) => {
          this.saving = false;
          this.applyAsset(res.asset);
        },
        error: (err) => {
          this.saving = false;
          this.errorLogger.logError(err, 'Failed to update asset');
        },
      });
  }

  // Soft-archive (Task 10): keeps history, moves the asset to the Archived view.
  protected archive(): void {
    this.remove(false);
  }

  // Hard-delete (Task 10): destroys the record; requires explicit confirmation.
  protected async confirmHardDelete(): Promise<void> {
    const alert = await this.alertCtrl.create({
      header: 'Delete asset permanently?',
      message:
        'This removes the asset and its history for good. This cannot be undone.',
      buttons: [
        { text: 'Cancel', role: 'cancel' },
        {
          text: 'Delete',
          role: 'destructive',
          handler: () => this.remove(true),
        },
      ],
    });
    await alert.present();
  }

  private remove(hardDelete: boolean): void {
    const assetID = this.$assetID();
    if (!assetID || !this.space?.id) {
      return;
    }
    this.assetService
      .removeAsset({ spaceID: this.space.id, assetID, hardDelete })
      .subscribe({
        next: () => {
          this.navigateForwardToSpacePage('assets').catch(
            this.errorLogger.logErrorHandler('Failed to navigate to assets'),
          );
        },
        error: this.errorLogger.logErrorHandler('Failed to remove asset'),
      });
  }

  // Transfer flow (Task 11): opens the destination-space picker dialog.
  protected async transfer(): Promise<void> {
    const assetID = this.$assetID();
    if (!assetID || !this.space?.id) {
      return;
    }
    const modal = await this.modalCtrl.create({
      component: TransferAssetComponent,
      componentProps: { spaceID: this.space.id, assetID },
    });
    await modal.present();
    const { data } = await modal.onDidDismiss();
    if (data) {
      this.loadAsset(assetID);
    }
  }
}
