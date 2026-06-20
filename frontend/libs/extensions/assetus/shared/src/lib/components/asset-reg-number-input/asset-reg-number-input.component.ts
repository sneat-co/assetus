import {
  Component,
  EventEmitter,
  Input,
  OnChanges,
  Output,
  SimpleChanges,
  ViewChild,
} from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import {
  IonButton,
  IonButtons,
  IonIcon,
  IonInput,
  IonItem,
  IonLabel,
} from '@ionic/angular/standalone';
import { ISpaceContext } from '@sneat/space-models';

// Ported from legacy ext-assetus-components (legacy assetus components lib).
@Component({
  selector: 'assetus-asset-reg-number',
  templateUrl: 'asset-reg-number-input.component.html',
  imports: [
    ReactiveFormsModule,
    IonItem,
    IonInput,
    IonButtons,
    IonButton,
    IonLabel,
    IonIcon,
  ],
})
export class AssetRegNumberInputComponent implements OnChanges {
  @Input({ required: true }) space?: ISpaceContext;
  @Input({ required: true }) assetID?: string;
  @Input({ required: true }) countyID?: string;
  @Input({ required: true }) regNumber?: string = '';
  @Input() hideSaveButton = false;
  @Input() placeholder = '';

  @Output() regNumberChange = new EventEmitter<string>();

  @Input() isSkipped = false;
  @Output() isSkippedChange = new EventEmitter<boolean>();

  @ViewChild(IonInput, { static: true }) regNumberInput!: IonInput;

  protected validatedRegNumber?: string;

  protected isSaving = false;
  protected readonly regNumberControl = new FormControl('');

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['regNumber'] && !this.regNumberControl.dirty) {
      this.regNumberControl.setValue(this.regNumber || '');
    }
  }

  protected get showSave(): boolean {
    return (
      !this.hideSaveButton &&
      this.regNumberControl.dirty &&
      (this.regNumber || '') !== this.regNumberControl.value
    );
  }

  protected get isValidated(): boolean {
    return this.validatedRegNumber === this.regNumberControl.value?.trim();
  }
  protected validate(): void {
    this.validatedRegNumber = this.regNumberControl.value?.trim();
    this.skipOrNext();
  }

  protected skipOrNext(): void {
    const regNumber = this.regNumberControl.value?.trim();
    if (regNumber) {
      this.regNumberChange.emit(regNumber);
    }
    this.isSkippedChange.emit();
  }

  protected submit(): void {
    const space = this.space;

    if (!space?.id || !this.assetID) {
      return;
    }

    // The reg-number is persisted via the parent's full-asset save: on the live
    // assetus backend the registration number lives in the vehicle's typed
    // `extra` and update_asset is a full-asset update (name/category/condition/
    // visibility required), so a standalone reg-number-only HTTP update — the
    // dead legacy `assets/update_asset` contract — no longer exists. We emit the
    // value upward (the parent vehicle-card folds it into the asset and saves)
    // and mark the control pristine to reflect the committed value.
    this.regNumberChange.emit(this.regNumberControl.value || '');
    this.regNumberControl.markAsPristine();
  }

  public focusToRegNumberInput(): void {
    this.regNumberInput
      .setFocus()
      .catch((e) => console.error('Failed to focus to reg number input', e));
  }
}
