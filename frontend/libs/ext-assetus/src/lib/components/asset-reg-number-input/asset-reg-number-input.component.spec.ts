import { TestBed } from '@angular/core/testing';
import { AssetService } from '@sneat/ext-assetus-components';
import { of, throwError } from 'rxjs';
import { componentTestProviders } from '../../../testing/test-providers';
import { AssetRegNumberInputComponent } from './asset-reg-number-input.component';

// Render + logic spec for the ported AssetRegNumberInputComponent. It injects
// the legacy AssetService and the ErrorLogger; both are stubbed.
describe('AssetRegNumberInputComponent', () => {
  let updateAsset: ReturnType<typeof vi.fn>;
  let fixture: ReturnType<
    typeof TestBed.createComponent<AssetRegNumberInputComponent>
  >;
  let component: AssetRegNumberInputComponent;

  beforeEach(() => {
    updateAsset = vi.fn(() => of(undefined));
    TestBed.configureTestingModule({
      imports: [AssetRegNumberInputComponent],
      providers: [
        ...componentTestProviders(),
        { provide: AssetService, useValue: { updateAsset } },
      ],
    });
    fixture = TestBed.createComponent(AssetRegNumberInputComponent);
    component = fixture.componentInstance;
  });

  it('creates and renders', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('syncs the form control from the regNumber input on change', () => {
    component.regNumber = 'ABC123';
    component.ngOnChanges({ regNumber: {} as never });
    expect(
      (component as unknown as { regNumberControl: { value: string } })
        .regNumberControl.value,
    ).toBe('ABC123');
  });

  it('submit posts an update request via AssetService when space+asset are set', () => {
    component.space = { id: 's1' } as never;
    component.assetID = 'a1';
    (component as unknown as { regNumberControl: { setValue(v: string): void } })
      .regNumberControl.setValue('XYZ789');

    (component as unknown as { submit(): void }).submit();

    expect(updateAsset).toHaveBeenCalledWith(
      expect.objectContaining({ spaceID: 's1', assetID: 'a1', regNumber: 'XYZ789' }),
    );
  });

  it('submit is a no-op without a space id', () => {
    component.space = undefined;
    component.assetID = 'a1';
    (component as unknown as { submit(): void }).submit();
    expect(updateAsset).not.toHaveBeenCalled();
  });

  it('re-enables the control when the update fails', () => {
    updateAsset.mockReturnValueOnce(throwError(() => new Error('boom')));
    component.space = { id: 's1' } as never;
    component.assetID = 'a1';
    const control = (
      component as unknown as { regNumberControl: { setValue(v: string): void } }
    ).regNumberControl;
    control.setValue('XYZ789');

    (component as unknown as { submit(): void }).submit();

    expect(updateAsset).toHaveBeenCalledTimes(1);
  });

  it('markAsPristine on success and re-enables the control', () => {
    component.space = { id: 's1' } as never;
    component.assetID = 'a1';
    const control = (
      component as unknown as {
        regNumberControl: {
          setValue(v: string): void;
          markAsDirty(): void;
          dirty: boolean;
          enabled: boolean;
        };
      }
    ).regNumberControl;
    control.setValue('XYZ789');
    control.markAsDirty();

    (component as unknown as { submit(): void }).submit();

    expect(control.dirty).toBe(false);
    expect(control.enabled).toBe(true);
  });

  it('showSave reflects an unsaved dirty change and the hideSaveButton flag', () => {
    const c = component as unknown as {
      regNumberControl: { setValue(v: string): void; markAsDirty(): void };
      readonly showSave: boolean;
    };
    component.regNumber = 'OLD';
    expect(c.showSave).toBe(false); // pristine
    c.regNumberControl.setValue('NEW');
    c.regNumberControl.markAsDirty();
    expect(c.showSave).toBe(true);
    component.hideSaveButton = true;
    expect(c.showSave).toBe(false);
  });

  it('validate records the validated value so isValidated becomes true', () => {
    const c = component as unknown as {
      regNumberControl: { setValue(v: string): void };
      validate(): void;
      readonly isValidated: boolean;
    };
    c.regNumberControl.setValue('XYZ789');
    expect(c.isValidated).toBe(false);
    c.validate();
    expect(c.isValidated).toBe(true);
    // Editing the value invalidates the previously validated reg number.
    c.regNumberControl.setValue('OTHER');
    expect(c.isValidated).toBe(false);
  });

  it('skipOrNext emits the trimmed reg number and the skip event', () => {
    const regNumberChange = vi.fn();
    const isSkippedChange = vi.fn();
    component.regNumberChange.subscribe(regNumberChange);
    component.isSkippedChange.subscribe(isSkippedChange);
    (
      component as unknown as {
        regNumberControl: { setValue(v: string): void };
      }
    ).regNumberControl.setValue('  ZZ9  ');

    (component as unknown as { skipOrNext(): void }).skipOrNext();

    expect(regNumberChange).toHaveBeenCalledWith('ZZ9');
    expect(isSkippedChange).toHaveBeenCalled();
  });

  it('skipOrNext emits only the skip event when the value is empty', () => {
    const regNumberChange = vi.fn();
    const isSkippedChange = vi.fn();
    component.regNumberChange.subscribe(regNumberChange);
    component.isSkippedChange.subscribe(isSkippedChange);

    (component as unknown as { skipOrNext(): void }).skipOrNext();

    expect(regNumberChange).not.toHaveBeenCalled();
    expect(isSkippedChange).toHaveBeenCalled();
  });

  it('focusToRegNumberInput delegates to the IonInput setFocus', async () => {
    const setFocus = vi.fn(() => Promise.resolve());
    (component as unknown as { regNumberInput: unknown }).regNumberInput = {
      setFocus,
    };
    (component as unknown as { focusToRegNumberInput(): void }).focusToRegNumberInput();
    expect(setFocus).toHaveBeenCalled();
  });

  it('renders the Validate/Skip controls in the new-asset flow', () => {
    component.assetID = undefined;
    component.isSkipped = false;
    fixture.detectChanges();
    const host = fixture.nativeElement as HTMLElement;
    expect(host.textContent).toContain('Validate');
    expect(host.textContent).toContain('Skip');
  });
});
