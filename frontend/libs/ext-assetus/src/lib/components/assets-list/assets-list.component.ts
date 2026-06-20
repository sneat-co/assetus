import {
  Component,
  Input,
  OnChanges,
  SimpleChanges,
  inject,
} from '@angular/core';
import {
  IonBadge,
  IonButton,
  IonButtons,
  IonIcon,
  IonItem,
  IonLabel,
  IonSpinner,
  ModalController,
} from '@ionic/angular/standalone';
import { IIdAndBrief } from '@sneat/core';
import {
  ASSET_SERVICE,
  AssetCategory,
  IAssetDbo,
  IAssetService,
} from '@sneat/extension-assetus-contract';
import { ErrorLogger, IErrorLogger } from '@sneat/core';
import { ISpaceContext } from '@sneat/space-models';
import { SpaceNavService } from '@sneat/space-services';
import { MileAgeDialogComponent } from '../mileage-dialog/mileage-dialog.component';

// Ported from legacy ext-assetus-components (legacy assetus components lib).
@Component({
  selector: 'assetus-assets-list',
  templateUrl: './assets-list.component.html',
  providers: [ModalController],
  imports: [
    IonItem,
    IonSpinner,
    IonLabel,
    IonButton,
    IonIcon,
    IonBadge,
    IonButtons,
  ],
})
export class AssetsListComponent implements OnChanges {
  private readonly errorLogger = inject<IErrorLogger>(ErrorLogger);
  private readonly assetService: IAssetService = inject(ASSET_SERVICE);
  private readonly spaceNavService = inject(SpaceNavService);
  private readonly modalCtrl = inject(ModalController);

  protected assets?: IIdAndBrief<IAssetDbo>[];
  protected mileAgeAsset?: IIdAndBrief<IAssetDbo>;

  @Input() allAssets?: IIdAndBrief<IAssetDbo>[];
  @Input({ required: true }) space?: ISpaceContext;
  @Input() assetType?: AssetCategory;
  @Input() filter = '';

  @Input() sorter: (
    a: IIdAndBrief<IAssetDbo>,
    b: IIdAndBrief<IAssetDbo>,
  ) => number = () => {
    return 0;
  };

  protected deletingIDs: string[] = [];

  ngOnChanges(changes: SimpleChanges): void {
    const { allAssets, assetType, filter } = this;
    if (!allAssets) {
      this.assets = undefined;
      return;
    }
    if (!allAssets.length) {
      this.assets = [];
      return;
    }
    const f = filter?.toLowerCase();
    if (!allAssets || (!filter && !assetType)) {
      this.assets = [...allAssets];
    } else {
      this.assets = allAssets?.filter(
        (asset) =>
          (!assetType || asset?.brief?.category === assetType) &&
          (!filter || asset?.brief?.name?.toLowerCase().includes(f) || -1),
      );
    }
    this.assets = this.assets ? [...this.assets].sort(this.sorter) : this.assets;
    console.log(
      'AssetsListComponent.ngOnChanges =>',
      changes,
      this.assetType,
      this.space,
      'allAssets:',
      this.allAssets,
      'filtered assets:',
      this.assets,
    );
  }

  protected goAsset(asset: IIdAndBrief<IAssetDbo>): void {
    if (!asset) {
      return;
    }
    if (!this.space) {
      this.errorLogger.logError(
        'can not navigate to asset page without team context',
      );
      return;
    }
    this.spaceNavService
      .navigateForwardToSpacePage(this.space, `asset/${asset.id}`, {
        state: { asset },
      })
      .catch(
        this.errorLogger.logErrorHandler('failed to navigate to asset page'),
      );
  }

  protected delete(event: Event, asset: IIdAndBrief<IAssetDbo>): void {
    event.stopPropagation();
    event.preventDefault();
    const { id, brief } = asset;
    this.deletingIDs.push(id);
    const deleteCompleted = () =>
      (this.deletingIDs = this.deletingIDs.filter((v) => v !== id));
    setTimeout(() => {
      if (
        !confirm(
          `Are you sure you want to delete this asset?

       ID: ${id}
       Title: ${brief?.name}

       This operation can not be undone.`,
        )
      ) {
        deleteCompleted();
        return;
      }
      this.assetService
        .removeAsset({ spaceID: this.space?.id || '', assetID: asset.id })
        .subscribe({
          next: () => {
            this.assets = this.assets?.filter((a) => a.id !== id);
          },
          error: this.errorLogger.logErrorHandler(
            'failed to delete an asset with ID=' + id,
          ),
          complete: deleteCompleted,
        });
    }, 1);
  }

  protected async addNewMilesAndFuel(
    event: Event,
    asset: IIdAndBrief<IAssetDbo>,
  ) {
    event.stopPropagation();
    event.preventDefault();

    const modal = await this.modalCtrl.create({
      component: MileAgeDialogComponent,
      componentProps: { asset },
    });
    await modal.present();
  }
}
