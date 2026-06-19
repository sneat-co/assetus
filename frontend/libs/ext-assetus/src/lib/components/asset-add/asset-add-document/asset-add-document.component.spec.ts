import { TestBed } from '@angular/core/testing';
import { WritableSignal } from '@angular/core';
import { EMPTY } from 'rxjs';
import { AssetService } from '@sneat/ext-assetus-components';
import { spacePageTestProviders } from '../../../../testing/test-providers';
import { AssetAddDocumentComponent } from './asset-add-document.component';

// Render + logic spec for the ported AssetAddDocumentComponent, which extends
// the legacy AddAssetBaseComponent.
describe('AssetAddDocumentComponent', () => {
  let fixture: ReturnType<
    typeof TestBed.createComponent<AssetAddDocumentComponent>
  >;
  let component: AssetAddDocumentComponent;
  let createAsset: ReturnType<typeof vi.fn>;

  // Seeds the read-only `space` getter via the protected $spaceRef signal.
  const seedSpace = (id = 's1'): void => {
    (
      component as unknown as { $spaceRef: WritableSignal<{ id: string }> }
    ).$spaceRef.set({ id });
  };

  const documentAssetInput = () =>
    ({
      id: 'a1',
      space: { id: 's1' },
      dbo: { category: 'vehicle', extra: {} },
    }) as never;

  beforeEach(() => {
    createAsset = vi.fn().mockReturnValue(EMPTY);
    TestBed.configureTestingModule({
      imports: [AssetAddDocumentComponent],
      providers: [
        ...spacePageTestProviders(),
        { provide: AssetService, useValue: { createAsset } },
      ],
    });
    fixture = TestBed.createComponent(AssetAddDocumentComponent);
    component = fixture.componentInstance;
  });

  it('creates and renders with a populated document asset', () => {
    component.documentAsset = documentAssetInput();
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('maps the document type and extraType onto the current document asset', () => {
    component.documentAsset = documentAssetInput();
    (component as unknown as { documentType: string }).documentType = 'car';
    (
      component as unknown as { onVehicleTypeChanged(): void }
    ).onVehicleTypeChanged();
    expect(
      (component as unknown as { documentAsset: { dbo: { extraType: string } } })
        .documentAsset.dbo.extraType,
    ).toBe('document');
  });

  it('onVehicleTypeChanged is a no-op when there is no asset dbo', () => {
    component.documentAsset = undefined;
    expect(() =>
      (
        component as unknown as { onVehicleTypeChanged(): void }
      ).onVehicleTypeChanged(),
    ).not.toThrow();
    expect(component.documentAsset).toBeUndefined();
  });

  it('onAssetChanged is a tolerated no-op', () => {
    expect(() =>
      (
        component as unknown as { onAssetChanged(a: unknown): void }
      ).onAssetChanged({} as never),
    ).not.toThrow();
  });

  it('ngOnChanges seeds a draft document when space arrives', () => {
    seedSpace();
    component.ngOnChanges({ space: {} as never });
    expect(component.documentAsset?.dbo?.extraType).toBe('document');
    expect(component.documentAsset?.space?.id).toBe('s1');
  });

  it('ngOnChanges ignores change sets without a space key', () => {
    component.ngOnChanges({ documentAsset: {} as never });
    expect(component.documentAsset).toBeUndefined();
  });

  it('formatDate returns an empty string for a missing value', () => {
    expect(
      (component as unknown as { formatDate(v?: string): string }).formatDate(
        undefined,
      ),
    ).toBe('');
  });

  it('formatDate returns an empty string for an array value', () => {
    expect(
      (
        component as unknown as { formatDate(v?: string | string[]): string }
      ).formatDate(['2026-03-15']),
    ).toBe('');
  });

  it('formatDate formats an ISO date', () => {
    expect(
      (component as unknown as { formatDate(v?: string): string }).formatDate(
        '2026-03-15',
      ),
    ).toContain('2026');
  });

  it('submitDocumentForm throws when there is no document type', () => {
    seedSpace();
    component.documentAsset = documentAssetInput();
    expect(() =>
      (
        component as unknown as { submitDocumentForm(): void }
      ).submitDocumentForm(),
    ).toThrow();
  });

  it('submitDocumentForm throws when there is no asset', () => {
    seedSpace();
    component.documentAsset = undefined;
    (component as unknown as { documentType: string }).documentType = 'car';
    expect(() =>
      (
        component as unknown as { submitDocumentForm(): void }
      ).submitDocumentForm(),
    ).toThrow();
  });

  it('submitDocumentForm creates the asset and sets isSubmitting', () => {
    seedSpace();
    component.documentAsset = documentAssetInput();
    (component as unknown as { documentType: string }).documentType = 'car';
    (component as unknown as { yearOfBuild: string }).yearOfBuild = '2019';
    (
      component as unknown as { submitDocumentForm(): void }
    ).submitDocumentForm();
    expect(createAsset).toHaveBeenCalledOnce();
    expect(
      (component as unknown as { isSubmitting: boolean }).isSubmitting,
    ).toBe(true);
    const request = createAsset.mock.calls[0][1] as {
      asset: { status: string; yearOfBuild: number };
    };
    expect(request.asset.status).toBe('active');
    expect(request.asset.yearOfBuild).toBe(2019);
  });

  it('submitDocumentForm omits yearOfBuild when it is blank', () => {
    seedSpace();
    component.documentAsset = documentAssetInput();
    (component as unknown as { documentType: string }).documentType = 'car';
    (
      component as unknown as { submitDocumentForm(): void }
    ).submitDocumentForm();
    const request = createAsset.mock.calls[0][1] as {
      asset: { yearOfBuild?: number };
    };
    expect(request.asset.yearOfBuild).toBeUndefined();
  });
});
