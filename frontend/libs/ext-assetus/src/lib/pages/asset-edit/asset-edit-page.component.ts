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
import { take, takeUntil } from 'rxjs/operators';
import {
  AssetCategory,
  AssetCondition,
  AssetVisibility,
  categoryOptions,
  conditionOptions,
  IAssetDbo,
  visibilityOptions,
} from '@sneat/extension-assetus-contract';
import { AssetService, AssetusCoreServicesModule } from '../../services';
import { TransferAssetComponent } from '../transfer/transfer-asset.component';

// Asset manage/edit page (route `asset/:assetID/edit`). Loads the asset ONCE for
// prefill (one-shot, so in-progress edits aren't clobbered by live updates),
// lets a member edit name/description/category/condition/visibility and save via
// updateAsset(). Also hosts the destructive actions: transfer, archive
// (soft-delete) and confirmed hard-delete.
@Component({
  selector: 'assetus-asset-edit-page',
  templateUrl: './asset-edit-page.component.html',
  imports: [
    FormsModule,
    AssetusCoreServicesModule,
    SpaceServiceModule,
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
    IonInput,
    IonTextarea,
    IonSelect,
    IonSelectOption,
  ],
  providers: [
    { provide: ClassName, useValue: 'AssetEditPageComponent' },
    SpaceComponentBaseParams,
  ],
})
export class AssetEditPageComponent extends SpaceBaseComponent {
  private readonly assetService = inject(AssetService);
  private readonly modalCtrl = inject(ModalController);
  private readonly alertCtrl = inject(AlertController);

  protected readonly categoryOptions = categoryOptions;
  protected readonly conditionOptions = conditionOptions;
  protected readonly visibilityOptions = visibilityOptions;

  protected readonly $assetID = signal<string | undefined>(undefined);
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
    // One-shot read for prefill: take(1) so subsequent live updates don't clobber
    // the user's in-progress edits.
    this.assetService
      .watchAssetByID(this.space, assetID)
      .pipe(take(1), takeUntil(this.destroyed$))
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
    this.name = asset.name;
    this.description = asset.description ?? '';
    this.category = asset.category;
    this.condition = asset.condition;
    this.visibility = asset.visibility;
  }

  private backToView(): void {
    const assetID = this.$assetID();
    if (!assetID) {
      return;
    }
    this.navigateForwardToSpacePage(`asset/${assetID}`).catch(
      this.errorLogger.logErrorHandler('Failed to navigate to asset'),
    );
  }

  protected cancel(): void {
    this.backToView();
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
        next: () => {
          this.saving = false;
          this.backToView();
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
