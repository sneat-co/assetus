import { TestBed } from '@angular/core/testing';
import { componentTestProviders } from '../../../testing/test-providers';
import { MakeModelCardComponent } from './make-model-card.component';

// Render + logic spec for the ported MakeModelCardComponent.
describe('MakeModelCardComponent', () => {
  let fixture: ReturnType<
    typeof TestBed.createComponent<MakeModelCardComponent>
  >;
  let component: MakeModelCardComponent;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [MakeModelCardComponent],
      providers: componentTestProviders(),
    });
    fixture = TestBed.createComponent(MakeModelCardComponent);
    component = fixture.componentInstance;
  });

  it('creates and renders with a non-empty list of makes', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
    expect(component.makes.length).toBeGreaterThan(0);
  });

  const isKnownMake = (): boolean =>
    (component as unknown as { isKnownMake(): boolean }).isKnownMake();
  const isKnownModel = (): boolean =>
    (component as unknown as { isKnownModel(): boolean }).isKnownModel();
  const onMakeChanged = (e: Event): void =>
    (component as unknown as { onMakeChanged(e: Event): void }).onMakeChanged(e);
  const onModelChanged = (e: Event): void =>
    (component as unknown as { onModelChanged(e: Event): void }).onModelChanged(
      e,
    );

  it('emits the make and clears the model on a known make change', () => {
    component.make = component.makes[0].id;
    component.model = 'A4';
    const makeEmitted: string[] = [];
    const modelEmitted: string[] = [];
    component.makeChange.subscribe((m) => makeEmitted.push(m));
    component.modelChange.subscribe((m) => modelEmitted.push(m));

    onMakeChanged(new Event('change'));

    expect(makeEmitted).toEqual([component.make]);
    expect(component.model).toBe('');
    expect(modelEmitted).toContain('');
    // A known make repopulates the models list from its catalogue entry.
    expect(component.models.length).toBeGreaterThan(0);
  });

  it('empties the models list when the make is unknown', () => {
    component.make = 'NotARealCarMake';
    const makeEmitted: string[] = [];
    component.makeChange.subscribe((m) => makeEmitted.push(m));

    onMakeChanged(new Event('change'));

    expect(component.models).toEqual([]);
    expect(makeEmitted).toEqual(['NotARealCarMake']);
  });

  it('onModelChanged emits the current model', () => {
    component.model = 'A6';
    const modelEmitted: string[] = [];
    component.modelChange.subscribe((m) => modelEmitted.push(m));

    onModelChanged(new Event('change'));

    expect(modelEmitted).toEqual(['A6']);
  });

  it('isKnownMake reflects whether the make is in the catalogue', () => {
    component.make = undefined;
    expect(isKnownMake()).toBe(false);
    component.make = 'NotARealCarMake';
    expect(isKnownMake()).toBe(false);
    component.make = component.makes[0].id;
    expect(isKnownMake()).toBe(true);
  });

  it('isKnownModel matches by id or title case-insensitively', () => {
    component.model = undefined;
    expect(isKnownModel()).toBe(false);
    component.models = [{ id: 'A4', title: 'A4' }];
    component.model = 'a4';
    expect(isKnownModel()).toBe(true);
    component.model = 'Z9';
    expect(isKnownModel()).toBe(false);
  });

});
