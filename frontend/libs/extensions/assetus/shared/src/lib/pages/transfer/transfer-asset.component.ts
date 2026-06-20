import { Component, Input, OnInit, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import {
  IonButton,
  IonButtons,
  IonContent,
  IonHeader,
  IonItem,
  IonLabel,
  IonList,
  IonSelect,
  IonSelectOption,
  IonTitle,
  IonToolbar,
  ModalController,
} from '@ionic/angular/standalone';
import { ErrorLogger, IErrorLogger, IIdAndBrief } from '@sneat/core';
import { IUserSpaceBrief } from '@sneat/auth-models';
import { SneatUserService } from '@sneat/auth-core';
import { zipMapBriefsWithIDs } from '@sneat/space-models';
import {
  ASSET_SERVICE,
  IAssetService,
} from '@sneat/extension-assetus-contract';

// Transfer flow (Task 11): pick a destination space and reassign the asset to
// it via AssetService.transferAsset(). The backend records the append-only
// Transferred history event with prior/new owner.
@Component({
  selector: 'assetus-transfer-asset',
  templateUrl: './transfer-asset.component.html',
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
    IonSelect,
    IonSelectOption,
  ],
})
export class TransferAssetComponent implements OnInit {
  private readonly modalCtrl = inject(ModalController);
  private readonly assetService: IAssetService = inject(ASSET_SERVICE);
  private readonly userService = inject(SneatUserService);
  private readonly errorLogger = inject<IErrorLogger>(ErrorLogger);

  @Input({ required: true }) spaceID!: string;
  @Input({ required: true }) assetID!: string;

  protected readonly $spaces = signal<readonly IIdAndBrief<IUserSpaceBrief>[]>(
    [],
  );
  protected toSpaceID?: string;
  protected transferring = false;

  ngOnInit(): void {
    this.userService.userState.subscribe({
      next: (state) => {
        const spaces = state?.record
          ? zipMapBriefsWithIDs(state.record.spaces) || []
          : [];
        // A transfer must change owner, so exclude the current owning space.
        this.$spaces.set(spaces.filter((s) => s.id !== this.spaceID));
      },
      error: this.errorLogger.logErrorHandler('Failed to load spaces'),
    });
  }

  protected transfer(): void {
    if (!this.toSpaceID) {
      return;
    }
    this.transferring = true;
    this.assetService
      .transferAsset({
        spaceID: this.spaceID,
        assetID: this.assetID,
        toSpaceID: this.toSpaceID,
      })
      .subscribe({
        next: (res) => {
          this.transferring = false;
          this.dismiss(res.owner.spaceID);
        },
        error: (err) => {
          this.transferring = false;
          this.errorLogger.logError(err, 'Failed to transfer asset');
        },
      });
  }

  protected cancel(): void {
    this.dismiss();
  }

  private dismiss(transferredToSpaceID?: string): void {
    this.modalCtrl
      .dismiss(transferredToSpaceID, transferredToSpaceID ? 'transferred' : 'cancel')
      .catch(this.errorLogger.logError);
  }
}
