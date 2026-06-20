import { Component, Input, OnInit, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  IonButton,
  IonButtons,
  IonContent,
  IonHeader,
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
import { ErrorLogger, IErrorLogger } from '@sneat/core';
import {
  AssetCategory,
  AssetCondition,
  AssetVisibility,
  categoryOptions,
  conditionOptions,
  defaultVisibilityForSpaceType,
  visibilityOptions,
  ICreateAssetRequest,
} from '@sneat/extension-assetus-contract';
import { AssetService } from '../../services';

// Create-asset dialog (Task 8). Collects name, category, condition and a
// visibility that defaults to the owning space's default but can be overridden.
// Optional metadata (description/location) is included to exercise the backend
// contract. Calls AssetService.createAsset() and dismisses with the new id.
@Component({
  selector: 'assetus-new-asset-dialog',
  templateUrl: './new-asset-dialog.component.html',
  imports: [
    FormsModule,
    IonHeader,
    IonToolbar,
    IonTitle,
    IonButtons,
    IonButton,
    IonContent,
    IonList,
    IonItem,
    IonLabel,
    IonInput,
    IonTextarea,
    IonSelect,
    IonSelectOption,
  ],
})
export class NewAssetDialogComponent implements OnInit {
  private readonly modalCtrl = inject(ModalController);
  private readonly assetService = inject(AssetService);
  private readonly errorLogger = inject<IErrorLogger>(ErrorLogger);

  @Input({ required: true }) spaceID!: string;
  @Input() spaceType?: string;

  protected readonly categoryOptions = categoryOptions;
  protected readonly conditionOptions = conditionOptions;
  protected readonly visibilityOptions = visibilityOptions;

  protected name = '';
  protected category: AssetCategory = 'other';
  protected condition: AssetCondition = 'good';
  protected visibility: AssetVisibility = 'private';
  protected description = '';
  protected location = '';
  protected creating = false;

  ngOnInit(): void {
    // Default visibility to the space default; the user can still override it.
    this.visibility = defaultVisibilityForSpaceType(this.spaceType);
  }

  protected create(): void {
    const name = this.name.trim();
    if (!name) {
      return;
    }
    const request: ICreateAssetRequest = {
      spaceID: this.spaceID,
      name,
      category: this.category,
      condition: this.condition,
      visibility: this.visibility,
    };
    if (this.description.trim()) {
      request.description = this.description.trim();
    }
    if (this.location.trim()) {
      request.location = this.location.trim();
    }
    this.creating = true;
    this.assetService.createAsset(request).subscribe({
      next: (res) => {
        this.creating = false;
        this.dismiss(res.id);
      },
      error: (err) => {
        this.creating = false;
        this.errorLogger.logError(err, 'Failed to create asset');
      },
    });
  }

  protected cancel(): void {
    this.dismiss();
  }

  private dismiss(createdAssetID?: string): void {
    this.modalCtrl
      .dismiss(createdAssetID, createdAssetID ? 'created' : 'cancel')
      .catch(this.errorLogger.logError);
  }
}
