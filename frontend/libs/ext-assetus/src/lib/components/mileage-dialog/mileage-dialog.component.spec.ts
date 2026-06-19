import { TestBed } from '@angular/core/testing';
import { ModalController } from '@ionic/angular/standalone';
import { AssetService } from '@sneat/ext-assetus-components';
import { of } from 'rxjs';
import { componentTestProviders } from '../../../testing/test-providers';
import { MileAgeDialogComponent } from './mileage-dialog.component';

// Render + logic spec for the ported MileAgeDialogComponent. It injects the
// legacy AssetService, the ModalController and the ErrorLogger; all stubbed.
describe('MileAgeDialogComponent', () => {
  let addVehicleRecord: ReturnType<typeof vi.fn>;
  let dismiss: ReturnType<typeof vi.fn>;
  let fixture: ReturnType<
    typeof TestBed.createComponent<MileAgeDialogComponent>
  >;
  let component: MileAgeDialogComponent;

  beforeEach(() => {
    addVehicleRecord = vi.fn(() => of('rec1'));
    dismiss = vi.fn(() => Promise.resolve(true));
    TestBed.configureTestingModule({
      imports: [MileAgeDialogComponent],
      providers: [
        ...componentTestProviders(),
        { provide: AssetService, useValue: { addVehicleRecord } },
        { provide: ModalController, useValue: { dismiss } },
      ],
    });
    fixture = TestBed.createComponent(MileAgeDialogComponent);
    component = fixture.componentInstance;
  });

  it('creates and renders', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('submit throws when the asset id is missing', () => {
    expect(() => component.submit()).toThrowError(/assetId is not set/);
  });

  it('submit throws when the space id is missing', () => {
    component.asset = { id: 'a1', brief: {} } as never;
    expect(() => component.submit()).toThrowError(/spaceId is not set/);
  });

  it('submit posts a vehicle record and dismisses the modal', () => {
    component.asset = { id: 'a1', brief: {} } as never;
    component.space = { id: 's1', brief: {} } as never;
    component.submit();
    expect(addVehicleRecord).toHaveBeenCalledWith(
      expect.objectContaining({ spaceID: 's1', assetID: 'a1' }),
    );
    expect(dismiss).toHaveBeenCalled();
  });

  it('cancel dismisses the modal', () => {
    component.cancel();
    expect(dismiss).toHaveBeenCalled();
  });
});
